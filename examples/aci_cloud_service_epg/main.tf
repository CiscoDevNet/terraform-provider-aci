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

resource "aci_tenant" "azure_cloud_tenant_tf_test" {
  name = "azure_terraform_test_tenant_svc_epg"
}

resource "aci_cloud_applicationcontainer" "azure_cloud_app_tf_test" {
  tenant_dn = aci_tenant.azure_cloud_tenant_tf_test.id
  name      = "azure_terraform_test_app_svc_epg"
}

resource "aci_vrf" "azure_cloud_vrf_tf_test" {
  tenant_dn = aci_tenant.azure_cloud_tenant_tf_test.id
  name      = "azure_terraform_test_vrf_svc_epg"
}

resource "aci_cloud_service_epg" "azure_cloud_svc_epg_tf_test_1" {
  cloud_applicationcontainer_dn = aci_cloud_applicationcontainer.azure_cloud_app_tf_test.id
  name                          = "azure_terraform_test_cloud_svc_epg_1"
  access_type                   = "Public"
  deployment_type               = "CloudNative"
  cloud_service_epg_type        = "Azure-SqlServer"
  relation_cloud_rs_cloud_epg_ctx   = aci_vrf.azure_cloud_vrf_tf_test.id
}

resource "aci_cloud_context_profile" "azure_cloud_ctxt_profile_tf_test" {
  tenant_dn = aci_tenant.azure_cloud_tenant_tf_test.id
  name = "azure_terraform_test_ctxt_profile_svc_epg"
  primary_cidr = "7.1.0.0/16"
  region = "westus"
  cloud_vendor = "azure"
  relation_cloud_rs_to_ctx = aci_vrf.azure_cloud_vrf_tf_test.id
}

resource "aci_cloud_cidr_pool" "azure_cloud_cidr_pool_tf_test" {
  cloud_context_profile_dn = aci_cloud_context_profile.azure_cloud_ctxt_profile_tf_test.id
  addr = "7.1.0.0/16"
}

resource "aci_cloud_subnet" "azure_cloud_subnet_tf_test" {
  cloud_cidr_pool_dn = aci_cloud_cidr_pool.azure_cloud_cidr_pool_tf_test.id
  ip = "7.1.0.0/24"
}

resource "aci_cloud_service_epg" "azure_cloud_svc_epg_tf_test_2" {
  cloud_applicationcontainer_dn = aci_cloud_applicationcontainer.azure_cloud_app_tf_test.id
  name                          = "azure_terraform_test_cloud_svc_epg_2"
  access_type                   = "Private"
  deployment_type               = "CloudNativeManaged"
  cloud_service_epg_type        = "Azure-AksCluster"
  relation_cloud_rs_cloud_epg_ctx   = aci_vrf.azure_cloud_vrf_tf_test.id
}

resource "aci_cloud_service_endpoint_selector" "azure_cloud_svc_ep_selector_tf_test_1" {
  cloud_service_epg_dn = aci_cloud_service_epg.azure_cloud_svc_epg_tf_test_2.id
  name = "azure_terraform_test_cloud_svc_ep_selector_1"
  match_expression = "IP=='7.1.0.0/24'"
}

resource "aci_cloud_service_epg" "azure_cloud_svc_epg_tf_test_3" {
  cloud_applicationcontainer_dn = aci_cloud_applicationcontainer.azure_cloud_app_tf_test.id
  name                          = "azure_terraform_test_cloud_svc_epg_3"
  access_type                   = "Private"
  deployment_type               = "Third-party"
  cloud_service_epg_type        = "Custom"
  relation_cloud_rs_cloud_epg_ctx   = aci_vrf.azure_cloud_vrf_tf_test.id
}

resource "aci_cloud_private_link_label" "azure_cloud_private_link_tf_test" {
  cloud_service_epg_dn = aci_cloud_service_epg.azure_cloud_svc_epg_tf_test_3.id
  name = "azure_terraform_test_cloud_svc_private_link_label"
}

resource "aci_cloud_service_endpoint_selector" "azure_cloud_svc_ep_selector_tf_test_2" {
  cloud_service_epg_dn = aci_cloud_service_epg.azure_cloud_svc_epg_tf_test_3.id
  name = "azure_terraform_test_cloud_svc_ep_selector_2"
  match_expression = "URL=='https://cisco.com'"
}