---
layout: "aci"
page_title: "ACI: aci_l2_outside"
sidebar_current: "docs-aci-data-source-l2_outside"
description: |-
  Data source for ACI L2 Outside
---

# aci_l2_outside #
Data source for ACI L2 Outside

## Example Usage ##

```hcl
data "aci_l2_outside" "example" {
  tenant_dn  = "${aci_tenant.example.id}"
  name  = "example"
}
```


## Argument Reference ##
* `tenant_dn` - (Required) Distinguished name of parent Tenant object.
* `name` - (Required) name of Object l2_outside.



## Attribute Reference

* `id` - Attribute id set to the Dn of the L2 Outside.
* `annotation` - (Optional) annotation for object l2_outside.
* `name_alias` - (Optional) name_alias for object l2_outside.
* `target_dscp` - (Optional) target dscp
