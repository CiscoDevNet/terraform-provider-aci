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

resource "aci_tenant" "footenant" {
  name        = "terraform_tenant"
  description = "tenant created while acceptance testing"
}

resource "aci_aaa_domain" "foosecurity_domain" {
  name        = "aaa_domain_1"
  description = "from terraform"
  annotation  = "aaa_domain_tag"
  name_alias  = "example"
}

data "aci_aaa_domain" "example6" {
  name = aci_aaa_domain.foosecurity_domain.name
}

output "name6" {
  value = data.aci_aaa_domain.example6
}

resource "aci_aaa_domain_relationship" "foodomain_relationship_parent_object" {
  name      = resource.aci_aaa_domain.foosecurity_domain.name
  parent_dn = aci_tenant.footenant.id
}

output "t_name6" {
  value = resource.aci_aaa_domain_relationship.foodomain_relationship_parent_object.name
}
