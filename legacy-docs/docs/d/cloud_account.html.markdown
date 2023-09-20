---
subcategory: "Cloud"
layout: "aci"
page_title: "ACI: aci_cloud_account"
sidebar_current: "docs-aci-data-source-cloud-account"
description: |-
  Data source for Cloud Network Controller Cloud Account
---

# aci_cloud_account #

Data source for Cloud Network Controller Cloud Account
Note: This data source is supported in Cloud Network Controller only.

## API Information ##

* `Class` - cloudAccount
* `Distinguished Name` - uni/tn-{tenant_name}/act-[{account_id}]-vendor-{vendor}

## GUI Information ##

* `Location` - Cloud Network Controller -> Application Management -> Tenants -> {tenant_name}



## Example Usage ##

```hcl
data "aci_cloud_account" "example" {
  tenant_dn  = aci_tenant.example.id
  account_id  = "azure_account_id"
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
* `access_type` - (Optional) Authentication type for the Cloud Account (managed=no credentials required (IAM), credentials=using accessKeys).
