---
layout: "aci"
page_title: "ACI: aci_logical_node_profile"
sidebar_current: "docs-aci-resource-logical_node_profile"
description: |-
  Manages ACI Logical Node Profile
---

# aci_logical_node_profile

Manages ACI Logical Node Profile

## Example Usage

```hcl
	resource "aci_logical_node_profile" "foological_node_profile" {
		l3_outside_dn = aci_l3_outside.example.id
		description   = "sample logical node profile"
		name          = "demo_node"
		annotation    = "tag_node"
		name_alias    = "alias_node"
		tag           = "black"
		target_dscp   = "unspecified"
	  }
```

## Argument Reference

- `l3_outside_dn` - (Required) Distinguished name of parent L3-outside object.
- `name` - (Required) Name of Object logical node profile.
- `annotation` - (Optional) Annotation for object logical node profile.
- `description` - (Optional) Description for object logical node profile.
- `config_issues` - (Optional) Bitmask representation of the configuration issues found during the endpoint group deployment.\[Read-only\] expected values are a list of "none", "node-path-misconfig", "routerid-not-changable-with-mcast" or "loopback-ip-missing".
- `name_alias` - (Optional) Name alias for object logical node profile.
- `tag` - (Optional) Specifies the color of a policy label. Allowed values are "black", "navy", "dark-blue", "medium-blue", "blue", "dark-green", "green", "teal", "dark-cyan", "deep-sky-blue", "dark-turquoise", "medium-spring-green", "lime", "spring-green", "aqua", "cyan", "midnight-blue", "dodger-blue", "light-sea-green", "forest-green", "sea-green", "dark-slate-gray", "lime-green", "medium-sea-green", "turquoise", "royal-blue", "steel-blue", "dark-slate-blue", "medium-turquoise", "indigo", "dark-olive-green", "cadet-blue", "cornflower-blue", "medium-aquamarine", "dim-gray", "slate-blue", "olive-drab", "slate-gray", "light-slate-gray", "medium-slate-blue", "lawn-green", "chartreuse", "aquamarine", "maroon", "purple", "olive", "gray", "sky-blue", "light-sky-blue", "blue-violet", "dark-red", "dark-magenta", "saddle-brown", "dark-sea-green", "light-green", "medium-purple", "dark-violet", "pale-green", "dark-orchid", "yellow-green", "sienna", "brown", "dark-gray", "light-blue", "green-yellow", "pale-turquoise", "light-steel-blue", "powder-blue", "fire-brick", "dark-goldenrod", "medium-orchid", "rosy-brown", "dark-khaki", "silver", "medium-violet-red", "indian-red", "peru", "chocolate", "tan", "light-gray", "thistle", "orchid", "goldenrod", "pale-violet-red", "crimson", "gainsboro", "plum", "burlywood", "light-cyan", "lavender", "dark-salmon", "violet", "pale-goldenrod", "light-coral", "khaki", "alice-blue", "honeydew", "azure", "sandy-brown", "wheat", "beige", "white-smoke", "mint-cream", "ghost-white", "salmon", "antique-white", "linen", "light-goldenrod-yellow", "old-lace", "red", "fuchsia", "magenta", "deep-pink", "orange-red", "tomato", "hot-pink", "coral", "dark-orange", "light-salmon", "orange", "light-pink", "pink", "gold", "peachpuff", "navajo-white", "moccasin", "bisque", "misty-rose", "blanched-almond", "papaya-whip", "lavender-blush", "seashell", "cornsilk", "lemon-chiffon", "floral-white", "snow", "yellow", "light-yellow", "ivory" and "white".
- `target_dscp` - (Optional) Node level DSCP value. Allowed values are "CS0", "CS1", "AF11", "AF12", "AF13", "CS2", "AF21", "AF22", "AF23", "CS3", "AF31", "AF32", "AF33", "CS4", "AF41", "AF42", "AF43", "CS5", "VA", "EF", "CS6", "CS7" and "unspecified". Default is "unspecified".

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Logical Node Profile.

## Importing

An existing Logical Node Profile can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_logical_node_profile.example <Dn>
```
