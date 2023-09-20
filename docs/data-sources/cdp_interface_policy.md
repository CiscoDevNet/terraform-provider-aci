---
subcategory: "Access Policies"
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
* `name` - (Required) Name of Object cdp interface policy.



## Attribute Reference

* `id` - Attribute id set to the Dn of the CDP Interface Policy.
* `admin_st` - (Optional) Administrative state
* `annotation` - (Optional) Annotation for object cdp interface policy.
* `name_alias` - (Optional) Name alias for object cdp interface policy.
* `description` - (Optional) Description for object cdp interface policy.
