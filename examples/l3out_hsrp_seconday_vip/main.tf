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

resource "aci_l3out_hsrp_secondary_vip" "example" {

  l3out_hsrp_interface_group_dn = aci_l3out_hsrp_interface_group.example.id
  ip                            = "10.0.0.1"
  annotation                    = "example"
  config_issues                 = "GroupMac-Conflicts-Other-Group"
  name_alias                    = "example"
  description                   = "from terraform"
}
