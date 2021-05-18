provider "aci" {
  username = ""
  password = ""
  url      = ""
  insecure = true
}

resource "aci_l3out_ospf_interface_profile" "example" {
  logical_interface_profile_dn = "${aci_logical_interface_profile.example.id}"
  description                  = "from terraform"
  annotation                   = "example"
  auth_key                     = "example"
  auth_key_id                  = "255"
  auth_type                    = "simple"
  name_alias                   = "example"
  relation_ospf_rs_if_pol      = "${aci_ospf_interface_policy.example.id}"
}