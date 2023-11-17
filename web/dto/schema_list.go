package dto

type SchemaListItem struct {
	ID         uint64 `json:"ID" db:"id"`
	DatabaseID uint64
	Name       string `json:"name" db:"name"`
	OwnerID    uint64 `json:"ownerID" db:"owner_id"`
	OwnerName  string `json:"ownerName" db:"owner_name"`
}
