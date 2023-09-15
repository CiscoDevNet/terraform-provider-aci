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

resource "aci_recurring_window" "example" {
  scheduler_dn      = aci_trigger_scheduler.example.id
  name              = "example"
  concur_cap        = "unlimited"
  day               = "every-day"
  hour              = "0"
  minute            = "0"
  node_upg_interval = "0"
  proc_break        = "none"
  proc_cap          = "unlimited"
  time_cap          = "unlimited"
  annotation        = "Example"
}