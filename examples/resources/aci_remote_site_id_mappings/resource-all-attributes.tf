
resource "aci_remote_site_id_mappings" "full_example_associated_site" {
  parent_dn         = aci_associated_site.example.id
  annotation        = "annotation"
  description       = "description_1"
  name              = "name_1"
  name_alias        = "name_alias_1"
  owner_key         = "owner_key_1"
  owner_tag         = "owner_tag_1"
  remote_vrf_pc_tag = "any"
  remote_pc_tag     = "16386"
  site_id           = "100"
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
