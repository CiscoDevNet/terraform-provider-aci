resource "aci_tenant" "tenant_for_entry" {
  name        = "tenant_for_entry"
  description = "This tenant is created by terraform ACI provider"
}

resource "aci_filter" "filter_for_entry" {
  tenant_dn   = aci_tenant.tenant_for_entry.id
  name        = "filter_for_entry"
  description = "This filter is created by terraform ACI provider."
}

resource "aci_entry" "demoentry" {
  filter_dn   = aci_filter.filter_for_entry.id
  name        = "test_tf_entry"
  description = "This entry is created by terraform ACI provider"
}
