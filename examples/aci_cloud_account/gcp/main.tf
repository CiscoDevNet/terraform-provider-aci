terraform {
  required_providers {
    aci = {
      source = "ciscodevnet/aci"
    }
  }
}

# GCP cloud APIC
provider "aci" {
  username = ""
  password = ""
  url      = ""
  insecure = true
}

# 1. Creating a GCP cloud account
resource "aci_tenant" "gcp_cloud_tenant" {
  description = "sample aci_tenant from terraform"
  name        = "gcp_test_tenant"
}

# access_type = "managed" 
resource "aci_cloud_account" "foo_cloud_account" {
  tenant_dn   = aci_tenant.gcp_cloud_tenant.id
  access_type = "credentials" //credentials or managed
  account_id  = "example_id"
  vendor      = "gcp"
}


# access_type = "credentials" (unmanaged)
resource "aci_tenant" "gcp_tenant" {
  description = "sample aci_tenant from terraform"
  name        = "gcp_tenant"
}

resource "aci_cloud_account" "gcp_cloud_account" {
  tenant_dn   = aci_cloud_credentials.gcp_credentials.tenant_dn
  access_type = "credentials" //credentials or managed
  account_id  = "example_id"
  vendor      = "gcp"
  cloud_credentials_dn = aci_cloud_credentials.gcp_credentials.id
}

resource "aci_cloud_credentials" "gcp_credentials" {
  tenant_dn   = aci_tenant.gcp_tenant.id
  key_id = "gcp_cred_id"
  name      = "gcp_test_cred"
  email= "anvjain@cisco.com"
}