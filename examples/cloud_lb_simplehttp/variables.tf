variable "name" {
  description = "Name to assign to the cloud boad balancer"
  type        = string
}

variable "epg_dn" {
  description = "DN for the target EPG where the loadbalancer should send traffic"
  type        = string
}

variable "contract_subject_dn" {
  description = "DN for the contract subject where the loadbalancer should be attached"
  type        = string
}

variable "subnet_a_dn" {
  description = "Dn for the ACI subnet A"
  type        = string
}

variable "subnet_b_dn" {
  description = "Dn for the ACI subnet B"
  type        = string
}

variable "listenerPort" {
  description = "External TCP port number for the http listener (ex. 80)"
  type        = string
}

variable "hostPort" {
  description = "Internal TCP port number for the http server (ex. 80)"
  type        = string
}
