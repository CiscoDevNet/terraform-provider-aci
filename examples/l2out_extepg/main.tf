
provider "aci" {
  username = ""
  password = ""
  url      = ""
  insecure = true
}


resource "aci_l2out_extepg" "example" {

  l2_outside_dn  = aci_l2_outside.example.id
  name  = "demo_ext_epg"
  annotation  = "example"
  exception_tag  = "example"
  flood_on_encap = "disabled"
  match_t = "All"
  name_alias  = "example"
  pref_gr_memb = "exclude"
  prio = "level1"
  target_dscp = "AF11"

}