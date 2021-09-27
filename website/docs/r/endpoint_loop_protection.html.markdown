---
layout: "aci"
page_title: "ACI: aci_endpoint_loop_protection"
sidebar_current: "docs-aci-resource-endpoint_loop_protection"
description: |-
  Manages ACI Endpoint Loop Protection
---

# aci_endpoint_loop_protection #

Manages ACI Endpoint Loop Protection

## API Information ##

* `Class` - epLoopProtectP
* `Distinguished Named` - uni/infra/epLoopProtectP-{name}

## GUI Information ##

* `Location` - System -> System Settings -> Endpoint Controls -> Ep Loop Protection -> Policy


## Example Usage ##

```hcl
resource "aci_endpoint_loop_protection" "example" {
  
  action = ["port-disable"]
  admin_st = "disabled"
  annotation = "orchestrator:terraform"
  loop_detect_intvl = "60"
  loop_detect_mult = "4"

}
```
## NOTE ##
User can use resource of type aci_endpoint_loop_protection to change configuration of object Endpoint Loop Protection. User cannot create more than one instances of object Endpoint Loop Protection.

## Argument Reference ##

* `annotation` - (Optional) Annotation of object Endpoint Loop Protection Policy.
* `action` - (Optional) Action.Sets the action to take when a loop is detected. Allowed values are "bd-learn-disable", "port-disable". Type: List.
* `admin_st` - (Optional) Admin State.The administrative state of the object or policy. Allowed values are "disabled", "enabled". Type: String.
* `loop_detect_intvl` - (Optional) Loop Detection Interval.Sets the loop detection interval, which specifies the time to detect a loop. Allowed range is "30"-"300". Type: String.
* `loop_detect_mult` - (Optional) Loop Detection Multiplier.Sets the loop detection multiplication factor, which is the number of times a single Endpoint moves between ports within the Detection interval. Allowed range is "1"-"255". Type: String


## Importing ##

An existing EPLoopProtectionPolicy can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_endpoint_loop_protection.example <Dn>
```