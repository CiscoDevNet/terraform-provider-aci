
resource "aci_epg_useg_vm_attribute" "full_example_epg_useg_block_statement" {
  parent_dn   = aci_epg_useg_block_statement.example.id
  annotation  = "annotation"
  category    = "all_category"
  description = "description_1"
  label_name  = "label_name"
  name        = "vm_attribute"
  name_alias  = "name_alias_1"
  operator    = "contains"
  owner_key   = "owner_key_1"
  owner_tag   = "owner_tag_1"
  type        = "domain"
  value       = "default_value"
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

resource "aci_epg_useg_vm_attribute" "full_example_epg_useg_sub_block_statement" {
  parent_dn   = aci_epg_useg_sub_block_statement.example.id
  annotation  = "annotation"
  category    = "all_category"
  description = "description_1"
  label_name  = "label_name"
  name        = "vm_attribute"
  name_alias  = "name_alias_1"
  operator    = "contains"
  owner_key   = "owner_key_1"
  owner_tag   = "owner_tag_1"
  type        = "domain"
  value       = "default_value"
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
