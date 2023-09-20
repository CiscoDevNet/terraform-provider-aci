---
subcategory: "Cloud"
layout: "aci"
page_title: "ACI: aci_tenant_to_cloud_account"
sidebar_current: "docs-aci-data-source-tenant-to-cloud-account"
description: |-
  Data source for ACI Tenant to Cloud Account association
---

# aci_tenant_to_cloud_account #

Data source for ACI Tenant to Cloud Account association
Note: This data source is supported in Cloud Network Controller only.

## API Information ##

* `Class` - fvRsCloudAccount
* `Distinguished Name` - uni/tn-{tenant_name}/rsCloudAccount

## GUI Information ##

* `Location` - Cloud Network Controller -> Application Management -> Tenants -> {tenant_name}



## Example Usage ##

```hcl
data "aci_tenant_to_cloud_account" "example" {
  tenant_dn  = aci_tenant.example.id
  cloud_account_dn = aci_cloud_account.example.id
}
```

## Argument Reference ##

* `tenant_dn` - (Required) Distinguished name of the parent Tenant object.
* `cloud_account_dn` - (Optional) The distinguished name of the target Cloud Account object.

## Attribute Reference ##
* `id` - Attribute id set to Dn of the Tenant to Cloud Account association object.
* `annotation` - (Optional) Annotation of the Tenant to Cloud Account association object.
* `name_alias` - (Optional) Name Alias of the Tenant to Cloud Account association object.

