resource "aci_tenant" "terraform_tenant" {
  name        = "tf_tenant"
  description = "This tenant is created by terraform"
}

resource "aci_application_profile" "terraform_ap" {
  tenant_dn = aci_tenant.terraform_tenant.id
  name      = "tf_ap"
}

# create ESG
resource "aci_endpoint_security_group" "terraform_esg" {
  application_profile_dn = aci_application_profile.terraform_ap.id
  name                   = "tf_esg"
}

# create selector
resource "aci_endpoint_security_group_tag_selector" "terraform_tag_selector" {
  endpoint_security_group_dn = aci_endpoint_security_group.terraform_esg.id
  match_key                  = "example"
  match_value                = "example1"
  value_operator             = "equals"
  description                = "Tag Selector"
}


