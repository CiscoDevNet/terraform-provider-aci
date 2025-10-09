
resource "aci_vmm_uplink_policy" "full_example_vmm_uplink_container" {
  parent_dn   = aci_vmm_uplink_container.example.id
  annotation  = "annotation"
  name_alias  = "name_alias_1"
  uplink_id   = "1"
  uplink_name = "uplink_name_1"
  annotations = [
    {
      key   = "key_0"
      value = "value_1"
    }
  ]
  tags = [
    {
      key   = "key_0"
      value = "value_1"
    }
  ]
}
