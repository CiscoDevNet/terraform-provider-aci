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

resource "aci_maintenance_policy" "example" {
  name                   = "mnt_policy"
  admin_st               = "triggered"
  description            = "from terraform"
  annotation             = "example"
  graceful               = "yes"
  ignore_compat          = "yes"
  internal_label         = "example"
  name_alias             = "example"
  notif_cond             = "notifyOnlyOnFailures"
  run_mode               = "pauseOnlyOnFailures"
  version                = "n9000-15.0(1k)"
  version_check_override = "trigger-immediate"
}