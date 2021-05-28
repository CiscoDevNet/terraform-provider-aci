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
resource "aci_miscabling_protocol_interface_policy" "test_miscable" {
  description = "example description"
  name        = "demo_mcpol"
  admin_st    = "enabled"
  annotation  = "tag_mcpol"
  name_alias  = "alias_mcpol"
}
