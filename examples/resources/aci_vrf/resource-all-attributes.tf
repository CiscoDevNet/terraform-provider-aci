
resource "aci_vrf" "full_example_tenant" {
  parent_dn                             = aci_tenant.example.id
  annotation                            = "annotation"
  bd_enforcement                        = "no"
  description                           = "description_1"
  ip_data_plane_learning                = "disabled"
  name                                  = "test_name"
  name_alias                            = "name_alias_1"
  owner_key                             = "owner_key_1"
  owner_tag                             = "owner_tag_1"
  policy_control_enforcement_direction  = "egress"
  policy_control_enforcement_preference = "enforced"
  relation_to_bgp_timers = {
    annotation      = "annotation_1"
    bgp_timers_name = aci_bgp_timers.example.name
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
  relation_to_monitoring_policy = {
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
  relation_to_bgp_address_family_contexts = [
    {
      address_family                  = "ipv4-ucast"
      annotation                      = "annotation_1"
      bgp_address_family_context_name = aci_bgp_address_family_context.example.name
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
  relation_to_eigrp_address_family_contexts = [
    {
      address_family                    = "ipv4-ucast"
      annotation                        = "annotation_1"
      eigrp_address_family_context_name = aci_eigrp_address_family_context.example.name
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
  relation_to_end_point_retention_policy = {
    annotation                      = "annotation_1"
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
  relation_to_l3out_route_tag_policy = {
    annotation                  = "annotation_1"
    l3out_route_tag_policy_name = aci_l3out_route_tag_policy.example.name
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
  relation_to_address_family_ospf_timers = [
    {
      address_family   = "ipv4-ucast"
      annotation       = "annotation_1"
      ospf_timers_name = aci_ospf_timers.example.name
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
  relation_to_wan_vpn = {
    annotation = "annotation_1"
    target_dn  = "uni/tn-test_tenant/sdwanvpncont/sdwanvpnentry-sdwanvpn_1"
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
  relation_to_ospf_timers = {
    annotation       = "annotation_1"
    ospf_timers_name = aci_ospf_timers.example.name
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
