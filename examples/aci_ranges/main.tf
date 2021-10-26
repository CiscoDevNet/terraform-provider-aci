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

resource "aci_vlan_pool" "example" {
  name        = "example"
  description = "From Terraform"
  alloc_mode  = "static"
  annotation  = "example"
  name_alias  = "example"
}
resource "aci_ranges" "example" {
  vlan_pool_dn = aci_vlan_pool.example.id
  from         = "vlan-1"
  description  = "From Terraform"
  to           = "vlan-2"
  alloc_mode   = "inherit"
  annotation   = "example"
  name_alias   = "name_alias"
  role         = "external"
}
