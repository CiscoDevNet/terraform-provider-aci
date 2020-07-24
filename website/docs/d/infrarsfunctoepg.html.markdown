---
layout: "aci"
page_title: "ACI: aci_epgs_using_function"
sidebar_current: "docs-aci-data-source-epgs_using_function"
description: |-
  Data source for ACI EPGs Using Function
---

# aci_epgs_using_function #
Data source for ACI EPGs Using Function

## Example Usage ##

```hcl

data "aci_epgs_using_function" "example" {
  access_generic_dn   = "${aci_access_generic.example.id}"
  t_dn                = "example"
}

```

## Argument Reference ##
* `access_generic_dn` - (Required) Distinguished name of parent AccessGeneric object.
* `t_dn` - (Required) tDn of Object epgs_using_function.



## Attribute Reference

* `id` - Attribute id set to the Dn of the EPGs Using Function.
* `annotation` - (Optional) annotation for object epgs_using_function.
* `encap` - (Optional) vlan number encap
* `instr_imedcy` - (Optional) instrumentation immediacy
* `mode` - (Optional) bgp domain mode
* `primary_encap` - (Optional) primary_encap for object epgs_using_function.
