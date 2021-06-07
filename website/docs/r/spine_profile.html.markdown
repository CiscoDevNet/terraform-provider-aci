---
layout: "aci"
page_title: "ACI: aci_spine_profile"
sidebar_current: "docs-aci-resource-spine_profile"
description: |-
  Manages ACI Spine Profile
---

# aci_spine_profile #
Manages ACI Spine Profile

## Example Usage ##

```hcl

resource "aci_spine_profile" "example" {
  name        = "spine_profile_1"
  description = "from terraform"
  annotation  = "spine profile"
  name_alias  = "check"
}

```


## Argument Reference ##
* `name` - (Required) Name of Object Spine Profile.
* `description` - (Optional) Description for object Spine Profile.
* `annotation` - (Optional) Annotation for object Spine Profile.
* `name_alias` - (Optional) Name alias for object Spine Profile.
* `relation_infra_rs_sp_acc_port_p` - (Optional) Relation to class infraSpAccPortP. Cardinality - N_TO_M. Type - Set of String.
                


## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Spine Profile.

## Importing ##

An existing Spine Profile can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_spine_profile.example <Dn>
```