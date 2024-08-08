
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
  epg_useg_block_statement = [
    {
      annotation  = "annotation_1"
      description = "description_1"
      match       = "all"
      name        = "criterion"
      name_alias  = "name_alias_1"
      owner_key   = "owner_key_1"
      owner_tag   = "owner_tag_1"
      precedence  = "1"
      scope       = "scope-bd"
    }
  ]
  relation_to_application_epg_monitoring_policy = [
    {
      annotation             = "annotation_1"
      monitoring_policy_name = aci_monitoring_policy.example.name
    }
  ]
  relation_to_bridge_domain = [
    {
      annotation         = "annotation_1"
      bridge_domain_name = aci_bridge_domain.example.name
    }
  ]
  relation_to_consumed_contracts = [
    {
      annotation    = "annotation_1"
      priority      = "level1"
      contract_name = aci_contract.example.name
    }
  ]
  relation_to_imported_contracts = [
    {
      annotation             = "annotation_1"
      priority               = "level1"
      imported_contract_name = aci_imported_contract.example.name
    }
  ]
  relation_to_custom_qos_policy = [
    {
      annotation             = "annotation_1"
      custom_qos_policy_name = aci_custom_qos_policy.example.name
    }
  ]
  relation_to_domains = [
    {
      annotation                    = "annotation_1"
      binding_type                  = "dynamicBinding"
      class_preference              = "encap"
      custom_epg_name               = "custom_epg_name_1"
      delimiter                     = "@"
      encapsulation                 = "vlan-100"
      encapsulation_mode            = "auto"
      epg_cos                       = "Cos0"
      epg_cos_pref                  = "disabled"
      deployment_immediacy          = "immediate"
      ipam_dhcp_override            = "10.0.0.2"
      ipam_enabled                  = "yes"
      ipam_gateway                  = "10.0.0.1"
      lag_policy_name               = "lag_policy_name_1"
      netflow_direction             = "both"
      enable_netflow                = "disabled"
      number_of_ports               = "1"
      port_allocation               = "elastic"
      primary_encapsulation         = "vlan-200"
      primary_encapsulation_inner   = "vlan-300"
      resolution_immediacy          = "immediate"
      secondary_encapsulation_inner = "vlan-400"
      switching_mode                = "AVE"
      target_dn                     = "uni/vmmp-VMware/dom-domain_1"
      untagged                      = "no"
      vnet_only                     = "no"
    }
  ]
  relation_to_data_plane_policing_policy = [
    {
      annotation                      = "annotation_1"
      data_plane_policing_policy_name = aci_data_plane_policing_policy.example.name
    }
  ]
  relation_to_fibre_channel_paths = [
    {
      annotation  = "annotation_1"
      description = "description_1"
      target_dn   = "topology/pod-1/paths-101/pathep-[eth1/1]"
      vsan        = "vsan-10"
      vsan_mode   = "native"
    }
  ]
  relation_to_intra_epg_contracts = [
    {
      annotation    = "annotation_1"
      contract_name = aci_contract.example.name
    }
  ]
  relation_to_static_leafs = [
    {
      annotation           = "annotation_1"
      description          = "description_1"
      encapsulation        = "vlan-100"
      deployment_immediacy = "immediate"
      mode                 = "native"
      target_dn            = "topology/pod-1/node-101"
    }
  ]
  relation_to_static_paths = [
    {
      annotation            = "annotation_1"
      description           = "description_1"
      encapsulation         = "vlan-202"
      deployment_immediacy  = "immediate"
      mode                  = "native"
      primary_encapsulation = "vlan-203"
      target_dn             = "topology/pod-1/paths-101/pathep-[eth1/1]"
    }
  ]
  relation_to_taboo_contracts = [
    {
      annotation          = "annotation_1"
      taboo_contract_name = aci_taboo_contract.example.name
    }
  ]
  relation_to_provided_contracts = [
    {
      annotation     = "annotation_1"
      match_criteria = "All"
      priority       = "level1"
      contract_name  = aci_contract.example.name
    }
  ]
  relation_to_contract_masters = [
    {
      annotation = "annotation_1"
      target_dn  = aci_application_epg.test_application_epg_0.id
    }
  ]
  relation_to_trust_control_policy = [
    {
      annotation                = "annotation_1"
      trust_control_policy_name = aci_trust_control_policy.example.name
    }
  ]
  associated_site = [
    {
      annotation  = "annotation_1"
      description = "description_1"
      name        = "name_1"
      name_alias  = "name_alias_1"
      owner_key   = "owner_key_1"
      owner_tag   = "owner_tag_1"
      site_id     = "0"
    }
  ]
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
