---
subcategory: "Networking"
layout: "aci"
page_title: "ACI: aci_subnet"
sidebar_current: "docs-aci-resource-subnet"
description: |-
  Manages ACI Subnet
---

# aci_subnet

Manages ACI Subnet

## API Information
Class - fvSubnet
- Distinguished Name - uni/tn-{tenant_name}/BD-{bd_name}/subnet-[{subnet_ip}]
- Distinguished Name - uni/tn-{tenant_name}/ap-{ap_name}/epg-{epg_name}/subnet-[{subnet_ip}]

## GUI Information
- Location - Tenant > Networking > Bridge Domains > Subnets
- Location - Tenant > Application Profiles > Application EPGs > Subnets

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

	# Create EP Reachability - Under AP -> EPG Subnet
	resource "aci_subnet" "foo_epg_subnet_next_hop_addr" {
		parent_dn     = aci_application_epg.foo_epg.id
		ip            = "10.0.3.29/32"
		scope         = ["private"]
		description   = "This subject is created by terraform"
		ctrl          = ["no-default-gateway"]
		preferred     = "no"
		virtual       = "yes"
		next_hop_addr = "10.0.3.30"
	}

	# Create Anycast MAC - Under AP -> EPG Subnet
	resource "aci_subnet" "foo_epg_subnet_anycast_mac" {
		parent_dn   = aci_application_epg.foo_epg.id
		ip          = "10.0.3.29/32"
		scope       = ["private"]
		description = "This subject is created by terraform"
		ctrl        = ["no-default-gateway"]
		preferred   = "no"
		virtual     = "yes"
		anycast_mac = "F0:1F:20:34:89:AB"
	}

	# Create MSNLB in IGMP mode - Under AP -> EPG Subnet
	resource "aci_subnet" "foo_epg_subnet_msnlb_mcast_igmp" {
		parent_dn   = aci_application_epg.foo_epg.id
		ip          = "10.0.3.29/32"
		scope       = ["private"]
		description = "This subject is created by terraform"
		ctrl        = ["no-default-gateway"]
		preferred   = "no"
		virtual     = "yes"
		msnlb = {
			mode  = "mode-mcast-igmp"
			group = "224.0.0.1"
			mac   = "00:00:00:00:00:00"
		}
	}

	# Create MSNLB in static multicast mode - Under AP -> EPG Subnet
	resource "aci_subnet" "foo_epg_subnet_msnlb_mcast_static" {
		parent_dn   = aci_application_epg.foo_epg.id
		ip          = "10.0.3.29/32"
		scope       = ["private"]
		description = "This subject is created by terraform"
		ctrl        = ["no-default-gateway"]
		preferred   = "no"
		virtual     = "yes"
		msnlb = {
			mode  = "mode-mcast--static"
			group = ""
			mac   = "03:1F:20:34:89:AA"
		}
	}

	# Create MSNLB in unicast mode - Under AP -> EPG Subnet
	resource "aci_subnet" "foo_epg_subnet_msnlb_mode_uc" {
		parent_dn   = aci_application_epg.foo_epg.id
		ip          = "10.0.3.29/32"
		scope       = ["private"]
		description = "This subject is created by terraform"
		ctrl        = ["no-default-gateway"]
		preferred   = "no"
		virtual     = "yes"
		msnlb = {
			mode  = "mode-uc"
			group = "0.0.0.0"
			mac   = "00:1F:20:34:89:AA"
		}
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
- `next_hop_addr` - (Optional) EP Reachability of the Application EPGs Subnet object. Type - String.
- `msnlb` - (Optional) A block representing MSNLB of the Application EPGs Subnet object. Type: Block.
   - `mode` - Mode of the MSNLB object, Allowed values are "mode-mcast--static", "mode-uc" and "mode-mcast-igmp". Default is "mode-uc".
   - `group` - The IGMP mode group IP address of the MSNLB object, must be a valid multicast IP address.
   - `mac` - MAC address of the unicast and static multicast mode of the MSNLB object. The valid static multicast MAC address format is `03:XX:XX:XX:XX:XX`.
- `anycast_mac` - Anycast MAC of the Application EPGs Subnet object. Type - String.
## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Subnet.

## Importing

An existing Subnet can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_subnet.example "<Dn>"
```
