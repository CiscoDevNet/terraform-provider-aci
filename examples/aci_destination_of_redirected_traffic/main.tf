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

resource "aci_destination_of_redirected_traffic" "example" {
  service_redirect_policy_dn = aci_service_redirect_policy.example.id
  ip                         = "1.2.3.4"
  mac                        = "12:25:56:98:45:74"
  ip2                        = "10.20.30.40"
  dest_name                  = "last"
  pod_id                     = "5"
  annotation                 = "load_traffic_dest"
  description                = "From Terraform"
  name_alias                 = "load_traffic_dest"
}