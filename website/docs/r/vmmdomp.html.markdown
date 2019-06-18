---
layout: "aci"
page_title: "ACI: aci_vmm_domain"
sidebar_current: "docs-aci-resource-vmm_domain"
description: |-
  Manages ACI VMM Domain
---

# aci_vmm_domain #
Manages ACI VMM Domain

## Example Usage ##

```hcl
	resource "aci_vmm_domain" "foovmm_domain" {
		provider_profile_dn = "${aci_provider_profile.example.id}"
		description         = "%s"
		name                = "demo_domp"
		access_mode         = "read-write"
		annotation          = "tag_dom"
		arp_learning        = "disabled"
		ave_time_out        = "30"
		config_infra_pg     = "no"
		ctrl_knob           = "epDpVerify"
		delimiter           = ";"
		enable_ave          = "no"
		enable_tag          = "no"
		encap_mode          = "unknown"
		enf_pref            = "hw"
		ep_inventory_type   = "on-link"
		ep_ret_time         = "0"
		hv_avail_monitor    = "no"
		mcast_addr          = "224.0.1.0"
		mode                = "default"
		name_alias          = "alias_dom"
		pref_encap_mode     = "unspecified"
	} 
```
## Argument Reference ##
* `provider_profile_dn` - (Required) Distinguished name of parent ProviderProfile object.
* `name` - (Required) name of Object vmm_domain.
* `access_mode` - (Optional) access_mode for object vmm_domain. Allowed values are "read-write" and "read-only". Default is "read-write".
* `annotation` - (Optional) annotation for object vmm_domain.
* `arp_learning` - (Optional) Enable/Disable arp learning for AVS Domain. Allowed values are "enabled" and "disabled". Default value is "disabled".
* `ave_time_out` - (Optional) ave_time_out for object vmm_domain. Allowed value range is "101" - "3001".
* `config_infra_pg` - (Optional) Flag to enable config_infra_pg for object vmm_domain. Allowed values are "yes" and "no". Default is "no".
* `ctrl_knob` - (Optional) Type pf control knob to use. Allowed values are "none" and "epDpVerify".
* `delimiter` - (Optional) delimiter for object vmm_domain.
* `enable_ave` - (Optional) Flag to enable ave for object vmm_domain. Allowed values are "yes" and "no". Default is "no".
* `enable_tag` - (Optional) Flag enable tagging for object vmm_domain. Allowed values are "yes" and "no". Default is "no".
* `encap_mode` - (Optional)  The layer 2 encapsulation protocol to use with the virtual switch. Allowed values are "unknown", "vlan" and "vxlan". Default is "unknown".
* `enf_pref` - (Optional) The switching enforcement preference. This determines whether switches can be done within the virtual switch (Local Switching) or whether all switched traffic must go through the fabric (No Local Switching). Allowed values are "hw", "sw" and "unknown". Default is "hw".
* `ep_inventory_type` - (Optional) Determines which end point inventory_type to use for object vmm_domain. Allowed values are "none" and "on-link". Default is "on-link".
* `ep_ret_time` - (Optional) end point retention time for object vmm_domain. Allowed value range is "1" - "6001". Default value is "0".
* `hv_avail_monitor` - (Optional) Flag to enable hv_avail_monitor for object vmm_domain. Allowed values are "yes" and "no". Default is "no".
* `mcast_addr` - (Optional) The multicast address of the VMM domain profile.
* `mode` - (Optional) The switch to be used for the domain profile. Allowed values are "default", "n1kv", "unknown", "ovs", "k8s", "rhev", "cf" and "openshift". Default is "default".
* `name_alias` - (Optional) name_alias for object vmm_domain.
* `pref_encap_mode` - (Optional) The preferred encapsulation mode for object vmm_domain. Allowed values are "unspecified", "vlan" and "vxlan". Default is "unspecified".

* `relation_vmm_rs_pref_enhanced_lag_pol` - (Optional) Relation to class lacpEnhancedLagPol. Cardinality - N_TO_ONE. Type - String.
                
* `relation_infra_rs_vlan_ns` - (Optional) Relation to class fvnsVlanInstP. Cardinality - N_TO_ONE. Type - String.
                
* `relation_vmm_rs_dom_mcast_addr_ns` - (Optional) Relation to class fvnsMcastAddrInstP. Cardinality - N_TO_ONE. Type - String.
                
* `relation_vmm_rs_default_cdp_if_pol` - (Optional) Relation to class cdpIfPol. Cardinality - N_TO_ONE. Type - String.
                
* `relation_vmm_rs_default_lacp_lag_pol` - (Optional) Relation to class lacpLagPol. Cardinality - N_TO_ONE. Type - String.
                
* `relation_infra_rs_vlan_ns_def` - (Optional) Relation to class fvnsAInstP. Cardinality - N_TO_ONE. Type - String.
                
* `relation_infra_rs_vip_addr_ns` - (Optional) Relation to class fvnsAddrInst. Cardinality - N_TO_ONE. Type - String.
                
* `relation_vmm_rs_default_lldp_if_pol` - (Optional) Relation to class lldpIfPol. Cardinality - N_TO_ONE. Type - String.
                
* `relation_vmm_rs_default_stp_if_pol` - (Optional) Relation to class stpIfPol. Cardinality - N_TO_ONE. Type - String.
                
* `relation_infra_rs_dom_vxlan_ns_def` - (Optional) Relation to class fvnsAInstP. Cardinality - N_TO_ONE. Type - String.
                
* `relation_vmm_rs_default_fw_pol` - (Optional) Relation to class nwsFwPol. Cardinality - N_TO_ONE. Type - String.
                
* `relation_vmm_rs_default_l2_inst_pol` - (Optional) Relation to class l2InstPol. Cardinality - N_TO_ONE. Type - String.
                


## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the VMM Domain.

## Importing ##

An existing VMM Domain can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_vmm_domain.example <Dn>
```