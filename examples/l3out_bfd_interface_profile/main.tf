
provider "aci" {
  username = ""
  password = ""
  url      = ""
  insecure = true
}

resource "aci_l3out_bfd_interface_profile" "example" {
  logical_interface_profile_dn = aci_logical_interface_profile.example.id
  annotation                   = "example"
  description                  = "from terraform"
  key                          = "example"
  key_id                       = "25"
  interface_profile_type       = "sha1"
}
