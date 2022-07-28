---
subcategory: "Cloud"
layout: "aci"
page_title: "ACI: aci_account"
sidebar_current: "docs-aci-data-source-account"
description: |-
  Data source for ACI Account
---

# aci_cloud_account #

Data source for ACI Cloud Account


## API Information ##

* `Class` - cloudAccount
* `Distinguished Name` - uni/tn-{tenant_name}/act-[{account_id}]-vendor-{vendor}

## GUI Information ##

* `Location` - Azure APIC -> Application Management -> Tenants



## Example Usage ##

```hcl
data "aci_account" "example" {
  tenant_dn  = aci_tenant.example.id
  account_id  = "example"
  vendor  = "example"
}
```

## Argument Reference ##

* `tenant_dn` - (Required) Distinguished name of the parent Tenant object.
* `account_id` - (Required) ID of the Cloud Account object.
* `vendor` - (Required) Vendor of the Cloud Account object.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the Cloud Account.
* `annotation` - (Optional) Annotation of the Cloud Account object.
* `name_alias` - (Optional) Name Alias of the Cloud Account object.
* `access_type` - (Optional) Account acccess type. How to authenticate to the account (managed=no credentials required (IAM), credentials=using accessKeys)
