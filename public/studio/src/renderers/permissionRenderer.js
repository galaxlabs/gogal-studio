import { renderTable } from "../components/dataTable.js";

function boolValue(row, normalKey, permKey) {
  return row[normalKey] ?? row[permKey] ?? false;
}

export function renderPermissions(perms) {
  return renderTable(
    [
      { key: "role", label: "Role" },
      { key: "permlevel", label: "Level" },
      {
        key: "read",
        label: "Read",
        render: (row) => boolValue(row, "read", "read_perm")
      },
      {
        key: "write",
        label: "Write",
        render: (row) => boolValue(row, "write", "write_perm")
      },
      {
        key: "create",
        label: "Create",
        render: (row) => boolValue(row, "create", "create_perm")
      },
      {
        key: "delete",
        label: "Delete",
        render: (row) => boolValue(row, "delete", "delete_perm")
      }
    ],
    perms
  );
}