---
layout: "aci"
page_title: "ACI: aci_leaf_access_port_policy_group"
sidebar_current: "docs-aci-resource-leaf_access_port_policy_group"
description: |-
  Manages ACI Leaf Access Port Policy Group
---

# aci_leaf_access_port_policy_group #
Manages ACI Leaf Access Port Policy Group

## Example Usage ##

```hcl
	resource "aci_leaf_access_port_policy_group" "fooleaf_access_port_policy_group" {
		description = "From Terraform"
		name        = "demo_access_port"
		annotation  = "tag_ports"
		name_alias  = "name_alias"
	} 
```
## Argument Reference ##
* `name` - (Required) Name of Object leaf access port policy group.
* `description` - (Optional) Description for object leaf access port policy group.
* `annotation` - (Optional) Annotation for object leaf access port policy group.
* `name_alias` - (Optional) Name alias for object leaf access port policy group.

* `relation_infra_rs_span_v_src_grp` - (Optional) Relation to class spanVSrcGrp. Cardinality - N_TO_M. Type - Set of String.
                
* `relation_infra_rs_stormctrl_if_pol` - (Optional) Relation to class stormctrlIfPol. Cardinality - N_TO_ONE. Type - String.
                
* `relation_infra_rs_poe_if_pol` - (Optional) Relation to class poeIfPol. Cardinality - N_TO_ONE. Type - String.
                
* `relation_infra_rs_lldp_if_pol` - (Optional) Relation to class lldpIfPol. Cardinality - N_TO_ONE. Type - String.
                
* `relation_infra_rs_macsec_if_pol` - (Optional) Relation to class macsecIfPol. Cardinality - N_TO_ONE. Type - String.
                
* `relation_infra_rs_qos_dpp_if_pol` - (Optional) Relation to class qosDppPol. Cardinality - N_TO_ONE. Type - String.
                
* `relation_infra_rs_h_if_pol` - (Optional) Relation to class fabricHIfPol. Cardinality - N_TO_ONE. Type - String.
                
* `relation_infra_rs_netflow_monitor_pol` - (Optional) Relation to class netflowMonitorPol. Cardinality - N_TO_M. Type - Set of Map.
	- `flt_type` - (Required) Netflow IP filter type. Allowed values: "ce", "ipv4", "ipv6". 
	- `target_dn` - (Required) Distinguished name of target Netflow Monitor object.
                
* `relation_infra_rs_l2_port_auth_pol` - (Optional) Relation to class l2PortAuthPol. Cardinality - N_TO_ONE. Type - String.
                
* `relation_infra_rs_mcp_if_pol` - (Optional) Relation to class mcpIfPol. Cardinality - N_TO_ONE. Type - String.
                
* `relation_infra_rs_l2_port_security_pol` - (Optional) Relation to class l2PortSecurityPol. Cardinality - N_TO_ONE. Type - String.
                
* `relation_infra_rs_copp_if_pol` - (Optional) Relation to class coppIfPol. Cardinality - N_TO_ONE. Type - String.
                
* `relation_infra_rs_span_v_dest_grp` - (Optional) Relation to class spanVDestGrp. Cardinality - N_TO_M. Type - Set of String.
                
* `relation_infra_rs_dwdm_if_pol` - (Optional) Relation to class dwdmIfPol. Cardinality - N_TO_ONE. Type - String.
                
* `relation_infra_rs_qos_pfc_if_pol` - (Optional) Relation to class qosPfcIfPol. Cardinality - N_TO_ONE. Type - String.
                
* `relation_infra_rs_qos_sd_if_pol` - (Optional) Relation to class qosSdIfPol. Cardinality - N_TO_ONE. Type - String.
                
* `relation_infra_rs_mon_if_infra_pol` - (Optional) Relation to class monInfraPol. Cardinality - N_TO_ONE. Type - String.
                
* `relation_infra_rs_fc_if_pol` - (Optional) Relation to class fcIfPol. Cardinality - N_TO_ONE. Type - String.
                
* `relation_infra_rs_qos_ingress_dpp_if_pol` - (Optional) Relation to class qosDppPol. Cardinality - N_TO_ONE. Type - String.
                
* `relation_infra_rs_cdp_if_pol` - (Optional) Relation to class cdpIfPol. Cardinality - N_TO_ONE. Type - String.
                
* `relation_infra_rs_l2_if_pol` - (Optional) Relation to class l2IfPol. Cardinality - N_TO_ONE. Type - String.
                
* `relation_infra_rs_stp_if_pol` - (Optional) Relation to class stpIfPol. Cardinality - N_TO_ONE. Type - String.
                
* `relation_infra_rs_qos_egress_dpp_if_pol` - (Optional) Relation to class qosDppPol. Cardinality - N_TO_ONE. Type - String.
                
* `relation_infra_rs_att_ent_p` - (Optional) Relation to class infraAttEntityP. Cardinality - N_TO_ONE. Type - String.
                
* `relation_infra_rs_l2_inst_pol` - (Optional) Relation to class l2InstPol. Cardinality - N_TO_ONE. Type - String.
                


## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Leaf Access Port Policy Group.

## Importing ##

An existing Leaf Access Port Policy Group can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_leaf_access_port_policy_group.example <Dn>
```