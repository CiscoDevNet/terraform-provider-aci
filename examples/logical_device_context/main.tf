provider "aci" {
  username = ""
  password = ""
  url      = ""
  insecure = true
}

resource "aci_logical_device_context" "example" {
  tenant_dn         = aci_tenant.tenentcheck.id
  ctrct_name_or_lbl = "default"
  graph_name_or_lbl = "any"
  node_name_or_lbl  = "N1"
  context           = "ctx1"
  description       = "from terraform"
  name_alias        = "example"
}
