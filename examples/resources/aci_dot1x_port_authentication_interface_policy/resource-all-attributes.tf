
resource "aci_dot1x_port_authentication_interface_policy" "full_example" {
  admin_state = "disabled"
  annotation  = "annotation"
  description = "description_1"
  host_mode   = "multi-auth"
  name        = "test_name"
  name_alias  = "name_alias_1"
  owner_key   = "owner_key_1"
  owner_tag   = "owner_tag_1"
  configuration = {
    annotation                        = "annotation_1"
    authentication_mode               = "bypass"
    maximum_reauthentication_requests = "5"
    maximum_requests                  = "4"
    reauthenticate                    = "no"
    reauthentication_period           = "3000"
    server_timeout                    = "20"
    supplicant_timeout                = "25"
    transmit_period                   = "35"
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
