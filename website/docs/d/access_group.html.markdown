---
layout: "aci"
page_title: "ACI: aci_access_group"
sidebar_current: "docs-aci-data-source-access_group"
description: |-
  Data source for ACI Access Group
---

# aci_access_group

Data source for ACI Access Group

## Example Usage

```hcl

data "aci_access_group" "example" {
  access_port_selector_dn  = aci_access_port_selector.example.id
}

```

## Argument Reference

- `access_port_selector_dn` - (Required) Distinguished name of parent Access Port Selector object.

## Attribute Reference

- `id` - Attribute id set to the Dn of the Access Access Group.
- `annotation` - (Optional) Annotation for object access group.
- `fex_id` - (Optional) Interface policy group FEX ID.
- `tdn` - (Optional) Interface policy group's target dn.
