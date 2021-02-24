resource "aci_rest" "rest_cons_if" {
  path       = "api/node/mo/${aci_endpoint_security_group.terraform_esg.id}/rsconsIf-interface.json"
  class_name = "fvRsConsIf"

  content = {
    "tnVzCPIfName" = "interface"
  }
}
