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
	name       = "test_acc_tenant"
}

resource "aci_l3_outside" "example" {
		tenant_dn      = aci_tenant.example.id
		description    = "from terraform"
		name           = "demo_l3out"
		annotation     = "tag_l3out"
		enforce_rtctrl = "export"
		name_alias     = "alias_out"
		target_dscp    = "unspecified"
}
resource "aci_logical_node_profile" "foological_node_profile" {
	l3_outside_dn = aci_l3_outside.example.id
	description   = "From Terraform"
	name          = "demo_node"
	annotation    = "tag_node"
	config_issues = "none"
	name_alias    = "alias_node"
	tag           = "black"
	target_dscp   = "unspecified"
}	