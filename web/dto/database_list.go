package dto

type DatabaseListItem struct {
	ID              uint64 `json:"ID" db:"id"`
	Name            string `json:"name" db:"name"`
	OwnerID         uint64 `json:"ownerID" db:"owner_id"`
	OwnerName       string `json:"ownerName" db:"owner_name"`
	IsTemplate      bool   `json:"isTemplate" db:"is_template"`
	ConnectionLimit int64  `json:"connectionLimit" db:"connection_limit"`
	TablespaceID    uint64 `json:"tablespaceID" db:"tablespace_id"`
	TablespaceName  string `json:"tablespaceName" db:"tablespace_name"`
	Size            string `json:"size" db:"size"`
	Encoding        string `json:"encoding" db:"encoding"`
	Collation       string `json:"collation" db:"collation"`
	CharacterType   string `json:"characterType" db:"character_type"`
}
