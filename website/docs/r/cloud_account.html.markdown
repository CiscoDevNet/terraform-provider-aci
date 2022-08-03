---
subcategory: "Cloud"
layout: "aci"
page_title: "ACI: aci_cloud_account"
sidebar_current: "docs-aci-resource--cloud-account"
description: |-
  Manages ACI Cloud Account
---

# aci_cloud_account #

Manages ACI Cloud Account

## API Information ##

* `Class` - cloudAccount
* `Distinguished Name` - uni/tn-{tenant_name}/act-[{account_id}]-vendor-{vendor}

## GUI Information ##

* `Location` - Cloud APIC -> Application Management -> Tenants


## Example Usage ##

```hcl
resource "aci_cloud_account" "example" {
  tenant_dn  = aci_tenant.example.id
  account_id  = "example"
  vendor  = "example"
  access_type = "managed"
  account_id = "example"
  vendor = "azure"
  cloud_rs_account_to_access_policy = aci_resource.example.id
  cloud_rs_credentials = aci_resource.example.id
}
```

## Argument Reference ##

* `tenant_dn` - (Required) Distinguished name of the parent Tenant object.
* `account_id` - (Required) ID of the Cloud Account object.
* `vendor` - (Required) Vendor of the Cloud Account object.
* `name` - (Optional) Name of the Cloud Account object.
* `annotation` - (Optional) Annotation of the Cloud Account object.
* `access_type` - (Optional) Authentication to the Cloud Account (managed=no credentials required (IAM), credentials=using accessKeys). Allowed values are "credentials", "managed". Default value is "managed". Type: String.
* `account_id` - (Optional) Account id. An object identifier.
* `vendor` - (Optional) The Cloud Account vendor of the controller. Allowed values are "aws", "azure", "gcp", "unknown", and default value is "aws". Type: String.

* `relation_cloud_rs_account_to_access_policy` - (Optional) Represents the relation to a Relation to the Access policy to be used (class cloudAccessPolicy). Relation to CloudAccessPolicy cardianity is n-1 Type: String.


* `relation_cloud_rs_credentials` - (Optional) Represents the relation to a Credentials to use to manage the account (class cloudCredentials). If access type is credentials, relation to the credentials to use Type: String.



## Importing ##

An existing Cloud Account can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_cloud_account.example "<Dn>"
```