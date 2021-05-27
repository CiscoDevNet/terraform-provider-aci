---
layout: "aci"
page_title: "ACI: aci_cloud_aws_provider"
sidebar_current: "docs-aci-resource-cloud_aws_provider"
description: |-
  Manages ACI Cloud AWS Provider
---

# aci_cloud_aws_provider

Manages ACI Cloud AWS Provider
<b>Note: This resource is supported in Cloud APIC only.</b>

## Example Usage

```hcl
	resource "aci_cloud_aws_provider" "foocloud_aws_provider" {
		tenant_dn         = aci_tenant.footenant.id
		description       = "aws account config"
		access_key_id     = "access_key"
		account_id        = "acc_id"
		annotation        = "tag_aws"
		secret_access_key = "secret_key"
	}
```

## Argument Reference

- `tenant_dn` - (Required) Distinguished name of parent Tenant object.
- `access_key_id` - (Optional) access_key_id for the AWS account provided in the account id field.
- `account_id` - (Optional) AWS account-id to manage with cloud APIC.
- `description` - (Optional) Description for object cloud aws provider.
- `annotation` - (Optional) Annotation for object cloud aws provider.
- `email` - (Optional) Email address of the local user.
- `http_proxy` - (Optional) Http proxy for object cloud aws provider.
- `is_account_in_org` - (Optional) Flag to decide whether the account is in the organization or not.
  Allowed values: "no", "yes". Default value: "no".
- `is_trusted` - (Optional) Whether the account is trusted with Tenant infra account.
  Allowed values: "no", "yes". Default value: "no".
- `name_alias` - (Optional) Name alias for object cloud aws provider.
- `provider_id` - (Optional) Provider id for object cloud aws provider.
- `region` - (Optional) Which AWS region to manage. \[Supported only in Cloud APIC 4.2(x) or earlier\]
- `secret_access_key` - (Optional) Secret access key for the AWS account provided in the account id field.

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Cloud AWS Provider.

## Importing

An existing Cloud AWS Provider can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_cloud_aws_provider.example <Dn>
```
