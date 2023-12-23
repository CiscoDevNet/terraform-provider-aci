terraform {
  required_providers {
    aci = {
      source = "ciscodevnet/aci"
    }
  }
}

provider "aci" {
  username = ""
  password = ""
  url      = ""
  insecure = true
}

resource "aci_vrf_snmp_context_community" "example" {
  vrf_snmp_context_dn = aci_vrf_snmp_context.example.id
  name                = "test"
  description         = "From Terraform"
  annotation          = "Test_Annotation"
  name_alias          = "Test_name_alias"
}

resource "aci_vrf_snmp_context" "example" {
  vrf_dn = aci_vrf.example.id
  name   = "example"
}

resource "aci_tenant" "example" {
  name = "tenant_snmp"
}

resource "aci_vrf" "example" {
  tenant_dn = aci_tenant.example.id
  name      = "vrf_example"
}
