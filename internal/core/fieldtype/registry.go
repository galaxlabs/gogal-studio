package fieldtype

import "fmt"

type OptionMode string

const (
	OptionNone         OptionMode = "none"
	OptionLines        OptionMode = "lines"
	OptionDocType      OptionMode = "doctype"
	OptionChildDocType OptionMode = "child_doctype"
)

type ControlType string

const (
	ControlInput    ControlType = "input"
	ControlTextarea ControlType = "textarea"
	ControlEditor   ControlType = "editor"
	ControlSelect   ControlType = "select"
	ControlCheckbox ControlType = "checkbox"
	ControlLink     ControlType = "link"
	ControlTable    ControlType = "table"
	ControlFile     ControlType = "file"
	ControlImage    ControlType = "image"
	ControlCode     ControlType = "code"
	ControlLayout   ControlType = "layout"
)

type FieldTypeDef struct {
	Name             string      `json:"name"`
	Control          ControlType `json:"control"`
	SQLType          string      `json:"sql_type"`
	OptionMode       OptionMode  `json:"option_mode"`
	RequiresOptions  bool        `json:"requires_options"`
	IsLayout         bool        `json:"is_layout"`
	IsNumeric        bool        `json:"is_numeric"`
	IsDateTime       bool        `json:"is_datetime"`
	IsText           bool        `json:"is_text"`
	IsAttach         bool        `json:"is_attach"`
	DefaultColumns   int         `json:"default_columns"`
	DefaultLength    int         `json:"default_length"`
	DefaultPrecision int         `json:"default_precision"`
	Description      string      `json:"description"`
}

var registry = map[string]FieldTypeDef{
	"Data": {
		Name:           "Data",
		Control:        ControlInput,
		SQLType:        "TEXT",
		OptionMode:     OptionNone,
		DefaultColumns: 6,
		DefaultLength:  140,
		IsText:         true,
		Description:    "Short text input.",
	},
	"Password": {
		Name:           "Password",
		Control:        ControlInput,
		SQLType:        "TEXT",
		OptionMode:     OptionNone,
		DefaultColumns: 6,
		DefaultLength:  255,
		IsText:         true,
		Description:    "Password/secret field.",
	},
	"Text": {
		Name:           "Text",
		Control:        ControlTextarea,
		SQLType:        "TEXT",
		OptionMode:     OptionNone,
		DefaultColumns: 12,
		IsText:         true,
		Description:    "Multi-line text.",
	},
	"Small Text": {
		Name:           "Small Text",
		Control:        ControlTextarea,
		SQLType:        "TEXT",
		OptionMode:     OptionNone,
		DefaultColumns: 12,
		IsText:         true,
		Description:    "Small multi-line text.",
	},
	"Long Text": {
		Name:           "Long Text",
		Control:        ControlTextarea,
		SQLType:        "TEXT",
		OptionMode:     OptionNone,
		DefaultColumns: 12,
		IsText:         true,
		Description:    "Long multi-line text.",
	},
	"Text Editor": {
		Name:           "Text Editor",
		Control:        ControlEditor,
		SQLType:        "TEXT",
		OptionMode:     OptionNone,
		DefaultColumns: 12,
		IsText:         true,
		Description:    "Rich text editor field.",
	},
	"Code": {
		Name:           "Code",
		Control:        ControlCode,
		SQLType:        "TEXT",
		OptionMode:     OptionNone,
		DefaultColumns: 12,
		IsText:         true,
		Description:    "Code editor field.",
	},
	"JSON": {
		Name:           "JSON",
		Control:        ControlCode,
		SQLType:        "JSONB",
		OptionMode:     OptionNone,
		DefaultColumns: 12,
		Description:    "JSON data field.",
	},
	"Int": {
		Name:           "Int",
		Control:        ControlInput,
		SQLType:        "BIGINT",
		OptionMode:     OptionNone,
		DefaultColumns: 4,
		IsNumeric:      true,
		Description:    "Integer number.",
	},
	"Float": {
		Name:             "Float",
		Control:          ControlInput,
		SQLType:          "DOUBLE PRECISION",
		OptionMode:       OptionNone,
		DefaultColumns:   4,
		DefaultPrecision: 6,
		IsNumeric:        true,
		Description:      "Floating point number.",
	},
	"Currency": {
		Name:             "Currency",
		Control:          ControlInput,
		SQLType:          "NUMERIC(18,6)",
		OptionMode:       OptionNone,
		DefaultColumns:   4,
		DefaultPrecision: 6,
		IsNumeric:        true,
		Description:      "Currency/amount field.",
	},
	"Percent": {
		Name:             "Percent",
		Control:          ControlInput,
		SQLType:          "NUMERIC(18,6)",
		OptionMode:       OptionNone,
		DefaultColumns:   4,
		DefaultPrecision: 6,
		IsNumeric:        true,
		Description:      "Percentage field.",
	},
	"Check": {
		Name:           "Check",
		Control:        ControlCheckbox,
		SQLType:        "BOOLEAN",
		OptionMode:     OptionNone,
		DefaultColumns: 3,
		Description:    "Boolean checkbox.",
	},
	"Date": {
		Name:           "Date",
		Control:        ControlInput,
		SQLType:        "DATE",
		OptionMode:     OptionNone,
		DefaultColumns: 4,
		IsDateTime:     true,
		Description:    "Date field.",
	},
	"Datetime": {
		Name:           "Datetime",
		Control:        ControlInput,
		SQLType:        "TIMESTAMPTZ",
		OptionMode:     OptionNone,
		DefaultColumns: 4,
		IsDateTime:     true,
		Description:    "Date and time field.",
	},
	"Time": {
		Name:           "Time",
		Control:        ControlInput,
		SQLType:        "TIME",
		OptionMode:     OptionNone,
		DefaultColumns: 4,
		IsDateTime:     true,
		Description:    "Time field.",
	},
	"Select": {
		Name:            "Select",
		Control:         ControlSelect,
		SQLType:         "TEXT",
		OptionMode:      OptionLines,
		RequiresOptions: true,
		DefaultColumns:  6,
		DefaultLength:   140,
		IsText:          true,
		Description:     "Select field with newline-separated options.",
	},
	"Link": {
		Name:            "Link",
		Control:         ControlLink,
		SQLType:         "TEXT",
		OptionMode:      OptionDocType,
		RequiresOptions: true,
		DefaultColumns:  6,
		DefaultLength:   140,
		IsText:          true,
		Description:     "Link to another DocType.",
	},
	"Dynamic Link": {
		Name:            "Dynamic Link",
		Control:         ControlLink,
		SQLType:         "TEXT",
		OptionMode:      OptionDocType,
		RequiresOptions: true,
		DefaultColumns:  6,
		DefaultLength:   140,
		IsText:          true,
		Description:     "Dynamic link where target DocType comes from another field.",
	},
	"Table": {
		Name:            "Table",
		Control:         ControlTable,
		SQLType:         "",
		OptionMode:      OptionChildDocType,
		RequiresOptions: true,
		DefaultColumns:  12,
		Description:     "Child table field.",
	},
	"Attach": {
		Name:           "Attach",
		Control:        ControlFile,
		SQLType:        "TEXT",
		OptionMode:     OptionNone,
		DefaultColumns: 6,
		IsAttach:       true,
		Description:    "File attachment.",
	},
	"Attach Image": {
		Name:           "Attach Image",
		Control:        ControlImage,
		SQLType:        "TEXT",
		OptionMode:     OptionNone,
		DefaultColumns: 6,
		IsAttach:       true,
		Description:    "Image attachment.",
	},
	"Section": {
		Name:           "Section",
		Control:        ControlLayout,
		SQLType:        "",
		OptionMode:     OptionNone,
		IsLayout:       true,
		DefaultColumns: 12,
		Description:    "Section break/layout field.",
	},
	"Column": {
		Name:           "Column",
		Control:        ControlLayout,
		SQLType:        "",
		OptionMode:     OptionNone,
		IsLayout:       true,
		DefaultColumns: 6,
		Description:    "Column break/layout field.",
	},
}

func Get(name string) (FieldTypeDef, bool) {
	def, ok := registry[name]
	return def, ok
}

func MustGet(name string) (FieldTypeDef, error) {
	def, ok := Get(name)
	if !ok {
		return FieldTypeDef{}, fmt.Errorf("unknown field type: %s", name)
	}

	return def, nil
}

func Names() []string {
	names := make([]string, 0, len(registry))

	for name := range registry {
		names = append(names, name)
	}

	return names
}

func All() []FieldTypeDef {
	defs := make([]FieldTypeDef, 0, len(registry))

	for _, def := range registry {
		defs = append(defs, def)
	}

	return defs
}

func IsValid(name string) bool {
	_, ok := registry[name]
	return ok
}

func RequiresOptions(name string) bool {
	def, ok := registry[name]
	if !ok {
		return false
	}

	return def.RequiresOptions
}

func SQLType(name string) string {
	def, ok := registry[name]
	if !ok {
		return "TEXT"
	}

	return def.SQLType
}

func DefaultColumns(name string) int {
	def, ok := registry[name]
	if !ok || def.DefaultColumns == 0 {
		return 6
	}

	return def.DefaultColumns
}

func DefaultLength(name string) int {
	def, ok := registry[name]
	if !ok {
		return 0
	}

	return def.DefaultLength
}
