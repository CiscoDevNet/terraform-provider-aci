---
layout: "aci"
page_title: "ACI: aci_mcp_instance_policy"
sidebar_current: "docs-aci-resource-mcp_protocol_instance_policy"
description: |-
  Manages ACI MCP Instance Policy
---

# aci_mcp_instance_policy #

Manages ACI MCP Instance Policy

## API Information ##

* `Class` - mcpInstPol
* `Distinguished Named` - uni/infra/mcpInstP-{name}

## GUI Information ##

* `Location` - Fabric -> Access Policies -> Policies -> Global -> MCP Instance Policy Default


## Example Usage ##

```hcl
resource "aci_mcp_instance_policy" "example" {
  admin_st         = "disabled"
  annotation       = "orchestrator:terraform"
  name_alias       = "mcp_instance_alias"
  description      = "From Terraform"
  ctrl             = []
  init_delay_time  = "180"
  key              = "example"
  loop_detect_mult = "3"
  loop_protect_act = "port-disable"
  tx_freq          = "2"
  tx_freq_msec     = "0"
}
```

## NOTE ##
User can use resource of type aci_mcp_instance_policy to change configuration of object MCP Instance Policy. User cannot create more than one instances of object MCP Instance Policy.

## Argument Reference ##


* `annotation` - (Optional) Annotation of object MCP Instance Policy.
* `admin_st` - (Optional) Admin State.The administrative state of the object or policy. Allowed values are "disabled", "enabled". Type: String.
* `description` - (Optional) Description for object MCP Instance Policy. Type: String.
* `name_alias` - (Optional) Name Alias for object MCP Instance Policy. Type: String.
* `ctrl` - (Optional) Controls.The control state. Allowed values are "pdu-per-vlan", "stateful-ha",  Type: List.
* `init_delay_time` - (Optional) Init Delay Time. Allowed range is "0"-"1800", Type: String.
* `key` - (Required) Secret Key.The key or password used to uniquely identify this configuration object.
* `loop_detect_mult` - (Optional) Loop Detection Multiplier. Allowed range is "1"-"255", Type: String.
* `loop_protect_act` - (Optional) Loop Protection Action. Allowed values are "port-disable","none". Type: String.
* `tx_freq` - (Optional) Transmission Frequency.Sets the transmission frequency of the instance advertisements. Allowed range is "1"-"300", Type: String.(Note: For value less than "2", loop_protect_act attribute needs to be "port-disable")
* `tx_freq_msec` - (Optional) Transmission Frequency.Sets the transmission frequency of mcp advertisements in milliseconds Allowed range is "0"-"999", Type: String.(Note: For value "0" of tx_freq, range of tx_freq_msec is  "100"-"999")


## Importing ##

An existing Mis-cablingProtocolInstancePolicy can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_mcp_instance_policy.example <Dn>
```