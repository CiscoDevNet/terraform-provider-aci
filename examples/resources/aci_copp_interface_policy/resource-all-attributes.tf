
resource "aci_copp_interface_policy" "full_example" {
  annotation  = "annotation"
  description = "description_1"
  name        = "test_name"
  name_alias  = "name_alias_1"
  owner_key   = "owner_key_1"
  owner_tag   = "owner_tag_1"
  protocol_policies = [
    {
      annotation      = "annotation_1"
      burst           = "15"
      description     = "description_1"
      match_protocols = ["arp", "bfd"]
      name            = "name_0"
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
