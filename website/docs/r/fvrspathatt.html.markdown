---
layout: "aci"
page_title: "ACI: aci_static_path"
sidebar_current: "docs-aci-resource-static_path"
description: |-
  Manages ACI Static Path
---

# aci_epg_to_static_path #
Manages ACI Static Path

## Example Usage ##

```hcl
resource "aci_epg_to_static_path" "example" {

  application_epg_dn  = "${aci_application_epg.example.id}"

  tDn  = "example"
  annotation  = "example"
  encap  = "example"
  instr_imedcy  = "example"
  mode  = "example"
  primary_encap  = "example"
}
```
## Argument Reference ##
* `application_epg_dn` - (Required) Distinguished name of parent ApplicationEPG object.
* `tDn` - (Required) tDn of Object static_path.
* `annotation` - (Optional) annotation for object static_path.
* `encap` - (Optional) encapsulation
* `instr_imedcy` - (Optional) immediacy
* `mode` - (Optional) mode of the static association with the path
* `primary_encap` - (Optional) primary_encap for object static_path.



## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Static Path.

## Importing ##

An existing Static Path can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_static_path.example <Dn>
```