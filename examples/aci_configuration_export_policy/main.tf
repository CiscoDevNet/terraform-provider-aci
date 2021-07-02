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
resource "aci_configuration_export_policy" "example" {
  name                  = "example"
  description           = "from terraform"
  admin_st              = "untriggered"
  annotation            = "example"
  format                = "json"
  include_secure_fields = "yes"
  max_snapshot_count    = "10"
  name_alias            = "example"
  snapshot              = "yes"
  target_dn             = "uni/tn-test"
}