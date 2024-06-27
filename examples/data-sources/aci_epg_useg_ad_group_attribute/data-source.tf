
data "aci_epg_useg_ad_group_attribute" "example_epg_useg_block_statement" {
  parent_dn = aci_epg_useg_block_statement.example.id
  selector  = "adepg/authsvr-common-sg1-ISE_1/grpcont/dom-cisco.com/grp-Eng"
}
