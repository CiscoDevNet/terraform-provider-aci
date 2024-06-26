
resource "aci_remote_site_id_mappings" "full_example_associated_site" {
  parent_dn         = aci_associated_site.example.id
  annotation        = "annotation"
  description       = "description"
  name              = "name"
  name_alias        = "name_alias"
  owner_key         = "owner_key"
  owner_tag         = "owner_tag"
  remote_vrf_pc_tag = "remote_vrf_pc_tag"
  remote_pc_tag     = "remote_pc_tag"
  site_id           = "0"
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
