terraform {
  required_providers {
    aci = {
      source = "CiscoDevNet/aci"
    }
  }
}

provider "aci" {
  username = ""
  password = ""
  url      = ""
  insecure = true
}

resource "aci_spanning_tree_interface_policy" "example" {
  name        = "example"
  annotation  = "orchestrator:terraform"
  description = "from terraform"
  ctrl        = ["bpdu-guard"]
}

data "aci_spanning_tree_interface_policy" "example" {
  name  = "example"
}
 
output "name" {
  value = data.aci_spanning_tree_interface_policy.example
}