
resource "aci_certificate_authority" "full_example" {
  annotation        = "annotation"
  certificate_chain = <<EOT
-----BEGIN CERTIFICATE-----
MIICODCCAaGgAwIBAgIJAIt8XMntue0VMA0GCSqGSIb3DQEBCwUAMDQxDjAMBgNV
BAMMBUFkbWluMRUwEwYDVQQKDAxZb3VyIENvbXBhbnkxCzAJBgNVBAYTAlVTMCAX
DTE4MDEwOTAwNTk0NFoYDzIxMTcxMjE2MDA1OTQ0WjA0MQ4wDAYDVQQDDAVBZG1p
bjEVMBMGA1UECgwMWW91ciBDb21wYW55MQswCQYDVQQGEwJVUzCBnzANBgkqhkiG
9w0BAQEFAAOBjQAwgYkCgYEAohG/7axtt7CbSaMP7r+2mhTKbNgh0Ww36C7Ta14i
v+VmLyKkQHnXinKGhp6uy3Nug+15a+eIu7CrgpBVMQeCiWfsnwRocKcQJWIYDrWl
XHxGQn31yYKR6mylE7Dcj3rMFybnyhezr5D8GcP85YRPmwG9H2hO/0Y1FUnWu9Iw
AQkCAwEAAaNQME4wHQYDVR0OBBYEFD0jLXfpkrU/ChzRvfruRs/fy1VXMB8GA1Ud
IwQYMBaAFD0jLXfpkrU/ChzRvfruRs/fy1VXMAwGA1UdEwQFMAMBAf8wDQYJKoZI
hvcNAQELBQADgYEAOmvre+5tgZ0+F3DgsfxNQqLTrGiBgGCIymPkP/cBXXkNuJyl
3ac7tArHQc7WEA4U2R2rZbEq8FC3UJJm4nUVtCPvEh3G9OhN2xwYev79yt6pIn/l
KU0Td2OpVyo0eLqjoX5u2G90IBWzhyjFbo+CcKMrSVKj1YOdG0E3OuiJf00=
-----END CERTIFICATE-----
EOT
  certificate_usage = test_value
  description       = "description_1"
  name              = "test_name"
  name_alias        = "name_alias_1"
  owner_key         = "owner_key_1"
  owner_tag         = "owner_tag_1"
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

// This example is only applicable to Cisco Cloud Network Controller
resource "aci_certificate_authority" "full_example_tenant" {
  parent_dn         = aci_tenant.example.id
  annotation        = "annotation"
  certificate_chain = <<EOT
-----BEGIN CERTIFICATE-----
MIICODCCAaGgAwIBAgIJAIt8XMntue0VMA0GCSqGSIb3DQEBCwUAMDQxDjAMBgNV
BAMMBUFkbWluMRUwEwYDVQQKDAxZb3VyIENvbXBhbnkxCzAJBgNVBAYTAlVTMCAX
DTE4MDEwOTAwNTk0NFoYDzIxMTcxMjE2MDA1OTQ0WjA0MQ4wDAYDVQQDDAVBZG1p
bjEVMBMGA1UECgwMWW91ciBDb21wYW55MQswCQYDVQQGEwJVUzCBnzANBgkqhkiG
9w0BAQEFAAOBjQAwgYkCgYEAohG/7axtt7CbSaMP7r+2mhTKbNgh0Ww36C7Ta14i
v+VmLyKkQHnXinKGhp6uy3Nug+15a+eIu7CrgpBVMQeCiWfsnwRocKcQJWIYDrWl
XHxGQn31yYKR6mylE7Dcj3rMFybnyhezr5D8GcP85YRPmwG9H2hO/0Y1FUnWu9Iw
AQkCAwEAAaNQME4wHQYDVR0OBBYEFD0jLXfpkrU/ChzRvfruRs/fy1VXMB8GA1Ud
IwQYMBaAFD0jLXfpkrU/ChzRvfruRs/fy1VXMAwGA1UdEwQFMAMBAf8wDQYJKoZI
hvcNAQELBQADgYEAOmvre+5tgZ0+F3DgsfxNQqLTrGiBgGCIymPkP/cBXXkNuJyl
3ac7tArHQc7WEA4U2R2rZbEq8FC3UJJm4nUVtCPvEh3G9OhN2xwYev79yt6pIn/l
KU0Td2OpVyo0eLqjoX5u2G90IBWzhyjFbo+CcKMrSVKj1YOdG0E3OuiJf00=
-----END CERTIFICATE-----
EOT
  certificate_usage = test_value
  description       = "description_1"
  name              = "test_name"
  name_alias        = "name_alias_1"
  owner_key         = "owner_key_1"
  owner_tag         = "owner_tag_1"
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
