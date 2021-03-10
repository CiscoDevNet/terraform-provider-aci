provider "aci" {
  username = ""
  password = ""
  url      = ""
  insecure = true
}

resource "aci_leaf_breakout_port_group" "example" {
  name        = "first"
  annotation  = "example"
  description = "adfadf"
  brkout_map  = "100g-4x"
  name_alias  = "aliasing"
}
