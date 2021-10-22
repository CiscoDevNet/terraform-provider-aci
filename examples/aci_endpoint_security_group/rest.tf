resource "aci_rest" "rest_con_export" {
  path    = "api/node/mo/uni/tn-tf_tenant/cif-exported_contract.json"
  payload = <<EOF
  "vzCPIf": {
    "attributes": {
      "dn": "uni/tn-tf_tenant/cif-exported_contract",
      "name": "exported_contract"
    },
    "children": [
      {
        "vzRsIf": {
          "attributes": {
            "tDn": "uni/tn-tenant_export_cons/brc-exported_contract"
          },
          "children": []
        }
      }
    ]
  }
  EOF
}

resource "aci_rest" "rest_qos_custom_pol" {
  path       = "api/node/mo/${aci_tenant.terraform_tenant.id}/qoscustom-testpol.json"
  class_name = "qosCustomPol"

  content = {
    "name" = "testpol"
  }
}

resource "aci_rest" "rest_taboo_con" {
  path       = "api/node/mo/${aci_tenant.terraform_tenant.id}/taboo-testcon.json"
  class_name = "vzTaboo"

  content = {
    "name" = "testcon"
  }
}