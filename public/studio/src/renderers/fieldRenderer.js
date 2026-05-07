import { renderTable } from "../components/dataTable.js";

export function renderFields(fields) {
  return renderTable(
    [
      { key: "idx", label: "Idx" },
      {
        key: "fieldname",
        label: "Field",
        render: (row) => `<strong>${row.fieldname}</strong><br><small>${row.label}</small>`
      },
      {
        key: "fieldtype",
        label: "Type",
        render: (row) => `<span class="badge">${row.fieldtype}</span>`
      },
      { key: "options", label: "Options" },
      { key: "reqd", label: "Reqd" },
      { key: "hidden", label: "Hidden" },
      { key: "read_only", label: "Read Only" },
      { key: "in_list_view", label: "List" }
    ],
    fields
  );
}