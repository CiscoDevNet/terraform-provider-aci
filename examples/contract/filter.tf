// Add below filter and filter entries to contract web_contract in contract.tf file

resource "aci_filter" "allow_https" {
  tenant_dn = aci_tenant.tenant_for_contract.id
  name      = "allow_https"
}
resource "aci_filter" "allow_icmp" {
  tenant_dn = aci_tenant.tenant_for_contract.id
  name      = "allow_icmp"
}

resource "aci_filter_entry" "https" {
  name        = "https"
  filter_dn   = aci_filter.allow_https.id
  ether_t     = "ip"
  prot        = "tcp"
  d_from_port = "https"
  d_to_port   = "https"
  stateful    = "yes"
}

resource "aci_filter_entry" "icmp" {
  name      = "icmp"
  filter_dn = aci_filter.allow_icmp.id
  ether_t   = "ip"
  prot      = "icmp"
  stateful  = "yes"
}
