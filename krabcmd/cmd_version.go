package krabcmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/ohkrab/krab/krab"
)

// CmdVersion returns version information.
type CmdVersion struct {
}

type ResponseVersion struct {
	Name  string `json:"name"`
	Build string `json:"build"`
}

func (c *CmdVersion) Command() []string {
	return []string{"version"}
}

func (c *CmdVersion) HttpEndpoint() (method string, path string) {
	method = http.MethodGet
	path = fmt.Sprint("/", strings.Join(c.Command(), "/"))
	return
}

func (c *CmdVersion) Do(ctx context.Context, w io.Writer) error {
	return json.NewEncoder(w).Encode(ResponseVersion{
		Name:  fmt.Sprint(krab.InfoName, " ", krab.InfoVersion),
		Build: fmt.Sprint("Build ", krab.InfoCommit, " ", krab.InfoBuildDate),
	})
}
