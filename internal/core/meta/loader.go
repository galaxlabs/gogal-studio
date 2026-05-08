package meta

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DocTypeMeta struct {
	Name          string      `json:"name"`
	Label         string      `json:"label"`
	Module        string      `json:"module"`
	AppName       string      `json:"app_name"`
	TableName     string      `json:"table_name"`
	Autoname      string      `json:"autoname"`
	NamingRule    string      `json:"naming_rule"`
	TitleField    string      `json:"title_field"`
	SortField     string      `json:"sort_field"`
	SortOrder     string      `json:"sort_order"`
	DocumentType  string      `json:"document_type"`
	IsSingle      bool        `json:"is_single"`
	IsSubmittable bool        `json:"is_submittable"`
	IsChildTable  bool        `json:"is_child_table"`
	IsTree        bool        `json:"is_tree"`
	Fields        []FieldMeta `json:"fields"`
	Permissions   []PermMeta  `json:"permissions"`
}

type FieldMeta struct {
	Name               string `json:"name"`
	Fieldname          string `json:"fieldname"`
	Label              string `json:"label"`
	Fieldtype          string `json:"fieldtype"`
	Options            string `json:"options"`
	Reqd               bool   `json:"reqd"`
	Hidden             bool   `json:"hidden"`
	ReadOnly           bool   `json:"read_only"`
	InListView         bool   `json:"in_list_view"`
	InStandardFilter   bool   `json:"in_standard_filter"`
	SearchIndex        bool   `json:"search_index"`
	UniqueField        bool   `json:"unique_field"`
	NoCopy             bool   `json:"no_copy"`
	SetOnlyOnce        bool   `json:"set_only_once"`
	AllowOnSubmit      bool   `json:"allow_on_submit"`
	Permlevel          int    `json:"permlevel"`
	Columns            int    `json:"columns"`
	Length             int    `json:"length"`
	PrecisionValue     int    `json:"precision_value"`
	DefaultValue       string `json:"default_value"`
	Description        string `json:"description"`
	DependsOn          string `json:"depends_on"`
	MandatoryDependsOn string `json:"mandatory_depends_on"`
	ReadOnlyDependsOn  string `json:"read_only_depends_on"`
	Placeholder        string `json:"placeholder"`
	FetchFrom          string `json:"fetch_from"`
	ValidationRule     string `json:"validation_rule"`
	Idx                int    `json:"idx"`
}

type PermMeta struct {
	Name      string `json:"name"`
	Role      string `json:"role"`
	Permlevel int    `json:"permlevel"`
	Read      bool   `json:"read"`
	Write     bool   `json:"write"`
	Create    bool   `json:"create"`
	Delete    bool   `json:"delete"`
	Submit    bool   `json:"submit"`
	Cancel    bool   `json:"cancel"`
	Amend     bool   `json:"amend"`
	Print     bool   `json:"print"`
	Email     bool   `json:"email"`
	Export    bool   `json:"export"`
	Import    bool   `json:"import"`
	Share     bool   `json:"share"`
	Report    bool   `json:"report"`
	Idx       int    `json:"idx"`
}

type Loader struct {
	DB    *pgxpool.Pool
	mu    sync.RWMutex
	cache map[string]DocTypeMeta
}

func NewLoader(db *pgxpool.Pool) *Loader {
	return &Loader{
		DB:    db,
		cache: map[string]DocTypeMeta{},
	}
}

func (l *Loader) GetDocTypeMeta(ctx context.Context, doctypeName string) (DocTypeMeta, error) {
	l.mu.RLock()
	cached, ok := l.cache[doctypeName]
	l.mu.RUnlock()

	if ok {
		return cached, nil
	}

	m, err := l.LoadDocTypeMeta(ctx, doctypeName)
	if err != nil {
		return DocTypeMeta{}, err
	}

	l.mu.Lock()
	l.cache[doctypeName] = m
	l.mu.Unlock()

	return m, nil
}

func (l *Loader) ClearCache() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.cache = map[string]DocTypeMeta{}
}

func (l *Loader) LoadDocTypeMeta(ctx context.Context, doctypeName string) (DocTypeMeta, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	doc, err := l.loadDocType(ctx, doctypeName)
	if err != nil {
		return DocTypeMeta{}, err
	}

	fields, err := l.loadFields(ctx, doctypeName)
	if err != nil {
		return DocTypeMeta{}, err
	}

	perms, err := l.loadPermissions(ctx, doctypeName)
	if err != nil {
		return DocTypeMeta{}, err
	}

	doc.Fields = fields
	doc.Permissions = perms

	return doc, nil
}

func (l *Loader) loadDocType(ctx context.Context, doctypeName string) (DocTypeMeta, error) {
	var m DocTypeMeta

	err := l.DB.QueryRow(ctx, `
		SELECT
			name,
			COALESCE(label, ''),
			COALESCE(module, ''),
			COALESCE(app_name, ''),
			COALESCE(table_name, ''),
			COALESCE(autoname, ''),
			COALESCE(naming_rule, ''),
			COALESCE(title_field, ''),
			COALESCE(sort_field, ''),
			COALESCE(sort_order, ''),
			COALESCE(document_type, ''),
			is_single,
			is_submittable,
			is_child_table,
			is_tree
		FROM "tabDocType"
		WHERE name = $1
	`, doctypeName).Scan(
		&m.Name,
		&m.Label,
		&m.Module,
		&m.AppName,
		&m.TableName,
		&m.Autoname,
		&m.NamingRule,
		&m.TitleField,
		&m.SortField,
		&m.SortOrder,
		&m.DocumentType,
		&m.IsSingle,
		&m.IsSubmittable,
		&m.IsChildTable,
		&m.IsTree,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return DocTypeMeta{}, fmt.Errorf("doctype not found: %s", doctypeName)
		}

		return DocTypeMeta{}, err
	}

	return m, nil
}

func (l *Loader) loadFields(ctx context.Context, doctypeName string) ([]FieldMeta, error) {
	rows, err := l.DB.Query(ctx, `
		SELECT
			name,
			COALESCE(fieldname, ''),
			COALESCE(label, ''),
			COALESCE(fieldtype, ''),
			COALESCE(options, ''),
			reqd,
			hidden,
			read_only,
			in_list_view,
			in_standard_filter,
			search_index,
			unique_field,
			no_copy,
			set_only_once,
			allow_on_submit,
			permlevel,
			columns,
			length,
			precision_value,
			COALESCE(default_value, ''),
			COALESCE(description, ''),
			COALESCE(depends_on, ''),
			COALESCE(mandatory_depends_on, ''),
			COALESCE(read_only_depends_on, ''),
			COALESCE(placeholder, ''),
			COALESCE(fetch_from, ''),
			COALESCE(validation_rule, ''),
			idx
		FROM "tabDocField"
		WHERE parent = $1
		ORDER BY idx, name
	`, doctypeName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	fields := []FieldMeta{}

	for rows.Next() {
		var f FieldMeta

		if err := rows.Scan(
			&f.Name,
			&f.Fieldname,
			&f.Label,
			&f.Fieldtype,
			&f.Options,
			&f.Reqd,
			&f.Hidden,
			&f.ReadOnly,
			&f.InListView,
			&f.InStandardFilter,
			&f.SearchIndex,
			&f.UniqueField,
			&f.NoCopy,
			&f.SetOnlyOnce,
			&f.AllowOnSubmit,
			&f.Permlevel,
			&f.Columns,
			&f.Length,
			&f.PrecisionValue,
			&f.DefaultValue,
			&f.Description,
			&f.DependsOn,
			&f.MandatoryDependsOn,
			&f.ReadOnlyDependsOn,
			&f.Placeholder,
			&f.FetchFrom,
			&f.ValidationRule,
			&f.Idx,
		); err != nil {
			return nil, err
		}

		fields = append(fields, f)
	}

	return fields, rows.Err()
}

func (l *Loader) loadPermissions(ctx context.Context, doctypeName string) ([]PermMeta, error) {
	rows, err := l.DB.Query(ctx, `
		SELECT
			name,
			role,
			permlevel,
			"read",
			"write",
			create_perm,
			delete_perm,
			submit_perm,
			cancel_perm,
			amend_perm,
			print_perm,
			email_perm,
			export_perm,
			import_perm,
			share_perm,
			report_perm,
			idx
		FROM "tabDocPerm"
		WHERE parent = $1
		ORDER BY idx, name
	`, doctypeName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	perms := []PermMeta{}

	for rows.Next() {
		var p PermMeta

		if err := rows.Scan(
			&p.Name,
			&p.Role,
			&p.Permlevel,
			&p.Read,
			&p.Write,
			&p.Create,
			&p.Delete,
			&p.Submit,
			&p.Cancel,
			&p.Amend,
			&p.Print,
			&p.Email,
			&p.Export,
			&p.Import,
			&p.Share,
			&p.Report,
			&p.Idx,
		); err != nil {
			return nil, err
		}

		perms = append(perms, p)
	}

	return perms, rows.Err()
}
