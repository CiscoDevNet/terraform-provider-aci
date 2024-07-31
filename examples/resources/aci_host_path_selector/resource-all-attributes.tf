
resource "aci_host_path_selector" "full_example" {
  annotation  = "annotation"
  description = "description_1"
  name        = "host_path_selector"
  name_alias  = "name_alias_1"
  owner_key   = "owner_key_1"
  owner_tag   = "owner_tag_1"
  relation_to_host_paths = [
    {
      annotation = "annotation_1"
      target_dn  = "target_dn_0"
    }
  ]
  relation_to_access_base_group = [
    {
      annotation = "annotation_1"
      target_dn  = "target_dn_1"
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
