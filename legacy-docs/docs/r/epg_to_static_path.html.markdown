---
subcategory: "Application Management"
layout: "aci"
page_title: "ACI: aci_epg_to_static_path"
sidebar_current: "docs-aci-resource-aci_epg_to_static_path"
description: |-
  Manages ACI EPG to Static Path
---

# aci_epg_to_static_path #

Manages ACI EPG to Static Path

## Example Usage ##

```hcl
resource "aci_epg_to_static_path" "example" {
  application_epg_dn  = aci_application_epg.example.id
  tdn  = "topology/pod-1/paths-129/pathep-[eth1/3]"
  annotation = "annotation"
  encap  = "vlan-1000"
  instr_imedcy = "lazy"
  mode  = "regular"
  primary_encap ="vlan-500"
}
```

## Argument Reference ##

* `application_epg_dn` - (Required) Distinguished name of parent ApplicationEPG object.
* `encap` - (Required) Encapsulation of the Static Path.
* `tdn` - (Required) tdn of Object Static Path.
* `annotation` - (Optional) Annotation for object Static Path.
* `instr_imedcy` - (Optional) Immediacy of the Static Path.
Allowed values: "immediate", "lazy". Default value: "lazy"
* `mode` - (Optional) Mode of the static association with the path.
Allowed values: "regular", "native", "untagged". Default value: "regular"
* `primary_encap` - (Optional) Primary encap for object Static Path.

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Static Path.

## Importing ##

An existing Static Path can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_epg_to_static_path.example <Dn>
```
