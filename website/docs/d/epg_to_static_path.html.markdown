---
layout: "aci"
page_title: "ACI: aci_epg_to_static_path"
sidebar_current: "docs-aci-data-source-epg_to_static_path"
description: |-
  Data source for ACI EPG to Static Path
---

# aci_epg_to_static_path #
Data source for ACI EPG to Static Path

## Example Usage ##

```hcl
data "aci_epg_to_static_path" "example" {
  application_epg_dn  = aci_application_epg.example.id
  tdn  = "example"
}
```
## Argument Reference ##
* `application_epg_dn` - (Required) Distinguished name of parent ApplicationEPG object.
* `tdn` - (Required) tdn of Object Static Path.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Static Path.
* `annotation` - (Optional) Annotation for object Static Path.
* `encap` - (Optional) Encapsulation of the Static Path.
* `instr_imedcy` - (Optional) Immediacy of the Static Path.
* `mode` - (Optional) Mode of the static association with the path.
* `primary_encap` - (Optional) Primary encap for object Static Path.
