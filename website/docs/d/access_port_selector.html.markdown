---
subcategory: "Access Policies"
layout: "aci"
page_title: "ACI: aci_access_port_selector"
sidebar_current: "docs-aci-data-source-access_port_selector"
description: |-
  Data source for ACI Access Port Selector
---

# aci_access_port_selector

Data source for ACI Access Port Selector

## Example Usage

```hcl
data "aci_access_port_selector" "dev_acc_port_select" {
  leaf_interface_profile_dn  = aci_leaf_interface_profile.example.id
  name                       = "foo_acc_port_select"
  access_port_selector_type  = "ALL"
}
```

## Argument Reference

- `leaf_interface_profile_dn` - (Required) Distinguished name of parent Leaf Interface Profile object.
- `name` - (Required) Name of Object Access Port Selector.
- `access_port_selector_type` - (Required) The host port selector type. Allowed values are "ALL" and "range". Default is "ALL".

## Attribute Reference

- `id` - Attribute id set to the Dn of the Access Port Selector.
- `annotation` - (Optional) Annotation for object Access Port Selector.
- `description` - (Optional) Description for object Access Port Selector.
- `name_alias` - (Optional) Name alias for object Access Port Selector.
