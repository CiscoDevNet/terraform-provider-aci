---
subcategory: "Cloud"
layout: "aci"
page_title: "ACI: aci_cloud_account"
sidebar_current: "docs-aci-data-source-cloud-account"
description: |-
  Data source for ACI Cloud Account
---

# aci_cloud_account #

Data source for ACI Cloud Account


## API Information ##

* `Class` - cloudAccount
* `Distinguished Name` - uni/tn-{tenant_name}/act-[{account_id}]-vendor-{vendor}

## GUI Information ##

* `Location` - Cloud APIC -> Application Management -> Tenants



## Example Usage ##

```hcl
data "aci_cloud_account" "example" {
  tenant_dn  = aci_tenant.example.id
  account_id  = "example"
  vendor  = "azure"
}
```

## Argument Reference ##

* `tenant_dn` - (Required) Distinguished name of the parent Tenant object.
* `account_id` - (Required) ID of the Cloud Account object.
* `vendor` - (Required) Vendor of the Cloud Account object.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the Cloud Account.
* `name` - (Optional) Name of the Cloud Account object.
* `annotation` - (Optional) Annotation of the Cloud Account object.
* `name_alias` - (Optional) Name Alias of the Cloud Account object.
* `access_type` - (Optional) Authentication to the Cloud Account (managed=no credentials required (IAM), credentials=using accessKeys).
