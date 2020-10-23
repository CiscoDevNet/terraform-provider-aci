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
* `name` - (Required) name of Object physical_domain.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Physical Domain.
* `annotation` - (Optional) annotation for object physical_domain.
* `name_alias` - (Optional) name_alias for object physical_domain.
