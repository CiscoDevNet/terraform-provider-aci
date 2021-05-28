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

	resource "aci_vrf" "vrf1" {
		tenant_dn = "${aci_tenant.footenant.id}"
		name      = "acc-vrf"
	}

	resource "aci_cloud_context_profile" "foocloud_context_profile" {
		name 		             = "ctx_prof_cidr"
		description              = "cloud_context_profile"
		tenant_dn                = "${aci_tenant.footenant.id}"
		primary_cidr             = "10.230.231.1/16"
		region                   = "us-west-1"
		cloud_vendor			 = "aws"
		relation_cloud_rs_to_ctx = "${aci_vrf.vrf1.id}"
	}

resource "aci_cloud_cidr_pool" "foocloud_cidr_pool" {
		cloud_context_profile_dn = "${aci_cloud_context_profile.foocloud_context_profile.id}"
		description              = "cloud CIDR from terraform"
		addr                     = "10.0.1.10/28"
		annotation               = "tag_cidr"
		name_alias               = "name_alias"
		primary                  = "yes"
	}