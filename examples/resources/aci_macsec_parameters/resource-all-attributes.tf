
resource "aci_macsec_parameters" "full_example" {
  annotation             = "annotation"
  cipher_suite           = "gcm-aes-128"
  confidentiality_offset = "offset-50"
  description            = "description_1"
  key_server_priority    = "39"
  name                   = "test_name"
  name_alias             = "name_alias_1"
  owner_key              = "owner_key_1"
  owner_tag              = "owner_tag_1"
  window_size            = "200"
  security_policy        = "must-secure"
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
