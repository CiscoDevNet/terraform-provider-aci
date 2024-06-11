
resource "aci_epg_useg_criterion_identity_group_attribute" "full_example_epg_useg_criterion" {
  parent_dn   = aci_epg_useg_criterion.example.id
  annotation  = "annotation"
  description = "description"
  name        = "name"
  name_alias  = "name_alias"
  owner_key   = "owner_key"
  owner_tag   = "owner_tag"
  selector    = "adepg/authsvr-common-sg1-ISE_1/grpcont/dom-cisco.com/grp-Eng"
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
