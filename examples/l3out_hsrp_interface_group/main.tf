terraform {
  required_providers {
    aci = {
      source = "ciscodevnet/aci"
    }
  }
}

#configure provider with your cisco aci credentials.
provider "aci" {
  username = "" # <APIC username>
  password = "" # <APIC pwd>
  url      = "" # <cloud APIC URL>
  insecure = true
}

resource "aci_l3out_hsrp_interface_group" "example" {
  l3out_hsrp_interface_profile_dn = aci_l3out_hsrp_interface_profile.example.id
  name                            = "one"
  annotation                      = "example"
  description                     = "From terraform"
  config_issues                   = "none"
  group_af                        = "ipv4"
  group_id                        = "20"
  group_name                      = "test"
  ip                              = "10.22.30.40"
  ip_obtain_mode                  = "admin"
  mac                             = "02:10:45:00:00:56"
  name_alias                      = "example"
}