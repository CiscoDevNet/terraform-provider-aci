
data "aci_epg_useg_criterion_vm_attribute" "example_epg_useg_criterion" {
  parent_dn = aci_epg_useg_criterion.example.id
  name      = "vm_attribute"
}

data "aci_epg_useg_criterion_vm_attribute" "example_epg_useg_sub_criterion" {
  parent_dn = aci_epg_useg_sub_criterion.example.id
  name      = "vm_attribute"
}
