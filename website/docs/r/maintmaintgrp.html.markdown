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


  name  = "example"
  annotation  = "example"
  fwtype  = "example"
  name_alias  = "example"
  pod_maintenance_group_type  = "example"
}
```
## Argument Reference ##
* `name` - (Required) name of Object pod_maintenance_group.
* `annotation` - (Optional) annotation for object pod_maintenance_group.
* `fwtype` - (Optional) fwtype for object pod_maintenance_group.
* `name_alias` - (Optional) name_alias for object pod_maintenance_group.
* `pod_maintenance_group_type` - (Optional) component type

* `relation_maint_rs_mgrpp` - (Optional) Relation to class maintMaintP. Cardinality - N_TO_ONE. Type - String.
                


## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the POD Maintenance Group.

## Importing ##

An existing POD Maintenance Group can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_pod_maintenance_group.example <Dn>
```