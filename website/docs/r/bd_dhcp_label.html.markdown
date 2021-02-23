---
layout: "aci"
page_title: "ACI: aci_bd_dhcp_label"
sidebar_current: "docs-aci-resource-bd_dhcp_label"
description: |-
  Manages ACI BD DHCP Label
---

# aci_bd_dhcp_label

Manages ACI BD DHCP Label

## Example Usage

```hcl
resource "aci_bd_dhcp_label" "foo_bd_dhcp_label" {

  bridge_domain_dn  = "${aci_bridge_domain.foo_bridge_domain.id}"
  name  = "example"
  annotation  = "example"
  name_alias  = "example"
  owner  = "tenant"
  tag  = "aqua"
}
```

## Argument Reference

- `bridge_domain_dn` - (Required) Distinguished name of parent Bridge Domain object.
- `name` - (Required) The Bridge Domain DHCP label name. This name can be up to 64 alphanumeric characters.
- `annotation` - (Optional) Annotation for object BD DHCP Label.
- `name_alias` - (Optional) Name alias for object BD DHCP Label.
- `owner` - (Optional) Owner of the target relay servers.  
  Allowed values: "infra", "tenant". Default value: "infra".
- `tag` - (Optional) Label color.  
  Allowed values: "alice-blue", "antique-white", "aqua", "aquamarine", "azure", "beige", "bisque", "black", "blanched-almond", "blue", "blue-violet", "brown", "burlywood", "cadet-blue", "chartreuse", "chocolate", "coral", "cornflower-blue", "cornsilk", "crimson", "cyan", "dark-blue", "dark-cyan", "dark-goldenrod", "dark-gray", "dark-green", "dark-khaki", "dark-magenta", "dark-olive-green", "dark-orange", "dark-orchid", "dark-red", "dark-salmon", "dark-sea-green", "dark-slate-blue", "dark-slate-gray", "dark-turquoise", "dark-violet", "deep-pink", "deep-sky-blue", "dim-gray", "dodger-blue", "fire-brick", "floral-white", "forest-green", "fuchsia", "gainsboro", "ghost-white", "gold", "goldenrod", "gray", "green", "green-yellow", "honeydew", "hot-pink", "indian-red", "indigo", "ivory", "khaki", "lavender", "lavender-blush", "lawn-green", "lemon-chiffon", "light-blue", "light-coral", "light-cyan", "light-goldenrod-yellow", "light-gray", "light-green", "light-pink", "light-salmon", "light-sea-green", "light-sky-blue", "light-slate-gray", "light-steel-blue", "light-yellow", "lime", "lime-green", "linen", "magenta", "maroon", "medium-aquamarine", "medium-blue", "medium-orchid","medium-sea-green", "medium-slate-blue", "medium-spring-green", "medium-turquoise", "medium-violet-red", "midnight-blue", "mint-cream", "misty-rose", "moccasin", "navajo-white", "navy", "old-lace", "olive", "olive-drab", "orange", "orange-red", "orchid", "pale-goldenrod", "pale-green", "pale-turquoise", "pale-violet-red", "papaya-whip", "peachpuf", "peru", "pink", "plum", "powder-blue", "purple", "red", "rosy-brown", "royal-blue", "saddle-brown", "salmon", "sandy-brown", "sea-green", "seashell", "sienna", "silver", "sky-blue", "slate-blue", "slate-gray", "snow", "spring-green", "steel-blue", "tan", "teal", "thistle", "tomato", "turquoise", "violet", "wheat", "white", "white-smoke", "yellow", "yellow-green".

- `relation_dhcp_rs_dhcp_option_pol` - (Optional) Relation to class dhcpOptionPol. Cardinality - N_TO_ONE. Type - String.

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the BD DHCP Label.

## Importing

An existing BD DHCP Label can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_bd_dhcp_label.example <Dn>
```
