---
layout: "aci"
page_title: "ACI: aci_tenant"
sidebar_current: "docs-aci-data-source-tenant"
description: |-
  Data source for ACI Tenant
---

# aci_tenant #
Data source for ACI Tenant

## Example Usage ##

```hcl
data "aci_tenant" "example" {
  name  = "example"
}
```
## Argument Reference ##
* `name` - (Required) name of Object tenant.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Tenant.
* `annotation` - (Optional) annotation for object tenant.
* `name_alias` - (Optional) name_alias for object tenant.
