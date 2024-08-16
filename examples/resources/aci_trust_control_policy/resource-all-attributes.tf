
resource "aci_trust_control_policy" "full_example_tenant" {
  parent_dn         = aci_tenant.example.id
  annotation        = "annotation"
  description       = "description_1"
  has_dhcpv4_server = "no"
  has_dhcpv6_server = "no"
  has_ipv6_router   = "no"
  name              = "test_name"
  name_alias        = "name_alias_1"
  owner_key         = "owner_key_1"
  owner_tag         = "owner_tag_1"
  trust_arp         = "no"
  trust_nd          = "no"
  trust_ra          = "no"
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
