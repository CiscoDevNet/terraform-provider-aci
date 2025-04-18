
resource "aci_associated_site" "full_example_application_epg" {
  parent_dn   = aci_application_epg.example.id
  annotation  = "annotation"
  description = "description_1"
  name        = "name_1"
  name_alias  = "name_alias_1"
  owner_key   = "owner_key_1"
  owner_tag   = "owner_tag_1"
  site_id     = "0"
  remote_sites = [
    {
      annotation        = "annotation_1"
      description       = "description_1"
      name              = "name_1"
      name_alias        = "name_alias_1"
      owner_key         = "owner_key_1"
      owner_tag         = "owner_tag_1"
      remote_vrf_pc_tag = "any"
      remote_pc_tag     = "16386"
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

resource "aci_associated_site" "full_example_bridge_domain" {
  parent_dn   = aci_bridge_domain.example.id
  annotation  = "annotation"
  description = "description_1"
  name        = "name_1"
  name_alias  = "name_alias_1"
  owner_key   = "owner_key_1"
  owner_tag   = "owner_tag_1"
  site_id     = "0"
  remote_sites = [
    {
      annotation        = "annotation_1"
      description       = "description_1"
      name              = "name_1"
      name_alias        = "name_alias_1"
      owner_key         = "owner_key_1"
      owner_tag         = "owner_tag_1"
      remote_vrf_pc_tag = "any"
      remote_pc_tag     = "16386"
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
