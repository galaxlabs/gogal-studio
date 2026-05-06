package doctype

type DocPermJSON struct {
	Doctype   string   `json:"doctype"`
	Role      string   `json:"role"`
	PermLevel FlexInt  `json:"permlevel"`
	IfOwner   FlexBool `json:"if_owner"`

	Read       FlexBool `json:"read"`
	Write      FlexBool `json:"write"`
	Create     FlexBool `json:"create"`
	CreatePerm FlexBool `json:"create_perm"`
	Delete     FlexBool `json:"delete"`
	DeletePerm FlexBool `json:"delete_perm"`

	Submit     FlexBool `json:"submit"`
	Cancel     FlexBool `json:"cancel"`
	Amend      FlexBool `json:"amend"`
	Report     FlexBool `json:"report"`
	Export     FlexBool `json:"export"`
	Import     FlexBool `json:"import"`
	Share      FlexBool `json:"share"`
	Print      FlexBool `json:"print"`
	Email      FlexBool `json:"email"`
	Select     FlexBool `json:"select"`
	SelectPerm FlexBool `json:"select_perm"`
	Mask       FlexBool `json:"mask"`

	Idx FlexInt `json:"idx"`
}
