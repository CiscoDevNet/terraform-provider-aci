resource "aci_tenant" "tenant_for_entry" {
  name        = "tenant_for_entry"
  description = "This tenant is created by terraform ACI provider"
}

resource "aci_filter" "filter_for_entry" {
  tenant_dn   = "${aci_tenant.tenant_for_entry.id}"
  name        = "filter_for_entry"
  description = "This filter is created by terraform ACI provider."
}

resource "aci_filter_entry" "demoentry" {
  filter_dn     = "${aci_filter.filter_for_entry.id}"
  name          = "test_tf_entry"
  description   = "This entry is created by terraform ACI provider"
  apply_to_frag = "no"
  arp_opc       = "unspecified"
  d_from_port   = "unspecified"
  d_to_port     = "unspecified"
  ether_t       = "ip"
  icmpv4_t      = "unspecified"
  icmpv6_t      = "unspecified"
  match_dscp    = "AF11"
  prot          = "tcp"
  s_from_port   = "unspecified"
  s_to_port     = "unspecified"
  stateful      = "no"
  tcp_rules     = "ack"
}
