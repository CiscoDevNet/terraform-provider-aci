package provider

import (
	"fmt"
	"slices"
)

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
  name      = "test_name"
}

resource "aci_application_epg" "test" {
  application_profile_dn = aci_application_profile.test.id
  name                   = "test_name"
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

const testConfigVmmDomPMin = `
resource "aci_vmm_domain" "test" {
  provider_profile_dn = "uni/vmmp-VMware"
  name                = "test_vmm_domain"
}
`

const testConfigVmmVSwitchPolicyContMinDependencyWithVmmDomP = testConfigVmmDomPMin + `
resource "aci_vswitch_policy" "test" {
  vmm_domain_dn = aci_vmm_domain.test.id
}
`

const testConfigImportedVnsLDevVipWithFvTenant = `
resource "aci_tenant" "test_tenant_imported_device" {
  name = "test_tenant_imported_device"
}

resource "aci_physical_domain" "test" {
  name = "test"
}

resource "aci_l4_l7_device" "test_imported_device" {
  tenant_dn                            = aci_tenant.test_tenant_imported_device.id
  name                                 = "test_imported_device"
  active                               = "no"
  context_aware                        = "single-Context"
  device_type                          = "PHYSICAL"
  function_type                        = "GoTo"
  is_copy                              = "no"
  mode                                 = "legacy-Mode"
  promiscuous_mode                     = "no"
  service_type                         = "OTHERS"
  relation_vns_rs_al_dev_to_phys_dom_p = "uni/phys-test"
}
`

func GetParentDnForTesting(parentClassName, className string) string {
	parentClassDnMap := map[string]string{
		"fvTenant":             "uni/tn-test_name",
		"fvAp":                 "uni/tn-test_name/ap-test_name",
		"fvESg":                "uni/tn-test_name/ap-test_name/esg-test_name",
		"fvAEPg":               "uni/tn-test_name/ap-test_name/epg-test_name",
		"fvCrtrn":              "uni/tn-test_name/ap-test_name/epg-test_name/crtrn",
		"fvSCrtrn":             "uni/tn-test_name/ap-test_name/epg-test_name/crtrn/crtrn-sub_criterion",
		"fvCtx":                "uni/tn-test_name/ctx-test_name",
		"fvFBRGroup":           "uni/tn-test_name/ctx-test_name/fbrg-fallback_route_group",
		"fvBD":                 "uni/tn-test_name/BD-test_name",
		"fvSiteAssociated":     "uni/tn-test_name/BD-test_name/stAsc",
		"vzTaboo":              "uni/tn-test_name/taboo-test_name",
		"vzTSubj":              "uni/tn-test_name/taboo-test_name/tsubj-test_name",
		"fvTrackList":          "uni/tn-test_name/tracklist-test_name",
		"netflowMonitorPol":    "uni/tn-test_name/monitorpol-netfow_monitor",
		"qosCustomPol":         "uni/tn-test_name/qoscustom-test_name",
		"pimRouteMapPol":       "uni/tn-test_name/rtmap-test_name",
		"infraAccPortP":        "uni/infra/accportprof-test_name",
		"infraSpAccPortP":      "uni/infra/spaccportprof-test_name",
		"infraFexP":            "uni/infra/fexprof-test_name",
		"infraHPortS":          "uni/infra/accportprof-test_name/hports-test_name-typ-range",
		"infraSHPortS":         "uni/infra/spaccportprof-test_name/shports-test_name-typ-ALL",
		"infraAttEntityP":      "uni/infra/attentp-test_name",
		"l3extOut":             "uni/tn-test_name/out-test_l3_outside",
		"l3extLoopBackIfP":     "uni/tn-test_name/out-test_l3_outside/lnodep-logical_node_profile/rsnodeL3OutAtt-[topology/pod-2/node-2011]/lbp-[1.2.3.5]",
		"l3extInstP":           "uni/tn-test_name/out-test_l3_outside/instP-test_name",
		"l3extConsLbl":         "uni/tn-test_name/out-test_l3_outside/conslbl-test_name",
		"mgmtInstP":            "uni/tn-mgmt/extmgmt-default/instp-test_name",
		"vmmDomP":              "uni/vmmp-VMware/dom-test_vmm_domain",
		"vmmVSwitchPolicyCont": "uni/vmmp-VMware/dom-test_vmm_domain/vswitchpolcont",
		"vmmUplinkPCont":       "uni/vmmp-VMware/dom-test_vmm_domain/uplinkpcont",
		"vzAny":                "uni/tn-test_name/ctx-test_name/any",
	}

	classDnMap := map[string]string{
		"pkiTP":        "uni/userext/pkiext",
		"pkiKeyRing":   "uni/userext/pkiext",
		"l3extProvLbl": "uni/tn-infra/out-test_l3_outside",
	}

	// ignore l3extProvLbl because the l3extOut parent in tree requires to point the infra tenant
	ignoreClassName := []string{"l3extProvLbl"}
	if !slices.Contains(ignoreClassName, className) {
		if value, ok := parentClassDnMap[parentClassName]; ok {
			return fmt.Sprintf("%s/", value)
		}
	}

	if value, ok := classDnMap[className]; ok {
		return fmt.Sprintf("%s/", value)
	}

	return ""
}
