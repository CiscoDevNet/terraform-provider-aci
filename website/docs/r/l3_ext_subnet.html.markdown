---
layout: "aci"
page_title: "ACI: aci_l3_ext_subnet"
sidebar_current: "docs-aci-resource-l3_ext_subnet"
description: |-
  Manages ACI l3 extension subnet
---

# aci_l3_ext_subnet

Manages ACI l3 extension subnet

## Example Usage

```hcl

	resource "aci_l3_ext_subnet" "foosubnet" {
	  external_network_instance_profile_dn  = aci_external_network_instance_profile.example.id
	  description                           = "Sample L3 External subnet"
	  ip                                    = "10.0.3.28/27"
	  aggregate                             = "shared-rtctrl"
	  annotation                            = "tag_ext_subnet"
	  name_alias                            = "alias_ext_subnet"
	  scope                                 = ["import-security"]
	}

```

## Argument Reference

- `external_network_instance_profile_dn` - (Required) Distinguished name of parent External Network Instance Profile object.
- `ip` - (Required) IP address of Object l3 extension subnet.
- `aggregate` - (Optional) Aggregate Routes for l3 extension subnet. Allowed values are "import-rtctrl", "export-rtctrl", "shared-rtctrl" and "none". Multiple comma-delimited values are allowed. e.g., "export-rtctrl,import-rtctrl". 
- `annotation` - (Optional) Annotation for object l3 extension subnet.
- `description` - (Optional) Description for object l3 extension subnet.
- `name_alias` - (Optional) Name alias for object l3 extension subnet.
- `scope` - (Optional) The list of domain applicable to the capability. Allowed values are "import-rtctrl", "export-rtctrl", "import-security", "shared-security" and "shared-rtctrl". Default is "import-security".

- `relation_l3ext_rs_subnet_to_profile` - (Optional) Relation to class rtctrlProfile. Cardinality - N_TO_M. Type - Set of Map.
- `relation_l3ext_rs_subnet_to_rt_summ` - (Optional) Relation to class rtsumARtSummPol. Cardinality - N_TO_ONE. Type - String.

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the l3 extension subnet.

## Importing

An existing Subnet can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_l3_ext_subnet.example <Dn>
```
