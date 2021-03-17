---
layout: "aci"
page_title: "ACI: aci_pod_maintenance_group"
sidebar_current: "docs-aci-resource-pod_maintenance_group"
description: |-
  Manages ACI POD Maintenance Group
---

# aci_pod_maintenance_group #
Manages ACI POD Maintenance Group

## Example Usage ##

```hcl
resource "aci_pod_maintenance_group" "example" {
  name  = "mgmt"
  fwtype  = "controller"
  name_alias  = "aliasing"
  pod_maintenance_group_type  = "ALL"
}
```


## Argument Reference ##

* `name` - (Required) name of pod maintenance group object.
* `annotation` - (Optional) annotation for pod maintenance group object.
* `fwtype` - (Optional) fwtype for pod maintenance group object.
Allowed values: "controller", "switch", "catalog", "plugin","pluginPackage", "config", "vpod"
* `name_alias` - (Optional) name_alias for pod maintenance group object.
* `pod_maintenance_group_type` - (Optional) component type for pod maintenance group object. Allowed values are "range", "ALL" and "ALL_IN_POD". Default value is "range".  

* `relation_maint_rs_mgrpp` - (Optional) relation to class maintMaintP. Cardinality - N_TO_ONE. Type - String.
                


## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the POD Maintenance Group.

## Importing ##

An existing POD Maintenance Group can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_pod_maintenance_group.example <Dn>
```