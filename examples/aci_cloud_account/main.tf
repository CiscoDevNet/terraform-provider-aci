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

resource "aci_cloud_account" "azure-account" {
  tenant_dn         = aci_tenant.footenant.id
  access_type       = "managed"  //credentials or managed
  account_id        = "example_id"
  vendor        = "azure"
}

resource "aci_tenant" "cloud_tenant" {
  description = "sample aci_tenant from terraform"
  name        = "cloud_tenant"
}

resource "aci_cloud_account" "cloud_account_1" {
  tenant_dn         = aci_tenant.cloud_tenant.id
  access_type       = "managed"  //credentials or managed
  account_id        = "my_example_account_id"
  name              = "azure_cloud"
  vendor        = "azure"
}

data "aci_cloud_account" "example" {
  tenant_dn  = aci_cloud_account.cloud_account_1.tenant_dn
  account_id  = aci_cloud_account.cloud_account_1.account_id
  vendor  = aci_cloud_account.cloud_account_1.vendor
}

output "name" {
  value = data.aci_cloud_account.example
}