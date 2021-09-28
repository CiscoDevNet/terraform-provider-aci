---
layout: "aci"
page_title: "ACI: aci_endpoint_controls"
sidebar_current: "docs-aci-resource-endpoint_controls"
description: |-
  Manages ACI Endpoint Control 
---

# aci_endpoint_controls #
Manages ACI Endpoint Control 

## API Information ##
* `Class` - epControlP
* `Distinguished Named` - uni/infra/epCtrlP-{name}

## GUI Information ##
* `Location` - System -> System Settings -> Endpoint Controls -> Rouge EP Control -> Policy


## Example Usage ##
```hcl
resource "aci_endpoint_controls" "example" {
  admin_st = "disabled"
  annotation = "orchestrator:terraform"
  hold_intvl = "1800"
  rogue_ep_detect_intvl = "60"
  rogue_ep_detect_mult = "4"
  description = "from terraform"
  name_alias = "example_name_alias"
}
```

## NOTE ##
User can use resource of type `aci_endpoint_controls` to change configuration of object Endpoint Control. User cannot create more than one instances of object Endpoint Control.

## Argument Reference ##
* `annotation` - (Optional) Annotation of object Endpoint Control.
* `admin_st` - (Optional) The administrative state of  object Endpoint Control. Allowed values are "disabled" and "enabled".
* `hold_intvl` - (Optional) The period of time before declaring that the neighbor is down of object Endpoint Control. Allowed range: "300" - "3600".
* `rogue_ep_detect_intvl` - (Optional) Rogue Endpoint Detection Interval of object Endpoint Control. Allowed range: "30" - "3600".
* `rogue_ep_detect_mult` - (Optional) Rogue Endpoint Detection Multiplication Factor of object Endpoint Control. Allowed range is "2" - "65535". 
* `description` - (Optional) Description of object Endpoint Control.
* `name_alias` - (Optional) Name alias of object Endpoint Control.



## Importing ##

An existing EndpointControl can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_endpoint_controls.example <Dn>
```