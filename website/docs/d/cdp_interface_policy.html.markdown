---
layout: "aci"
page_title: "ACI: aci_cdp_interface_policy"
sidebar_current: "docs-aci-data-source-cdp_interface_policy"
description: |-
  Data source for ACI CDP Interface Policy
---

# aci_cdp_interface_policy #
Data source for ACI CDP Interface Policy

## Example Usage ##

```hcl
data "aci_cdp_interface_policy" "example" {


  name  = "example"
}
```
## Argument Reference ##
* `name` - (Required) name of Object cdp_interface_policy.



## Attribute Reference

* `id` - Attribute id set to the Dn of the CDP Interface Policy.
* `admin_st` - (Optional) administrative state
* `annotation` - (Optional) annotation for object cdp_interface_policy.
* `name_alias` - (Optional) name_alias for object cdp_interface_policy.
