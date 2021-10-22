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
		description = "Tenant created while acceptance testing"
		name        = "demo_tenant"
	}

	resource "aci_cloud_applicationcontainer" "foocloud_applicationcontainer" {
		tenant_dn   = aci_tenant.footenant.id
		name        = "demo_app"
		annotation  = "tag_app"
	}

	resource "aci_cloud_epg" "foocloud_epg" {
		cloud_applicationcontainer_dn = aci_cloud_applicationcontainer.foocloud_applicationcontainer.id
		name                          = "cloud_epg"
	}

	resource "aci_cloud_endpoint_selector" "foocloud_endpoint_selector" {
		cloud_epg_dn     = aci_cloud_epg.foocloud_epg.id
		description      = "sample ep selector"
		name             = "ep_select"
		annotation       = "tag_ep"
		match_expression = "custom:Name=='admin-ep2'"
		name_alias       = "alias_ep"
	}