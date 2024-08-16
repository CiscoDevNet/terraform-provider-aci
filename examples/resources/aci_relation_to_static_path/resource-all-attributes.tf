
resource "aci_relation_to_static_path" "full_example_application_epg" {
  parent_dn             = aci_application_epg.example.id
  annotation            = "annotation"
  description           = "description_1"
  encapsulation         = "vlan-201"
  deployment_immediacy  = "immediate"
  mode                  = "native"
  primary_encapsulation = "vlan-203"
  target_dn             = "topology/pod-1/paths-101/pathep-[eth1/1]"
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
