terraform {
  required_providers {
    aci = {
      source = "ciscodevnet/aci"
    }
  }
}

provider "aci" {
  username = "admin"
  password = "ins3965!"
  url      = "https://10.23.248.120"
  insecure = true
}

resource "aci_tenant" "tenant_for_match_rules" {
  name        = "tenant_for_match_rules"
  description = "This tenant is created by terraform ACI provider"
}

resource "aci_match_rule" "match_rule" {
  tenant_dn  = aci_tenant.tenant_for_match_rules.id
  name  = "match_rule"
}

resource "aci_match_community_terms" "community_terms" {
  match_rule_dn = aci_match_rule.match_rule.id
  name  = "community_terms"
  description = "This is community terms"

  // match_community_factors {
  //   community = "regular:as2-nn2:4:15"
  //   scope = "non-transitive"
  // }
  match_community_factors {
    community = "regular:as2-nn2:4:16"
    scope = "transitive"
  }
}

// data "aci_match_community_terms" "example" {
//   match_rule_dn = aci_match_community_terms.community_terms.match_rule_dn
//   name  = aci_match_community_terms.community_terms.name
// }

// output "name" {
//   value = data.aci_match_community_terms.example
// }