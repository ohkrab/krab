package krab

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2"
)

// Addr represents resource reference.
type Addr struct {
	Keyword string
	Labels  []string
}

// String returns full reference name including the keyword.
func (a Addr) String() string {
	return fmt.Sprintf("%s.%s", a.Keyword, a.OnlyRefNames())
}

// Absolute returns keyword and labels as a single slice.
func (a Addr) Absolute() []string {
	return append([]string{a.Keyword}, a.Labels...)
}

// OnlyRefNames returns reference name without the keyword.
func (a Addr) OnlyRefNames() string {
	return strings.Join(a.Labels, ".")
}

func parseTraversalToAddr(t hcl.Traversal) (Addr, error) {
	addr := Addr{
		Keyword: t.RootName(),
		Labels:  make([]string, 0),
	}

	for _, rel := range t[1:] {
		attr, ok := rel.(hcl.TraverseAttr)
		if !ok {
			return Addr{}, fmt.Errorf("Failed to parse hcl.Traversal to Addr")
		}

		addr.Labels = append(addr.Labels, attr.Name)
	}

	return addr, nil
}
