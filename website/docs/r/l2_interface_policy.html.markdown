---
layout: "aci"
page_title: "ACI: aci_l2_interface_policy"
sidebar_current: "docs-aci-resource-l2_interface_policy"
description: |-
  Manages ACI L2 Interface Policy
---

# aci_l2_interface_policy

Manages ACI L2 Interface Policy

## Example Usage

```hcl
resource "aci_l2_interface_policy" "example" {
  name        = "demo_l2_pol"
  annotation  = "tag_l2_pol"
  description = "from description"
  name_alias  = "alias_l2_pol"
  qinq        = "disabled"
  vepa        = "disabled"
  vlan_scope  = "global"
}
```

## Argument Reference

- `name` - (Required) Name of Object L2 interface policy.
- `annotation` - (Optional) Annotation for object L2 interface policy.
- `description` - (Optional) Description for object L2 interface policy.
- `name_alias` - (Optional) Name alias for object L2 interface policy.
- `qinq` - (Optional) Determines if QinQ is disabled or if the port should be considered a core or edge port. Allowed values are "disabled", "edgePort", "corePort" and "doubleQtagPort". Default is "disabled".
- `vepa` - (Optional) Determines if Virtual Ethernet Port Aggregator is disabled or enabled. Allowed values are "disabled" and "enabled". Default is "disabled".
- `vlan_scope` - (Optional) The scope of the VLAN. Allowed values are "global" and "portlocal". Default is "global".

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the L2 Interface Policy.

## Importing

An existing L2 Interface Policy can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import aci_l2_interface_policy.example <Dn>
```
