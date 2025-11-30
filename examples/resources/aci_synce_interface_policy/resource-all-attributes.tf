
resource "aci_synce_interface_policy" "full_example" {
  admin_state                    = "enabled"
  annotation                     = "annotation"
  description                    = "description_1"
  name                           = "test_name"
  name_alias                     = "name_alias_1"
  owner_key                      = "owner_key_1"
  owner_tag                      = "owner_tag_1"
  quality_level_options          = "op2g2"
  quality_receive_exact_value    = "fsync-ql-o2-g2-eec2"
  quality_receive_highest_value  = "fsync-ql-common-none"
  quality_receive_lowest_value   = "fsync-ql-common-none"
  quality_transmit_exact_value   = "fsync-ql-o2-g2-eec2"
  quality_transmit_highest_value = "fsync-ql-common-none"
  quality_transmit_lowest_value  = "fsync-ql-common-none"
  selection_input                = "yes"
  source_priority                = "10"
  synchronization_status_message = "yes"
  wait_to_restore_time           = "8"
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
