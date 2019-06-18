resource "aci_rest" "madebyresttf" {
  path       = "/api/node/mo.json"
  class_name = "fvTenant"

  content = {
    "dn"   = "uni/tn-tntestrest1"
    "name" = "tntestrest1"
  }
}
