
resource "aci_neighbor_discovery_interface_policy" "full_example_tenant" {
  parent_dn                      = aci_tenant.example.id
  annotation                     = "annotation"
  controller_state               = ["managed-cfg"]
  description                    = "description_1"
  hop_limit                      = "40"
  mtu                            = "8700"
  name                           = "test_name"
  name_alias                     = "name_alias_1"
  neighbor_solicitation_interval = "1500"
  neighbor_solicitation_retries  = "6"
  nud_retry_base                 = "2"
  nud_retry_interval             = "1300"
  nud_retry_max_attempts         = "5"
  owner_key                      = "owner_key_1"
  owner_tag                      = "owner_tag_1"
  router_advertisement_interval  = "500"
  router_advertisement_lifetime  = "1500"
  reachable_time                 = "2"
  retransmit_timer               = "2"
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
