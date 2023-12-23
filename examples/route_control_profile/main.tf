terraform {
  required_providers {
    aci = {
      source = "ciscodevnet/aci"
    }
  }
}

#configure provider with your cisco aci credentials.
provider "aci" {
  username = "" # <APIC username>
  password = "" # <APIC pwd>
  url      = "" # <cloud APIC URL>
  insecure = true
}

resource "aci_tenant" "tenant_for_route_control" {
  name        = "tenant_for_route_control"
  description = "This tenant is created by terraform ACI provider"
}

resource "aci_l3_outside" "example" {
  tenant_dn = aci_tenant.tenant_for_route_control.id
  name      = "example_l3out"
}

resource "aci_route_control_profile" "example" {
  parent_dn                  = aci_tenant.tenant_for_route_control.id
  name                       = "one"
  annotation                 = "example"
  description                = "from terraform"
  name_alias                 = "example"
  route_control_profile_type = "global"
}

resource "aci_route_control_profile" "example2" {
  parent_dn                  = aci_l3_outside.example.id
  name                       = "route_control_profile_1"
  annotation                 = "route_control_profile_tag"
  description                = "from terraform"
  name_alias                 = "example"
  route_control_profile_type = "global"
}

## DEPRECATED VERSION
resource "aci_bgp_route_control_profile" "example" {
  parent_dn                  = aci_tenant.tenant_for_route_control.id
  name                       = "bgp_route_control_profile_1"
  annotation                 = "bgp_route_control_profile_tag"
  description                = "from terraform"
  name_alias                 = "example"
  route_control_profile_type = "global"
}
