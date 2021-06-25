---
layout: "aci"
page_title: "ACI: aci_v_switch_policy_group"
sidebar_current: "docs-aci-resource-v_switch_policy_group"
description: |-
  Manages ACI VSwitch Policy Group
---

# aci_v_switch_policy_group #

Manages ACI VSwitch Policy Group

## API Information ##

* `Class` - vmmVSwitchPolicyCont
* `Distinguished Named` - uni/vmmp-{vendor}/dom-{name}/vswitchpolcont

## GUI Information ##

* `Location` - Virtual Networking -> {vendor} -> {domain_name} -> VSwitch Policy

## Example Usage ##

```hcl
resource "aci_vswitch_policy" "example" {
  vmm_domain_dn  = aci_vmm_domain.example.id
  annotation = "orchestrator:terraform"

  relation_vmm_rs_vswitch_exporter_pol {
    active_flow_time_out = "60"
    idle_flow_time_out = "15"
    sampling_rate = "0"
    target_dn = aci_resource.example.id
  }

  relation_vmm_rs_vswitch_override_cdp_if_pol = aci_resource.example.id

  relation_vmm_rs_vswitch_override_fw_pol = aci_resource.example.id

  relation_vmm_rs_vswitch_override_lacp_pol = aci_resource.example.id

  relation_vmm_rs_vswitch_override_lldp_if_pol = aci_resource.example.id

  relation_vmm_rs_vswitch_override_mcp_if_pol = aci_resource.example.id

  relation_vmm_rs_vswitch_override_mtu_pol = aci_resource.example.id

  relation_vmm_rs_vswitch_override_stp_pol = aci_resource.example.id
}
```

## Argument Reference ##

* `vmm_domain_dn` - (Required) Distinguished name of parent VMM Domain object.
* `annotation` - (Optional) Annotation of object VSwitch Policy Group.
* `description` - (Optional) Description of object VSwitch Policy Group.
* `name_alias` - (Optional) Name Alias of object VSwitch Policy Group.
* `relation_vmm_rs_vswitch_exporter_pol` - (Optional) A block representing the relation to a VMM Netflow Exporter Policy (class netflowVmmExporterPol). Type: Block.
  * `active_flow_time_out` - (Optional) The range of allowed values is "1" to "36001". Default value is "60".
  * `idle_flow_time_out` - (Optional) The range of allowed values is "1" to "6001". Default value is "15".
  * `sampling_rate` - (Optional) The range of allowed values is "1" to "10001". Default value is "0".
  * `target_dn` - (Required) The distinguished name of the target exporter policy. Type: String

* `relation_vmm_rs_vswitch_override_cdp_if_pol` - (Optional) Represents the relation to a CDP Interface Policy (class cdpIfPol). Relationship to policy providing physical configuration of the interfaces. Type: String.

* `relation_vmm_rs_vswitch_override_fw_pol` - (Optional) Represents the relation to a Firewall Policy (class nwsFwPol). Type: String.

* `relation_vmm_rs_vswitch_override_lacp_pol` - (Optional) Represents the relation to a LACP Lag Policy (class lacpLagPol). Type: String.

* `relation_vmm_rs_vswitch_override_lldp_if_pol` - (Optional) Represents the relation to a LLDP Interface Policy (class lldpIfPol). Type: String.

* `relation_vmm_rs_vswitch_override_mcp_if_pol` - (Optional) Represents the relation to a MCP Interface Policy (class mcpIfPol). Relationship to policy providing physical configuration of the interfaces Type: String.

* `relation_vmm_rs_vswitch_override_mtu_pol` - (Optional) Represents the relation to a MTU Policy (class l2InstPol). Relationship to policy providing physical configuration of the interfaces Type: String.

* `relation_vmm_rs_vswitch_override_stp_pol` - (Optional) Represents the relation to a STP Policy (class stpIfPol). Type: String.

## Importing ##

An existing VSwitchPolicyGroup can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import v_switch_policy_group.example <Dn>
```