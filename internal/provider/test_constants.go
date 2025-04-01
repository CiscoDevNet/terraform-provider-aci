package provider

const testConstantsConfigFvTenantMin = `
resource "aci_tenant" "test" {
  name = "test_name"
}
`

const testConfigL3extOutMin = testConstantsConfigFvTenantMin + `
resource "aci_vrf" "test" {
  tenant_dn = aci_tenant.test.id
  name      = "test_vrf"
}

resource "aci_l3_outside" "test" {
  tenant_dn = aci_tenant.test.id
  name      = "test_l3_outside"
  relation_l3ext_rs_ectx = aci_vrf.test.id
}
`

const testConfigFvAEPgMin = testConstantsConfigFvTenantMin + `
resource "aci_application_profile" "test" {
  tenant_dn = aci_tenant.test.id
  name      = "test_ap"
}

resource "aci_application_epg" "test" {
  application_profile_dn = aci_application_profile.test.id
  name                   = "test_epg"
}
`

const testConfigL3extOutMinDependencyWithFvTenant = testConstantsConfigFvTenantMin + `
resource "aci_vrf" "test" {
  tenant_dn = aci_tenant.test.id
  name      = "test_vrf"
}

resource "aci_l3_outside" "test" {
  tenant_dn = aci_tenant.test.id
  name      = "test_l3_outside"
  relation_l3ext_rs_ectx = aci_vrf.test.id
}
`

const testConfigFvCtxMinDependencyWithFvTenant = testConstantsConfigFvTenantMin + `
resource "aci_vrf" "test" {
  tenant_dn = aci_tenant.test.id
  name      = "test_vrf"
}
`

const testConfigFvTenantInfraMin = `
data "aci_tenant" "test" {
  name = "infra"
}
`

const testConfigL3extOutMinDependencyWithFvTenantInfra = testConfigFvTenantInfraMin + `
resource "aci_vrf" "test" {
  tenant_dn = data.aci_tenant.test.id
  name      = "test_vrf"
}

resource "aci_l3_outside" "test" {
  tenant_dn = data.aci_tenant.test.id
  name      = "test_l3_outside"
  relation_l3ext_rs_ectx = aci_vrf.test.id
}
`

const testConfigL3extLoopBackIfPMinDependencyWithL3extRsNodeL3OutAtt = testConfigL3extOutMin + `
resource "aci_logical_node_profile" "test" {
  l3_outside_dn = aci_l3_outside.test.id
  name          = "logical_node_profile"
} 

resource "aci_logical_node_to_fabric_node" "test" {
  logical_node_profile_dn = aci_logical_node_profile.test.id
  tdn                     = "topology/pod-2/node-2011"
  rtr_id                  = "1.2.3.4"
}

resource "aci_l3out_loopback_interface_profile" "test" {
  fabric_node_dn = aci_logical_node_to_fabric_node.test.id
  addr           = "1.2.3.5"
}
`

const testConstantsConfigFvTenantMinDependencyWithPkiTP = testConstantsConfigFvTenantMin + `
resource "aci_certificate_authority" "test" {
  parent_dn = aci_tenant.test.id
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
  name = "test_name"
}
`

const testConfigDefaultMinDependencyWithPkiTP = `
resource "aci_certificate_authority" "test" {
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
  name = "test_name"
}
`

const testConfigDataSourceSystem = `
data "aci_system" "version" {
  pod_id    = "1"
  system_id = "1"
}
`
