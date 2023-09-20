---
subcategory: "Cloud"
layout: "aci"
page_title: "ACI: aci_cloud_aws_provider"
sidebar_current: "docs-aci-data-source-cloud_aws_provider"
description: |-
  Data source for Cloud Network Controller Cloud AWS Provider
---

# aci_cloud_aws_provider #
Data source for Cloud Network Controller Cloud AWS Provider  
<b>Note: This resource is supported in Cloud Network Controller only.</b>
## Example Usage ##

```hcl
data "aci_cloud_aws_provider" "aws_prov" {
  tenant_dn  = aci_tenant.dev_tenant.id
}
```
## Argument Reference ##
* `tenant_dn` - (Required) Distinguished name of parent Tenant object.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Cloud AWS Provider.
* `access_key_id` - (Optional) access_key_id for the AWS account provided in the account id field.
* `account_id` - (Optional) AWS account-id to manage with Cloud Network Controller.
* `description` - (Optional) Description for object cloud aws provider.
* `annotation` - (Optional) Annotation for object cloud aws provider.
* `email` - (Optional) Email address of the local user.
* `http_proxy` - (Optional) Http proxy for object cloud aws provider.
* `is_account_in_org` - (Optional) Flag to decide whether the account is in the organization or not.
* `is_trusted` - (Optional) Whether the account is trusted with Tenant infra account.
* `name_alias` - (Optional) Name alias for object cloud aws provider.
* `provider_id` - (Optional) Provider id for object cloud aws provider.
* `region` - (Optional) Which AWS region to manage. \[Supported only in Cloud-APIC 4.2(x) or earlier\]
* `secret_access_key` - (Optional) Secret access key for the AWS account provided in the account id field.
