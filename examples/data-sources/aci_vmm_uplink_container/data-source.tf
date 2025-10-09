
data "aci_vmm_uplink_container" "example_vmm_domain" {
  parent_dn = aci_vmm_domain.example.id
}
