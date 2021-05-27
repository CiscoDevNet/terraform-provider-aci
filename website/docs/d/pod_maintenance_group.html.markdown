---
layout: "aci"
page_title: "ACI: aci_pod_maintenance_group"
sidebar_current: "docs-aci-data-source-pod_maintenance_group"
description: |-
  Data source for ACI POD Maintenance Group
---

# aci_pod_maintenance_group

Data source for ACI POD Maintenance Group

## Example Usage

```hcl
data "aci_pod_maintenance_group" "example" {
  name  = "example"
}
```

## Argument Reference

- `name` - (Required) The name for a set of nodes that a maintenance policy can be applied to. The maintenance policy determines the pre-defined action to take when there is a disruptive change made to the service profile associated with the node group.

## Attribute Reference

- `id` - attribute id set to the Dn of pod maintenance group object.
- `annotation` - (Optional) Annotation for pod maintenance group object.
- `description` - (Optional) Description for pod maintenance group object.
- `fwtype` - (Optional) The firmware type for pod maintenance group object. Allowed values are "catalog", "config", "controller", "plugin", "pluginPackage", "switch" and "vpod". Default value is "switch".
- `name_alias` - (Optional) Name alias for pod maintenance group object.
- `pod_maintenance_group_type` - (Optional) component type for pod maintenance group object. Allowed values are "range", "ALL" and "ALL_IN_POD". Default value is "range".
