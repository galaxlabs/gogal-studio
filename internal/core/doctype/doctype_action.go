package doctype

type DocTypeActionJSON struct {
	Doctype     string   `json:"doctype"`
	ActionName  string   `json:"action_name"`
	Label       string   `json:"label"`
	GroupName   string   `json:"group_name"`
	Group       string   `json:"group"`
	ActionType  string   `json:"action_type"`
	Action      string   `json:"action"`
	Handler     string   `json:"handler"`
	Route       string   `json:"route"`
	Method      string   `json:"method"`
	Permission  string   `json:"permission"`
	VisibleWhen string   `json:"visible_when"`
	Hidden      FlexBool `json:"hidden"`
	Custom      FlexBool `json:"custom"`
	Enabled     FlexBool `json:"enabled"`
	Idx         FlexInt  `json:"idx"`
}
