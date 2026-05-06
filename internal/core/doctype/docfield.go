package doctype

type DocFieldJSON struct {
	Fieldname string `json:"fieldname"`
	Label     string `json:"label"`
	Fieldtype string `json:"fieldtype"`
	Options   string `json:"options"`

	Required FlexBool `json:"required"`
	Reqd     FlexBool `json:"reqd"`

	Unique      FlexBool `json:"unique"`
	UniqueField FlexBool `json:"unique_field"`

	ReadOnly FlexBool `json:"read_only"`
	Hidden   FlexBool `json:"hidden"`

	InListView       FlexBool `json:"in_list_view"`
	InStandardFilter FlexBool `json:"in_standard_filter"`
	InFilter         FlexBool `json:"in_filter"`
	InGlobalSearch   FlexBool `json:"in_global_search"`
	InPreview        FlexBool `json:"in_preview"`
	InImportTemplate FlexBool `json:"in_import_template"`

	SearchIndex FlexBool `json:"search_index"`
	Sticky      FlexBool `json:"sticky"`

	Idx     FlexInt `json:"idx"`
	Columns FlexInt `json:"columns"`
	Length  FlexInt `json:"length"`

	Precision   string   `json:"precision"`
	NonNegative FlexBool `json:"non_negative"`

	DefaultValue string `json:"default"`
	Description  string `json:"description"`

	DependsOn            string   `json:"depends_on"`
	MandatoryDependsOn   string   `json:"mandatory_depends_on"`
	ReadOnlyDependsOn    string   `json:"read_only_depends_on"`
	Collapsible          FlexBool `json:"collapsible"`
	CollapsibleDependsOn string   `json:"collapsible_depends_on"`
	HideBorder           FlexBool `json:"hide_border"`

	FetchFrom    string   `json:"fetch_from"`
	FetchIfEmpty FlexBool `json:"fetch_if_empty"`

	Bold              FlexBool `json:"bold"`
	Translatable      FlexBool `json:"translatable"`
	AllowInQuickEntry FlexBool `json:"allow_in_quick_entry"`
	ShowOnTimeline    FlexBool `json:"show_on_timeline"`

	PrintHide          FlexBool `json:"print_hide"`
	PrintHideIfNoValue FlexBool `json:"print_hide_if_no_value"`
	ReportHide         FlexBool `json:"report_hide"`
	PrintWidth         string   `json:"print_width"`
	Width              string   `json:"width"`
	MaxHeight          string   `json:"max_height"`

	PermLevel             FlexInt  `json:"permlevel"`
	IgnoreUserPermissions FlexBool `json:"ignore_user_permissions"`
	AllowOnSubmit         FlexBool `json:"allow_on_submit"`
	AllowBulkEdit         FlexBool `json:"allow_bulk_edit"`

	NoCopy                    FlexBool `json:"no_copy"`
	SetOnlyOnce               FlexBool `json:"set_only_once"`
	RememberLastSelectedValue FlexBool `json:"remember_last_selected_value"`
	IgnoreXSSFilter           FlexBool `json:"ignore_xss_filter"`

	Alignment        string `json:"alignment"`
	Placeholder      string `json:"placeholder"`
	DocumentationURL string `json:"documentation_url"`
	OldFieldname     string `json:"oldfieldname"`
	OldFieldtype     string `json:"oldfieldtype"`

	HideDays    FlexBool `json:"hide_days"`
	HideSeconds FlexBool `json:"hide_seconds"`
	SortOptions FlexBool `json:"sort_options"`
	LinkFilters string   `json:"link_filters"`

	MakeAttachmentPublic   FlexBool `json:"make_attachment_public"`
	Mask                   FlexBool `json:"mask"`
	ButtonColor            string   `json:"button_color"`
	ShowDescriptionOnClick FlexBool `json:"show_description_on_click"`

	IsVirtual   FlexBool `json:"is_virtual"`
	NotNullable FlexBool `json:"not_nullable"`

	ValidationRule string `json:"validation_rule"`
}
