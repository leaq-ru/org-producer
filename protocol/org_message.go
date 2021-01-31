package protocol

type OrgMessage struct {
	INN           string  `json:"i"`
	EmployeeCount string  `json:"e"`
	OkvedOsn      Okved   `json:"o"`
	OkvedDop      []Okved `json:"od"`
}

type Okved struct {
	Name string `json:"n"`
	Code string `json:"c"`
	Ver  string `json:"v"`
}
