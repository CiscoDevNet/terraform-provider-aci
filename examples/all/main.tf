provider "aci" {
  username = ""
  password = ""
  url      = ""
  insecure = true
}


resource "aci_tenant" "madebytf2" {
  name        = "madebytf2"
  description = "This tenant is created by terraform"
}

resource "aci_app_profile" "madebytf2" {
  tenant_dn   = aci_tenant.madebytf2.id
  name        = "madebytf2"
  description = "This app profile is created by terraform"
}

resource "aci_epg" "madebytf2" {
  application_profile_dn = aci_app_profile.madebytf2.id
  name                   = "madebytf2"
  description            = "This epg is created by terraform"
}

resource "aci_bridge_domain" "madebytf2" {
  tenant_dn   = aci_tenant.madebytf2.id
  name        = "madebytf2"
  description = "This bridge domain is created by terraform"
  mac         = "00:22:BD:F8:19:FF"
}

resource "aci_contract" "madebytf2" {
  tenant_dn   = aci_tenant.madebytf2.id
  name        = "madebytfsdfsdf"
  description = "This contract is created by terraform"
  scope       = "context"
  dscp        = "VA"
}

resource "aci_subject" "madebytf2" {
  contract_dn = aci_contract.madebytf2.id
  name        = "madebytftt"
  description = "This subject is created by terraform"
}

resource "aci_subnet" "madebytf2" {
  name             = "10.0.3.28/27"
  bridge_domain_dn = aci_bridge_domain.madebytf2.id
  ip_address       = "10.0.3.28/27"
  scope            = ["private"]
  description      = "This subject is created by terraform"
}

resource "aci_filter" "madebytf2" {
  tenant_dn   = aci_tenant.madebytf2.id
  name        = "madebytf"
  description = "This filter is created by terraform"
}

resource "aci_entry" "madebytf2" {
  filter_dn   = aci_filter.madebytf2.id
  name        = "madebytf"
  description = "This entry is created by terraform"
}

resource "aci_rest" "madebyresttf" {
  path       = "/api/node/mo/uni/tn-tntest.json"
  class_name = "fvTenant"

  content = {
    "dn"   = "uni/tn-tntest"
    "name" = "tntest"
  }
}
