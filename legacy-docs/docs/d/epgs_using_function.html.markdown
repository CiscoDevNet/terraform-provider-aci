---
subcategory: "Access Policies"
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
  access_generic_dn   = aci_access_generic.example.id
  tdn                 = aci_application_epg.epg2.id
}

```

## Argument Reference ##
* `access_generic_dn` - (Required) Distinguished name of parent AccessGeneric object.
* `tdn` - (Required) tDn of Object EPGs Using Function.



## Attribute Reference

* `id` - Attribute id set to the Dn of the EPGs Using Function.
* `annotation` - (Optional) Annotation for object EPGs Using Function.
* `encap` - (Optional) Vlan number encap
* `instr_imedcy` - (Optional) Instrumentation immediacy
* `mode` - (Optional) Bgp domain mode
* `primary_encap` - (Optional) Primary encap for object EPGs Using Function.
