---
layout: "aci"
page_title: "ACI: aci_monitoring_policy"
sidebar_current: "docs-aci-data-source-monitoring_policy"
description: |-
  Data source for ACI Monitoring Policy
---

# aci_monitoring_policy #
Data source for ACI Monitoring Policy

## Example Usage ##

```hcl
data "aci_monitoring_policy" "example" {
  tenant_dn = "example"
  name  = "example"
}
```


## Argument Reference ##
* `tenant_dn` - (Required) tenant dn of object monitoring policy.
* `name` - (Required) name of object monitoring policy.


## Attribute Reference

* `id` - Attribute id set to the Dn of the object monitoring policy.
* `annotation` - (Optional) annotation for object monitoring policy.
* `name_alias` - (Optional) name_alias for object monitoring policy.
