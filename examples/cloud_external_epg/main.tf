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
		tenant_dn   = "${aci_tenant.footenant.id}"
		name        = "demo_app"
		annotation  = "tag_app"
	}

	resource "aci_cloud_external_epg" "foocloud_external_epg" {
		cloud_applicationcontainer_dn = "${aci_cloud_applicationcontainer.foocloud_applicationcontainer.id}"
		description                   = "sample cloud external epg"
		name                          = "cloud_ext_epg"
		annotation                    = "tag_ext_epg"
		exception_tag                 = "0"
		flood_on_encap                = "disabled"
		match_t                       = "All"
		name_alias                    = "alias_ext"
		pref_gr_memb                  = "exclude"
		prio                          = "unspecified"
		route_reachability            = "inter-site"
	}