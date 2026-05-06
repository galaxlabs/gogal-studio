package doctype

type DocTypeLinkJSON struct {
	Doctype        string   `json:"doctype"`
	LinkDoctype    string   `json:"link_doctype"`
	LinkFieldname  string   `json:"link_fieldname"`
	ParentDoctype  string   `json:"parent_doctype"`
	TableFieldname string   `json:"table_fieldname"`
	GroupName      string   `json:"group_name"`
	Group          string   `json:"group"`
	Hidden         FlexBool `json:"hidden"`
	IsChildTable   FlexBool `json:"is_child_table"`
	Custom         FlexBool `json:"custom"`
	Idx            FlexInt  `json:"idx"`
}
