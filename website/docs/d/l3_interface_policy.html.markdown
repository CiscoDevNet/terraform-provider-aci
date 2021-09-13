---
layout: "aci"
page_title: "ACI: aci_l3_interface_policy"
sidebar_current: "docs-aci-data-source-l3_interface_policy"
description: |-
  Data source for ACI L3 Interface Policy
---

# aci_l3_interface_policy #
Data source for ACI L3 Interface Policy

## API Information ##
* `Class` - l3IfPol
* `Distinguished Named` - uni/fabric/l3IfP-{name}

## GUI Information ##
* `Location` - Fabric -> Fabric Policies -> Policies -> Interface -> L3 Interface > Create L3 Interface Policy

## Example Usage ##

```hcl
data "aci_l3_interface_policy" "example" {
  name  = "example"
}
```

## Argument Reference ##
* `name` - (Required) Name of object L3 Interface Policy.

## Attribute Reference
* `id` - Attribute id set to the Dn of the L3 Interface Policy.
* `annotation` - (Optional) Annotation for object L3 Interface Policy.
* `bfd_isis` - (Optional) BFD ISIS Configuration for object L3 Interface Policy.
* `name_alias` - (Optional) Name alias for object L3 Interface Policy.
* `description` - (Optional) Description for object L3 Interface Policy.

