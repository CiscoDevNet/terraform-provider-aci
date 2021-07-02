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

resource "aci_firmware_policy" "example" {
  name  = "example"
  description = "from terraform"
  annotation  = "example"
  effective_on_reboot  = "no"
  ignore_compat  = "no"
  internal_label  = "example_policy"
  name_alias  = "example"
  version  = "n9000-14.2(3q)"
  version_check_override  = "untriggered"
}