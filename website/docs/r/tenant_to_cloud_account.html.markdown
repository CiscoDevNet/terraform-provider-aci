---
subcategory: "Cloud"
layout: "aci"
page_title: "ACI: aci_tenant_to_cloud_account"
sidebar_current: "docs-aci-resource-tenant-to-cloud-account"
description: |-
  Manages ACI Tenant to Cloud Account association
---

# aci_tenant_to_cloud_account #

Manages ACI Tenant to Cloud Account association

## API Information ##

* `Class` - fvRsCloudAccount
* `Distinguished Name` - uni/tn-{tenant_name}/rsCloudAccount

## GUI Information ##

* `Location` - Cloud APIC -> Application Management -> Tenants  -> {tenant_name}


## Example Usage ##

```hcl
resource "aci_tenant_to_cloud_account" "example" {
  tenant_dn  = aci_tenant.example.id
  cloud_account_dn = aci_cloud_account.example.id
}
```

## Argument Reference ##

* `tenant_dn` - (Required) Distinguished name of the parent Tenant object.
* `cloud_account_dn` - (Optional) The distinguished name of the target Cloud Account object.
* `annotation` - (Optional) Annotation of the Tenant to Cloud Account association object.
* `name_alias` - (Optional) Name Alias of the Tenant to Cloud Account association object.


## Importing ##

An existing Tenant to Cloud Account can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_tenant_to_cloud_account.example "<Dn>"
```