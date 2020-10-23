---
layout: "aci"
page_title: "ACI: aci_pod_maintenance_group"
sidebar_current: "docs-aci-data-source-pod_maintenance_group"
description: |-
  Data source for ACI POD Maintenance Group
---

# aci_pod_maintenance_group #
Data source for ACI POD Maintenance Group

## Example Usage ##

```hcl
data "aci_pod_maintenance_group" "example" {


  name  = "example"
}
```
## Argument Reference ##
* `name` - (Required) name of Object pod_maintenance_group.



## Attribute Reference

* `id` - Attribute id set to the Dn of the POD Maintenance Group.
* `annotation` - (Optional) annotation for object pod_maintenance_group.
* `fwtype` - (Optional) fwtype for object pod_maintenance_group.
* `name_alias` - (Optional) name_alias for object pod_maintenance_group.
* `pod_maintenance_group_type` - (Optional) component type
