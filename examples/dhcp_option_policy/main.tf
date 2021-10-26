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

resource "aci_dhcp_option_policy" "example" {

  tenant_dn  = aci_tenant.example.id
  name       = "example_dhcp_option_policy"
  annotation = "example"
  name_alias = "example"

  dhcp_option {
    name           = "example_one"
    annotation     = "annotation_one"
    data           = "data_one"
    dhcp_option_id = "1"
    name_alias     = "one"
  }
  dhcp_option {
    name           = "example_two"
    annotation     = "annotation_two"
    data           = "data_two"
    dhcp_option_id = "2"
    name_alias     = "two"
  }

}