---
layout: "aci"
page_title: "ACI: aci_span_sourcedestination_group_match_label"
sidebar_current: "docs-aci-data-source-span_sourcedestination_group_match_label"
description: |-
  Data source for ACI SPAN Source-destination Group Match Label
---

# aci_span_sourcedestination_group_match_label #
Data source for ACI SPAN Source-destination Group Match Label

## Example Usage ##

```hcl
data "aci_span_sourcedestination_group_match_label" "example" {
  span_source_group_dn  = aci_span_source_group.example.id
  name  = "example"
}
```
## Argument Reference ##
* `span_source_group_dn` - (Required) Distinguished name of parent SPANSourceGroup object.
* `name` - (Required) Name of Object SPANSourceGroup object.



## Attribute Reference

* `id` - Attribute id set to the Dn of the SPAN Source-destination Group Match Label.
* `annotation` - (Optional) Annotation for object SPANSourceGroup object.
* `description` - (Optional) Description for object SPANSourceGroup object.
* `name_alias` - (Optional) Name Alias for object SPANSourceGroup object.
* `tag` - (Optional) label color
