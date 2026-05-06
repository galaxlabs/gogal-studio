package doctype

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type DocTypeJSON struct {
	Name          string `json:"name"`
	Module        string `json:"module"`
	Label         string `json:"label"`
	TableName     string `json:"table_name"`
	IsCore        bool   `json:"is_core"`
	IsSingle      bool   `json:"is_single"`
	IsSubmittable bool   `json:"is_submittable"`
	IsTree        bool   `json:"is_tree"`
	IsChildTable  bool   `json:"is_child_table"`

	AllowImport  bool `json:"allow_import"`
	AllowExport  bool `json:"allow_export"`
	TrackChanges bool `json:"track_changes"`
	EditableGrid bool `json:"editable_grid"`
	QuickEntry   bool `json:"quick_entry"`

	Controller string `json:"controller"`
	Route      string `json:"route"`
	NamingRule string `json:"naming_rule"`
	TitleField string `json:"title_field"`
	ImageField string `json:"image_field"`
	SortField  string `json:"sort_field"`
	SortOrder  string `json:"sort_order"`

	Fields      []DocFieldJSON      `json:"fields"`
	Actions     []DocTypeActionJSON `json:"actions"`
	Links       []DocTypeLinkJSON   `json:"links"`
	Permissions []DocPermJSON       `json:"permissions"`
	States      []DocTypeStateJSON  `json:"states"`
}

type FlexBool bool

func (b *FlexBool) UnmarshalJSON(data []byte) error {
	var boolValue bool
	if err := json.Unmarshal(data, &boolValue); err == nil {
		*b = FlexBool(boolValue)
		return nil
	}

	var intValue int
	if err := json.Unmarshal(data, &intValue); err == nil {
		*b = FlexBool(intValue != 0)
		return nil
	}

	var stringValue string
	if err := json.Unmarshal(data, &stringValue); err == nil {
		switch stringValue {
		case "1", "true", "True", "TRUE", "yes", "Yes", "YES":
			*b = true
			return nil
		case "0", "false", "False", "FALSE", "no", "No", "NO", "":
			*b = false
			return nil
		default:
			return fmt.Errorf("invalid boolean string: %s", stringValue)
		}
	}

	return fmt.Errorf("invalid boolean value: %s", string(data))
}

func (b FlexBool) Bool() bool {
	return bool(b)
}

type FlexInt int

func (i *FlexInt) UnmarshalJSON(data []byte) error {
	var intValue int
	if err := json.Unmarshal(data, &intValue); err == nil {
		*i = FlexInt(intValue)
		return nil
	}

	var stringValue string
	if err := json.Unmarshal(data, &stringValue); err == nil {
		if stringValue == "" {
			*i = 0
			return nil
		}

		parsed, err := strconv.Atoi(stringValue)
		if err != nil {
			return err
		}

		*i = FlexInt(parsed)
		return nil
	}

	return fmt.Errorf("invalid integer value: %s", string(data))
}

func (i FlexInt) Int() int {
	return int(i)
}
