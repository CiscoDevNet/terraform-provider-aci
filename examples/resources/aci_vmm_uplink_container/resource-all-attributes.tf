
resource "aci_vmm_uplink_container" "full_example_vmm_domain" {
  parent_dn         = aci_vmm_domain.example.id
  annotation        = "annotation"
  name_alias        = "name_alias_1"
  number_of_uplinks = "2"
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
  uplink_policies = [
    {
      annotation  = "annotation_1"
      name_alias  = "name_alias_1"
      uplink_id   = "2"
      uplink_name = "uplink_name_2"
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
  ]
}
