export function isSystemField(fieldname) {
  return [
    "id",
    "name",
    "owner",
    "created_by",
    "modified_by",
    "created_at",
    "updated_at",
    "creation",
    "modified",
    "deleted_at",
    "idx",
    "docstatus",
    "version",
    "site_id",
    "app_id",
    "module_id"
  ].includes(fieldname);
}

export function isEditableField(field) {
  if (!field) return false;
  if (field.hidden) return false;
  if (field.read_only) return false;
  if (isSystemField(field.fieldname)) return false;
  return true;
}

export function getInputType(fieldtype) {
  switch (fieldtype) {
    case "Int":
    case "Float":
    case "Currency":
      return "number";
    case "Date":
      return "date";
    case "Datetime":
      return "datetime-local";
    case "Check":
      return "checkbox";
    default:
      return "text";
  }
}

export function normalizeValueForSave(field, value) {
  if (field.fieldtype === "Check") return Boolean(value);
  if (["Int"].includes(field.fieldtype)) return value === "" ? null : Number.parseInt(value, 10);
  if (["Float", "Currency"].includes(field.fieldtype)) return value === "" ? null : Number.parseFloat(value);
  return value;
}
