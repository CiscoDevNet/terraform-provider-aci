terraform {
  required_providers {
    aci = {
      source = "CiscoDevNet/aci"
    }
  }
}

provider "aci" {
  username = "" # <APIC username>
  password = "" # <APIC pwd>
  url      = "" # <APIC URL>
  insecure = true
}

resource "aci_tenant" "terraform_tenant" {
  name        = "tf_tenant"
  description = "This tenant is created by terraform"
}

resource "aci_ip_sla_monitoring_policy" "example" {
  tenant_dn             = aci_tenant.terraform_tenant.id
  name                  = "example"
  type_of_service       = "0"
  traffic_class_value   = "0"
  sla_detect_multiplier = "3"
  sla_frequency         = "60"
  sla_port              = "10"
  sla_type              = "tcp"
  threshold             = "900"
  timeout               = "900"
}
