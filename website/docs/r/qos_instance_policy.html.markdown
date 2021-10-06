---
layout: "aci"
page_title: "ACI: aci_qos_instance_policy"
sidebar_current: "docs-aci-resource-qos_instance_policy"
description: |-
  Manages ACI QOS Instance Policy
---

# aci_qos_instance_policy #

Manages ACI QOS Instance Policy

## API Information ##

* `Class` - qosInstPol
* `Distinguished Named` - uni/infra/qosinst-{name}


## Example Usage ##

```hcl
resource "aci_qos_instance_policy" "example" {
  name_alias            = "qos_instance_alias"
  description           = "From Terraform"
  etrap_age_timer       = "0" 
  etrap_bw_thresh       = "0"
  etrap_byte_ct         = "0"
  etrap_st              = "no"
  fabric_flush_interval = "500"
  fabric_flush_st       = "no"
  annotation            = "orchestrator:terraform"
  ctrl                  = "none"
  uburst_spine_queues   = "10"
  uburst_tor_queues     = "10"
}
```
## NOTE ##
User can use resource of type aci_qos_instance_policy to change configuration of object QOS Instance Policy. User cannot create more than one instances of object QOS Instance Policy.

## Argument Reference ##


* `annotation` - (Optional) Annotation of object QOS Instance Policy.
* `description` - (Optional) Description for object QOS Instance Policy. Type: String.
* `name_alias` - (Optional) Name Alias for object QOS Instance Policy. Type: String.
* `etrap_age_timer` - (Optional) E-trap flow age out timer. Min Allowed Value is "0".
* `etrap_bw_thresh` - (Optional) Track activeness of elephant flow. Min Allowed Value is "0".
* `etrap_byte_ct` - (Optional) E-trap elephant flow identifier. Min Allowed Value is "0".
* `etrap_st` - (Optional) E-trap enable knob. E-trap parameters. Allowed values are "no", "yes". Type: String.
* `fabric_flush_interval` - (Optional) Fabric Flush Interval in ms. Allowed range is "100"-"1000". Type: String.
* `fabric_flush_st` - (Optional) Fabric PFC Flush enable knob. Fabric Flush parameters Allowed values are "no", "yes". Type: String.
* `ctrl` - (Optional) Global Control Settings. The control state. Allowed values are "dot1p-preserve", "none". Type: String.
* `uburst_spine_queues` - (Optional) Micro burst spine queues percent. Allowed range is "0"-"100". Type: String
* `uburst_tor_queues` - (Optional) Micro burst tor queues percent. Allowed range is "0"-"100". Type: String


## Importing ##

An existing QOSInstancePolicy can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_qos_instance_policy.example <Dn>
```