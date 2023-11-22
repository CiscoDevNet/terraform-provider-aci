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

# VRF setup part
resource "aci_vrf" "vrf1" {
  tenant_dn = aci_tenant.tf_tenant.id
  name      = "vrf-1"
}

# AAA Domain setup part
resource "aci_aaa_domain" "aaa_domain_1" {
  name = "aaa_domain_1"
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

# Application Load Balancer
resource "aci_cloud_l4_l7_native_load_balancer" "cloud_native_alb" {
  tenant_dn = aci_tenant.tf_tenant.id
  name      = "cloud_native_alb"
  aaa_domain_dn = [
    aci_aaa_domain.aaa_domain_1.id
  ]
  relation_cloud_rs_ldev_to_cloud_subnet = [
    aci_cloud_subnet.cloud_subnet.id
  ]
  cloud_l4l7_load_balancer_type = "application"
}

# Third-Party Firewall
resource "aci_cloud_l4_l7_third_party_device" "cloud_third_party_fw" {
  tenant_dn    = aci_tenant.tf_tenant.id
  name         = "cloud_third_party_fw"
  service_type = "FW"

  aaa_domain_dn = [
    aci_aaa_domain.aaa_domain_1.id
  ]
  relation_cloud_rs_ldev_to_ctx = aci_vrf.vrf1.id

  interface_selectors {
    allow_all = "no"
    name      = "Interface_1"
    end_point_selectors {
      match_expression = "IP=='1.1.1.21/24'"
      name             = "Interface_1_ep_1"
    }
    end_point_selectors {
      match_expression = "custom:Name1=='admin-ep1'"
      name             = "Interface_1_ep_2"
    }
  }
  interface_selectors {
    allow_all = "no"
    name      = "Interface_2"
    end_point_selectors {
      match_expression = "IP=='1.1.1.21/24'"
      name             = "Interface_2_ep_1"
    }
    end_point_selectors {
      match_expression = "custom:Name1=='admin-ep1'"
      name             = "Interface_2_ep_2"
    }
  }
}

# Third-Party Load Balancer
resource "aci_cloud_l4_l7_third_party_device" "cloud_third_party_lb" {
  tenant_dn    = aci_tenant.tf_tenant.id
  name         = "cloud_third_party_lb"
  service_type = "ADC"

  aaa_domain_dn = [
    aci_aaa_domain.aaa_domain_1.id
  ]
  relation_cloud_rs_ldev_to_ctx = aci_vrf.vrf1.id

  interface_selectors {
    allow_all = "no"
    name      = "Interface_1"
    end_point_selectors {
      match_expression = "IP=='1.1.1.21/24'"
      name             = "Interface_1_ep_1"
    }
    end_point_selectors {
      match_expression = "custom:Name1=='admin-ep1'"
      name             = "Interface_1_ep_2"
    }
  }
}

# Service Graph Part
resource "aci_l4_l7_service_graph_template" "cloud_service_graph" {
  tenant_dn                         = aci_tenant.tf_tenant.id
  name                              = "cloud_service_graph"
  l4_l7_service_graph_template_type = "cloud"
}

resource "aci_function_node" "function_node_0" {
  l4_l7_service_graph_template_dn     = aci_l4_l7_service_graph_template.cloud_service_graph.id
  name                                = "N0"
  func_template_type                  = "ADC_ONE_ARM"
  managed                             = "yes"
  relation_vns_rs_node_to_cloud_l_dev = aci_cloud_l4_l7_native_load_balancer.cloud_native_alb.id
}

resource "aci_function_node" "function_node_1" {
  l4_l7_service_graph_template_dn     = aci_l4_l7_service_graph_template.cloud_service_graph.id
  name                                = "N2"
  func_template_type                  = "OTHER"
  managed                             = "no"
  relation_vns_rs_node_to_cloud_l_dev = aci_cloud_l4_l7_third_party_device.cloud_third_party_lb.id
}

resource "aci_function_node" "function_node_2" {
  l4_l7_service_graph_template_dn      = aci_l4_l7_service_graph_template.cloud_service_graph.id
  name                                 = "N1"
  func_template_type                   = "FW_ROUTED"
  managed                              = "no"
  relation_vns_rs_node_to_cloud_l_dev  = aci_cloud_l4_l7_third_party_device.cloud_third_party_fw.id
  l4_l7_device_interface_consumer_name = "Interface_1"
  l4_l7_device_interface_provider_name = "Interface_2"
}

resource "aci_connection" "consumer" {
  l4_l7_service_graph_template_dn = aci_l4_l7_service_graph_template.cloud_service_graph.id
  name                            = "CON0"
  adj_type                        = "L3"
  conn_dir                        = "consumer"
  conn_type                       = "external"
  direct_connect                  = "yes"
  unicast_route                   = "yes"
  relation_vns_rs_abs_connection_conns = [
    aci_l4_l7_service_graph_template.cloud_service_graph.term_cons_dn,
    aci_function_node.function_node_0.conn_consumer_dn,
  ]
}

resource "aci_connection" "consumer_provider_1" {
  l4_l7_service_graph_template_dn = aci_l4_l7_service_graph_template.cloud_service_graph.id
  name                            = "CON1"
  adj_type                        = "L3"
  conn_type                       = "external"
  direct_connect                  = "yes"
  unicast_route                   = "yes"
  relation_vns_rs_abs_connection_conns = [
    aci_function_node.function_node_1.conn_consumer_dn,
    aci_function_node.function_node_0.conn_provider_dn
  ]
}

resource "aci_connection" "consumer_provider_2" {
  l4_l7_service_graph_template_dn = aci_l4_l7_service_graph_template.cloud_service_graph.id
  name                            = "CON2"
  adj_type                        = "L3"
  conn_type                       = "external"
  direct_connect                  = "yes"
  unicast_route                   = "yes"
  relation_vns_rs_abs_connection_conns = [
    aci_function_node.function_node_1.conn_provider_dn,
    aci_function_node.function_node_2.conn_consumer_dn
  ]
}

resource "aci_connection" "provider" {
  l4_l7_service_graph_template_dn = aci_l4_l7_service_graph_template.cloud_service_graph.id
  name                            = "CON3"
  adj_type                        = "L3"
  conn_dir                        = "provider"
  conn_type                       = "external"
  direct_connect                  = "yes"
  unicast_route                   = "yes"
  relation_vns_rs_abs_connection_conns = [
    aci_l4_l7_service_graph_template.cloud_service_graph.term_prov_dn,
    aci_function_node.function_node_2.conn_provider_dn
  ]
}
