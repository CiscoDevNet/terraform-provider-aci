---
layout: "aci"
page_title: "ACI: aci_vpc_explicit_protection_group"
sidebar_current: "docs-aci-resource-vpc_explicit_protection_group"
description: |-
  Manages ACI VPC Explicit Protection Group
---

# aci_vpc_explicit_protection_group #
Manages ACI VPC Explicit Protection Group

## Example Usage ##

```hcl
resource "aci_vpc_explicit_protection_group" "example" {


  name  = "example"
  annotation  = "example"
  switch1 = "switch1 id"
  switch2 = "switch2 id"
  vpc_domain_policy = "test"
  vpc_explicit_protection_group_id  = "example"
}
```
## Argument Reference ##
* `name` - (Required) name of Object vpc_explicit_protection_group.
* `switch1` - (Required) Id of switch 1 to attach.
* `switch2` - (Required) Id of switch 2 to attach.
* `annotation` - (Optional) annotation for object vpc_explicit_protection_group.
* `vpc_explicit_protection_group_id` - (Optional) explicit protection group ID
* `vpc_domain_policy` - (Optional) VPC domain policy name.              


## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the VPC Explicit Protection Group.

## Importing ##

An existing VPC Explicit Protection Group can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_vpc_explicit_protection_group.example <Dn>
```