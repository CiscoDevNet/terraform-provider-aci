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

resource "aci_cloud_aws_provider" "foocloud_aws_provider" {
		tenant_dn         = "${aci_tenant.footenant.id}"
		description       = "aws account config"
		access_key_id     = "access_key"
		account_id        = "acc_id"
		annotation        = "tag_aws"
		secret_access_key = "secret_key"
	}