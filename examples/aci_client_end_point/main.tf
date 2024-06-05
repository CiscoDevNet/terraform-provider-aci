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

data "aci_client_end_point" "this" {
  ip                 = "10.1.1.33"
  allow_empty_result = true
}

output "mac_endpoints" {
  value = data.aci_client_end_point.this.fvcep_objects
}

data "aci_tenant" "tenant" {
  name = "tenant_name"
}

data "aci_application_profile" "application_profile" {
  tenant_dn = data.aci_tenant.tenant.id
  name      = "application_name"
}

data "aci_application_epg" "application_epg" {
  application_profile_dn = data.aci_application_profile.application_profile.id
  name                   = "epg_name"
}

# Filter based on Tenants

data "aci_client_end_point" "end_points" {
  filter_dn = data.aci_tenant.tenant.id
}

# Filter based on Application Profiles

data "aci_client_end_point" "end_points" {
  filter_dn = data.aci_application_profile.application_profile.id
}

# Filter based on EPGs

data "aci_client_end_point" "end_points" {
  filter_dn = data.aci_application_epg.application_epg.id
}