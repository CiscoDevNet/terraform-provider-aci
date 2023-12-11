---
subcategory: "L3Out"
layout: "aci"
page_title: "ACI: aci_l3_ext_subnet"
sidebar_current: "docs-aci-data-source-aci_l3_ext_subnet"
description: |-
  Data source for ACI External EPG Subnet
---

# aci_l3_ext_subnet #

Data source for ACI External EPG Subnet

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

* `external_network_instance_profile_dn` - (Required) Distinguished name of the parent External Network Instance Profile object.
* `ip` - (Required) IP address of the External EPG Subnet object.

## Attribute Reference ##

* `id` - Attribute id set to the Dn of the External EPG Subnet object.
* `aggregate` - (Optional) Aggregate Routes of the External EPG Subnet object.
* `annotation` - (Optional) Annotation of the External EPG Subnet object.
* `description` - (Optional) Description of the External EPG Subnet object.
* `name_alias` - (Optional) Name alias of the External EPG Subnet object.
* `scope` - (Optional) The list of domain applicable to the capability.
* `relation_l3ext_rs_subnet_to_profile` - (Optional) Relation to Route Control Profile (class rtctrlProfile). Type: Block.
	* `tn_rtctrl_profile_name` - **Deprecated** (Optional) Associates the External EPGs with the Route Control Profiles.
	* `tn_rtctrl_profile_dn` - (Optional) Associates the External EPGs with the Route Control Profiles.
	* `direction` - (Optional) Relation to configure route map for each BGP peer in the inbound and outbound directions.
* `relation_l3ext_rs_subnet_to_rt_summ` - (Optional) Relation to a Route Summarization Policy (class rtsumARtSummPol).