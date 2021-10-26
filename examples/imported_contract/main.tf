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

resource "aci_tenant" "terraform_ten" {
  name = "import_con"
}

resource "aci_imported_contract" "example" {

  tenant_dn         = aci_tenant.terraform_ten.id
  name              = "example"
  annotation        = "example"
  name_alias        = "example"
  relation_vz_rs_if = aci_tenant.terraform_ten.id
}
