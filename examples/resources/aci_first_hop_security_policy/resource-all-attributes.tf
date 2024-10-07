
resource "aci_first_hop_security_policy" "full_example_tenant" {
  parent_dn            = aci_tenant.example.id
  annotation           = "annotation"
  description          = "description_1"
  ip_inspection        = "disabled"
  name                 = "test_name"
  name_alias           = "name_alias_1"
  owner_key            = "owner_key_1"
  owner_tag            = "owner_tag_1"
  router_advertisement = "disabled"
  source_guard         = "disabled"
  route_advertisement_guard_policy = {
    annotation           = "annotation_1"
    description          = "description_1"
    managed_config_check = "no"
    managed_config_flag  = "no"
    max_hop_limit        = "10"
    max_router_pref      = "disabled"
    min_hop_limit        = "1"
    name                 = "name_1"
    name_alias           = "name_alias_1"
    other_config_check   = "no"
    other_config_flag    = "no"
    owner_key            = "owner_key_1"
    owner_tag            = "owner_tag_1"
  }
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
