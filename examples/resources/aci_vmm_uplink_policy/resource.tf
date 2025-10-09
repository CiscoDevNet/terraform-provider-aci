
resource "aci_vmm_uplink_policy" "example_vmm_uplink_container" {
  parent_dn   = aci_vmm_uplink_container.example.id
  uplink_id   = "1"
  uplink_name = "uplink_name_1"
}
