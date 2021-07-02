---
layout: "aci"
page_title: "ACI: aci_vpc_explicit_protection_group"
sidebar_current: "docs-aci-resource-vpc_explicit_protection_group"
description: |-
  Manages ACI VPC Explicit Protection Group
---

# aci_vpc_explicit_protection_group

Manages ACI VPC Explicit Protection Group

## Example Usage

```hcl
resource "aci_vpc_explicit_protection_group" "example" {
  name                              = "example"
  annotation                        = "tag_vpc"
  switch1                           = "switch1_id"
  switch2                           = "switch2_id"
  vpc_domain_policy                 = "test"
  vpc_explicit_protection_group_id  = "1"
}
```

## Argument Reference

- `name` - (Required) Name of Object VPC Explicit Protection Group.
- `switch1` - (Required) Node Id of switch 1 to attach.
- `switch2` - (Required) Node Id of switch 2 to attach.
- `annotation` - (Optional) Annotation for object VPC Explicit Protection Group.
- `vpc_domain_policy` - (Optional) VPC domain policy name.
- `vpc_explicit_protection_group_id` - (Optional) Explicit protection group ID. Integer values are allowed between 1-1000. default value is "0".

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the VPC Explicit Protection Group.

## Importing

An existing VPC Explicit Protection Group can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import aci_vpc_explicit_protection_group.example <Dn>
```
