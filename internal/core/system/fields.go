package system

// FieldDef is a portable description of a DocType field, independent of any
// specific storage or serialisation format.
type FieldDef struct {
	Fieldname          string
	Label              string
	Fieldtype          string
	Options            string
	Reqd               bool
	Hidden             bool
	ReadOnly           bool
	InListView         bool
	InStandardFilter   bool
	InFilter           bool
	InGlobalSearch     bool
	SearchIndex        bool
	UniqueField        bool
	NoCopy             bool
	SetOnlyOnce        bool
	AllowOnSubmit      bool
	Permlevel          int
	Columns            int
	Length             int
	PrecisionValue     int
	DefaultValue       string
	Description        string
	DependsOn          string
	MandatoryDependsOn string
	ReadOnlyDependsOn  string
	Placeholder        string
	FetchFrom          string
	ValidationRule     string
	Idx                int
}

// StandardFields returns the 7 system-level fields that every DocType carries
// in the database. They are intentionally omitted from JSON metadata files to
// keep the schema definitions clean.
//
// Order: name(1) → owner(2) → creation(3) → modified(4) →
//
//	modified_by(5) → docstatus(6) → idx(7)
func StandardFields() []FieldDef {
	return []FieldDef{
		{
			Fieldname:   "name",
			Label:       "Name",
			Fieldtype:   "Data",
			Hidden:      true,
			ReadOnly:    true,
			SearchIndex: true,
			UniqueField: true,
			NoCopy:      true,
			SetOnlyOnce: true,
			Columns:     6,
			Length:      140,
			Idx:         1,
			Description: "System document name / public identifier.",
		},
		{
			Fieldname:   "owner",
			Label:       "Owner",
			Fieldtype:   "Link",
			Options:     "User",
			Hidden:      true,
			ReadOnly:    true,
			InFilter:    true,
			SearchIndex: true,
			NoCopy:      true,
			Columns:     6,
			Length:      140,
			Idx:         2,
			Description: "User who owns the document.",
		},
		{
			Fieldname:   "creation",
			Label:       "Creation",
			Fieldtype:   "Datetime",
			Hidden:      true,
			ReadOnly:    true,
			InFilter:    true,
			NoCopy:      true,
			Columns:     4,
			Idx:         3,
			Description: "Document creation datetime.",
		},
		{
			Fieldname:   "modified",
			Label:       "Modified",
			Fieldtype:   "Datetime",
			Hidden:      true,
			ReadOnly:    true,
			InFilter:    true,
			SearchIndex: true,
			NoCopy:      true,
			Columns:     4,
			Idx:         4,
			Description: "Last modified datetime.",
		},
		{
			Fieldname:   "modified_by",
			Label:       "Modified By",
			Fieldtype:   "Link",
			Options:     "User",
			Hidden:      true,
			ReadOnly:    true,
			InFilter:    true,
			NoCopy:      true,
			Columns:     6,
			Length:      140,
			Idx:         5,
			Description: "User who last modified the document.",
		},
		{
			Fieldname:   "docstatus",
			Label:       "Docstatus",
			Fieldtype:   "Int",
			Hidden:      true,
			ReadOnly:    true,
			InFilter:    true,
			NoCopy:      true,
			Columns:     3,
			Idx:         6,
			Description: "0 = Draft, 1 = Submitted, 2 = Cancelled.",
		},
		{
			Fieldname:   "idx",
			Label:       "Index",
			Fieldtype:   "Int",
			Hidden:      true,
			ReadOnly:    true,
			NoCopy:      true,
			Columns:     3,
			Idx:         7,
			Description: "Sort order index.",
		},
	}
}

// IsSystemField reports whether fieldname is one of the 7 reserved system fields.
func IsSystemField(fieldname string) bool {
	switch fieldname {
	case "name", "owner", "creation", "modified", "modified_by", "docstatus", "idx":
		return true
	default:
		return false
	}
}

// SystemFieldCount returns the number of standard system fields.
func SystemFieldCount() int {
	return len(StandardFields())
}
