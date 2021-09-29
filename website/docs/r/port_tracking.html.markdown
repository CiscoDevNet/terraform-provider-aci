---
layout: "aci"
page_title: "ACI: aci_port_tracking"
sidebar_current: "docs-aci-resource-port_tracking"
description: |-
  Manages ACI Port Tracking
---

# aci_port_tracking #

Manages ACI Port Tracking

## API Information ##

* `Class` - infraPortTrackPol
* `Distinguished Named` - uni/infra/trackEqptFabP-{name}

## GUI Information ##

* `Location` - System -> System Settings -> Port Tracking


## Example Usage ##

```hcl
resource "aci_port_tracking" "example" {

    admin_st           = "off"
    annotation         = "orchestrator:terraform"
    delay              = "120"
    include_apic_ports = "no"
    minlinks           = "0"
    name_alias         = "port_tracking_alias"
    description        = "From Terraform"
}
```
## NOTE ##
User can use resource of type aci_port_tracking to change configuration of object Port Tracking. User cannot create more than one instances of object Port Tracking.

## Argument Reference ##


* `annotation` - (Optional) Annotation of object Port Tracking. Type: String.
* `name_alias` - (Optional) Name Alias of object Port Tracking. Type: String.
* `description` - (Optional) Description of object Port Tracking. Type: String.
* `admin_st` - (Optional) Port Tracking State.The administrative state of the object or policy. Allowed values are "off", "on". Type: String.
* `delay` - (Optional) Delay Timeout.The administrative port delay. Allowed range is "1"-"300". Type: String.
* `include_apic_ports` - (Optional) Include APIC Ports when port tracking is triggered. Allowed values are "no", "yes". Type: String.
* `minlinks` - (Optional) Minimum links left up before trigger. Allowed range is "0"-"48". Type: String.


## Importing ##

An existing PortTracking can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_port_tracking.example <Dn>
```