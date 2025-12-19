
resource "aci_mcp_interface_policy" "full_example" {
  admin_state                         = "enabled"
  annotation                          = "annotation"
  description                         = "description_1"
  grace_period_seconds                = "3"
  grace_period_milliseconds           = "138"
  maximum_number_of_vlan              = "129"
  strict_mode                         = "on"
  mcp_pdu_per_vlan                    = "off"
  name                                = "test_name"
  name_alias                          = "name_alias_1"
  owner_key                           = "owner_key_1"
  owner_tag                           = "owner_tag_1"
  initial_delay_time                  = "1750"
  transmission_frequency_seconds      = "10"
  transmission_frequency_milliseconds = "800"
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
