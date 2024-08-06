
resource "aci_access_interface_override" "full_example" {
  annotation  = "annotation"
  description = "description_1"
  name        = "host_path_selector"
  name_alias  = "name_alias_1"
  owner_key   = "owner_key_1"
  owner_tag   = "owner_tag_1"
  relation_to_host_paths = [
    {
      annotation = "annotation_1"
      target_dn  = "topology/pod-1/paths-101/pathep-[eth1/1]"
    }
  ]
  relation_to_access_interface_policy_group = [
    {
      annotation = "annotation_1"
      target_dn  = "uni/infra/funcprof/accportgrp-access_interface_policy_group"
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
