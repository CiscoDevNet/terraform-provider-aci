---
layout: "aci"
page_title: "ACI: aci_v_switch_policy_group"
sidebar_current: "docs-aci-resource-v_switch_policy_group"
description: |-
  Manages ACI VSwitch Policy Group
---

# vSwitch_Policy #

Manages ACI VSwitch Policy Group

## API Information ##

* `Class` - vmmVSwitchPolicyCont
* `Distinguished Named` - uni/vmmp-{vendor}/dom-{name}/vswitchpolcont

## GUI Information ##

* `Location` - Virtual Networking -> VMM Domain -> VSwitchPolicy


## Example Usage ##

```hcl
resource "aci_v_switch_policy_group" "example" {
  vmm_domain_dn  = aci_vmm_domain.example.id
  annotation = "orchestrator:terraform"

  vmm_rs_vswitch_exporter_pol {
    active_flow_time_out = "60"
    idle_flow_time_out = "15"
    sampling_rate = "0"
    target_dn = aci_resource.example.id
  }

  vmm_rs_vswitch_override_cdp_if_pol = aci_resource.example.id

  vmm_rs_vswitch_override_fw_pol = aci_resource.example.id

  vmm_rs_vswitch_override_lacp_pol = aci_resource.example.id

  vmm_rs_vswitch_override_lldp_if_pol = aci_resource.example.id

  vmm_rs_vswitch_override_mcp_if_pol = aci_resource.example.id

  vmm_rs_vswitch_override_mtu_pol = aci_resource.example.id

  vmm_rs_vswitch_override_stp_pol = aci_resource.example.id
}
```

## Argument Reference ##

* `vmm_domain_dn` - (Required) Distinguished name of parent VMMDomain object.

* `annotation` - (Optional) Annotation of object VSwitch Policy Group.


* `relation_vmm_rs_vswitch_exporter_pol` - (Optional) A block representing the relation to a Relation to VMM Netflow Exporter Policy (class netflowVmmExporterPol). Relationship to VMM Netflow Exporter Policy Type: Block.
  * `active_flow_time_out` - (Optional) activeFlowTimeOut.  default value is "60".
  * `idle_flow_time_out` - (Optional) idleFlowTimeOut.  default value is "15".
  * `sampling_rate` - (Optional) samplingRate.  default value is "0".
  * `target_dn` - (Required) The distinguished name of the target. Type: String


* `relation_vmm_rs_vswitch_override_cdp_if_pol` - (Optional) Represents the relation to a Relation to CDP Interface Policy (class cdpIfPol). Relationship to policy providing physical configuration of the interfaces Type: String.


* `relation_vmm_rs_vswitch_override_fw_pol` - (Optional) Represents the relation to a Relation to Firewall Policy (class nwsFwPol). Relationship to nws fw policy Type: String.


* `relation_vmm_rs_vswitch_override_lacp_pol` - (Optional) Represents the relation to a Relation to LACP Lag Policy (class lacpLagPol). Relationship to lacp pc policy Type: String.


* `relation_vmm_rs_vswitch_override_lldp_if_pol` - (Optional) Represents the relation to a Relation to LLDP Interface Policy (class lldpIfPol).  Type: String.


* `relation_vmm_rs_vswitch_override_mcp_if_pol` - (Optional) Represents the relation to a Relation to MCP Interface Policy (class mcpIfPol). Relationship to policy providing physical configuration of the interfaces Type: String.


* `relation_vmm_rs_vswitch_override_mtu_pol` - (Optional) Represents the relation to a Relation to MTU Policy (class l2InstPol). Relationship to policy providing physical configuration of the interfaces Type: String.


* `relation_vmm_rs_vswitch_override_stp_pol` - (Optional) Represents the relation to a Relation to STP Policy (class stpIfPol). Relationship to stp policy Type: String.



## Importing ##

An existing VSwitchPolicyGroup can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import v_switch_policy_group.example <Dn>
```