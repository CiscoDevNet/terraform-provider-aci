---
layout: "aci"
page_title: "ACI: aci_spine_switch_policy_group"
sidebar_current: "docs-aci-resource-spine_switch_policy_group"
description: |-
  Manages ACI Spine Switch Policy Group
---

# aci_spine_switch_policy_group #
Manages ACI Spine Switch Policy Group

## API Information ##
* `Class` - infraSpineAccNodePGrp
* `Distinguished Named` - uni/infra/funcprof/spaccnodepgrp-{name}

## GUI Information ##
* `Location` - Fabric -> Access Policies -> Switches -> Spine Switches -> Policy Groups -> Create Spine Switch Policy Group


## Example Usage ##
```hcl
resource "aci_spine_switch_policy_group" "example" {
  name  = "example"
  annotation = "orchestrator:terraform"
}
```

## Argument Reference ##
* `name` - (Required) Name of object Spine Switch Policy Group.
* `annotation` - (Optional) Annotation of object Spine Switch Policy Group.
* `relation_infra_rs_iacl_spine_profile` - (Optional) Represents the relation to a CoPP Prefilter Profile for Spines (class iaclSpineProfile). Relationship the CoPP Prefilter Spine profile to be applied on spines Type: String.
* `relation_infra_rs_spine_bfd_ipv4_inst_pol` - (Optional) Represents the relation to a BFD Ipv4 Instance Policy (class bfdIpv4InstPol). Relationship to BFD Ipv4 Instance Policy Type: String.
* `relation_infra_rs_spine_bfd_ipv6_inst_pol` - (Optional) Represents the relation to a BFD Ipv6 Instance Policy (class bfdIpv6InstPol). Relationship to BFD Ipv6 Instance Policy Type: String.
* `relation_infra_rs_spine_copp_profile` - (Optional) Represents the relation to a CoPP Profile for Spines (class coppSpineProfile). Relationship the CoPP profile to be applied on spines Type: String.
* `relation_infra_rs_spine_p_grp_to_cdp_if_pol` - (Optional) Represents the relation to a Relation to cdp interface policy for mgmt port (class cdpIfPol). Relationship to cdp interface policy for mgmt port Type: String.
* `relation_infra_rs_spine_p_grp_to_lldp_if_pol` - (Optional) Represents the relation to a Relation to lldp interface policy for mgmt port (class lldpIfPol). Relationship to lldp interface policy for mgmt port Type: String.
* `name_alias` - (Optional) Name alias for object Spine Switch Policy Group. 
* `description` - (Optional) Description for object Spine Switch Policy Group.


## Importing ##

An existing SpineSwitchPolicyGroup can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import spine_switch_policy_group.example <Dn>
```