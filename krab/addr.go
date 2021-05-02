package krab

import "fmt"

type Addr struct {
	Keyword string
	Type    string `hcl:"type,label"`
	Name    string `hcl:"name,label"`
}

func (a *Addr) String() string {
	if a.Type == "" {
		return fmt.Sprintf("%s.%s", a.Keyword, a.Name)
	} else {
		return fmt.Sprintf("%s.%s.%s", a.Keyword, a.Type, a.Name)
	}
}
