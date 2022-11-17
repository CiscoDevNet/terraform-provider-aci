terraform {
  required_providers {
    aci = {
      source = "ciscodevnet/aci"
    }
  }
}

provider "aci" {
  username = "" # <APIC username>
  password = "" # <APIC pwd>
  url      = "" # <cloud APIC URL>
  insecure = true
}


resource "aci_tenant" "src_tenant" {
  name = "src_tenant"
}

resource "aci_tenant" "dst_tenant" {
  name = "dst_tenant"
}

resource "aci_vrf" "src_vrf" {
  tenant_dn = aci_tenant.src_tenant.id
  name      = "src_vrf"
}

resource "aci_vrf" "dst_vrf1" {
  tenant_dn = aci_tenant.dst_tenant.id
  name      = "dst_vrf1"
}

resource "aci_vrf" "dst_vrf2" {
  tenant_dn = aci_tenant.dst_tenant.id
  name      = "dst_vrf2"
}

resource "aci_cloud_context_profile" "src_vrf_ctx_prof" {
  name                     = "src_vrf_ctx_prof"
  tenant_dn                = aci_tenant.src_tenant.id
  primary_cidr             = "10.10.10.1/24"
  region                   = "asia-east1"
  cloud_vendor             = "gcp"
  relation_cloud_rs_to_ctx = aci_vrf.src_vrf.id
}

resource "aci_cloud_context_profile" "dst_vrf1_ctx_prof" {
  name                     = "dst_vrf1_ctx_prof"
  tenant_dn                = aci_tenant.dst_tenant.id
  primary_cidr             = "10.10.10.2/24"
  region                   = "asia-east1"
  cloud_vendor             = "gcp"
  relation_cloud_rs_to_ctx = aci_vrf.dst_vrf1.id
}

resource "aci_cloud_context_profile" "dst_vrf2_ctx_prof" {
  name                     = "dst_vrf2_ctx_prof"
  tenant_dn                = aci_tenant.dst_tenant.id
  primary_cidr             = "10.10.10.3/24"
  region                   = "asia-east1"
  cloud_vendor             = "gcp"
  relation_cloud_rs_to_ctx = aci_vrf.dst_vrf2.id
}

# Only for the Cloud APIC Version >= 25.0
resource "aci_cloud_vrf_leak_routes" "cloud_internal_leak_routes" {
  vrf_dn = aci_vrf.src_vrf.id
  leak_to {
    vrf_dn = aci_vrf.dst_vrf1.id
  }
  leak_to {
    vrf_dn = aci_vrf.dst_vrf2.id
  }
}
