---
layout: "aci"
page_title: "ACI: aci_span_sourcedestination_group_match_label"
sidebar_current: "docs-aci-resource-span_sourcedestination_group_match_label"
description: |-
  Manages ACI SPAN Source-destination Group Match Label
---

# aci_span_sourcedestination_group_match_label #

Manages ACI SPAN Source-destination Group Match Label

## Example Usage ##

```hcl
resource "aci_span_sourcedestination_group_match_label" "example" {
span_source_group_dn  = aci_span_source_group.example.id
name                  = "example"
annotation            = "tag_label"
name_alias            = "alias_label"
tag                   = "label_color"
}
```

## Argument Reference ##

* `span_source_group_dn` - (Required) Distinguished name of parent SPANSourceGroup object.
* `name` - (Required) name of Object span_sourcedestination_group_match_label.
* `annotation` - (Optional)
* `name_alias` - (Optional)
* `tag` - (Optional) label color.
Allowed values: "black", "navy", "dark-blue", "medium-blue", "blue", "dark-green", "green", "teal", "dark-cyan", "deep-sky-blue", "dark-turquoise", "medium-spring-green", "lime", "spring-green", "aqua", "cyan", "midnight-blue", "dodger-blue", "light-sea-green", "forest-green", "sea-green", "dark-slate-gray", "lime-green", "medium-sea-green", "turquoise", "royal-blue", "steel-blue", "dark-slate-blue", "medium-turquoise", "indigo", "dark-olive-green", "cadet-blue", "cornflower-blue", "medium-aquamarine", "dim-gray", "slate-blue", "olive-drab", "slate-gray", "light-slate-gray", "medium-slate-blue", "lawn-green", "chartreuse", "aquamarine", "maroon", "purple", "olive", "gray", "sky-blue", "light-sky-blue", "blue-violet", "dark-red", "dark-magenta", "saddle-brown", "dark-sea-green", "light-green", "medium-purple", "dark-violet", "pale-green", "dark-orchid", "yellow-green", "sienna", "brown", "dark-gray", "light-blue", "green-yellow", "pale-turquoise", "light-steel-blue", "powder-blue", "fire-brick", "dark-goldenrod", "medium-orchid", "rosy-brown", "dark-khaki", "silver", "medium-violet-red", "indian-red", "peru", "chocolate", "tan", "light-gray", "thistle", "orchid", "goldenrod", "pale-violet-red", "crimson", "gainsboro", "plum", "burlywood", "light-cyan", "lavender", "dark-salmon", "violet", "pale-goldenrod", "light-coral", "khaki", "alice-blue", "honeydew", "azure", "sandy-brown", "wheat", "beige", "white-smoke", "mint-cream", "ghost-white", "salmon", "antique-white", "linen", "light-goldenrod-yellow", "old-lace", "red", "fuchsia", "magenta", "deep-pink", "orange-red", "tomato", "hot-pink", "coral", "dark-orange", "light-salmon", "orange", "light-pink", "pink", "gold", "peachpuff", "navajo-white", "moccasin", "bisque", "misty-rose", "blanched-almond", "papaya-whip", "lavender-blush", "seashell", "cornsilk", "lemon-chiffon", "floral-white", "snow", "yellow", "light-yellow", "ivory", "white"

## Attribute Reference ##

The only attribute that this resource exports is the `id`, which is set to the
Dn of the SPAN Source-destination Group Match Label.

## Importing ##

An existing SPAN Source-destination Group Match Label can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import aci_span_sourcedestination_group_match_label.example <Dn>
```
