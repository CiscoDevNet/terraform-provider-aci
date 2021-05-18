---
layout: "aci"
page_title: "ACI: aci_l2_outside"
sidebar_current: "docs-aci-resource-l2_outside"
description: |-
  Manages ACI L2 Outside
---

# aci_l2_outside

Manages ACI L2 Outside

## Example Usage

```hcl
resource "aci_l2_outside" "example" {

  tenant_dn  = "${aci_tenant.example.id}"
  name  = "example"
  annotation  = "example"
  name_alias  = "example"
  target_dscp = "AF11"

}
```

## Argument Reference

- `tenant_dn` - (Required) Distinguished name of parent tenant object.
- `name` - (Required) Name of object l2 outside.
- `annotation` - (Optional) Annotation for object l2 outside.

- `name_alias` - (Optional) Name alias for object l2 outside.

- `target_dscp` - (Optional) Target dscp.  
  Allowed values: "AF11", "AF12", "AF13", "AF21", "AF22", "AF23", "AF31", "AF32", "AF33", "AF41", "AF42", "AF43", "CS0", "CS1", "CS2", "CS3", "CS4", "CS5", "CS6", "CS7", "EF", "VA", "unspecified". Default value: "unspecified".

- `relation_l2ext_rs_e_bd` - (Optional) Relation to class fvBD. Cardinality - N_TO_ONE. Type - String.
- `relation_l2ext_rs_l2_dom_att` - (Optional) Relation to class l2extDomP. Cardinality - N_TO_ONE. Type - String.

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the L2 Outside.

## Importing

An existing L2 Outside can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_l2_outside.example <Dn>
```
