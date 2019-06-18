variable "bd_subnet" {
  type    = "string"
  default = "1.1.1.1/24"
}

variable "provider_profile_dn" {
  default = "uni/vmmp-VMware"
}

variable "vmm_domain" {
  default = "ESX0-leaf102"
}