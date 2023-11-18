package views

const (
	NavNone     = iota
	NavDatabase = iota + 1
)

type LayoutInfo struct {
	Blank    bool
	Nav      int
	Database string
	Footer   string
}
