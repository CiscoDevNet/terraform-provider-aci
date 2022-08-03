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

resource "aci_cloud_account" "azure_account" {
  tenant_dn   = aci_tenant.footenant.id
  access_type = "managed"
  account_id  = "my_example_account_id"
  vendor      = "azure"
}

# Own Subscription (attaching the subscription id from the same tenant)
resource "aci_tenant_to_cloud_account" "tenant_to_account" {
  tenant_dn        = aci_tenant.footenant.id
  cloud_account_dn = aci_cloud_account.azure_account.id
}

data "aci_tenant_to_cloud_account" "example" {
  tenant_dn        = aci_tenant_to_cloud_account.tenant_to_account.tenant_dn
  cloud_account_dn = aci_tenant_to_cloud_account.tenant_to_account.cloud_account_dn
}

output "name" {
  value = data.aci_tenant_to_cloud_account.example
}

# Shared Subscription (attaching the subscription id from another exiting tenant)
resource "aci_tenant" "cloud_tenant" {
  description = "sample aci_tenant from terraform"
  name        = "cloud_tenant"
}

resource "aci_cloud_account" "cloud_account" {
  tenant_dn   = aci_tenant.cloud_tenant.id
  access_type = "managed"
  account_id  = "example_account_id"
  name        = "azure_cloud"
  vendor      = "azure"
}

resource "aci_tenant_to_cloud_account" "new_tenant_to_account" {
  tenant_dn        = aci_tenant.cloud_tenant.id
  cloud_account_dn = aci_cloud_account.azure_account.id
}