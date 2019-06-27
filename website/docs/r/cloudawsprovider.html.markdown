---
layout: "aci"
page_title: "ACI: aci_cloud_aws_provider"
sidebar_current: "docs-aci-resource-cloud_aws_provider"
description: |-
  Manages ACI Cloud AWS Provider
---

# aci_cloud_aws_provider #
Manages ACI Cloud AWS Provider
<b>Note: This resource is supported in Cloud APIC only.</b>
## Example Usage ##

```hcl
	resource "aci_cloud_aws_provider" "foocloud_aws_provider" {
		tenant_dn         = "${aci_tenant.footenant.id}"
		description       = "aws account config"
		access_key_id     = "access_key"
		account_id        = "acc_id"
		annotation        = "tag_aws"
		region            = "us-west-2"
		secret_access_key = "secret_key"
	}
```
## Argument Reference ##
* `tenant_dn` - (Required) Distinguished name of parent Tenant object.
* `access_key_id` - (Optional) access_key_id for the AWS account provided in the account_id field.
* `account_id` - (Optional) AWS account-id to manage with cloud APIC.
* `annotation` - (Optional) annotation for object cloud_aws_provider.
* `email` - (Optional) email address of the local user.
* `http_proxy` - (Optional) http_proxy for object cloud_aws_provider.
* `is_account_in_org` - (Optional) Flag to decide whether the account is in the organization or not.
* `is_trusted` - (Optional) Whether the account is trusted with Tenant infra account.
* `name_alias` - (Optional) name_alias for object cloud_aws_provider.
* `provider_id` - (Optional) provider_id for object cloud_aws_provider.
* `region` - (Optional) which AWS region to manage.
* `secret_access_key` - (Optional) secret_access_key for the AWS account provided in the account_id field.



## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Cloud AWS Provider.

## Importing ##

An existing Cloud AWS Provider can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_cloud_aws_provider.example <Dn>
```