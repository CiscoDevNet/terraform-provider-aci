---
layout: "aci"
page_title: "ACI: aci_fc_domain"
sidebar_current: "docs-aci-data-source-fc_domain"
description: |-
  Data source for ACI FC Domain
---

# aci_fc_domain #
Data source for ACI FC Domain

## Example Usage ##

```hcl
data "aci_fc_domain" "example" {


  name  = "example"
}
```
## Argument Reference ##
* `name` - (Required) name of Object fc_domain.



## Attribute Reference

* `id` - Attribute id set to the Dn of the FC Domain.
* `annotation` - (Optional) annotation for object fc_domain.
* `name_alias` - (Optional) name_alias for object fc_domain.
