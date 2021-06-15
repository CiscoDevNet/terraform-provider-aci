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

resource "aci_logical_interface_context" "logical_interface_context_one" {
	logical_device_context_dn  = aci_logical_device_context.ldev_ctx.id
	annotation  = "first logical interface"
	conn_name_or_lbl  = "first_interface_ctx"
	l3_dest  = "no"
	name_alias  = "first_interface_ctx"
	permit_log  = "no"
}