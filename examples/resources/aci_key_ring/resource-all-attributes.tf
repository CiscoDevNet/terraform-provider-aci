
resource "aci_key_ring" "full_example" {
  admin_state           = "completed"
  annotation            = "annotation"
  certificate           = <<EOT
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
  description           = "description_1"
  elliptic_curve        = "none"
  key                   = <<EOT
-----BEGIN PRIVATE KEY-----
MIICdwIBADANBgkqhkiG9w0BAQEFAASCAmEwggJdAgEAAoGBAKIRv+2sbbewm0mj
D+6/tpoUymzYIdFsN+gu02teIr/lZi8ipEB514pyhoaerstzboPteWvniLuwq4KQ
VTEHgoln7J8EaHCnECViGA61pVx8RkJ99cmCkepspROw3I96zBcm58oXs6+Q/BnD
/OWET5sBvR9oTv9GNRVJ1rvSMAEJAgMBAAECgYByu3QO0qF9h7X3JEu0Ld4cKBnB
giQ2uJC/et7KxIJ/LOvw9GopBthyt27KwG1ntBkJpkTuAaQHkyNns7vLkNB0S0IR
+owVFEcKYq9VCHTaiQU8TDp24gN+yPTrpRuH8YhDVq5SfVdVuTMgHVQdj4ya4VlF
Gj+a7+ipxtGiLsVGrQJBAM7p0Fm0xmzi+tBOASUAcVrPLcteFIaTBFwfq16dm/ON
00Khla8Et5kMBttTbqbukl8mxFjBEEBlhQqb6EdQQ0sCQQDIhHx1a9diG7y/4DQA
4KvR3FCYwP8PBORlSamegzCo+P1OzxiEo0amX7yQMA5UyiP/kUsZrme2JBZgna8S
p4R7AkEAr7rMhSOPUnMD6V4WgsJ5g1Jp5kqkzBaYoVUUSms5RASz4+cwJVCwTX91
Y1jcpVIBZmaaY3a0wrx13ajEAa0dOQJBAIpjnb4wqpsEh7VpmJqOdSdGxb1XXfFQ
sA0T1OQYqQnFppWwqrxIL+d9pZdiA1ITnNqyvUFBNETqDSOrUHwwb2cCQGArE+vu
ffPUWQ0j+fiK+covFG8NL7H+26NSGB5+Xsn9uwOGLj7K/YT6CbBtr9hJiuWjM1Al
0V4ltlTuu2mTMaw=
-----END PRIVATE KEY-----
EOT
  key_type              = "RSA"
  modulus               = "mod2048"
  name                  = "test_name"
  name_alias            = "name_alias_1"
  owner_key             = "owner_key_1"
  owner_tag             = "owner_tag_1"
  regenerate            = "no"
  certificate_authority = "test_name"
  annotations = [
    {
      key   = "key_0"
      value = "value_1"
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
  ]
  tags = [
    {
      key   = "key_0"
      value = "value_1"
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
  ]
}

// This example is only applicable to Cisco Cloud Network Controller
resource "aci_key_ring" "full_example_tenant" {
  parent_dn             = aci_tenant.example.id
  admin_state           = "completed"
  annotation            = "annotation"
  certificate           = <<EOT
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
  description           = "description_1"
  elliptic_curve        = "none"
  key                   = <<EOT
-----BEGIN PRIVATE KEY-----
MIICdwIBADANBgkqhkiG9w0BAQEFAASCAmEwggJdAgEAAoGBAKIRv+2sbbewm0mj
D+6/tpoUymzYIdFsN+gu02teIr/lZi8ipEB514pyhoaerstzboPteWvniLuwq4KQ
VTEHgoln7J8EaHCnECViGA61pVx8RkJ99cmCkepspROw3I96zBcm58oXs6+Q/BnD
/OWET5sBvR9oTv9GNRVJ1rvSMAEJAgMBAAECgYByu3QO0qF9h7X3JEu0Ld4cKBnB
giQ2uJC/et7KxIJ/LOvw9GopBthyt27KwG1ntBkJpkTuAaQHkyNns7vLkNB0S0IR
+owVFEcKYq9VCHTaiQU8TDp24gN+yPTrpRuH8YhDVq5SfVdVuTMgHVQdj4ya4VlF
Gj+a7+ipxtGiLsVGrQJBAM7p0Fm0xmzi+tBOASUAcVrPLcteFIaTBFwfq16dm/ON
00Khla8Et5kMBttTbqbukl8mxFjBEEBlhQqb6EdQQ0sCQQDIhHx1a9diG7y/4DQA
4KvR3FCYwP8PBORlSamegzCo+P1OzxiEo0amX7yQMA5UyiP/kUsZrme2JBZgna8S
p4R7AkEAr7rMhSOPUnMD6V4WgsJ5g1Jp5kqkzBaYoVUUSms5RASz4+cwJVCwTX91
Y1jcpVIBZmaaY3a0wrx13ajEAa0dOQJBAIpjnb4wqpsEh7VpmJqOdSdGxb1XXfFQ
sA0T1OQYqQnFppWwqrxIL+d9pZdiA1ITnNqyvUFBNETqDSOrUHwwb2cCQGArE+vu
ffPUWQ0j+fiK+covFG8NL7H+26NSGB5+Xsn9uwOGLj7K/YT6CbBtr9hJiuWjM1Al
0V4ltlTuu2mTMaw=
-----END PRIVATE KEY-----
EOT
  key_type              = "RSA"
  modulus               = "mod2048"
  name                  = "test_name"
  name_alias            = "name_alias_1"
  owner_key             = "owner_key_1"
  owner_tag             = "owner_tag_1"
  regenerate            = "no"
  certificate_authority = "test_name"
  annotations = [
    {
      key   = "key_0"
      value = "value_1"
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
  ]
  tags = [
    {
      key   = "key_0"
      value = "value_1"
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
  ]
}
