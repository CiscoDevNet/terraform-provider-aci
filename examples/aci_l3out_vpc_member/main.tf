
terraform {
  required_providers {
    aci = {
      source = "ciscodevnet/aci"
    }
  }
}

provider "aci" {
  username = ""
  password = ""
  url      = ""
  insecure = true
}

resource "aci_l3out_vpc_member" "example" {
  leaf_port_dn = aci_l3out_path_attachment.example.id
  side         = "A"
  addr         = "10.0.0.1/16"
  annotation   = "example"
  ipv6_dad     = "enabled"
  ll_addr      = "::"
  description  = "from terraform"
  name_alias   = "example"
}