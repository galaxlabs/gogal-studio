export const FIELD_TYPES = {
  Data: {
    control: "input",
    inputType: "text",
    optionsMode: "none",
    placeholder: "Short text",
    defaultColumns: 6
  },

  Text: {
    control: "textarea",
    rows: 3,
    optionsMode: "none",
    placeholder: "Text",
    defaultColumns: 12
  },

  "Small Text": {
    control: "textarea",
    rows: 2,
    optionsMode: "none",
    placeholder: "Small text",
    defaultColumns: 12
  },

  "Long Text": {
    control: "textarea",
    rows: 6,
    optionsMode: "none",
    placeholder: "Long text",
    defaultColumns: 12
  },

  "Text Editor": {
    control: "editor",
    optionsMode: "none",
    placeholder: "Rich text editor",
    defaultColumns: 12
  },

  Int: {
    control: "input",
    inputType: "number",
    optionsMode: "none",
    placeholder: "0",
    defaultColumns: 4
  },

  Float: {
    control: "input",
    inputType: "number",
    step: "any",
    optionsMode: "none",
    placeholder: "0.00",
    defaultColumns: 4
  },

  Currency: {
    control: "input",
    inputType: "number",
    step: "any",
    optionsMode: "none",
    placeholder: "0.00",
    defaultColumns: 4
  },

  Check: {
    control: "checkbox",
    optionsMode: "none",
    defaultColumns: 3
  },

  Date: {
    control: "input",
    inputType: "date",
    optionsMode: "none",
    defaultColumns: 4
  },

  Datetime: {
    control: "input",
    inputType: "datetime-local",
    optionsMode: "none",
    defaultColumns: 4
  },

  Time: {
    control: "input",
    inputType: "time",
    optionsMode: "none",
    defaultColumns: 4
  },

  Select: {
    control: "select",
    optionsMode: "lines",
    placeholder: "One option per line",
    defaultColumns: 6
  },

  Link: {
    control: "link",
    optionsMode: "doctype",
    placeholder: "Target DocType name",
    defaultColumns: 6
  },

  Table: {
    control: "table",
    optionsMode: "child_doctype",
    placeholder: "Child Table DocType name",
    defaultColumns: 12
  },

  Attach: {
    control: "file",
    optionsMode: "none",
    defaultColumns: 6
  },

  "Attach Image": {
    control: "image",
    optionsMode: "none",
    defaultColumns: 6
  },

  JSON: {
    control: "code",
    language: "json",
    optionsMode: "none",
    defaultColumns: 12
  },

  Code: {
    control: "code",
    language: "text",
    optionsMode: "none",
    defaultColumns: 12
  },

  Section: {
    control: "layout",
    layoutType: "section",
    optionsMode: "none",
    defaultColumns: 12
  },

  Column: {
    control: "layout",
    layoutType: "column",
    optionsMode: "none",
    defaultColumns: 6
  }
};

export function getFieldTypeConfig(fieldtype) {
  return FIELD_TYPES[fieldtype] || FIELD_TYPES.Data;
}

export function getFieldTypeNames() {
  return Object.keys(FIELD_TYPES);
}