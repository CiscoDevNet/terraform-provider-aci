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

resource "aci_tenant" "tf_tenant" {
  name = "tf_tenant"
}

# Cloud Subnet setup part
resource "aci_vrf" "vrf1" {
  tenant_dn = aci_tenant.tf_tenant.id
  name      = "vrf-1"
}

resource "aci_cloud_context_profile" "ctx1" {
  name                     = "tf_ctx1"
  tenant_dn                = aci_tenant.tf_tenant.id
  primary_cidr             = "10.1.0.0/16"
  region                   = "westus"
  cloud_vendor             = "azure"
  relation_cloud_rs_to_ctx = aci_vrf.vrf1.id
  hub_network              = "uni/tn-infra/gwrouterp-default"
}

resource "aci_cloud_cidr_pool" "cloud_cidr_pool" {
  cloud_context_profile_dn = aci_cloud_context_profile.ctx1.id
  addr                     = "10.1.0.0/16"
}

data "aci_cloud_provider_profile" "cloud_profile" {
  vendor = "azure"
}

data "aci_cloud_providers_region" "cloud_region" {
  cloud_provider_profile_dn = data.aci_cloud_provider_profile.cloud_profile.id
  name                      = "westus"
}

data "aci_cloud_availability_zone" "region_availability_zone" {
  cloud_providers_region_dn = data.aci_cloud_providers_region.cloud_region.id
  name                      = "default"
}

resource "aci_cloud_subnet" "cloud_subnet" {
  cloud_cidr_pool_dn = aci_cloud_cidr_pool.cloud_cidr_pool.id
  ip                 = "10.1.1.0/24"
  usage              = "gateway"
  zone               = data.aci_cloud_availability_zone.region_availability_zone.id
  scope              = ["shared", "private", "public"]
}

# AAA Domain setup part
resource "aci_aaa_domain" "aaa_domain_1" {
  name = "aaa_domain_1"
}

resource "aci_aaa_domain" "aaa_domain_2" {
  name = "aaa_domain_2"
}

# Application Load Balancer
resource "aci_cloud_l4_l7_native_load_balancer" "cloud_native_alb" {
  tenant_dn = aci_tenant.tf_tenant.id
  name      = "cloud_native_alb"
  aaa_domain_dn = [
    aci_aaa_domain.aaa_domain_1.id,
    aci_aaa_domain.aaa_domain_2.id
  ]
  relation_cloud_rs_ldev_to_cloud_subnet = [
    aci_cloud_subnet.cloud_subnet.id
  ]
  active_active                 = "no"
  allow_all                     = "no"
  auto_scaling                  = "no"
  context_aware                 = "multi-Context"
  device_type                   = "CLOUD"
  function_type                 = "GoTo"
  is_copy                       = "no"
  is_instantiation              = "no"
  is_static_ip                  = "no"
  managed                       = "no"
  mode                          = "legacy-Mode"
  promiscuous_mode              = "no"
  scheme                        = "internal"
  size                          = "medium"
  sku                           = "standard"
  service_type                  = "NATIVELB"
  target_mode                   = "primary"
  trunking                      = "no"
  cloud_l4l7_load_balancer_type = "application"
  instance_count                = "2"
  max_instance_count            = "10"
  min_instance_count            = "5"
}

# Network Load Balancer
resource "aci_cloud_l4_l7_native_load_balancer" "cloud_native_nlb" {
  tenant_dn                     = aci_tenant.tf_tenant.id
  name                          = "cloud_native_nlb"
  service_type                  = "NATIVELB"
  cloud_l4l7_load_balancer_type = "network"
  aaa_domain_dn = [
    aci_aaa_domain.aaa_domain_1.id
  ]
  relation_cloud_rs_ldev_to_cloud_subnet = [
    aci_cloud_subnet.cloud_subnet.id
  ]
}
