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

resource "aci_tenant" "foo_tenant" {
  description = "sample aci_tenant from terraform"
  name        = "terraform_test_tenant"
}

resource "aci_bfd_multihop_interface_policy" "foo_bfd_multihop_interface_policy" {
  tenant_dn             = aci_tenant.foo_tenant.id
  name                  = "bfd_mh_interface_policy"
  admin_state           = "disabled"
  detection_multiplier  = "5"
  min_transmit_interval = "251"
  min_receive_interval  = "253"
}

data "aci_bfd_multihop_interface_policy" "data_bfd_multihop_interface_policy" {
  tenant_dn = aci_bfd_multihop_interface_policy.foo_bfd_multihop_interface_policy.tenant_dn
  name      = aci_bfd_multihop_interface_policy.foo_bfd_multihop_interface_policy.name
}

output "bfd_multihop_interface_policy" {
  value = data.aci_bfd_multihop_interface_policy.data_bfd_multihop_interface_policy
}