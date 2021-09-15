---
layout: "aci"
page_title: "ACI: aci_access_switch_policy_group"
sidebar_current: "docs-aci-resource-access_switch_policy_group"
description: |-
  Manages ACI Access Switch Policy Group
---

# aci_access_switch_policy_group #

Manages ACI Access Switch Policy Group

## API Information ##

* `Class` - infraAccNodePGrp
* `Distinguished Named` - uni/infra/funcprof/accnodepgrp-{name}

## GUI Information ##

* `Location` - Fabric -> Access Policies -> Switches -> Leaf Switches -> Policy Groups -> Create Access Switch Policy Group


## Example Usage ##

```hcl
resource "aci_access_switch_policy_group" "example" {
  name  = "example"
  annotation = "example"
  description = "example"
  name_alias = "example"
}
```

## Argument Reference ##


* `name` - (Required) Name of object Access Switch Policy Group.
* `annotation` - (Optional) Annotation of object Access Switch Policy Group.
* `relation_infra_rs_bfd_ipv4_inst_pol` - (Optional) Represents the relation to a BFD Ipv4 Instance Policy (class bfdIpv4InstPol). Relationship to BFD Ipv4 Instance Policy Type: String.
* `relation_infra_rs_bfd_ipv6_inst_pol` - (Optional) Represents the relation to a BFD Ipv6 Instance Policy (class bfdIpv6InstPol). Relationship to BFD Ipv6 Instance Policy Type: String.
* `relation_infra_rs_bfd_mh_ipv4_inst_pol` - (Optional) Represents the relation to a MH BFD Ipv4 Instance Policy (class bfdMhIpv4InstPol). Relationship to MH BFD Ipv4 Instance Policy Type: String.
* `relation_infra_rs_bfd_mh_ipv6_inst_pol` - (Optional) Represents the relation to a MH BFD Ipv6 Instance Policy (class bfdMhIpv6InstPol). Relationship to MH BFD Ipv6 Instance Policy Type: String.
* `relation_infra_rs_equipment_flash_config_pol` - (Optional) Represents the relation to a Flash Configuration Policy (class equipmentFlashConfigPol). Relation to equipmentFlashConfigInstPol Type: String.
* `relation_infra_rs_fc_fabric_pol` - (Optional) Represents the relation to a Fibre Channel Fabric Level Policy (class fcFabricPol). Relation to fcInstPol Type: String.
* `relation_infra_rs_fc_inst_pol` - (Optional) Represents the relation to a Fibre Channel Instance Policy (class fcInstPol). Relation to fcInstPol Type: String.
* `relation_infra_rs_iacl_leaf_profile` - (Optional) Represents the relation to a CoPP Prefilter Profile for Leafs (class iaclLeafProfile). Relationship the CoPP Prefilter Leaf profile to be applied on leafs Type: String.
* `relation_infra_rs_l2_node_auth_pol` - (Optional) Represents the relation to a Node Authentication (802.1x) policy (class l2NodeAuthPol). Relation to l2NodeAuthPol Type: String.
* `relation_infra_rs_leaf_copp_profile` - (Optional) Represents the relation to a CoPP Profile for Leafs (class coppLeafProfile). Relationship the CoPP profile to be applied on leafs Type: String.
* `relation_infra_rs_leaf_p_grp_to_cdp_if_pol` - (Optional) Represents the relation to a Relation to cdp interface policy for mgmt port (class cdpIfPol). Relationship to cdp interface policy for mgmt port Type: String.
* `relation_infra_rs_leaf_p_grp_to_lldp_if_pol` - (Optional) Represents the relation to a Relation to lldp interface policy for mgmt port (class lldpIfPol). Relationship to lldp interface policy for mgmt port Type: String.
* `relation_infra_rs_mon_node_infra_pol` - (Optional) Represents the relation to a Relation to Access Monitoring Policy (class monInfraPol). A source relation to the monitoring policy model. Type: String.
* `relation_infra_rs_mst_inst_pol` - (Optional) Represents the relation to a MST Instance Policy (class stpInstPol). A source relation to a spanning tree protocol policy. Type: String.
* `relation_infra_rs_netflow_node_pol` - (Optional) Represents the relation to a Netflow Node Policy (class netflowNodePol). Relationship to Netflow Node Policy Type: String.
* `relation_infra_rs_poe_inst_pol` - (Optional) Represents the relation to a POE Node Policy (class poeInstPol). Relationship to POE Node Policy Type: String.
* `relation_infra_rs_topoctrl_fast_link_failover_inst_pol` - (Optional) Represents the relation to a Fast Link Failover Instance Policy (class topoctrlFastLinkFailoverInstPol). Relation to topoctrlFastLinkFailoverPol Type: String.
* `relation_infra_rs_topoctrl_fwd_scale_prof_pol` - (Optional) Represents the relation to a Forwarding Scale Profile Policy (class topoctrlFwdScaleProfilePol). Relation to topoctrlFwdScaleProfilePol Type: String.
* `name_alias` - (Optional) Name alias for object Access Switch Policy Group. 
* `description` - (Optional) Description for object Access Switch Policy Group.


## Importing ##

An existing AccessSwitchPolicyGroup can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_access_switch_policy_group.example <Dn>
```