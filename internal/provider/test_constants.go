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
