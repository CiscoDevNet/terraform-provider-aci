
resource "aci_key_ring" "full_example" {
  admin_state = "completed"
  annotation  = "annotation"
  cert        = "-----BEGIN CERTIFICATE-----\nMIICODCCAaGgAwIBAgIJAIt8XMntue0VMA0GCSqGSIb3DQEBCwUAMDQxDjAMBgNV\nBAMMBUFkbWluMRUwEwYDVQQKDAxZb3VyIENvbXBhbnkxCzAJBgNVBAYTAlVTMCAX\nDTE4MDEwOTAwNTk0NFoYDzIxMTcxMjE2MDA1OTQ0WjA0MQ4wDAYDVQQDDAVBZG1p\nbjEVMBMGA1UECgwMWW91ciBDb21wYW55MQswCQYDVQQGEwJVUzCBnzANBgkqhkiG\n9w0BAQEFAAOBjQAwgYkCgYEAohG/7axtt7CbSaMP7r+2mhTKbNgh0Ww36C7Ta14i\nv+VmLyKkQHnXinKGhp6uy3Nug+15a+eIu7CrgpBVMQeCiWfsnwRocKcQJWIYDrWl\nXHxGQn31yYKR6mylE7Dcj3rMFybnyhezr5D8GcP85YRPmwG9H2hO/0Y1FUnWu9Iw\nAQkCAwEAAaNQME4wHQYDVR0OBBYEFD0jLXfpkrU/ChzRvfruRs/fy1VXMB8GA1Ud\nIwQYMBaAFD0jLXfpkrU/ChzRvfruRs/fy1VXMAwGA1UdEwQFMAMBAf8wDQYJKoZI\nhvcNAQELBQADgYEAOmvre+5tgZ0+F3DgsfxNQqLTrGiBgGCIymPkP/cBXXkNuJyl\n3ac7tArHQc7WEA4U2R2rZbEq8FC3UJJm4nUVtCPvEh3G9OhN2xwYev79yt6pIn/l\nKU0Td2OpVyo0eLqjoX5u2G90IBWzhyjFbo+CcKMrSVKj1YOdG0E3OuiJf00=\n-----END CERTIFICATE-----\n"
  description = "description"
  ecc_curve   = "none"
  key         = "-----BEGIN PRIVATE KEY-----\nMIICdwIBADANBgkqhkiG9w0BAQEFAASCAmEwggJdAgEAAoGBAKIRv+2sbbewm0mj\nD+6/tpoUymzYIdFsN+gu02teIr/lZi8ipEB514pyhoaerstzboPteWvniLuwq4KQ\nVTEHgoln7J8EaHCnECViGA61pVx8RkJ99cmCkepspROw3I96zBcm58oXs6+Q/BnD\n/OWET5sBvR9oTv9GNRVJ1rvSMAEJAgMBAAECgYByu3QO0qF9h7X3JEu0Ld4cKBnB\ngiQ2uJC/et7KxIJ/LOvw9GopBthyt27KwG1ntBkJpkTuAaQHkyNns7vLkNB0S0IR\n+owVFEcKYq9VCHTaiQU8TDp24gN+yPTrpRuH8YhDVq5SfVdVuTMgHVQdj4ya4VlF\nGj+a7+ipxtGiLsVGrQJBAM7p0Fm0xmzi+tBOASUAcVrPLcteFIaTBFwfq16dm/ON\n00Khla8Et5kMBttTbqbukl8mxFjBEEBlhQqb6EdQQ0sCQQDIhHx1a9diG7y/4DQA\n4KvR3FCYwP8PBORlSamegzCo+P1OzxiEo0amX7yQMA5UyiP/kUsZrme2JBZgna8S\np4R7AkEAr7rMhSOPUnMD6V4WgsJ5g1Jp5kqkzBaYoVUUSms5RASz4+cwJVCwTX91\nY1jcpVIBZmaaY3a0wrx13ajEAa0dOQJBAIpjnb4wqpsEh7VpmJqOdSdGxb1XXfFQ\nsA0T1OQYqQnFppWwqrxIL+d9pZdiA1ITnNqyvUFBNETqDSOrUHwwb2cCQGArE+vu\nffPUWQ0j+fiK+covFG8NL7H+26NSGB5+Xsn9uwOGLj7K/YT6CbBtr9hJiuWjM1Al\n0V4ltlTuu2mTMaw=\n-----END PRIVATE KEY-----\n"
  key_type    = "RSA"
  modulus     = "mod1024"
  name        = "test_name"
  name_alias  = "name_alias"
  owner_key   = "owner_key"
  owner_tag   = "owner_tag"
  regen       = "no"
  tp          = "test_name"
  annotations = [
    {
      key   = "annotations_1"
      value = "value_1"
    }
  ]
}

resource "aci_key_ring" "full_example_tenant" {
  parent_dn   = aci_tenant.example.id
  admin_state = "completed"
  annotation  = "annotation"
  cert        = "-----BEGIN CERTIFICATE-----\nMIICODCCAaGgAwIBAgIJAIt8XMntue0VMA0GCSqGSIb3DQEBCwUAMDQxDjAMBgNV\nBAMMBUFkbWluMRUwEwYDVQQKDAxZb3VyIENvbXBhbnkxCzAJBgNVBAYTAlVTMCAX\nDTE4MDEwOTAwNTk0NFoYDzIxMTcxMjE2MDA1OTQ0WjA0MQ4wDAYDVQQDDAVBZG1p\nbjEVMBMGA1UECgwMWW91ciBDb21wYW55MQswCQYDVQQGEwJVUzCBnzANBgkqhkiG\n9w0BAQEFAAOBjQAwgYkCgYEAohG/7axtt7CbSaMP7r+2mhTKbNgh0Ww36C7Ta14i\nv+VmLyKkQHnXinKGhp6uy3Nug+15a+eIu7CrgpBVMQeCiWfsnwRocKcQJWIYDrWl\nXHxGQn31yYKR6mylE7Dcj3rMFybnyhezr5D8GcP85YRPmwG9H2hO/0Y1FUnWu9Iw\nAQkCAwEAAaNQME4wHQYDVR0OBBYEFD0jLXfpkrU/ChzRvfruRs/fy1VXMB8GA1Ud\nIwQYMBaAFD0jLXfpkrU/ChzRvfruRs/fy1VXMAwGA1UdEwQFMAMBAf8wDQYJKoZI\nhvcNAQELBQADgYEAOmvre+5tgZ0+F3DgsfxNQqLTrGiBgGCIymPkP/cBXXkNuJyl\n3ac7tArHQc7WEA4U2R2rZbEq8FC3UJJm4nUVtCPvEh3G9OhN2xwYev79yt6pIn/l\nKU0Td2OpVyo0eLqjoX5u2G90IBWzhyjFbo+CcKMrSVKj1YOdG0E3OuiJf00=\n-----END CERTIFICATE-----\n"
  description = "description"
  ecc_curve   = "none"
  key         = "-----BEGIN PRIVATE KEY-----\nMIICdwIBADANBgkqhkiG9w0BAQEFAASCAmEwggJdAgEAAoGBAKIRv+2sbbewm0mj\nD+6/tpoUymzYIdFsN+gu02teIr/lZi8ipEB514pyhoaerstzboPteWvniLuwq4KQ\nVTEHgoln7J8EaHCnECViGA61pVx8RkJ99cmCkepspROw3I96zBcm58oXs6+Q/BnD\n/OWET5sBvR9oTv9GNRVJ1rvSMAEJAgMBAAECgYByu3QO0qF9h7X3JEu0Ld4cKBnB\ngiQ2uJC/et7KxIJ/LOvw9GopBthyt27KwG1ntBkJpkTuAaQHkyNns7vLkNB0S0IR\n+owVFEcKYq9VCHTaiQU8TDp24gN+yPTrpRuH8YhDVq5SfVdVuTMgHVQdj4ya4VlF\nGj+a7+ipxtGiLsVGrQJBAM7p0Fm0xmzi+tBOASUAcVrPLcteFIaTBFwfq16dm/ON\n00Khla8Et5kMBttTbqbukl8mxFjBEEBlhQqb6EdQQ0sCQQDIhHx1a9diG7y/4DQA\n4KvR3FCYwP8PBORlSamegzCo+P1OzxiEo0amX7yQMA5UyiP/kUsZrme2JBZgna8S\np4R7AkEAr7rMhSOPUnMD6V4WgsJ5g1Jp5kqkzBaYoVUUSms5RASz4+cwJVCwTX91\nY1jcpVIBZmaaY3a0wrx13ajEAa0dOQJBAIpjnb4wqpsEh7VpmJqOdSdGxb1XXfFQ\nsA0T1OQYqQnFppWwqrxIL+d9pZdiA1ITnNqyvUFBNETqDSOrUHwwb2cCQGArE+vu\nffPUWQ0j+fiK+covFG8NL7H+26NSGB5+Xsn9uwOGLj7K/YT6CbBtr9hJiuWjM1Al\n0V4ltlTuu2mTMaw=\n-----END PRIVATE KEY-----\n"
  key_type    = "RSA"
  modulus     = "mod1024"
  name        = "test_name"
  name_alias  = "name_alias"
  owner_key   = "owner_key"
  owner_tag   = "owner_tag"
  regen       = "no"
  tp          = "test_name"
  annotations = [
    {
      key   = "annotations_1"
      value = "value_1"
    }
  ]
}
