resource "aci_rest" "madebyresttf" {
  path       = "/api/node/mo.json"
  class_name = "fvTenant"

  content = {
    "dn"   = "uni/tn-tntestrest1"
    "name" = "tntestrest1"
  }
}

# resource "aci_rest" "test_rel" {
#   path       = "/api/node/mo/uni/tn-Tenant10/out-L3out_OSPF/rsl3DomAtt.json"
#   class_name = "l3extRsL3DomAtt"
#   content = {
#     "tDn" = "uni/l3dom-L3-out-Domain"
#   }
# }

resource "aci_rest" "madebyrest_yaml" {
  path    = "/api/mo/uni.json"
  payload = <<EOF
{
        "fvTenant": {
          "attributes": {
            "name": "Sales",
            "descr": "Sales department json"
          }
        }
}
  EOF
}
