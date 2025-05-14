
resource "aci_external_epg" "full_example_l3_outside" {
  parent_dn              = aci_l3_outside.example.id
  annotation             = "annotation"
  description            = "description_1"
  contract_exception_tag = "contract_exception_tag_1"
  flood_in_encapsulation = "disabled"
  match_criteria         = "All"
  name                   = "test_name"
  name_alias             = "name_alias_1"
  intra_epg_isolation    = "enforced"
  preferred_group_member = "exclude"
  priority               = "level1"
  target_dscp            = "AF11"
  relation_to_consumed_contracts = [
    {
      annotation    = "annotation_1"
      priority      = "level1"
      contract_name = aci_contract.example.name
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
  relation_to_imported_contracts = [
    {
      annotation             = "annotation_1"
      priority               = "level1"
      imported_contract_name = aci_imported_contract.example.name
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
  relation_to_custom_qos_policy = {
    annotation             = "annotation_1"
    custom_qos_policy_name = aci_custom_qos_policy.example.name
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
  relation_to_intra_epg_contracts = [
    {
      annotation    = "annotation_1"
      contract_name = aci_contract.example.name
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
  relation_to_taboo_contracts = [
    {
      annotation          = "annotation_1"
      taboo_contract_name = aci_taboo_contract.example.name
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
  relation_to_provided_contracts = [
    {
      annotation     = "annotation_1"
      match_criteria = "All"
      priority       = "level1"
      contract_name  = aci_contract.example.name
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
  relation_to_contract_masters = [
    {
      annotation = "annotation_1"
      target_dn  = aci_external_epg.test_external_epg_0.id
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
  relation_to_route_control_profiles = [
    {
      annotation                 = "annotation_1"
      direction                  = "export"
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
