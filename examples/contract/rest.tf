resource "aci_rest" "rest_abs_graph" {
  path       = "api/node/mo/${aci_tenant.tenant_for_contract.id}/AbsGraph-testgraph.json"
  class_name = "vnsAbsGraph"

  content = {
    "name" = "testgraph"
  }
}
