# provider "aci" {
#   username = ""
#   password = ""
#   url      = ""
#   insecure = true
# }

provider "aci" {
  username    = "nirav"
  private_key = "/Users/crest/go/src/github.com/ciscoecosystem/certdir/.key"
  cert_name   = ""
  url         = ""
  insecure    = true
}


resource "aci_tenant" "demotenant" {
  name                          = "tf_test_tenant"
  description                   = "This tenant is created by terraform"
  relation_fv_rs_tn_deny_rule   = ["${aci_filter.deny_rule_filter1.id}", "${aci_filter.deny_rule_filter2.id}"] # Relation to vzFilter class. Cardinality - N_TO_M.
  relation_fv_rs_tenant_mon_pol = "${aci_rest.rest_mon_epg_pol.id}"                                            # Relation to monEPGPol class. Cardinality - N_TO_ONE.
}

resource "aci_tenant" "test_tenant" {
  name        = "tf_test_rel_tenant"
  description = "This tenant is created by terraform"
}
