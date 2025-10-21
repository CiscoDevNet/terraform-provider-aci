
data "aci_imported_logical_device" "example_tenant" {
  parent_dn      = aci_tenant.example.id
  logical_device = aci_l4_l7_device.example_in_another_tenant.id
}
