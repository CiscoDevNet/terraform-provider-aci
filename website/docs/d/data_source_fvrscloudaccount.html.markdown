---
layout: "aci"
page_title: "ACI: aci_tenanttoaccountassociation"
sidebar_current: "docs-aci-data-source-tenanttoaccountassociation"
description: |-
  Data source for ACI Tenant to account association
---

# aci_tenanttoaccountassociation #

Data source for ACI Tenant to account association


## API Information ##

* `Class` - fvRsCloudAccount
* `Distinguished Name` - uni/tn-{name}/rsCloudAccount

## GUI Information ##

* `Location` - 



## Example Usage ##

```hcl
data "aci_tenanttoaccountassociation" "example" {
  tenant_dn  = aci_tenant.example.id
}
```

## Argument Reference ##

* `tenant_dn` - (Required) Distinguished name of parent Tenant object.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the Tenant to account association.
* `annotation` - (Optional) Annotation of object Tenant to account association.
* `name_alias` - (Optional) Name Alias of object Tenant to account association.
* `t_dn` - (Optional) Target-dn. The distinguished name of the target.
