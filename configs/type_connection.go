package configs

import (
	"net/url"

	"github.com/ohkrab/krab/addrs"
)

type Connection struct {
	addrs.Addr
	SourceInfo

	Uri url.URL `hcl:"uri"`
}
