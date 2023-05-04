resource "aci_tenant" "tenant_for_filter" {
  name        = "_ACI-BENCHMARK_FOR_FILTER"
  description = "This tenant is created by terraform ACI provider"
}

resource "aci_filter" "bd_flood_filter" {
  tenant_dn   = aci_tenant.tenant_for_filter.id
  name        = "test_bdfloodfilter"
  description = "This filter is created by terraform ACI provider."
}

resource "aci_filter" "bd_flood_filter2" {
  tenant_dn   = aci_tenant.tenant_for_filter.id
  name        = "test_bdfloodfilter2"
  description = "This filter is created by terraform ACI provider."
}

resource "aci_filter" "bd_flood_filter3" {
  tenant_dn   = aci_tenant.tenant_for_filter.id
  name        = "test_bdfloodfilter3"
  description = "This filter is created by terraform ACI provider."
}
