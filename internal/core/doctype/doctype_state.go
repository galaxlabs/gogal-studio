package doctype

type DocTypeStateJSON struct {
	Doctype string   `json:"doctype"`
	Title   string   `json:"title"`
	Color   string   `json:"color"`
	Custom  FlexBool `json:"custom"`
	Idx     FlexInt  `json:"idx"`
}
