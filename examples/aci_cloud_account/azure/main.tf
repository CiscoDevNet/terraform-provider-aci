terraform {
  required_providers {
    aci = {
      source = "ciscodevnet/aci"
    }
  }
}

# Azure cloud APIC
provider "aci" {
  username = ""
  password = ""
  url      = ""
  insecure = true
}

# 1. Creating a cloud account
resource "aci_tenant" "footenant" {
  description = "sample aci_tenant from terraform"
  name        = "test_tenant"
}

# TESTING access_type = "managed" 
resource "aci_cloud_account" "foo_cloud_account" {
  tenant_dn   = aci_tenant.footenant.id
  access_type = "managed" //credentials or managed
  account_id  = "example_id"
  vendor      = "azure"
}

# 2. Associating the tenant and cloud account it with subscription

# 2.1 Own Subscription (attaching the subscription id from the same tenant)
resource "aci_tenant" "azure_tenant" {
  description = "sample aci_tenant from terraform"
  name        = "test_tenant"
}

resource "aci_cloud_account" "azure_account" {
  tenant_dn   = aci_tenant.azure_tenant.id
  access_type = "managed"
  account_id  = "example_id"
  vendor      = "azure"
}

# Own Subscription
resource "aci_tenant_to_cloud_account" "tenant_to_azure_account" {
  tenant_dn        = aci_tenant.azure_tenant.id
  cloud_account_dn = aci_cloud_account.azure_account.id
}

# 2.1 Shared Subscription (attaching the subscription id from another exiting tenant)
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
# relation_cloud_rs_account_to_access_policy -> is only available in in the most recent version of cloud APIC (>25.0.3)
}

# Shared Subscription
resource "aci_tenant_to_cloud_account" "new_tenant_to_account" {
  tenant_dn        = aci_tenant.cloud_tenant.id
  cloud_account_dn = aci_cloud_account.azure_account.id
}

# DATASOURCE
data "aci_cloud_account" "aci_cloud_account_data" {
  tenant_dn  = aci_cloud_account.cloud_account.tenant_dn
  account_id = aci_cloud_account.cloud_account.account_id
  vendor     = aci_cloud_account.cloud_account.vendor
}

output "aci_cloud_account_output" {
  value = data.aci_cloud_account.aci_cloud_account_data
}

data "aci_tenant_to_cloud_account" "aci_tenant_to_cloud_account_data" {
  tenant_dn        = aci_tenant_to_cloud_account.new_tenant_to_account.tenant_dn
  cloud_account_dn = aci_tenant_to_cloud_account.new_tenant_to_account.cloud_account_dn
}

output "aci_tenant_to_cloud_account_output" {
  value = data.aci_tenant_to_cloud_account.aci_tenant_to_cloud_account_data
}


# TESTING access_type = "credentials"
resource "aci_tenant" "azure_cloud_tenant" {
  description = "sample aci_tenant from terraform"
  name        = "azure_cloud_tenant"
}

resource "aci_cloud_ad" "azure_ad" {
  tenant_dn   = aci_tenant.azure_cloud_tenant.id
  active_directory_id = "azure_ad_id"
  # name      = "azure_ad"
}

resource "aci_cloud_credentials" "azure_credentials" {
  tenant_dn   = aci_cloud_ad.azure_ad.tenant_dn
  key_id = "azure_cred_id"
  name      = "test_cred"
  # key = "secretkey"  # client secret
  relation_cloud_rs_ad= aci_cloud_ad.azure_ad.id
}

resource "aci_cloud_account" "azure_cloud_account" {
  depends_on               = [aci_cloud_credentials.azure_credentials]
  tenant_dn   = aci_cloud_credentials.azure_credentials.tenant_dn
  access_type = "credentials" 
  account_id  = "example_id"
  vendor      = "azure"
  cloud_credentials_dn = aci_cloud_credentials.azure_credentials.id
}

# tenant_to_cloud_account_association
resource "aci_tenant_to_cloud_account" "tenant_to_azure_cloud_account" {
  tenant_dn        = aci_tenant.azure_cloud_tenant.id
  cloud_account_dn = aci_cloud_account.azure_cloud_account.id
}
