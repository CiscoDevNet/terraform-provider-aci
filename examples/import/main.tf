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

resource "aci_tenant" "test" {
  #name = "1234"
}

# resource "aci_app_profile" "testap" {}

# resource "aci_bridge_domain" "testbd" {}

# resource "aci_contract" "testcontract" {}

# resource "aci_subnet" "testsubnet" {
#   name             = "10.0.3.30/30"
#   bridge_domain_dn = "uni/tn-Girish/BD-testingbd"
#   scope            = ["private"]
#   ip_address       = "10.0.3.30/30"
# }

# resource "aci_subject" "testsubject" {}


# resource "aci_filter" "testfilter" {}


# resource "aci_entry" "testentry" {}


# resource "aci_epg" "testepg" {}

resource "aci_vrf" "imp_vrf" {

}
