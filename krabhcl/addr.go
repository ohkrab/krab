package krabhcl

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

func AddrFromStrings(s []string) Addr {
	a := Addr{}
	for i := 0; i < len(s); i++ {
		a.Keyword = s[i]
		break
	}
	a.Labels = s[1:]
	return a
}

// Equal compares if other Addr is the same.
func (a Addr) Equal(o Addr) bool {
	if a.Keyword != o.Keyword {
		return false
	}

	if len(a.Labels) != len(o.Labels) {
		return false
	}

	for i := range a.Labels {
		if a.Labels[i] != o.Labels[i] {
			return false
		}
	}

	return true
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

func ParseTraversalToAddr(t hcl.Traversal) (Addr, error) {
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
