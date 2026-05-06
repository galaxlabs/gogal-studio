package doctype

type ModuleJSON struct {
	Name        string   `json:"name"`
	Label       string   `json:"label"`
	Description string   `json:"description"`
	Version     string   `json:"version"`
	IsCore      FlexBool `json:"is_core"`
	Enabled     FlexBool `json:"enabled"`
	DoctypePath string   `json:"doctype_path"`
}
