
resource "aci_relation_to_domain" "full_example_application_epg" {
  parent_dn                     = aci_application_epg.example.id
  annotation                    = "annotation"
  binding_type                  = "dynamicBinding"
  class_preference              = "encap"
  custom_epg_name               = "custom_epg_name_1"
  delimiter                     = "@"
  encapsulation                 = "vlan-100"
  encapsulation_mode            = "auto"
  epg_cos                       = "Cos0"
  epg_cos_pref                  = "disabled"
  deployment_immediacy          = "immediate"
  ipam_dhcp_override            = "10.0.0.2"
  ipam_enabled                  = "yes"
  ipam_gateway                  = "10.0.0.1"
  lag_policy_name               = "lag_policy_name_1"
  netflow_direction             = "both"
  enable_netflow                = "disabled"
  number_of_ports               = "1"
  port_allocation               = "elastic"
  primary_encapsulation         = "vlan-200"
  primary_encapsulation_inner   = "vlan-300"
  resolution_immediacy          = "immediate"
  secondary_encapsulation_inner = "vlan-400"
  switching_mode                = "AVE"
  target_dn                     = "uni/vmmp-VMware/dom-domain_2"
  untagged                      = "no"
  annotations = [
    {
      key   = "key_0"
      value = "value_1"
    }
  ]
  tags = [
    {
      key   = "key_0"
      value = "value_1"
    }
  ]
}
