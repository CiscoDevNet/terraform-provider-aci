---
layout: "aci"
page_title: "ACI: aci_physical_domain"
sidebar_current: "docs-aci-data-source-physical_domain"
description: |-
  Data source for ACI Physical Domain
---

# aci_physical_domain #
Data source for ACI Physical Domain

## Example Usage ##

```hcl
data "aci_physical_domain" "example" {
  name  = "example"
}
```
## Argument Reference ##
* `name` - (Required) Name of object physical domain..



## Attribute Reference

* `id` - Attribute id set to the Dn of the Physical Domain.
* `annotation` - (Optional) Annotation for object physical domain.
* `name_alias` - (Optional) Name alias for object physical domain.