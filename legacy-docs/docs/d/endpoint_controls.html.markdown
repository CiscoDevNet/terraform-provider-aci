---
subcategory: "System Settings"
layout: "aci"
page_title: "ACI: aci_endpoint_controls"
sidebar_current: "docs-aci-data-source-endpoint_controls"
description: |-
  Data source for ACI Endpoint Control 
---

# aci_endpoint_controls #
Data source for ACI Endpoint Control 


## API Information ##
* `Class` - epControlP
* `Distinguished Name` - uni/infra/epCtrlP-{name}

## GUI Information ##
* `Location` - System -> System Settings -> Endpoint Controls -> Rouge EP Control -> Policy

## Example Usage ##
```hcl
data "aci_endpoint_controls" "example" {}
```

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the Endpoint Control.
* `annotation` - (Optional) Annotation of object Endpoint Control.
* `name_alias` - (Optional) Name Alias of object Endpoint Control.
* `admin_st` - (Optional) The administrative state of object Endpoint Control.
* `hold_intvl` - (Optional) The period of time before declaring that the neighbor is down of object Endpoint Control.
* `rogue_ep_detect_intvl` - (Optional)  Rogue Endpoint Detection Interval of object Endpoint Control.
* `rogue_ep_detect_mult` - (Optional)  Rogue Endpoint Detection Multiplication Factor of object Endpoint Control.
* `description` - (Optional) Description of object Endpoint Control.
