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
	Convert()
}

func (d *Domain) Convert() {
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
