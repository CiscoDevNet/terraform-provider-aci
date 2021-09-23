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

resource "aci_vrf_snmp_context_community" "example" {
	vrf_snmp_context_dn = aci_vrf_snmp_context.test.id
	name = "test"
	description = "From Terraform"
	annotation = "Test_Annotation"
	name_alias = "Test_name_alias"
}
