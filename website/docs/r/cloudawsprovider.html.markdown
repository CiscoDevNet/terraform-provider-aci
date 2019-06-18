---
layout: "aci"
page_title: "ACI: aci_cloud_aws_provider"
sidebar_current: "docs-aci-resource-cloud_aws_provider"
description: |-
  Manages ACI Cloud AWS Provider
---

# aci_cloud_aws_provider #
Manages ACI Cloud AWS Provider
Note: This resource is supported in Cloud APIC only.
## Example Usage ##

```hcl
resource "aci_cloud_aws_provider" "example" {

  tenant_dn  = "${aci_tenant.example.id}"
  access_key_id  = "example"
  account_id  = "example"
  annotation  = "example"
  email  = "example"
  http_proxy  = "example"
  is_account_in_org  = "example"
  is_trusted  = "example"
  name_alias  = "example"
  provider_id  = "example"
  region  = "example"
  secret_access_key  = "example"
}
```
## Argument Reference ##
* `tenant_dn` - (Required) Distinguished name of parent Tenant object.
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



## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Cloud AWS Provider.

## Importing ##

An existing Cloud AWS Provider can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_cloud_aws_provider.example <Dn>
```