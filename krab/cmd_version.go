package krab

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ohkrab/krab/krabhcl"
)

// CmdVersion returns version information.
type CmdVersion struct{}

// ResponseVersion json
type ResponseVersion struct {
	Name  string `json:"name"`
	Build string `json:"build"`
}

func (c *CmdVersion) Addr() krabhcl.Addr { return krabhcl.Addr{Keyword: "version", Labels: []string{}} }

func (c *CmdVersion) Name() []string { return []string{"version"} }

func (c *CmdVersion) HttpMethod() string { return http.MethodGet }

func (c *CmdVersion) Do(ctx context.Context, o CmdOpts) (any, error) {
	return ResponseVersion{
		Name:  fmt.Sprint(InfoName, " ", InfoVersion),
		Build: fmt.Sprint("Build ", InfoCommit, " ", InfoBuildDate),
	}, nil
}
