
resource "aci_relation_to_static_leaf" "full_example_application_epg" {
  parent_dn            = aci_application_epg.example.id
  annotation           = "annotation"
  description          = "description_1"
  encapsulation        = "vlan-101"
  deployment_immediacy = "immediate"
  mode                 = "native"
  target_dn            = "topology/pod-1/node-101"
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
