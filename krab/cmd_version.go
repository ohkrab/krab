package krab

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// CmdVersion returns version information.
type CmdVersion struct{}

// ResponseVersion json
type ResponseVersion struct {
	Name  string `json:"name"`
	Build string `json:"build"`
}

func (c *CmdVersion) Name() []string { return []string{"version"} }

func (c *CmdVersion) HttpMethod() string { return http.MethodGet }

func (c *CmdVersion) Do(ctx context.Context, w io.Writer) error {
	return json.NewEncoder(w).Encode(ResponseVersion{
		Name:  fmt.Sprint(InfoName, " ", InfoVersion),
		Build: fmt.Sprint("Build ", InfoCommit, " ", InfoBuildDate),
	})
}
