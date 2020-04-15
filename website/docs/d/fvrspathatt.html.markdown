---
layout: "aci"
page_title: "ACI: aci_static_path"
sidebar_current: "docs-aci-data-source-static_path"
description: |-
  Data source for ACI Static Path
---

# aci_epg_to_static_path #
Data source for ACI Static Path

## Example Usage ##

```hcl
data "aci_epg_to_static_path" "example" {

  application_epg_dn  = "${aci_application_epg.example.id}"

  tDn  = "example"
}
```
## Argument Reference ##
* `application_epg_dn` - (Required) Distinguished name of parent ApplicationEPG object.
* `tDn` - (Required) tDn of Object static_path.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Static Path.
* `annotation` - (Optional) annotation for object static_path.
* `encap` - (Optional) encapsulation
* `instr_imedcy` - (Optional) immediacy
* `mode` - (Optional) mode of the static association with the path
* `primary_encap` - (Optional) primary_encap for object static_path.
