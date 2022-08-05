resource "aci_x509_certificate" "test_cert" {
  local_user_dn = "uni/userext/user-admin"
  description   = "From Terraform"

  name       = "test_1"
  annotation = "x509_certificate_tag"
  name_alias = "alias_name"
  data       = file("x509_certificate.cert")
}

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