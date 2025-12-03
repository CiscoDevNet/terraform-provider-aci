
resource "aci_poe_interface_policy" "full_example" {
  admin_state        = "disabled"
  annotation         = "annotation"
  description        = "description_1"
  maximum_power      = "15400"
  host_mode          = "never"
  name               = "test_name"
  name_alias         = "name_alias_1"
  owner_key          = "owner_key_1"
  owner_tag          = "owner_tag_1"
  poe_vlan           = "vlan-200"
  policing_action    = "log"
  port_priority_high = "yes"
  relation_to_application_epg = {
    annotation = "annotation_1"
    target_dn  = aci_application_epg.test_application_epg_0.id
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
