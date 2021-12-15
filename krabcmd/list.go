package krabcmd

type Registry []Action

func (r *Registry) Register(a Action) {
	*r = append(*r, a)
}
