package dto

type TableListItem struct {
	DatabaseID     uint64
	Schema         string  `db:"schema_name"`
	Name           string  `db:"name"`
	OwnerName      string  `db:"owner_name"`
	TablespaceName string  `db:"tablespace_name"`
	RLS            bool    `db:"rls"`
	Size           string  `db:"size"`
	SizePercent    float64 `db:"size_percent"`
}
