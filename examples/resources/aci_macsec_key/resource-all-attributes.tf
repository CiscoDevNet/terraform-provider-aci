
resource "aci_macsec_key" "full_example_macsec_key_chain" {
  parent_dn      = aci_macsec_key_chain.example.id
  annotation     = "annotation"
  description    = "description_1"
  end_time       = "infinite"
  key_name       = "aa"
  name           = "name_1"
  name_alias     = "name_alias_1"
  owner_key      = "owner_key_1"
  owner_tag      = "owner_tag_1"
  pre_shared_key = "123456789a223456789a323456789abc"
  start_time     = "2025-11-28T03:12:09.452-08:00"
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
