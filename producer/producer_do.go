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
	res, err := http.Get("https://www.nalog.ru/opendata/7707329152-rsmp/")
	if err != nil {
		return err
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return err
	}

	var dlLink string
	for _, node := range doc.Find("td").Nodes {
		if node.Data == "Гиперссылка (URL) на набор" {
			for _, attr := range node.NextSibling.FirstChild.Attr {
				if attr.Key == "href" {
					dlLink = attr.Val
					break
				}
			}
		}
	}

	if dlLink == "" {
		return errors.New("xml download link not found")
	}

	resZipXml, err := http.Get(dlLink)
	if err != nil {
		return err
	}
	defer resZipXml.Body.Close()

	const filename = "archive.zip"

	outFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, resZipXml.Body)
	if err != nil {
		return err
	}

	reader, err := zip.OpenReader(filename)
	if err != nil {
		return err
	}

	totalProduced := 0
	go func() {
		for ctx.Err() != nil {
			p.logger.Debug().Int("totalProduced", totalProduced).Send()
			time.Sleep(10 * time.Second)
		}
	}()

	for _, f := range reader.File {
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
