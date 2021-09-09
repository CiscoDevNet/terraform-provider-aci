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

resource "aci_tenant" "footenant" {
	description = "sample aci_tenant from terraform"
	name        = "demo_tenant"
	annotation  = "tag_tenant"
	name_alias  = "alias_tenant"
}

resource "aci_bfd_interface_policy" "example" {
  tenant_dn = aci_tenant.footenant.id
  name = "example"
  admin_st = "enabled"
  annotation  = "example"
  ctrl = "opt-subif"
  detect_mult  = "3"
  echo_admin_st = "disabled"
  echo_rx_intvl  = "50"
  min_rx_intvl  = "50"
  min_tx_intvl  = "50"
  name_alias  = "example"
  description = "example"
}