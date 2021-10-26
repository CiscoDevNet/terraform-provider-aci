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

resource "aci_hsrp_interface_policy" "example" {
  tenant_dn    = aci_tenant.tenentcheck.id
  name         = "one"
  annotation   = "example"
  description  = "from terraform"
  ctrl         = "bia"
  delay        = "10"
  name_alias   = "example"
  reload_delay = "10"
}