variable "vsphere_server" {
  default = ""
}

variable "vsphere_user" {
  default = ""
}

variable "vsphere_password" {
  default = ""
}

variable "vsphere_datacenter" {
  default = "ESX0"
}

variable "aci_vm1_address" {
  default = "1.1.1.10"
}

variable "aci_vm2_address" {
  default = "1.1.1.11"
}

variable "aci_vm1_name" {
  default = "aci-tf-test1"
}

variable "aci_vm2_name" {
  default = "aci-tf-test2"
}

variable "gateway" {
  default = "1.1.1.1"
}

variable "domain_name" {
  default = ""
}

variable "vsphere_template" {
  default = "ubuntu-1404-template"
}

variable "folder" {
  default = "CLEUR-workshop"
}

variable "dns_list" {
  default = ["172.23.136.143", "172.23.136.144"]
}

variable "dns_search" {
  default = ["cisco.com"]
}

variable "vsphere_host_name" {
  default = "10.23.239.30"
}

variable "vsphere_datastore" {
  default = "datastore1"
}

variable "aci_tenant_name" {
  default = ""
}

variable "aci_epg_name" {
  default = ""
}

variable "aci_application_profile_name" {
  default = ""
}
