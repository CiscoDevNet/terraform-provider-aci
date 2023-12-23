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

resource "aci_isis_domain_policy" "example" {
  annotation          = "orchestrator:terraform"
  mtu                 = "1492"
  redistrib_metric    = "63"
  description         = "from terraform"
  name_alias          = "example_name_alias"
  lsp_fast_flood      = "disabled"
  lsp_gen_init_intvl  = "50"
  lsp_gen_max_intvl   = "8000"
  lsp_gen_sec_intvl   = "50"
  spf_comp_init_intvl = "50"
  spf_comp_max_intvl  = "8000"
  spf_comp_sec_intvl  = "50"
  isis_level_name     = "example"
}