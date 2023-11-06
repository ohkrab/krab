package dto

type ActionListItem struct {
	Namespace   string
	Name        string
	Description string
	Transaction bool
	Arguments   []*ActionListItemArgument
}

type ActionListItemArgument struct {
	Name        string
	Type        string
	Description string
}
