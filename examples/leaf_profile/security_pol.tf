resource "aci_port_security_policy" "test_sec_pol" {
    description = "From Terraform"
		name        = "test_sec_pol"
		annotation  = "tag_port_pol"
		maximum     = "12"
		name_alias  = "alias_port_pol"
		timeout     = "60"
		violation   = "protect"
  
}
