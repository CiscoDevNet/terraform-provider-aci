---
layout: "aci"
page_title: "ACI: aci_cloud_aws_provider"
sidebar_current: "docs-aci-data-source-cloud_aws_provider"
description: |-
  Data source for ACI Cloud AWS Provider
---

# aci_cloud_aws_provider #
Data source for ACI Cloud AWS Provider
Note: This resource is supported in Cloud APIC only.
## Example Usage ##

```hcl
data "aci_cloud_aws_provider" "example" {

  tenant_dn  = "${aci_tenant.example.id}"
}
```
## Argument Reference ##
* `tenant_dn` - (Required) Distinguished name of parent Tenant object.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Cloud AWS Provider.
* `access_key_id` - (Optional) access_key_id for object cloud_aws_provider.
* `account_id` - (Optional) account_id for object cloud_aws_provider.
* `annotation` - (Optional) annotation for object cloud_aws_provider.
* `email` - (Optional) email address of the local user
* `http_proxy` - (Optional) http_proxy for object cloud_aws_provider.
* `is_account_in_org` - (Optional) is_account_in_org for object cloud_aws_provider.
* `is_trusted` - (Optional) is_trusted for object cloud_aws_provider.
* `name_alias` - (Optional) name_alias for object cloud_aws_provider.
* `provider_id` - (Optional) provider_id for object cloud_aws_provider.
* `region` - (Optional) region for object cloud_aws_provider.
* `secret_access_key` - (Optional) secret_access_key for object cloud_aws_provider.
