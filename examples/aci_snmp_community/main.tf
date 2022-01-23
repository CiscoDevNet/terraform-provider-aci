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

resource "aci_snmp_community" "public" {
  parent_dn   = "uni/fabric/snmppol-default"
  name        = "public"
  description = "from terraform"
  annotation  = "aci_snmp_community"
  name_alias  = "example"
}

resource "aci_snmp_community" "public_vrf" {
  parent_dn   = aci_vrf_snmp_context.example.id
  name        = "public"
  description = "from terraform"
  annotation  = "aci_snmp_community"
  name_alias  = "example"
}

resource "aci_vrf_snmp_context" "example" {
  vrf_dn     = aci_vrf.example.id
  name       = "example"
}

resource "aci_tenant" "example" {
  name = "tenant_snmp"
}

resource "aci_vrf" "example" {
  tenant_dn = aci_tenant.example.id
  name      = "vrf_example"
}
