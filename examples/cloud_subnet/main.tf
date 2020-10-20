provider "aci" {
  username = ""
  password = ""
  url      = ""
  insecure = true
}

resource "aci_tenant" "tenentcheck" {
  description = "cloud resources"
  name        = "cloudTenant"
  annotation  = "cloudTenant"
  name_alias  = "cloud_tenant"
}

resource "aci_vrf" "vrf1" {
  tenant_dn = "${aci_tenant.tenentcheck.id}"
  name      = "vrf-1"
}

resource "aci_cloud_context_profile" "ctx1" {
    name 		                = "check"
	description                 = "cloud_context_profile created while acceptance testing"
	tenant_dn                   = "${aci_tenant.tenentcheck.id}"
	primary_cidr                = "10.0.0.0/16"
	region                      = "us-west-1"
	cloud_vendor			    = "aws"
	relation_cloud_rs_to_ctx    = "${aci_vrf.vrf1.id}"
    hub_network                 =  "uni/tn-infra/gwrouterp-default"
}

resource "aci_cloud_cidr_pool" "cloud_cidr_pool" {
    cloud_context_profile_dn = "${aci_cloud_context_profile.ctx1.id}"
    addr                     = "10.0.0.0/16"
}

resource "aci_cloud_subnet" "cloud_subnet" {
  	cloud_cidr_pool_dn            = "${aci_cloud_cidr_pool.cloud_cidr_pool.id}"
  	ip                            = "10.0.1.0/24"
    usage = "gateway"
	zone = "uni/clouddomp/provp-aws/region-us-west-1/zone-us-west-1b"
}