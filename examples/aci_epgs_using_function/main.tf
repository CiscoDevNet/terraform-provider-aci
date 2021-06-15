
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





resource "aci_tenant" "terraform_tenant" {
    name        = "tf_tenant"
    description = "This tenant is created by terraform"
}

resource "aci_application_profile" "terraform_ap" {
    tenant_dn  = aci_tenant.terraform_tenant.id
    name       = "tf_ap"
}

resource "aci_application_epg" "terraform_epg" {
    application_profile_dn  = aci_application_profile.terraform_ap.id
    name                    = "tf_epg"
}

resource "aci_attachable_access_entity_profile" "fooattachable_access_entity_profile" {
	description = "From Terraform"
	name        = "demo_entity_prof"
	annotation  = "tag_entity"
	name_alias  = "Name_Alias"
}

resource "aci_access_generic" "example" {
    attachable_access_entity_profile_dn  = aci_attachable_access_entity_profile.fooattachable_access_entity_profile.id
    description = "from terraform"
    name  = "default"
    annotation  = "access_generic_tag"
    name_alias  = "access_generic"
}

resource "aci_epgs_using_function" "example" {
  access_generic_dn   = aci_access_generic.example.id
  tdn                 = aci_application_epg.terraform_epg.id
  annotation          = "annotation"
  encap               = "vlan-5"
  instr_imedcy        = "lazy"
  mode                = "regular"
  primary_encap       = "vlan-7"
}
