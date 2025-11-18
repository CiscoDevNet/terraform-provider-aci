
data "aci_vmm_credential" "example_vmm_domain" {
  parent_dn = aci_vmm_domain.example.id
  name      = "test_name"
}
