---
layout: "aci"
page_title: "ACI: aci_l2_domain"
sidebar_current: "docs-aci-data-source-l2_domain"
description: |-
  Data source for ACI L2 Domain
---

# aci_l2_domain #
Data source for ACI L2 Domain

## Example Usage ##

```hcl
data "aci_l2_domain" "fool2_domain" {
  name  = "l2_domain_1"
}
```
## Argument Reference ##
* `name` - (Required) Name of object L2 Domain.

## Attribute Reference

* `id` - Attribute id set to the Dn of the L2 Domain.
* `annotation` - (Optional) Annotation for object L2 Domain.
* `name_alias` - (Optional) Name alias for object L2 Domain.
