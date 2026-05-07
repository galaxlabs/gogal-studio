window.GOGAL_BUILDER_POLICY = {
  doctypeCreateFields: [
    "name",
    "module",
    "is_single",
    "is_submittable",
    "is_child_table",
    "editable_grid"
  ],

  doctypePropertyFields: {
    basic: [
      "name",
      "label",
      "module",
      "table_name"
    ],

    behavior: [
      "is_single",
      "is_submittable",
      "is_child_table",
      "editable_grid",
      "quick_entry",
      "track_changes"
    ],

    naming: [
      "naming_rule",
      "title_field",
      "sort_field",
      "sort_order"
    ],

    importExport: [
      "allow_import",
      "allow_export",
      "allow_rename"
    ],

    advanced: [
      "controller",
      "route",
      "is_tree",
      "is_virtual",
      "description"
    ]
  },

  // Never show these on main form canvas.
  doctypeMetaOnlyFields: [
    "module",
    "is_submittable",
    "istable",
    "is_child_table",
    "issingle",
    "is_single",
    "editable_grid",
    "quick_entry",
    "track_changes",
    "track_seen",
    "track_views",
    "custom",
    "beta",
    "is_virtual",
    "queue_in_background",
    "engine",
    "migration_hash",
    "row_format",
    "permissions",
    "actions",
    "links",
    "states",
    "fields",
    "form_builder",
    "form_builder_tab",
    "settings_tab",
    "connections_tab",
    "web_view",
    "email_settings_sb",
    "advanced",
    "json_hash",
    "source_path",
    "created_at",
    "updated_at",
    "status"
  ],

  // System fields should not appear as user-created business fields.
  systemFieldnames: [
    "name",
    "owner",
    "creation",
    "modified",
    "modified_by",
    "docstatus",
    "idx",
    "parent",
    "parenttype",
    "parentfield",
    "doctype"
  ],

  // Internal DocField properties. Keep these hidden from normal property panel.
  docfieldAdvancedOnlyFields: [
    "oldfieldname",
    "oldfieldtype",
    "permlevel",
    "ignore_user_permissions",
    "ignore_xss_filter",
    "report_hide",
    "print_hide",
    "print_hide_if_no_value",
    "print_width",
    "width",
    "max_height",
    "documentation_url",
    "link_filters",
    "make_attachment_public",
    "button_color",
    "show_description_on_click",
    "is_virtual",
    "not_nullable",
    "mask",
    "sticky"
  ],

  docfieldVisiblePropertyFields: {
    basic: [
      "label",
      "fieldname",
      "fieldtype",
      "options"
    ],

    validation: [
      "reqd",
      "required",
      "unique",
      "unique_field",
      "read_only",
      "hidden",
      "default"
    ],

    listSearch: [
      "in_list_view",
      "in_standard_filter",
      "in_filter",
      "in_global_search",
      "search_index"
    ],

    layout: [
      "columns",
      "description",
      "placeholder",
      "depends_on",
      "mandatory_depends_on",
      "read_only_depends_on"
    ]
  },

  fieldTypes: [
    "Data",
    "Small Text",
    "Text",
    "Long Text",
    "Int",
    "Float",
    "Currency",
    "Check",
    "Date",
    "Datetime",
    "Time",
    "Select",
    "Link",
    "Table",
    "Attach",
    "Attach Image",
    "JSON",
    "Code",
    "Section Break",
    "Column Break",
    "Tab Break",
    "Button",
    "HTML"
  ]
};
