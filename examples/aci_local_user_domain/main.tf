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

resource "aci_aaa_domain" "foosecurity_domain" {
  name        = "aaa_domain_local_user"
  description = "from terraform"
  annotation  = "aaa_domain_tag"
  name_alias  = "example"
}

resource "aci_local_user" "foolocal_user" {
  name        = "user_with_domain"
  pwd         = "Password123!"
  description = "This user is created by terraform"
}

resource "aci_user_security_domain" "foouser_security_domain" {
  local_user_dn  = aci_local_user.foolocal_user.id
  name  = aci_aaa_domain.foosecurity_domain.name
  annotation = "orchestrator:terraform"
  name_alias = "example_name_alias"
  description = "from Terraform"
}

resource "aci_user_security_domain_role" "foouser_domain_role" {
    user_domain_dn  = aci_user_security_domain.foouser_security_domain.id
    annotation     = "orchestrator:terraform"
    name            = "example"
    priv_type       = "readPriv"
    name_alias      = "user_role_alias"
    description     = "From Terraform"
}
