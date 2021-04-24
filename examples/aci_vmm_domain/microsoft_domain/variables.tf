variable "microsoft_domain" {
  default = "uni/vmmp-Microsoft"
}
// Microsoft vmm domain resources
variable "vmm_domain" {
  default = "microsoft_domain"
}

variable "aci_vmm_controller" {
  default = "microsoft_vmm_controller"
}


