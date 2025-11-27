
resource "aci_copp_interface_protocol_policy" "full_example_copp_interface_policy" {
  parent_dn       = aci_copp_interface_policy.example.id
  annotation      = "annotation"
  burst           = "15"
  description     = "description_1"
  match_protocols = ["arp", "bfd"]
  name            = "test_name"
  name_alias      = "name_alias_1"
  owner_key       = "owner_key_1"
  owner_tag       = "owner_tag_1"
  rate            = "15"
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
