terraform {
  required_providers {
    aci = {
      source = "ciscodevnet/aci"
    }
  }
}

provider "aci" {
  username = ""
  password = ""
  url      = ""
  insecure = true
}

resource "aci_tenant" "example" {
  name = "test_acc_tenant"
}

resource "aci_logical_device_context" "example" {
  tenant_dn                          = aci_tenant.example.id
  ctrct_name_or_lbl                  = "default"
  graph_name_or_lbl                  = "any"
  node_name_or_lbl                   = "N1"
  context                            = "ctx1"
  description                        = "from terraform"
  name_alias                         = "example"
  relation_vns_rs_l_dev_ctx_to_l_dev = "uni/tn-test_acc_tenant/lDevVip-LoadBalancer01"
}
