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

resource "aci_fabric_node_member" "example" {
  name        = "test"
  serial      = "1"
  annotation  = "example"
  description = "from terraform"
  ext_pool_id = "0"
  fabric_id   = "1"
  name_alias  = "example"
  node_id     = "201"
  node_type   = "unspecified"
  pod_id      = "1"
  role        = "unspecified"
}