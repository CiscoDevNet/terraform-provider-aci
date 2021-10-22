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
resource "aci_interface_fc_policy" "test_fc" {
  name         = "tf_pol"
  description  = "From Terraform"
  annotation   = "tag_if_policy"
  fill_pattern = "IDLE"
  name_alias   = "demo_alias"
  port_mode    = "f"
  rx_bb_credit = "64"
  speed        = "auto"
  trunk_mode   = "auto"

}
