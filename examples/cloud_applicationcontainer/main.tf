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

resource "aci_cloud_applicationcontainer" "foo_clou_app" {
  tenant_dn   = "${aci_tenant.footenant.id}"
  name        = "demo_cloud_app"
  description = "aci_cloud_applicationcontainer from terraform"
  annotation  = "tag_cloud_app"
  name_alias  = "alias_app"
}