---
layout: "aci"
page_title: "ACI: aci_subnet"
sidebar_current: "docs-aci-resource-subnet"
description: |-
  Manages ACI Subnet
---

# aci_subnet

Manages ACI Subnet

## Example Usage

```hcl

	resource "aci_subnet" "foosubnet" {
		parent_dn 		 = aci_bridge_domain.bd_for_subnet.id
		description      = "subnet"
		ip               = "10.0.3.28/27"
		annotation       = "tag_subnet"
		ctrl             = ["querier", "nd"]
		name_alias       = "alias_subnet"
		preferred        = "no"
		scope            = ["private", "shared"]
		virtual          = "yes"
	}

```

## Argument Reference

- `parent_dn` - (Required) Distinguished name of parent object.
- `ip` - (Required) The IP address and mask of the default gateway.
- `annotation` - (Optional) Annotation for object subnet.
- `description` - (Optional) Description for object subnet.
- `ctrl` - (Optional) The list of subnet control state. The control can be specific protocols applied to the subnet such as IGMP Snooping. Allowed values are "unspecified", "querier", "nd" and "no-default-gateway". Default is "nd". NOTE: "unspecified" should't be used along with other values.
- `name_alias` - (Optional) Name alias for object subnet.
- `preferred` - (Optional) Indicates if the subnet is preferred (primary) over the available alternatives. Only one preferred subnet is allowed. Allowed values are "yes" and "no". Default is "no".
- `scope` - (Optional) The List of network visibility of the subnet. Allowed values are "private", "public" and "shared". Default is "private".
- `virtual` - (Optional) Treated as virtual IP address. Used in case of BD extended to multiple sites. Allowed values are "yes" and "no". Default is "no".

- `relation_fv_rs_bd_subnet_to_out` - (Optional) Relation to class l3extOut. Cardinality - N_TO_M. Type - Set of String.
- `relation_fv_rs_nd_pfx_pol` - (Optional) Relation to class ndPfxPol. Cardinality - N_TO_ONE. Type - String.
- `relation_fv_rs_bd_subnet_to_profile` - (Optional) Relation to class rtctrlProfile. Cardinality - N_TO_ONE. Type - String.

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Subnet.

## Importing

An existing Subnet can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_subnet.example <Dn>
```
