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
resource "aci_tenant" "example" {
  name = "test_acc_tenant"
}

resource "aci_l3_outside" "fool3_outside" {
  tenant_dn      = aci_tenant.example.id
  description    = "from terraform"
  name           = "demo_l3out"
  annotation     = "tag_l3out"
  enforce_rtctrl = ["export"]
  name_alias     = "alias_out"
  target_dscp    = "unspecified"
}