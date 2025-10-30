
resource "aci_imported_logical_device" "full_example_tenant" {
  parent_dn      = aci_tenant.example.id
  annotation     = "annotation"
  description    = "description_1"
  logical_device = aci_l4_l7_device.example_in_another_tenant.id
  name           = "name_1"
  name_alias     = "name_alias_1"
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
