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

resource "aci_tenant" "tenant_for_match_rules" {
  name        = "tenant_for_match_rules"
  description = "This tenant is created by terraform ACI provider"
}

resource "aci_match_rule" "match_rule" {
  tenant_dn = aci_tenant.tenant_for_match_rules.id
  name      = "match_rule"
}

resource "aci_match_regex_community_terms" "regex_community_terms" {
  match_rule_dn  = aci_match_rule.match_rule.id
  name           = "regex_community_terms"
  community_type = "regular"
  description    = "This is regex community terms"
  regex          = ".*"
}

data "aci_match_regex_community_terms" "example" {
  match_rule_dn  = aci_match_regex_community_terms.regex_community_terms.match_rule_dn
  name           = aci_match_regex_community_terms.regex_community_terms.name
  community_type = aci_match_regex_community_terms.regex_community_terms.community_type
}

output "name" {
  value = data.aci_match_regex_community_terms.example
}