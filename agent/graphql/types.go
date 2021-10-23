package schema

type Migration struct {
	RefName     string `json:"refName"`
	Version     string `json:"version"`
	Transaction bool   `json:"transaction"`
	Up          string `json:"up"`
	Down        string `json:"down"`
}

type MigrationSet struct {
	RefName    string      `json:"refName"`
	Schema     string      `json:"schema"`
	Migrations []Migration `json:"migrations,omitempty"`
}
