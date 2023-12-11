---
subcategory: "Access Policies"
layout: "aci"
page_title: "ACI: aci_epgs_using_function"
sidebar_current: "docs-aci-resource-aci_epgs_using_function"
description: |-
  Manages ACI EPGs Using Function
---

# aci_epgs_using_function #

Manages ACI EPGs Using Function

## Example Usage ##

```hcl
resource "aci_epgs_using_function" "example" {
  access_generic_dn   = aci_access_generic.example.id
  tdn                 = aci_application_epg.epg2.id
  annotation          = "annotation"
  encap               = "vlan-5"
  instr_imedcy        = "lazy"
  mode                = "regular"
  primary_encap       = "vlan-7"
}
```

## Argument Reference ##

* `access_generic_dn` - (Required) Distinguished name of parent AccessGeneric object.
* `tdn` - (Required) tDn of Object EPGs Using Function.
* `encap` - (Required) Vlan number encap. 
* `annotation` - (Optional) annotation for object EPGs Using Function.
* `instr_imedcy` - (Optional) Instrumentation immediacy.
Allowed values: "immediate", "lazy". Default value: "lazy".
* `mode` - (Optional) Bgp domain mode.
Allowed values: "regular", "native", "untagged". Default value: "regular"
* `primary_encap` - (Optional) Primary encap for object EPGs Using Function.

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the EPGs Using Function.

## Importing ##

An existing EPGs Using Function can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_ep_gs_using_function.example <Dn>
```
