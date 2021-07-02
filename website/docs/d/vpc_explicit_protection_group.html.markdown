---
layout: "aci"
page_title: "ACI: aci_vpc_explicit_protection_group"
sidebar_current: "docs-aci-data-source-vpc_explicit_protection_group"
description: |-
  Data source for ACI VPC Explicit Protection Group
---

# aci_vpc_explicit_protection_group

Data source for ACI VPC Explicit Protection Group

## Example Usage

```hcl
data "aci_vpc_explicit_protection_group" "example" {
  name  = "example"
}
```

## Argument Reference

- `name` - (Required) name of Object vpc_explicit_protection_group.

## Attribute Reference

- `id` - Attribute id set to the Dn of the VPC Explicit Protection Group.
- `annotation` - (Optional) Annotation for object VPC Explicit Protection Group.
- `vpc_explicit_protection_group_id` - (Optional) explicit protection group ID
