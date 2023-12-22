---
subcategory: "System Settings"
layout: "aci"
page_title: "ACI: aci_endpoint_controls"
sidebar_current: "docs-aci-resource-aci_endpoint_controls"
description: |-
  Manages ACI Endpoint Control 
---

# aci_endpoint_controls #
Manages ACI Endpoint Control 

## API Information ##
* `Class` - epControlP
* `Distinguished Name` - uni/infra/epCtrlP-{name}

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
Users can use the resource of type `aci_endpoint_controls` to change the configuration of the object Endpoint Control. Users cannot create more than one instance of object Endpoint Control.

## Argument Reference ##
* `annotation` - (Optional) Annotation of object Endpoint Control.
* `admin_st` - (Optional) The administrative state of  object Endpoint Control. Allowed values are "disabled" and "enabled".
* `hold_intvl` - (Optional) The period of time before declaring that the neighbor is down of object Endpoint Control. Allowed range: "300" - "3600".
* `rogue_ep_detect_intvl` - (Optional) Rogue Endpoint Detection Interval of object Endpoint Control. Allowed range: "30" - "3600".
* `rogue_ep_detect_mult` - (Optional) Rogue Endpoint Detection Multiplication Factor of object Endpoint Control. Allowed range is "2" - "65535". 
* `description` - (Optional) Description of object Endpoint Control.
* `name_alias` - (Optional) Name alias of object Endpoint Control.



## Importing ##

An existing Endpoint Control can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_endpoint_controls.example <Dn>
```