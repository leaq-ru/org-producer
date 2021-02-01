package producer

import (
	"archive/zip"
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/nnqq/scr-org-producer/protocol"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func (p Producer) Do(ctx context.Context) error {
	p.logger.Debug().Msg("get nalog.ru HTML start...")
	res, err := http.Get("https://www.nalog.ru/opendata/7707329152-rsmp/")
	if err != nil {
		return err
	}
	p.logger.Debug().Msg("get nalog.ru HTML done")
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return err
	}

	var dlLink string
	doc.Find("td").EachWithBreak(func(_ int, s *goquery.Selection) bool {
		if s.Text() == "Гиперссылка (URL) на набор" {
			attr, ok := s.Next().Children().Attr("href")
			if ok && attr != "" {
				dlLink = attr
				return false
			}
		}
		return true
	})

	if dlLink == "" {
		return errors.New("zip download link not found")
	}

	p.logger.Debug().Str("dlLink", dlLink).Msg("got zip download url")

	p.logger.Debug().Msg("get zip reader start...")
	resZipXml, err := http.Get(dlLink)
	if err != nil {
		return err
	}
	p.logger.Debug().Msg("get zip reader done")
	defer resZipXml.Body.Close()

	const filename = "archive.zip"

	p.logger.Debug().Msg("loading zip to file start...")
	outFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resZipXml.Body)
	if err != nil {
		return err
	}
	p.logger.Debug().Msg("loading zip to file done")

	p.logger.Debug().Msg("zip.OpenReader start...")
	reader, err := zip.OpenReader(filename)
	if err != nil {
		return err
	}
	p.logger.Debug().Msg("zip.OpenReader done")

	totalProduced := 0
	go func() {
		for ctx.Err() == nil {
			p.logger.Debug().Int("totalProduced", totalProduced).Send()
			time.Sleep(10 * time.Second)
		}
	}()

	lastLoopIndex, err := p.stateModel.GetLastLoopIndex(ctx)
	if err != nil {
		return err
	}
	p.logger.Debug().Uint32("lastLoopIndex", lastLoopIndex).Msg("value from GetLastLoopIndex")

	go func() {
		for ctx.Err() == nil {
			e := p.stateModel.Commit(ctx, lastLoopIndex)
			if e != nil {
				p.logger.Error().Err(e).Send()
			}
			time.Sleep(time.Second)
		}
	}()

	lastIndex := lastLoopIndex
	if len(reader.File) < int(lastLoopIndex) {
		lastIndex = uint32(len(reader.File))
	}

	for i, f := range reader.File[lastIndex:] {
		xmlItem, errOpen := f.Open()
		if errOpen != nil {
			return errOpen
		}

		bytes, errReadAll := ioutil.ReadAll(xmlItem)
		if errReadAll != nil {
			return errReadAll
		}

		var body xmlFile
		err = xml.Unmarshal(bytes, &body)
		if err != nil {
			return err
		}

		for _, d := range body.Docs {
			msg := protocol.OrgMessage{
				EmployeeCount: d.EmployeeCount,
				OkvedOsn: protocol.Okved{
					Name: d.Okved.OkvedOsn.Name,
					Code: d.Okved.OkvedOsn.Code,
					Ver:  d.Okved.OkvedOsn.Ver,
				},
			}

			if d.Ind.INN != "" {
				msg.INN = d.Ind.INN
			}
			if d.Legal.INN != "" {
				msg.INN = d.Legal.INN
			}

			for _, od := range d.Okved.OkvedDop {
				msg.OkvedDop = append(msg.OkvedDop, protocol.Okved{
					Name: od.Name,
					Code: od.Code,
					Ver:  od.Ver,
				})
			}

			bytesMsg, errMarshal := json.Marshal(msg)
			if errMarshal != nil {
				return errMarshal
			}

			errPublish := p.stanConn.Publish(subject, bytesMsg)
			if errPublish != nil {
				return errPublish
			}
			totalProduced += 1
			lastLoopIndex = uint32(i)
		}
	}
	return nil
}

const subject = "org"

type xmlFile struct {
	XMLName xml.Name `xml:"Файл"`
	Docs    []xmlDoc `xml:"Документ"`
}

type xmlDoc struct {
	EmployeeCount string `xml:"ССЧР,attr"`
	Ind           ind    `xml:"ИПВклМСП"`
	Legal         legal  `xml:"ОргВклМСП"`
	Okved         okved  `xml:"СвОКВЭД"`
}

type ind struct {
	INN string `xml:"ИННФЛ,attr"`
}

type legal struct {
	INN string `xml:"ИННЮЛ,attr"`
}

type okved struct {
	OkvedOsn okvedBody   `xml:"СвОКВЭДОсн"`
	OkvedDop []okvedBody `xml:"СвОКВЭДДоп"`
}

type okvedBody struct {
	Code string `xml:"КодОКВЭД,attr"`
	Name string `xml:"НаимОКВЭД,attr"`
	Ver  string `xml:"ВерсОКВЭД,attr"`
}
