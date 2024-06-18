
resource "aci_endpoint_security_group" "full_example_application_profile" {
  parent_dn              = aci_application_profile.example.id
  annotation             = "annotation"
  description            = "description"
  exception_tag          = "exception_tag"
  match_criteria         = "All"
  name                   = "test_name"
  name_alias             = "name_alias"
  intra_esg_isolation    = "enforced"
  preferred_group_member = "exclude"
  admin_state            = "no"
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
  relation_to_intra_epg_contracts = [
    {
      annotation    = "annotation_1"
      contract_name = aci_contract.example.name
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
  relation_to_vrf = [
    {
      annotation = "annotation_1"
      vrf_name   = aci_vrf.example.name
    }
  ]
  relation_to_contract_masters = [
    {
      annotation = "annotation_1"
      target_dn  = "target_dn_0"
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
