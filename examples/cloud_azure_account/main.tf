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
  name        = "test_tenant"
}

resource "aci_cloud_azure_account" "azure-account" {
  tenant_dn         = aci_tenant.footenant.id
  access_type       = "managed"  //credentials or managed
  account_id        = "my_example_account_id"
  # name              = "azure_cloud_account"
  vendor        = "azure"
}