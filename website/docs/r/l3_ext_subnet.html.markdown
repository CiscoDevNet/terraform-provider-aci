---
subcategory: "L3Out"
layout: "aci"
page_title: "ACI: aci_l3_ext_subnet"
sidebar_current: "docs-aci-resource-l3_ext_subnet"
description: |-
  Manages ACI l3 extension subnet
---

# aci_l3_ext_subnet #

Manages ACI l3 extension subnet

## API Information ##

* `Class` - l3extSubnet
* `Distinguished Name` - uni/tn-{tenant}/out-{L3Out}/instP-{external EPG}/extsubnet-[{ip}]

## GUI Information ##

* `Location` - Tenant -> Networking -> L3Outs -> External EPGs -> Subnets -> Route Control Profile

## Example Usage ##

```hcl

	resource "aci_l3_ext_subnet" "foosubnet" {
	  external_network_instance_profile_dn  = aci_external_network_instance_profile.example.id
	  description                           = "Sample L3 External subnet"
	  ip                                    = "10.0.3.28/27"
	  aggregate                             = "shared-rtctrl"
	  annotation                            = "tag_ext_subnet"
	  name_alias                            = "alias_ext_subnet"
	  scope                                 = ["import-rtctrl", "export-rtctrl", "import-security"]
	  relation_l3ext_rs_subnet_to_profile {
		tn_rtctrl_profile_dn  = aci_route_control_profile.bgp_route_control_profile.id
		direction = "import"
	  }
	  relation_l3ext_rs_subnet_to_profile {
		tn_rtctrl_profile_dn  = aci_route_control_profile.bgp_route_control_profile_2.id
		direction = "export"
	  }
	}

```

## Argument Reference ##

* `external_network_instance_profile_dn` - (Required) Distinguished name of the parent External Network Instance Profile object.
* `ip` - (Required) ip of the L3 Extension Subnet object.
* `aggregate` - (Optional) Aggregate Routes of the L3 Extension Subnet object. Allowed values are "import-rtctrl", "export-rtctrl", "shared-rtctrl" and "none".
* `annotation` - (Optional) annotation of the L3 Extension Subnet object.
* `description` - (Optional) Description of the L3 Extension Subnet object.
* `name_alias` - (Optional) name_alias of the L3 Extension Subnet object.
* `scope` - (Optional) The list of domain applicable to the capability. Allowed values are "import-rtctrl", "export-rtctrl", "import-security", "shared-security" and "shared-rtctrl". Default is "import-security".
* `relation_l3ext_rs_subnet_to_profile` - (Optional) Relation to Route Control Profile (class rtctrlProfile). Type: Block.
	* `tn_rtctrl_profile_name` - **Deprecated** (Optional) Associates the External EPGs with the Route Control Profiles.
	* `tn_rtctrl_profile_dn` - (Optional) Associates the External EPGs with the Route Control Profiles.
	* `direction` - (Required) Relation to configure route map for each BGP peer in the inbound and outbound directions.
* `relation_l3ext_rs_subnet_to_rt_summ` - (Optional) Relation to class rtsumARtSummPol.

## Importing ##

An existing L3 Extension Subnet can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_l3_ext_subnet.example <Dn>
```
