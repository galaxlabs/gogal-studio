package doctype

type ModuleDefJSON struct {
	Name        string   `json:"name"`
	Label       string   `json:"label"`
	Description string   `json:"description"`
	Enabled     FlexBool `json:"enabled"`
	IsCore      FlexBool `json:"is_core"`
	Version     string   `json:"version"`
	Icon        string   `json:"icon"`
	Color       string   `json:"color"`
	SortOrder   FlexInt  `json:"sort_order"`
}
