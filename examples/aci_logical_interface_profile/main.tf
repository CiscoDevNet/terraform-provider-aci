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
  name        = "l3out_tf_tenant"
  annotation  = "tag_tenant"
  name_alias  = "alias_tenant"
}

resource "aci_vrf" "vrf1" {
  tenant_dn = aci_tenant.footenant.id
  name      = "vrf1"
}

resource "aci_l3_outside" "fool3_outside" {
  tenant_dn              = aci_tenant.footenant.id
  description            = "sample aci_l3_outside"
  name                   = "demo_l3out"
  annotation             = "tag_l3out"
  enforce_rtctrl         = ["export", "import"]
  name_alias             = "alias_out"
  target_dscp            = "unspecified"
  relation_l3ext_rs_ectx = aci_vrf.vrf1.id
}

resource "aci_logical_node_profile" "foological_node_profile" {
  l3_outside_dn = aci_l3_outside.fool3_outside.id
  description   = "sample logical node profile"
  name          = "demo_node"
  annotation    = "tag_node"
  config_issues = "none"
  name_alias    = "alias_node"
  tag           = "black"
  target_dscp   = "unspecified"
}

resource "aci_igmp_interface_policy" "example_igmp" {
  tenant_dn                 = aci_tenant.footenant.id
  name                      = "example_igmp"
  group_timeout             = "260"
  last_member_count         = "2"
  last_member_response_time = "1"
  querier_timeout           = "255"
  query_interval            = "125"
  robustness_variable       = "2"
  response_interval         = "10"
  startup_query_count       = "2"
  startup_query_interval    = "31"
  version                   = "v2"
}

resource "aci_pim_interface_policy" "example_ip" {
  tenant_dn                  = aci_tenant.footenant.id
  name                       = "example_ip"
  designated_router_delay    = "3"
  designated_router_priority = "1"
  hello_interval             = "30000"
  join_prune_interval        = "60"
  control_state              = ["border"]
}

resource "aci_pim_interface_policy" "example_ipv6" {
  tenant_dn                  = aci_tenant.footenant.id
  name                       = "example_ipv6"
  designated_router_delay    = "3"
  designated_router_priority = "1"
  hello_interval             = "30000"
  join_prune_interval        = "60"
  control_state              = ["border"]
}

resource "aci_logical_interface_profile" "foological_interface_profile" {
  logical_node_profile_dn               = aci_logical_node_profile.foological_node_profile.id
  description                           = "aci_logical_interface_profile from terraform"
  name                                  = "demo_int_prof"
  annotation                            = "tag_prof"
  name_alias                            = "alias_prof"
  prio                                  = "unspecified"
  tag                                   = "black"
  relation_l3ext_rs_pim_ip_if_pol       = aci_pim_interface_policy.example_ip.id
  relation_l3ext_rs_pim_ipv6_if_pol     = aci_pim_interface_policy.example_ipv6.id
  relation_l3ext_rs_igmp_if_pol         = aci_igmp_interface_policy.example_igmp.id
  relation_l3ext_rs_egress_qos_dpp_pol  = "uni/tn-l3out_tf_tenant/qosdpppol-egress_data_plane"
  relation_l3ext_rs_ingress_qos_dpp_pol = "uni/tn-l3out_tf_tenant/qosdpppol-ingress_data_plane"
  relation_l3ext_rs_l_if_p_cust_qos_pol = "uni/tn-l3out_tf_tenant/qoscustom-qos"
  relation_l3ext_rs_nd_if_pol           = "uni/tn-l3out_tf_tenant/ndifpol-nd"
  relation_l3ext_rs_l_if_p_to_netflow_monitor_pol {
    tn_netflow_monitor_pol_dn = "uni/tn-l3out_tf_tenant/monitorpol-netflow"
    flt_type                  = "ipv4"
  }
}
