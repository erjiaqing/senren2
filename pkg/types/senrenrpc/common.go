package senrenrpc

import (
	"github.com/erjiaqing/senren2/pkg/db"
)

type SuccessError struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

type SuccessErrorOnly struct {
	SuccessError `json:"result"`
}

type Domain string

type HasDomain interface {
	ConvertDomain()
	GetDomain() string
	SetDomain(string)
}

func (d *Domain) ConvertDomain() {
	if *d == "" {
		*d = "0000000000000000"
		return
	}

	if len(*d) == 16 {
		return
	}

	dom := db.DB.QueryRow("SELECT `uid` FROM `domain_info` WHERE short_name = ?", d)
	if err := dom.Scan(d); err != nil {
		panic(err)
		// not exist domain id
	}
}

func (d *Domain) GetDomain() string {
	return string(*d)
}

func (d *Domain) SetDomain(s string) {
	*d = Domain(s)
	d.ConvertDomain()
}
