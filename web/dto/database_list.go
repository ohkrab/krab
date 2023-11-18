package dto

type DatabaseListItem struct {
	ID              uint64  `db:"id"`
	Name            string  `db:"name"`
	OwnerID         uint64  `db:"owner_id"`
	OwnerName       string  `db:"owner_name"`
	IsTemplate      bool    `db:"is_template"`
	ConnectionLimit int64   `db:"connection_limit"`
	TablespaceID    uint64  `db:"tablespace_id"`
	TablespaceName  string  `db:"tablespace_name"`
	Size            string  `db:"size"`
	SizePercent     float64 `db:"size_percent"`
	Encoding        string  `db:"encoding"`
	Collation       string  `db:"collation"`
	CharacterType   string  `db:"character_type"`
	CanConnect      bool
}
