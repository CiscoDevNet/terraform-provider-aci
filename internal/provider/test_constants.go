package provider

const testConfigFvTenantMin = `
resource "aci_tenant" "test" {
  name = "test_tenant"
}

`

const testConfigL3extOutMin = testConfigFvTenantMin + `
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

const testConfigFvApMinDependencyWithFvTenant = testConfigFvTenantMin + `
resource "aci_application_profile" "test" {
  tenant_dn = aci_tenant.test.id
  name      = "test_ap"
}
`

const testConfigFvAEPgMinDependencyWithFvAp = testConfigFvApMinDependencyWithFvTenant + `
resource "aci_application_epg" "test" {
  application_profile_dn = aci_application_profile.test.id
  name                   = "test_epg"
}
`

const testConfigFvAEPgMin = testConfigFvTenantMin + `
resource "aci_application_profile" "test" {
  tenant_dn = aci_tenant.test.id
  name      = "test_ap"
}

resource "aci_application_epg" "test" {
  application_profile_dn = aci_application_profile.test.id
  name                   = "test_epg"
}
`

const testConfigL3extOutMinDependencyWithFvTenant = testConfigFvTenantMin + `
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

const testConfigFvCtxMinDependencyWithFvTenant = testConfigFvTenantMin + `
resource "aci_vrf" "test" {
  tenant_dn = aci_tenant.test.id
  name      = "test_vrf"
}
`
const testConfigFvAEPgMinDependencyWithFvTenant = testConfigFvAEPgMin

const testConfigFvSiteAssociatedMinDependencyWithFvCtx = testConfigFvTenantMin + `
resource "aci_vrf" "test" {
  tenant_dn = aci_tenant.test.id
  name      = "test_vrf"
}

resource "aci_associated_site" "test" {
  parent_dn = aci_vrf.test.id
  site_id = "102"
  name   = "test_associated_site"
}
`

const testConfigFvBDMinDependencyWithFvAp = testConfigFvTenantMin + `
resource "aci_bridge_domain" "test" {
  tenant_dn = aci_tenant.test.id
  name      = "test_bd"
}
`

const testConfigFvBDMinDependencyWithFvTenant = testConfigFvBDMinDependencyWithFvAp

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

const testConfigFvBDMinDependencyWithFvTenant = testConfigFvTenantMin + `
resource "aci_bridge_domain" "test" {
  tenant_dn = aci_tenant.test.id
  name      = "test_bd"
}
`

const testConfigFvBDMin = testConfigFvBDMinDependencyWithFvTenant
