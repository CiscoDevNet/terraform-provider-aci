provider "aci" {
  username = ""
  password = ""
  url      = ""
  insecure = true
}

resource "aci_epg_to_contract" "example_consumer" {
  application_epg_dn = "uni/tn-nirav/ap-nkap/epg-nkepg"
  contract_dn        = "uni/tn-nirav/brc-nkcon"
  contract_type      = "provider"
}


resource "aci_epg_to_contract" "example_provider" {
  application_epg_dn = "uni/tn-nirav/ap-nkap/epg-nkepg"
  contract_dn        = "uni/tn-nirav/brc-nkcon"
  contract_type      = "consumer"
  prio = "level1"
}

data "aci_epg_to_contract" "example" {
    application_epg_dn = "uni/tn-nirav/ap-nkap/epg-nkepg"
    contract_name  = "nkcon"
    contract_type = "consumer"
}