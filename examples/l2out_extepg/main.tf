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


resource "aci_tenant" "foo_tenant" {
  name        = "tenant_1"
  description = "This tenant is created by terraform"
}

resource "aci_l2_outside" "fool2_outside" {
  tenant_dn   = aci_tenant.foo_tenant.id
  description = "from terraform"
  name        = "l2_outside_1"
  annotation  = "l2_outside_tag"
  name_alias  = "example"
  target_dscp = "AF11"
}

resource "aci_l2out_extepg" "fool2out_extepg" {
  l2_outside_dn  = aci_l2_outside.fool2_outside.id
  description    = "from terraform"
  name           = "l2out_extepg_1"
  annotation     = "l2out_extepg_tag"
  exception_tag  = "example"
  flood_on_encap = "disabled"
  match_t        = "All"
  name_alias     = "example"
  pref_gr_memb   = "exclude"
  prio           = "level1"
  target_dscp    = "AF11"
}

data "aci_l2out_extepg" "example8" {
  l2_outside_dn = aci_l2_outside.fool2_outside.id
  name          = aci_l2out_extepg.fool2out_extepg.name
}

output "name8" {
  value = data.aci_l2out_extepg.example8
}