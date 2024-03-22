
resource "aci_communication_policy" "full_example" {
  annotation                         = "annotation"
  description                        = "description"
  name                               = "test_name"
  name_alias                         = "name_alias"
  owner_key                          = "owner_key"
  owner_tag                          = "owner_tag"
  strict_security_on_apic_oob_subnet = "no"
  http_ssl_configuration = [
    {
      access_control_allow_credential = "disabled"
      access_control_allow_origins    = "access_control_allow_origins_1"
      admin_st                        = "disabled"
      annotation                      = "annotation_1"
      cli_only_mode                   = "disabled"
      client_cert_auth_state          = "disabled"
      description                     = "description_1"
      dh_param                        = "1024"
      global_throttle_rate            = "global_throttle_rate_1"
      global_throttle_st              = "disabled"
      global_throttle_unit            = "global_throttle_unit_1"
      max_request_status_count        = "max_request_status_count_1"
      name                            = "name_1"
      name_alias                      = "name_alias_1"
      node_exporter                   = "disabled"
      port                            = "port_1"
      referer                         = "referer_1"
      server_header                   = "disabled"
      ssl_protocols                   = "TLSv1"
      throttle_rate                   = "throttle_rate_1"
      throttle_st                     = "disabled"
      visore_access                   = "disabled"
    }
  ]
  annotations = [
    {
      key   = "annotations_1"
      value = "value_1"
    }
  ]
}
