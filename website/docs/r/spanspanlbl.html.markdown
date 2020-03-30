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

  span_source_group_dn  = "${aci_span_source_group.example.id}"

  name  = "example"
  annotation  = "example"
  name_alias  = "example"
  tag  = "example"
}
```
## Argument Reference ##
* `span_source_group_dn` - (Required) Distinguished name of parent SPANSourceGroup object.
* `name` - (Required) name of Object span_sourcedestination_group_match_label.
* `annotation` - (Optional) 
* `name_alias` - (Optional) 
* `tag` - (Optional) label color



## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the SPAN Source-destination Group Match Label.

## Importing ##

An existing SPAN Source-destination Group Match Label can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_span_sourcedestination_group_match_label.example <Dn>
```