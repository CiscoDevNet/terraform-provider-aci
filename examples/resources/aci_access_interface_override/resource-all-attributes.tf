
resource "aci_access_interface_override" "full_example" {
  annotation  = "annotation"
  description = "description_1"
  name        = "host_path_selector"
  name_alias  = "name_alias_1"
  owner_key   = "owner_key_1"
  owner_tag   = "owner_tag_1"
  relation_to_host_path = {
    annotation = "annotation_0"
    target_dn  = "target_dn_0"
  }
  relation_to_access_interface_policy_group = {
    annotation = "annotation_0"
    target_dn  = "target_dn_0"
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




