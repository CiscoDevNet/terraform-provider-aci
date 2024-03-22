
resource "aci_application_epg" "full_example_application_profile" {
  parent_dn              = aci_application_profile.example.id
  annotation             = "annotation"
  description            = "description_1"
  contract_exception_tag = "contract_exception_tag_1"
  flood_in_encapsulation = "disabled"
  forwarding_control     = "proxy-arp"
  has_multicast_source   = "no"
  useg_epg               = "no"
  match_criteria         = "All"
  name                   = "test_name"
  name_alias             = "name_alias_1"
  intra_epg_isolation    = "enforced"
  preferred_group_member = "exclude"
  priority               = "level1"
  admin_state            = "no"
  epg_useg_block_statement = {
    annotation  = "annotation_0"
    description = "description_0"
    match       = "match_0"
    name        = "name_0"
    name_alias  = "name_alias_0"
    owner_key   = "owner_key_0"
    owner_tag   = "owner_tag_0"
    precedence  = "precedence_0"
    scope       = "scope_0"
  }
  relation_to_application_epg_monitoring_policy = {
    annotation             = "annotation_0"
    monitoring_policy_name = aci_monitoring_policy.example.name
  }
  relation_to_bridge_domain = {
    annotation         = "annotation_0"
    bridge_domain_name = aci_bridge_domain.example.name
  }
  relation_to_consumed_contracts = [
    {
      annotation    = "annotation_0"
      priority      = "priority_0"
      contract_name = aci_contract.example.name
      annotations = [
        {
          key   = "key_0"
          value = "value_0"
        }
      ]
      tags = [
        {
          key   = "key_0"
          value = "value_0"
        }
      ]
    }
  ]
  relation_to_imported_contracts = [
    {
      annotation             = "annotation_0"
      priority               = "priority_0"
      imported_contract_name = aci_imported_contract.example.name
      annotations = [
        {
          key   = "key_0"
          value = "value_0"
        }
      ]
      tags = [
        {
          key   = "key_0"
          value = "value_0"
        }
      ]
    }
  ]
  relation_to_custom_qos_policy = {
    annotation             = "annotation_0"
    custom_qos_policy_name = aci_custom_qos_policy.example.name
  }
  relation_to_domains = [
    {
      annotation                    = "annotation_0"
      binding_type                  = "binding_type_0"
      class_preference              = "class_preference_0"
      custom_epg_name               = "custom_epg_name_0"
      delimiter                     = "delimiter_0"
      encapsulation                 = "encapsulation_0"
      encapsulation_mode            = "encapsulation_mode_0"
      epg_cos                       = "epg_cos_0"
      epg_cos_pref                  = "epg_cos_pref_0"
      deployment_immediacy          = "deployment_immediacy_0"
      ipam_dhcp_override            = "10.0.0.2"
      ipam_enabled                  = "yes"
      ipam_gateway                  = "10.0.0.1"
      lag_policy_name               = "lag_policy_name_0"
      netflow_direction             = "netflow_direction_0"
      enable_netflow                = "enable_netflow_0"
      number_of_ports               = "number_of_ports_0"
      port_allocation               = "port_allocation_0"
      primary_encapsulation         = "primary_encapsulation_0"
      primary_encapsulation_inner   = "primary_encapsulation_inner_0"
      resolution_immediacy          = "resolution_immediacy_0"
      secondary_encapsulation_inner = "secondary_encapsulation_inner_0"
      switching_mode                = "switching_mode_0"
      untagged                      = "untagged_0"
      annotations = [
        {
          key   = "key_0"
          value = "value_0"
        }
      ]
      tags = [
        {
          key   = "key_0"
          value = "value_0"
        }
      ]
    }
  ]
  relation_to_data_plane_policing_policy = {
    annotation                      = "annotation_0"
    data_plane_policing_policy_name = aci_data_plane_policing_policy.example.name
  }
  relation_to_fibre_channel_paths = [
    {
      annotation  = "annotation_0"
      description = "description_0"
      vsan        = "vsan_0"
      vsan_mode   = "vsan_mode_0"
      annotations = [
        {
          key   = "key_0"
          value = "value_0"
        }
      ]
      tags = [
        {
          key   = "key_0"
          value = "value_0"
        }
      ]
    }
  ]
  relation_to_intra_epg_contracts = [
    {
      annotation    = "annotation_0"
      contract_name = aci_contract.example.name
      annotations = [
        {
          key   = "key_0"
          value = "value_0"
        }
      ]
      tags = [
        {
          key   = "key_0"
          value = "value_0"
        }
      ]
    }
  ]
  relation_to_static_leafs = [
    {
      annotation           = "annotation_0"
      description          = "description_0"
      encapsulation        = "encapsulation_0"
      deployment_immediacy = "deployment_immediacy_0"
      mode                 = "mode_0"
      annotations = [
        {
          key   = "key_0"
          value = "value_0"
        }
      ]
      tags = [
        {
          key   = "key_0"
          value = "value_0"
        }
      ]
    }
  ]
  relation_to_static_paths = [
    {
      annotation            = "annotation_0"
      description           = "description_0"
      encapsulation         = "encapsulation_0"
      deployment_immediacy  = "deployment_immediacy_0"
      mode                  = "mode_0"
      primary_encapsulation = "primary_encapsulation_0"
      annotations = [
        {
          key   = "key_0"
          value = "value_0"
        }
      ]
      tags = [
        {
          key   = "key_0"
          value = "value_0"
        }
      ]
    }
  ]
  relation_to_taboo_contracts = [
    {
      annotation          = "annotation_0"
      taboo_contract_name = aci_taboo_contract.example.name
      annotations = [
        {
          key   = "key_0"
          value = "value_0"
        }
      ]
      tags = [
        {
          key   = "key_0"
          value = "value_0"
        }
      ]
    }
  ]
  relation_to_provided_contracts = [
    {
      annotation     = "annotation_0"
      match_criteria = "match_criteria_0"
      priority       = "priority_0"
      contract_name  = aci_contract.example.name
      annotations = [
        {
          key   = "key_0"
          value = "value_0"
        }
      ]
      tags = [
        {
          key   = "key_0"
          value = "value_0"
        }
      ]
    }
  ]
  relation_to_contract_masters = [
    {
      annotation = "annotation_0"
      annotations = [
        {
          key   = "key_0"
          value = "value_0"
        }
      ]
      tags = [
        {
          key   = "key_0"
          value = "value_0"
        }
      ]
    }
  ]
  relation_to_trust_control_policy = {
    annotation                = "annotation_0"
    trust_control_policy_name = aci_trust_control_policy.example.name
  }
  annotations = [
    {
      key   = "key_0"
      value = "value_0"
    }
  ]
  tags = [
    {
      key   = "key_0"
      value = "value_0"
    }
  ]
}




