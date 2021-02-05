
provider "aci" {
  username = ""
  password = ""
  url      = ""
  insecure = true
}

resource "aci_tenant" "tenentcheck" {
  description = "hello"
  name        = "check_tenantnk"
  annotation  = "tag_tenant"
  name_alias  = "alias_tenant"
}

resource "aci_contract" "democontract" {
  tenant_dn   = aci_tenant.tenentcheck.id
  name        = "test_tf_contract"
  description = "ThiscontractiscreatedbyterraformACIprovider"
  scope       = "context"
  target_dscp = "VA"
  prio        = "unspecified"
  filter  {
    filter_name = "abcd"
    description = "first filter from contract resource"
    annotation  = "tag_filter"
    name_alias  = "abcd"
    filter_entry {
      description       = "hello"
      filter_entry_name = "check_entry1"
      d_from_port       = "80"
      ether_t           = "ipv4"
      prot              = "tcp"
    }
    filter_entry  {
      description       = "hello world"
      filter_entry_name = "check_entry2"
      d_from_port       = "443"
      ether_t           = "ipv4"
      prot              = "tcp"
    }
  }
}

