---
layout: "aci"
page_title: "ACI: aci_l3_outside"
sidebar_current: "docs-aci-data-source-l3_outside"
description: |-
  Data source for ACI L3 Outside
---

# aci_l3_outside #
Data source for ACI L3 Outside

## Example Usage ##

```hcl
data "aci_l3_outside" "dev_l3_out" {
  tenant_dn  = "${aci_tenant.dev_tenant.id}"
  name       = "foo_l3_out"
}
```
## Argument Reference ##
* `tenant_dn` - (Required) Distinguished name of parent Tenant object.
* `name` - (Required) name of Object l3_outside.



## Attribute Reference

* `id` - Attribute id set to the Dn of the L3 Outside.
* `description`- (Optional) Description for object l3_outside.
* `annotation` - (Optional) annotation for object l3_outside.
* `enforce_rtctrl` - (Optional) enforce route control type
* `name_alias` - (Optional) name_alias for object l3_outside.
* `target_dscp` - (Optional) The target differentiated services code point (DSCP) of the path attached to the layer 3 outside profile.
