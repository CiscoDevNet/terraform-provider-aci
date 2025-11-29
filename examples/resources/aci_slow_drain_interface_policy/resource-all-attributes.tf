
resource "aci_slow_drain_interface_policy" "full_example" {
  annotation                   = "annotation"
  congestion_clear_action      = "err-disable"
  congestion_detect_multiplier = "15"
  description                  = "description_1"
  flush_admin_state            = "disabled"
  flush_timeout                = "300"
  name                         = "test_name"
  name_alias                   = "name_alias_1"
  owner_key                    = "owner_key_1"
  owner_tag                    = "owner_tag_1"
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
