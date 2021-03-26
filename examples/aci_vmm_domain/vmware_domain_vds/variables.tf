variable "vds" {
  default = "uni/vmmp-VMware"
}

variable "vmm_domain" {
  default = "ESX0-leaf102-vds"
}

variable "aci_vmm_controller" {
  default = "vmware_vds_controlller"
}

variable "aci_vmm_credential" {
  default = "vmware_vds_credential"
}
