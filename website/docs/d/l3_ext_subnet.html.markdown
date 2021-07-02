---
layout: "aci"
page_title: "ACI: aci_l3_ext_subnet"
sidebar_current: "docs-aci-data-source-l3_ext_subnet"
description: |-
  Data source for ACI l3 extension subnet
---

# aci_l3_ext_subnet

Data source for ACI l3 extension subnet

## Example Usage

```hcl

data "aci_l3_ext_subnet" "example" {
  external_network_instance_profile_dn  = aci_external_network_instance_profile.example.id
  ip                                    = "10.0.3.28/27"
}

```

## Argument Reference

- `external_network_instance_profile_dn` - (Required) Distinguished name of parent External Network Instance Profile object.
- `ip` - (Required) IP address of Object l3 extension subnet.

## Attribute Reference

- `id` - Attribute id set to the Dn of the l3 extension subnet.
- `aggregate` - (Optional) Aggregate Routes for l3 extension subnet.
- `annotation` - (Optional) Annotation for object l3 extension subnet.
- `description` - (Optional) Description for object l3 extension subnet.
- `name_alias` - (Optional) Name alias for object l3 extension subnet.
- `scope` - (Optional) The list of domain applicable to the capability.
