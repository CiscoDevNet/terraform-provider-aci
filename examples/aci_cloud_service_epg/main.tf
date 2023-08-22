terraform {
  required_providers {
    aci = {
      source = "ciscodevnet/aci"
    }
  }
}

provider "aci" {
  username = "ansible_github_ci"
  password = "sJ94G92#8dq2hx*K4qh"
  url      = "https://20.245.236.136"
  insecure = true
}

resource "aci_tenant" "azure_cloud_tenant_tf_test" {
  name = "azure_terraform_test_tenant_svc_epg"
}

resource "aci_cloud_applicationcontainer" "azure_cloud_app_tf_test" {
  tenant_dn = aci_tenant.azure_cloud_tenant_tf_test.id
  name = "azure_terraform_test_app_svc_epg"
}

resource "aci_cloud_service_epg" "azure_cloud_svc_epg_tf_test"{
  cloud_applicationcontainer_dn = aci_cloud_applicationcontainer.azure_cloud_app_tf_test.id
  name =  "azure_terraform_test_cloud_svc_epg"
  access_type = "Public"
  deployment_type = "CloudNative"
  cloud_service_epg_type = "Azure-SqlServer"
}