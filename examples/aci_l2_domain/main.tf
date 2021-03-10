
provider "aci" {
  username = ""
  password = ""
  url      = ""
  insecure = true
}

resource "aci_l2_domain" "example" {
  name  = "demo_l2_domain"
  annotation  = "l2_domain for demo"
  name_alias  = "l2_domain"
}