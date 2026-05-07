export function renderDocTypeDetails(dt) {
  return `
    <div class="kv"><span>Name</span><strong>${dt.name}</strong></div>
    <div class="kv"><span>Module</span><strong>${dt.module}</strong></div>
    <div class="kv"><span>App</span><strong>${dt.app_name}</strong></div>
    <div class="kv"><span>Single</span><strong>${dt.is_single}</strong></div>
    <div class="kv"><span>Child Table</span><strong>${dt.is_child_table}</strong></div>
    <div class="kv"><span>Submittable</span><strong>${dt.is_submittable}</strong></div>
    <div class="kv"><span>Tree</span><strong>${dt.is_tree}</strong></div>
  `;
}