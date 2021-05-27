---
layout: "aci"
page_title: "ACI: aci_l3_outside"
sidebar_current: "docs-aci-resource-l3_outside"
description: |-
  Manages ACI L3 Outside
---

# aci_l3_outside

Manages ACI L3 Outside

## Example Usage

```hcl
	resource "aci_l3_outside" "fool3_outside" {
		tenant_dn      = aci_tenant.dev_tenant.id
		description    = "from terraform"
		name           = "demo_l3out"
		annotation     = "tag_l3out"
		enforce_rtctrl = ["export", "import"]
		name_alias     = "alias_out"
		target_dscp    = "unspecified"
	}
```

## Argument Reference

- `tenant_dn` - (Required) Distinguished name of parent Tenant object.
- `name` - (Required) Name of Object l3 outside.
- `description`- (Optional) Description for object l3 outside.
- `annotation` - (Optional) Annotation for object l3 outside.
- `enforce_rtctrl` - (Optional) Enforce route control type. Allowed values are "import" and "export". Default is "export". (Multiple Comma-Delimited values are allowed. E.g., "import,export").
- `name_alias` - (Optional) Name alias for object l3 outside.
- `target_dscp` - (Optional) The target differentiated services code point (DSCP) of the path attached to the layer 3 outside profile. Allowed values are "CS0", "CS1", "AF11", "AF12", "AF13", "CS2", "AF21", "AF22", "AF23", "CS3", "AF31", "AF32", "AF33", "CS4", "AF41", "AF42", "AF43", "CS5", "VA", "EF", "CS6", "CS7" and "unspecified". Default is "unspecified".

- `relation_l3ext_rs_dampening_pol` - (Optional) Relation to class rtctrlProfile. Cardinality - N_TO_M. Type - Set of Map.
- `relation_l3ext_rs_ectx` - (Optional) Relation to class fvCtx. Cardinality - N_TO_ONE. Type - String.
- `relation_l3ext_rs_out_to_bd_public_subnet_holder` - (Optional) Relation to class fvBDPublicSubnetHolder. Cardinality - N_TO_M. Type - Set of String.
- `relation_l3ext_rs_interleak_pol` - (Optional) Relation to class rtctrlProfile. Cardinality - N_TO_ONE. Type - String.
- `relation_l3ext_rs_l3_dom_att` - (Optional) Relation to class extnwDomP. Cardinality - N_TO_ONE. Type - String.

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the L3 Outside.

## Importing

An existing L3 Outside can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_l3_outside.example <Dn>
```
