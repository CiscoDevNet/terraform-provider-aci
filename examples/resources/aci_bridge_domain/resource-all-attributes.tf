
resource "aci_bridge_domain" "full_example_tenant" {
  parent_dn                          = aci_tenant.example.id
  optimize_wan_bandwidth             = "no"
  annotation                         = "annotation"
  arp_flooding                       = "no"
  description                        = "description_1"
  enable_rogue_exception_mac         = "no"
  clear_remote_mac_entries           = "no"
  endpoint_move_detection_mode       = "garp"
  advertise_host_routes              = "no"
  enable_intersite_bum_traffic       = "no"
  intersite_l2_stretch               = "no"
  ip_learning                        = "no"
  pim_ipv6                           = "no"
  limit_ip_learn_to_subnets          = "yes"
  link_local_ipv6_address            = "fe80::1"
  custom_mac_address                 = "00:22:BD:F8:19:FE"
  drop_arp_with_multicast_smac       = "no"
  pim                                = "no"
  multi_destination_flooding         = "bd-flood"
  name                               = "test_name"
  name_alias                         = "name_alias_1"
  owner_key                          = "owner_key_1"
  owner_tag                          = "owner_tag_1"
  service_bd_routing_disable         = "no"
  bridge_domain_type                 = "fc"
  unicast_routing                    = "no"
  l2_unknown_unicast_flooding        = "proxy"
  l3_unknown_multicast_flooding      = "flood"
  ipv6_l3_unknown_multicast_flooding = "flood"
  virtual_mac_address                = "00:22:BD:F8:19:FB"
  legacy_mode = {
    annotation    = "annotation_1"
    description   = "description_1"
    encapsulation = "vlan-100"
    name          = "name_1"
    name_alias    = "name_alias_1"
    owner_key     = "owner_key_1"
    owner_tag     = "owner_tag_1"
    annotations = [
      {
        key   = "key_0"
        value = "value_1"
      }
    ]
    tags = [
      {
        key   = "key_0"
        value = "value_1"
      }
    ]
  }
  rogue_coop_exceptions = [
    {
      annotation  = "annotation_1"
      description = "description_1"
      mac         = "00:00:00:00:00:00"
      name        = "name_1"
      name_alias  = "name_alias_1"
      annotations = [
        {
          key   = "key_0"
          value = "value_1"
        }
      ]
      tags = [
        {
          key   = "key_0"
          value = "value_1"
        }
      ]
    }
  ]
  relation_to_monitor_policy = {
    annotation             = "annotation_1"
    monitoring_policy_name = aci_monitoring_policy.example.name
    annotations = [
      {
        key   = "key_0"
        value = "value_1"
      }
    ]
    tags = [
      {
        key   = "key_0"
        value = "value_1"
      }
    ]
  }
  relation_to_first_hop_security_policy = {
    annotation                     = "annotation_1"
    first_hop_security_policy_name = aci_first_hop_security_policy.example.name
    annotations = [
      {
        key   = "key_0"
        value = "value_1"
      }
    ]
    tags = [
      {
        key   = "key_0"
        value = "value_1"
      }
    ]
  }
  relation_to_neighbor_discovery_interface_policy = {
    annotation                               = "annotation_1"
    neighbor_discovery_interface_policy_name = aci_neighbor_discovery_interface_policy.example.name
    annotations = [
      {
        key   = "key_0"
        value = "value_1"
      }
    ]
    tags = [
      {
        key   = "key_0"
        value = "value_1"
      }
    ]
  }
  relation_to_netflow_monitor_policies = [
    {
      annotation                  = "annotation_1"
      filter_type                 = "ce"
      netflow_monitor_policy_name = aci_netflow_monitor_policy.example.name
      annotations = [
        {
          key   = "key_0"
          value = "value_1"
        }
      ]
      tags = [
        {
          key   = "key_0"
          value = "value_1"
        }
      ]
    }
  ]
  relation_to_l3_outsides = [
    {
      annotation      = "annotation_1"
      l3_outside_name = aci_l3_outside.example.name
      annotations = [
        {
          key   = "key_0"
          value = "value_1"
        }
      ]
      tags = [
        {
          key   = "key_0"
          value = "value_1"
        }
      ]
    }
  ]
  relation_to_route_control_profile = {
    annotation                 = "annotation_1"
    l3_outside_name            = aci_l3_outside.example.name
    route_control_profile_name = aci_route_control_profile.example.name
    annotations = [
      {
        key   = "key_0"
        value = "value_1"
      }
    ]
    tags = [
      {
        key   = "key_0"
        value = "value_1"
      }
    ]
  }
  relation_to_dhcp_relay_policy = {
    annotation             = "annotation_1"
    dhcp_relay_policy_name = aci_dhcp_relay_policy.example.name
    annotations = [
      {
        key   = "key_0"
        value = "value_1"
      }
    ]
    tags = [
      {
        key   = "key_0"
        value = "value_1"
      }
    ]
  }
  relation_to_end_point_retention_policy = {
    annotation                      = "annotation_1"
    resolve_action                  = "inherit"
    end_point_retention_policy_name = aci_end_point_retention_policy.example.name
    annotations = [
      {
        key   = "key_0"
        value = "value_1"
      }
    ]
    tags = [
      {
        key   = "key_0"
        value = "value_1"
      }
    ]
  }
  relation_to_vrf = {
    annotation = "annotation_1"
    vrf_name   = aci_vrf.example.name
    annotations = [
      {
        key   = "key_0"
        value = "value_1"
      }
    ]
    tags = [
      {
        key   = "key_0"
        value = "value_1"
      }
    ]
  }
  relation_to_igmp_snooping_policy = {
    annotation                = "annotation_1"
    igmp_snooping_policy_name = aci_igmp_snooping_policy.example.name
    annotations = [
      {
        key   = "key_0"
        value = "value_1"
      }
    ]
    tags = [
      {
        key   = "key_0"
        value = "value_1"
      }
    ]
  }
  relation_to_mld_snooping_policy = {
    annotation               = "annotation_1"
    mld_snooping_policy_name = aci_mld_snooping_policy.example.name
    annotations = [
      {
        key   = "key_0"
        value = "value_1"
      }
    ]
    tags = [
      {
        key   = "key_0"
        value = "value_1"
      }
    ]
  }
  annotations = [
    {
      key   = "key_0"
      value = "value_1"
    }
  ]
  tags = [
    {
      key   = "key_0"
      value = "value_1"
    }
  ]
}
