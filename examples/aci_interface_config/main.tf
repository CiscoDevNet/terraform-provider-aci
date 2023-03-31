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

resource "aci_leaf_access_port_policy_group" "leaf_access_port" {
  name = "leaf_access_port"
}

# Create Access Port
resource "aci_interface_config" "access_port_config_1001" {
  node         = 1001
  interface    = "1/1/1"
  port_type    = "access"
  policy_group = aci_leaf_access_port_policy_group.leaf_access_port.id
  description  = "Access Port created using TF"
}

# Create Access Port - Breakout interface
resource "aci_interface_config" "access_port_config_1002_brkout" {
  node        = 1002
  interface   = "1/1/1"
  port_type   = "access"
  breakout    = "100g-4x"
  description = "Breakout an Access Port using TF"
}

# Create Fabric Port Configuration
resource "aci_interface_config" "fabric_port_config" {
  node        = 1003
  interface   = "2/2/2"
  port_type   = "fabric"
  description = "Fabric Port created using TF"
}
