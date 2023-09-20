---
subcategory: "Firmware"
layout: "aci"
page_title: "ACI: aci_pod_maintenance_group"
sidebar_current: "docs-aci-resource-pod_maintenance_group"
description: |-
  Manages ACI POD Maintenance Group
---

# aci_pod_maintenance_group

Manages ACI POD Maintenance Group

## Example Usage

```hcl
resource "aci_pod_maintenance_group" "example" {
  name  = "mgmt"
  fwtype  = "controller"
  description = "from terraform"
  name_alias  = "aliasing"
  pod_maintenance_group_type  = "ALL"
  annotation = "aci_pod_maintenance_group_annotation"
}
```

## Argument Reference

- `name` - (Required) The name for a set of nodes that a maintenance policy can be applied to. The maintenance policy determines the pre-defined action to take when there is a disruptive change made to the service profile associated with the node group.
- `annotation` - (Optional) Annotation for pod maintenance group object.
- `description` - (Optional) Description for pod maintenance group object.
- `fwtype` - (Optional) The firmware type for pod maintenance group object. Allowed values are "catalog", "config", "controller", "plugin", "pluginPackage", "switch" and "vpod". Default value is "switch".
- `name_alias` - (Optional) Name alias for pod maintenance group object.
- `pod_maintenance_group_type` - (Optional) component type for pod maintenance group object. Allowed values are "range", "ALL" and "ALL_IN_POD". Default value is "range".
- `relation_maint_rs_mgrpp` - (Optional) relation to class maintMaintP. Cardinality - N_TO_ONE. Type - String.

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the POD Maintenance Group.

## Importing

An existing POD Maintenance Group can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_pod_maintenance_group.example <Dn>
```
