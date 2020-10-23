---
layout: "aci"
page_title: "ACI: aci_filter"
sidebar_current: "docs-aci-data-source-filter"
description: |-
  Data source for ACI Filter
---

# aci_filter #
Data source for ACI Filter

## Example Usage ##

```hcl
data "aci_filter" "allow_http" {
  tenant_dn  = "${aci_tenant.dev_tenant.id}"
  name       = "http_flt"
}
```
## Argument Reference ##
* `tenant_dn` - (Required) Distinguished name of parent Tenant object.
* `name` - (Required) name of Object filter.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Filter.
* `annotation` - (Optional) annotation for object filter.
* `name_alias` - (Optional) name_alias for object filter.
