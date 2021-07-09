---
layout: "aci"
page_title: "ACI: aci_epgs_using_function"
sidebar_current: "docs-aci-resource-epgs_using_function"
description: |-
  Manages ACI EPGs Using Function
---

# aci_epgs_using_function #

Manages ACI EPGs Using Function

## Example Usage ##

```hcl
resource "aci_epgs_using_function" "example" {
  access_generic_dn = "${aci_access_generic.example.id}"
  tdn               = "${aci_application_epg.epg2.id}"
  annotation        = "example"
  encap             = "vlan-5"
  instr_imedcy      = "example"
  mode              = "example"
  primary_encap     = "example"
}
```

## Argument Reference ##

* `access_generic_dn` - (Required) Distinguished name of parent AccessGeneric object.
* `tdn` - (Required) tDn of Object epgs_using_function.
* `encap` - (Required) vlan number encap.
* `annotation` - (Optional) annotation for object epgs_using_function.
* `instr_imedcy` - (Optional) instrumentation immediacy.
Allowed values: "immediate", "lazy"
* `mode` - (Optional) bgp domain mode.
Allowed values: "regular", "native", "untagged"
* `primary_encap` - (Optional) primary_encap for object epgs_using_function.

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the EPGs Using Function.

## Importing ##

An existing EPGs Using Function can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_ep_gs_using_function.example <Dn>
```
