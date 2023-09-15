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

resource "aci_tenant" "test-tenant" {
  name        = "test-tenant-a"
  description = "This tenant is created by terraform"
}

resource "aci_service_redirect_policy" "test-srp" {
  tenant_dn   = aci_tenant.test-tenant-a.id
  name        = "test-srp-a"
  description = "This service redirect policy is created by terraform"
}

data "aci_client_end_point" "check" {
  allow_empty_result = true
  name               = "foobar"
}

resource "aci_destination_of_redirected_traffic" "example" {
  service_redirect_policy_dn = aci_service_redirect_policy.test-srp.id
  ip                         = data.aci_client_end_point.check.fvcep_objects.0.ip
  mac                        = data.aci_client_end_point.check.fvcep_objects.0.mac
}
