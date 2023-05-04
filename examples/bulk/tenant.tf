data "aci_tenant" "tenant_for_contract" {
  name        = "common"
}

resource "aci_tenant" "tenant_for_benchmark" {
  name        = "_ACI-BENCHMARK"
  description = "This tenant is created by terraform ACI provider"
}
