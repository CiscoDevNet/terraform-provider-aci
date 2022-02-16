---
subcategory: "Access Policies"
layout: "aci"
page_title: "ACI: aci_mcp_instance_policy"
sidebar_current: "docs-aci-data-source-mcp_instance_policy"
description: |-
  Data source for ACI MCP Instance Policy
---

# aci_mcp_instance_policy #

Data source for ACI MCP Instance Policy


## API Information ##

* `Class` - mcpInstPol
* `Distinguished Named` - uni/infra/mcpInstP-{name}

## GUI Information ##

* `Location` - Fabric -> Access Policies -> Policies -> Global -> MCP Instance Policy Default



## Example Usage ##

```hcl
data "aci_mcp_instance_policy" "example" {}
```

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the MCP Instance Policy.
* `annotation` - (Optional) Annotation of object MCP Instance Policy.
* `name_alias` - (Optional) Name Alias of object MCP Instance Policy.
* `description` - (Optional) Description of Object MCP Instance Policy.
* `admin_st` - (Optional) Admin State. The administrative state of the object or policy.
* `ctrl` - (Optional) Controls. The control state.
* `init_delay_time` - (Optional) Init Delay Time. 
* `loop_detect_mult` - (Optional) Loop Detection Multiplier. 
* `loop_protect_act` - (Optional) Loop Protection Action. 
* `tx_freq` - (Optional) Transmission Frequency. Sets the transmission frequency of the instance advertisements.
* `tx_freq_msec` - (Optional) Transmission Frequency. Sets the transmission frequency of mcp advertisements in milliseconds
