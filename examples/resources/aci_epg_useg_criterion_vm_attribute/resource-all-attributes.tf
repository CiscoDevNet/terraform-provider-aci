
resource "aci_epg_useg_criterion_vm_attribute" "full_example_epg_useg_criterion" {
  parent_dn   = aci_epg_useg_criterion.example.id
  annotation  = "annotation"
  category    = "all_category"
  description = "description"
  label_name  = "label_name"
  name        = "vm_attribute"
  name_alias  = "name_alias"
  operator    = "contains"
  owner_key   = "owner_key"
  owner_tag   = "owner_tag"
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

resource "aci_epg_useg_criterion_vm_attribute" "full_example_epg_useg_sub_criterion" {
  parent_dn   = aci_epg_useg_sub_criterion.example.id
  annotation  = "annotation"
  category    = "all_category"
  description = "description"
  label_name  = "label_name"
  name        = "vm_attribute"
  name_alias  = "name_alias"
  operator    = "contains"
  owner_key   = "owner_key"
  owner_tag   = "owner_tag"
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
