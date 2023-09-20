---
subcategory: "Virtual Networking"
layout: "aci"
page_title: "ACI: aci_vswitch_policy"
sidebar_current: "docs-aci-resource-vswitch_policy"
description: |-
  Manages ACI VSwitch Policy
---

# aci_vswitch_policy

Manages ACI VSwitch Policy

## API Information

- `Class` - vmmVSwitchPolicyCont
- `Distinguished Name` - uni/vmmp-{vendor}/dom-{name}/vswitchpolcont

## GUI Information

- `Location` - Virtual Networking -> {vendor} -> {domain_name} -> VSwitch Policy

## Example Usage

```hcl
resource "aci_vswitch_policy" "example" {
  vmm_domain_dn  = aci_vmm_domain.example.id
  annotation = "vswitch_policy_tag"
  description = "from terraform"
  name_alias = "vswitch_policy_alias"
  relation_vmm_rs_vswitch_exporter_pol {
    active_flow_time_out = "60"
    idle_flow_time_out = "15"
    sampling_rate = "0"
    target_dn = "uni/infra/vmmexporterpol-exporter_policy"
  }

  relation_vmm_rs_vswitch_override_cdp_if_pol = aci_cdp_interface_policy.example.id

  relation_vmm_rs_vswitch_override_fw_pol = aci_resource.example.id

  relation_vmm_rs_vswitch_override_lacp_pol = aci_lacp_policy.example.id

  relation_vmm_rs_vswitch_override_lldp_if_pol = aci_lldp_interface_policy.example.id

  relation_vmm_rs_vswitch_override_mcp_if_pol = aci_resource.example.id

  relation_vmm_rs_vswitch_override_mtu_pol = aci_resource.example.id

  relation_vmm_rs_vswitch_override_stp_pol = aci_spanning_tree_interface_policy.example.id
}
```

## Argument Reference

- `vmm_domain_dn` - (Required) Distinguished name of parent VMM Domain object.
- `annotation` - (Optional) Annotation of object VSwitch Policy Group.
- `description` - (Optional) Description of object VSwitch Policy Group.
- `name_alias` - (Optional) Name Alias of object VSwitch Policy Group.
- `relation_vmm_rs_vswitch_exporter_pol` - (Optional) A block representing the relation to a VMM Netflow Exporter Policy (class netflowVmmExporterPol). Type: Block.

  - `active_flow_time_out` - (Optional) The range of allowed values is "0" to "3600". Default value is "60".
  - `idle_flow_time_out` - (Optional) The range of allowed values is "0" to "600". Default value is "15".
  - `sampling_rate` - (Optional) The range of allowed values is "0" to "1000". Default value is "0".
  - `target_dn` - (Required) The distinguished name of the target exporter policy. Type: String

- `relation_vmm_rs_vswitch_override_cdp_if_pol` - (Optional) Represents the relation to a CDP Interface Policy (class cdpIfPol). Relationship to policy providing physical configuration of the interfaces. Type: String.

- `relation_vmm_rs_vswitch_override_fw_pol` - (Optional) Represents the relation to a Firewall Policy (class nwsFwPol). Type: String.

- `relation_vmm_rs_vswitch_override_lacp_pol` - (Optional) Represents the relation to a LACP Lag Policy (class lacpLagPol). Type: String.

- `relation_vmm_rs_vswitch_override_lldp_if_pol` - (Optional) Represents the relation to a LLDP Interface Policy (class lldpIfPol). Type: String.

- `relation_vmm_rs_vswitch_override_mcp_if_pol` - (Optional) Represents the relation to a MCP Interface Policy (class mcpIfPol). Relationship to policy providing physical configuration of the interfaces Type: String.

- `relation_vmm_rs_vswitch_override_mtu_pol` - (Optional) Represents the relation to a MTU Policy (class l2InstPol). Relationship to policy providing physical configuration of the interfaces Type: String.

- `relation_vmm_rs_vswitch_override_stp_pol` - (Optional) Represents the relation to a STP Policy (class stpIfPol). Type: String.

## Importing

An existing VSwitchPolicyGroup can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import vswitch_policy.example <Dn>
```
