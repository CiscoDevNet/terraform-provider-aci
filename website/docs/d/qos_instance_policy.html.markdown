---
layout: "aci"
page_title: "ACI: aci_qos_instance_policy"
sidebar_current: "docs-aci-data-source-qos_instance_policy"
description: |-
  Data source for ACI QOS Instance Policy
---

# aci_qos_instance_policy #

Data source for ACI QOS Instance Policy


## API Information ##

* `Class` - qosInstPol
* `Distinguished Named` - uni/infra/qosinst-{name}

## Example Usage ##

```hcl
data "aci_qos_instance_policy" "example" {}
```

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the QOS Instance Policy.
* `annotation` - (Optional) Annotation of object QOS Instance Policy.
* `name_alias` - (Optional) Name Alias of object QOS Instance Policy.
* `description` - (Optional) Description of object QOS Instance Policy
* `etrap_age_timer` - (Optional) E-trap flow age out timer. 
* `etrap_bw_thresh` - (Optional) Track activeness of elephant flow. 
* `etrap_byte_ct` - (Optional) E-trap elephant flow identifier. 
* `etrap_st` - (Optional) E-trap enable knob. E-trap parameters
* `fabric_flush_interval` - (Optional) Fabric Flush Interval in ms. 
* `fabric_flush_st` - (Optional) Fabric PFC Flush enable knob. Fabric Flush parameters
* `ctrl` - (Optional) Global Control Settings. The control state.
* `uburst_spine_queues` - (Optional) micro burst spine queues percent. Global microburst spine % queues
* `uburst_tor_queues` - (Optional) micro burst tor queues percent. Global microburst tor % queues
