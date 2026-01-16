
data "aci_vmm_domain" "example" {
  parent_dn = "uni/vmmp-VMware"
  name      = "test_name"
}
