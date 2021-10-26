terraform {
  required_providers {
    aci = {
      source = "ciscodevnet/aci"
    }
  }
}

provider "aci" {
  username = ""
  password = ""
  url      = ""
  insecure = true
}

resource "aci_tenant" "footenant" {
  description = "Tenant created while acceptance testing"
  name        = "demo_tenant"
}

resource "aci_cloud_applicationcontainer" "foocloud_applicationcontainer" {
  tenant_dn  = aci_tenant.footenant.id
  name       = "demo_app"
  annotation = "tag_app"
}

resource "aci_cloud_external_epg" "foocloud_external_epg" {
  cloud_applicationcontainer_dn = aci_cloud_applicationcontainer.foocloud_applicationcontainer.id
  name                          = "cloud_ext_epg"
}

resource "aci_cloud_endpoint_selectorfor_external_epgs" "foocloud_endpoint_selectorfor_external_epgs" {
  cloud_external_epg_dn = aci_cloud_external_epg.foocloud_external_epg.id
  subnet                = "0.0.0.0/0"
  description           = "sample external ep selector"
  name                  = "ext_ep_selector"
  annotation            = "tag_ext_selector"
  is_shared             = "yes"
  match_expression      = "custom:tag=='provbaz'"
  name_alias            = "alias_select"
}