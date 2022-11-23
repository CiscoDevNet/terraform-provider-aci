---
subcategory: "L3Out"
layout: "aci"
page_title: "ACI: aci_l3_ext_subnet"
sidebar_current: "docs-aci-data-source-l3_ext_subnet"
description: |-
  Data source for ACI l3 extension subnet
---

# aci_l3_ext_subnet #

Data source for ACI l3 extension subnet

## API Information ##

* `Class` - l3extSubnet
* `Distinguished Name` - uni/tn-{tenant}/out-{L3Out}/instP-{external EPG}/extsubnet-[{ip}]

## GUI Information ##

* `Location` - Tenant -> Networking -> L3Outs -> External EPGs -> Subnets -> Route Control Profile

## Example Usage ##

```hcl

data "aci_l3_ext_subnet" "example" {
  external_network_instance_profile_dn  = aci_external_network_instance_profile.example.id
  ip                                    = "10.0.3.28/27"
}

```

## Argument Reference ##

* `external_network_instance_profile_dn` - (Required) Distinguished name of parent External Network Instance Profile object.
* `ip` - (Required) IP address of Object l3 extension subnet.

## Attribute Reference ##

* `id` - Attribute id set to the Dn of the l3 extension subnet.
* `aggregate` - (Optional) Aggregate Routes for l3 extension subnet.
* `annotation` - (Optional) Annotation for object l3 extension subnet.
* `description` - (Optional) Description for object l3 extension subnet.
* `name_alias` - (Optional) Name alias for object l3 extension subnet.
* `scope` - (Optional) The list of domain applicable to the capability.
* `relation_l3ext_rs_subnet_to_profile` - (Optional) Relation to Route Control Profile (class rtctrlProfile). Type: Block.
	* `relation_l3ext_rs_subnet_to_profile.tn_rtctrl_profile_name` - **Deprecated** (Optional) Associates the external EPGs with the route control profiles.
	* `relation_l3ext_rs_subnet_to_profile.tn_rtctrl_profile_dn` - (Optional) Associates the external EPGs with the route control profiles.
	* `relation_l3ext_rs_subnet_to_profile.direction` - (Optional) Relation to configure route map for each BGP peer in the inbound and outbound directions.
* `relation_l3ext_rs_subnet_to_rt_summ` - (Optional) Relation to class rtsumARtSummPol.