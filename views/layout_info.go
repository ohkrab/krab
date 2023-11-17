package views

const (
	NavNone     = iota
	NavDatabase = iota + 1
)

type LayoutInfo struct {
	Nav int
	OID int
}
