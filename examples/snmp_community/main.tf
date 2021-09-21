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

resource "aci_snmp_community" "foosnmp_community" {
	name 		= "test"
	description = "From Terraform"
	vrf_dn = aci_vrf.test.id
	annotation = "Test_Annotation"
	name_alias = "Test_name_alias"
}
