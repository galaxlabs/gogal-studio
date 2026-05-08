import { renderTable } from "../components/table.js";

export function renderDocTypeDetails(doctype) {
  const dt = doctype || {};
  return renderTable(
    [
      { key: "label", label: "Property" },
      { key: "value", label: "Value" }
    ],
    [
      { label: "Name", value: dt.name },
      { label: "Module", value: dt.module },
      { label: "App", value: dt.app_name },
      { label: "Table", value: dt.table_name },
      { label: "Single", value: dt.is_single ?? false },
      { label: "Child Table", value: dt.is_child_table ?? false },
      { label: "Submittable", value: dt.is_submittable ?? false },
      { label: "Tree", value: dt.is_tree ?? false }
    ],
    "No DocType metadata found."
  );
}
