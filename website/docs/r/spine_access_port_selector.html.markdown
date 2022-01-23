---
subcategory: "Access Policies"
layout: "aci"
page_title: "ACI: aci_spine_access_port_selector"
sidebar_current: "docs-aci-resource-spine_access_port_selector"
description: |-
  Manages ACI Spine Access Port Selector
---

# aci_spine_access_port_selector #

Manages ACI Spine Access Port Selector

## API Information ##

* `Class` - infraSHPortS
* `Distinguished Named` - uni/infra/spaccportprof-{name}/shports-{name}-typ-{type}

## GUI Information ##

* `Location` - Fabric > Access Policies > Interfaces > Spine Interfaces > Profiles > {interface_profile}:{interface_selector}


## Example Usage ##

```hcl
resource "aci_spine_access_port_selector" "example" {
  spine_interface_profile_dn  = aci_spine_interface_profile.example.id
  name  = "example"
  spine_access_port_selector_type  = "ALL"
  annotation = "orchestrator:terraform"
  name_alias = "alias example"

  relation_infra_rs_sp_acc_grp = aci_resource.example.id
}
```

## Argument Reference ##

* `spine_interface_profile_dn` - (Required) Distinguished name of the parent Spine Interface Profile.
* `name` - (Required) Name of the Spine Access Port Selector.
* `spine_access_port_selector_type` - (Required) The type of Spine Access Port Selector. Allowed values are "ALL" and "range". Default is "ALL".
* `annotation` - (Optional) Annotation of the Spine Access Port Selector.
* `name_alias` - (Optional) Name Alias of the Spine Access Port Selector.

* `relation_infra_rs_sp_acc_grp` - (Optional) Represents the relation to a Spine Access Group (class infraSpAccGrp). Type: String.


## Importing ##

An existing Spine Access Port Selector can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_spine_access_port_selector.example <Dn>
```