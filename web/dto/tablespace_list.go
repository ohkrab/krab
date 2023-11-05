package dto

type TablespaceListItem struct {
	ID        uint64 `json:"ID" db:"id"`
	Name      string `json:"name" db:"name"`
	Size      string `json:"size" db:"size"`
	OwnerID   uint64 `json:"ownerID" db:"owner_id"`
	OwnerName string `json:"ownerName" db:"owner_name"`
	Location  string `json:"location" db:"location"`
}
