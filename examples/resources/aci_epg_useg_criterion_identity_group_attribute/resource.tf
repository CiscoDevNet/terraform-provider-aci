
resource "aci_epg_useg_criterion_identity_group_attribute" "example_epg_useg_criterion" {
  parent_dn = aci_epg_useg_criterion.example.id
  selector  = "adepg/authsvr-common-sg1-ISE_1/grpcont/dom-cisco.com/grp-Eng"
}
