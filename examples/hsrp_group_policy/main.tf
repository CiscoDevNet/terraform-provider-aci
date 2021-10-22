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

resource "aci_hsrp_group_policy" "example" {

  tenant_dn              = aci_tenant.example.id
  name                   = "example"
  annotation             = "example"
  ctrl                   = "preempt"
  hello_intvl            = "3000"
  hold_intvl             = "10000"
  key                    = "cisco"
  name_alias             = "example"
  preempt_delay_min      = "60"
  preempt_delay_reload   = "60"
  preempt_delay_sync     = "60"
  prio                   = "100"
  timeout                = "60"
  hsrp_group_policy_type = "md5"

}
