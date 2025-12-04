
resource "aci_link_level_interface_policy" "full_example" {
  annotation          = "annotation"
  auto_negotiation    = "off"
  description         = "description_1"
  port_delay          = "50"
  emi_retrain         = "enable"
  fec_mode            = "auto-fec"
  link_debounce       = "2001"
  name                = "test_name"
  name_alias          = "name_alias_1"
  owner_key           = "owner_key_1"
  owner_tag           = "owner_tag_1"
  physical_media_type = "sfp-10g-tx"
  speed               = "100G"
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
