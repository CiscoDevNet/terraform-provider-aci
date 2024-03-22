
resource "aci_management_access_policy" "full_example" {
  annotation                         = "annotation"
  description                        = "description_1"
  name                               = "test_name"
  name_alias                         = "name_alias_1"
  owner_key                          = "owner_key_1"
  owner_tag                          = "owner_tag_1"
  strict_security_on_apic_oob_subnet = "no"
  http_service = {
    allow_credentials        = "disabled"
    allow_origins            = "allow_origins_1"
    admin_state              = "disabled"
    annotation               = "annotation_1"
    cli_only_mode            = "disabled"
    description              = "description_1"
    global_throttle_rate     = "10000"
    global_throttle_state    = "disabled"
    global_throttle_unit     = "r/s"
    max_request_status_count = "0"
    name                     = "name_1"
    name_alias               = "name_alias_1"
    node_exporter            = "disabled"
    port                     = "80"
    redirect_state           = "disabled"
    server_header            = "disabled"
    login_throttle_rate      = "2"
    login_throttle_state     = "disabled"
    visore_access            = "disabled"
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
  http_ssl_configuration = {
    allow_credentials                       = "disabled"
    allow_origins                           = "allow_origins_1"
    admin_state                             = "enabled"
    annotation                              = "annotation_1"
    cli_only_mode                           = "disabled"
    client_certificate_authentication_state = "disabled"
    description                             = "description_1"
    dh_parameter                            = "1024"
    global_throttle_rate                    = "10000"
    global_throttle_state                   = "disabled"
    global_throttle_unit                    = "r/s"
    max_request_status_count                = "0"
    name                                    = "name_1"
    name_alias                              = "name_alias_1"
    node_exporter                           = "disabled"
    port                                    = "443"
    referer                                 = "referer_1"
    server_header                           = "disabled"
    ssl_protocols                           = ["TLSv1"]
    login_throttle_rate                     = "2"
    login_throttle_state                    = "disabled"
    visore_access                           = "disabled"
    certificate_authority = {
      annotation = "annotation_1"
      target_dn  = "uni/userext/pkiext/tp-test_name"
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
    key_ring = {
      annotation    = "annotation_1"
      key_ring_name = aci_key_ring.example.name
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
  ssh_access_via_web = {
    admin_state = "disabled"
    annotation  = "annotation_1"
    description = "description_1"
    name        = "name_1"
    name_alias  = "name_alias_1"
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
  ssh_service = {
    admin_state                   = "disabled"
    annotation                    = "annotation_1"
    description                   = "description_1"
    host_key_algorithms           = ["rsa-sha2-256"]
    kex_algorithms                = ["curve25519-sha256"]
    name                          = "name_1"
    name_alias                    = "name_alias_1"
    password_authentication_state = "disabled"
    port                          = "22"
    ssh_ciphers                   = ["aes192-ctr"]
    ssh_macs                      = ["hmac-sha1"]
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
  telnet_service = {
    admin_state = "disabled"
    annotation  = "annotation_1"
    description = "description_1"
    name        = "name_1"
    name_alias  = "name_alias_1"
    port        = "23"
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
