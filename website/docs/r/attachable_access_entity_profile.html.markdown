---
layout: "aci"
page_title: "ACI: aci_attachable_access_entity_profile"
sidebar_current: "docs-aci-resource-attachable_access_entity_profile"
description: |-
  Manages ACI Attachable Access Entity Profile
---

# aci_attachable_access_entity_profile #
Manages ACI Attachable Access Entity Profile

## Example Usage ##

```hcl
	resource "aci_attachable_access_entity_profile" "fooattachable_access_entity_profile" {
		description = "%s"
		name        = "demo_entity_prof"
		annotation  = "tag_entity"
		name_alias  = "%s"
	}
```
## Argument Reference ##
* `name` - (Required) name of Object attachable_access_entity_profile.
* `annotation` - (Optional) annotation for object attachable_access_entity_profile.
* `name_alias` - (Optional) name_alias for object attachable_access_entity_profile.

* `relation_infra_rs_dom_p` - (Optional) Relation to class infraADomP. Cardinality - N_TO_M. Type - Set of String.
                


## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Attachable Access Entity Profile.

## Importing ##

An existing Attachable Access Entity Profile can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_attachable_access_entity_profile.example <Dn>
```