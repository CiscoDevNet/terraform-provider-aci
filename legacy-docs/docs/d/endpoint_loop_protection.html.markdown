---
subcategory: "System Settings"
layout: "aci"
page_title: "ACI: aci_endpoint_loop_protection"
sidebar_current: "docs-aci-data-source-endpoint_loop_protection"
description: |-
  Data source for ACI Endpoint Loop Protection
---

# aci_endpoint_loop_protection #

Data source for ACI Endpoint Loop Protection


## API Information ##

* `Class` - epLoopProtectP
* `Distinguished Name` - uni/infra/epLoopProtectP-{name}

## GUI Information ##

* `Location` - System -> System Settings -> Endpoint Controls -> Ep Loop Protection -> Policy



## Example Usage ##

```hcl
data "aci_endpoint_loop_protection" "example" {}
```

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the Endpoint Loop Protection Policy.
* `annotation` - (Optional) Annotation of object Endpoint Loop Protection Policy.
* `name_alias` - (Optional) Name Alias of object Endpoint Loop Protection Policy.
* `description` - (Optional) Description of object Endpoint Loop Protection.
* `action` - (Optional) Action. Sets the action to take when a loop is detected.
* `admin_st` - (Optional) Admin State. The administrative state of the object or policy.
* `loop_detect_intvl` - (Optional) Loop Detection Interval. Sets the loop detection interval, which specifies the time to detect a loop.
* `loop_detect_mult` - (Optional) Loop Detection Multiplier. Sets the loop detection multiplication factor, which is the number of times a single Endpoint moves between ports within the Detection interval.
