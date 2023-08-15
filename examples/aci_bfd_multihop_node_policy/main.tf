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

resource "aci_tenant" "tenant_test_tf" {
  description = "sample aci_tenant from terraform"
  name        = "terraform_test_tenant_bfd_mh_node_policy"
}

resource "aci_bfd_multihop_node_policy" "bfd_mh_node_policy" {
  tenant_dn            = aci_tenant.tenant_test_tf.id
  description          = "sample aci_bfd_multihop_node_policy from terraform"
  name                 = "terraform_test_bfd_mh_node_policy"
  admin_state          = "enabled"
  detection_multiplier = "2"   # Values between 1 and 50. Default value set to 3.
  min_rx_interval      = "250" # Values between 250 and 999. Default value set to 250.
  min_tx_interval      = "500" # Values between 250 and 999. Default value set to 250.
}