---
layout: "aci"
page_title: "ACI: aci_access_port_selector"
sidebar_current: "docs-aci-resource-access_port_selector"
description: |-
  Manages ACI Access Port Selector
---

# aci_access_port_selector #
Manages ACI Access Port Selector

## Example Usage ##

```hcl
	resource "aci_access_port_selector" "fooaccess_port_selector" {
		leaf_interface_profile_dn = "${aci_leaf_interface_profile.example.id}"
		description               = "%s"
		name                      = "demo_port_selector"
		access_port_selector_type = "%s"
		annotation                = "tag_port_selector"
		name_alias                = "alias_port_selector"
	} 
```
## Argument Reference ##
* `leaf_interface_profile_dn` - (Required) Distinguished name of parent LeafInterfaceProfile object.
* `name` - (Required) name of Object access_port_selector.
* `access_port_selector_type` - (Required) The host port selector type.Allowed values are "ALL" and "range". Default is "ALL".
* `annotation` - (Optional) annotation for object access_port_selector.
* `name_alias` - (Optional) name_alias for object access_port_selector.
* `access_port_selector_type` - (Optional) host port selector type

* `relation_infra_rs_acc_base_grp` - (Optional) Relation to class infraAccBaseGrp. Cardinality - N_TO_ONE. Type - String.
                


## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Access Port Selector.

## Importing ##

An existing Access Port Selector can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_access_port_selector.example <Dn>
```