---
layout: "aci"
page_title: "ACI: aci_lacp_policy"
sidebar_current: "docs-aci-data-source-lacp_policy"
description: |-
  Data source for ACI LACP Policy
---

# aci_lacp_policy #
Data source for ACI LACP Policy

## Example Usage ##

```hcl
data "aci_lacp_policy" "dev_lacp_pol" {
  name  = "foo_lacp_pol"
}
```
## Argument Reference ##
* `name` - (Required) Name of Object lacp_policy.



## Attribute Reference

* `id` - Attribute id set to the Dn of the LACP Policy.
* `annotation` - (Optional) Annotation for object lacp_policy.
* `ctrl` - (Optional) LAG control properties
* `max_links` - (Optional) Maximum number of links.
* `min_links` - (Optional) Minimum number of links in port channel.
* `mode` - (Optional) Policy mode.
* `name_alias` - (Optional) Name_alias for object lacp_policy.
