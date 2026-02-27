
resource "aci_storm_control_interface_policy" "full_example" {
  annotation                              = "annotation"
  broadcast_burst_rate_packets_per_second = "11"
  broadcast_burst_rate_percentage         = "100.000000"
  broadcast_rate_percentage               = "100.000000"
  broadcast_rate_packets_per_second       = "10"
  burst_rate_packets_per_second           = "unspecified"
  burst_rate_percentage                   = "100.000000"
  description                             = "description_1"
  unicast_multicast_broadcast             = "Valid"
  multicast_burst_rate_packets_per_second = "13"
  multicast_burst_rate_percentage         = "100.000000"
  multicast_rate_percentage               = "100.000000"
  multicast_rate_packets_per_second       = "12"
  name                                    = "test_name"
  name_alias                              = "name_alias_1"
  owner_key                               = "owner_key_1"
  owner_tag                               = "owner_tag_1"
  rate_percentage                         = "100.000000"
  rate_packets_per_second                 = "unspecified"
  action                                  = "drop"
  soak_count                              = "3"
  unicast_burst_rate_packets_per_second   = "15"
  unicast_burst_rate_percentage           = "100.000000"
  unicast_rate_percentage                 = "100.000000"
  unicast_rate_packets_per_second         = "14"
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
