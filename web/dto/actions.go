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

type ActionForm struct {
	ExecutionID string
	Namespace   string
	Name        string
	Arguments   []*ActionFormArgument
}

type ActionFormArgument struct {
	Name        string
	Description string
	Value       string
}
