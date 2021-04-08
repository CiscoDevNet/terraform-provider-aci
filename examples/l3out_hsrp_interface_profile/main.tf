provider "aci" {
  username = ""
  password = ""
  url      = ""
  insecure = true
}

resource "aci_l3out_hsrp_interface_profile" "example" {

  logical_interface_profile_dn  = "${aci_logical_interface_profile.example.id}"
  annotation  = "example"
  name_alias  = "example"
  version = "v1"

}
