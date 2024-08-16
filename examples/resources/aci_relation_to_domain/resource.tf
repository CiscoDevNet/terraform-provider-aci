
resource "aci_relation_to_domain" "example_application_epg" {
  parent_dn = aci_application_epg.example.id
  target_dn = "uni/vmmp-VMware/dom-domain_2"
}
