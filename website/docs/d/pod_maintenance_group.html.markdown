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

* `id` - attribute id set to the Dn of pod maintenance group object.
* `annotation` - annotation for pod maintenance group object.
* `fwtype` - fwtype for pod maintenance group object.
* `name_alias` - name_alias for pod maintenance group object.
* `pod_maintenance_group_type` - component type for pod maintenance group object.
