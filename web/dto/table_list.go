package dto

type TableListItem struct {
	DatabaseName   string
	Schema         string  `db:"schema_name"`
	Name           string  `db:"name"`
	OwnerName      string  `db:"owner_name"`
	TablespaceName string  `db:"tablespace_name"`
	RLS            bool    `db:"rls"`
	Internal       bool    `db:"internal"`
	Size           string  `db:"size"`
	SizePercent    float64 `db:"size_percent"`
	EstimatedRows  int64   `db:"estimated_rows"`
}
