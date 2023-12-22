---
subcategory: "System Settings"
layout: "aci"
page_title: "ACI: aci_port_tracking"
sidebar_current: "docs-aci-resource-aci_port_tracking"
description: |-
  Manages ACI Port Tracking
---

# aci_port_tracking #

Manages ACI Port Tracking

## API Information ##

* `Class` - infraPortTrackPol
* `Distinguished Name` - uni/infra/trackEqptFabP-{name}

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
Users can use the resource of type aci_port_tracking to change the configuration of the object Port Tracking. Users cannot create more than one instance of object Port Tracking.

## Argument Reference ##


* `annotation` - (Optional) Annotation of object Port Tracking. Type: String.
* `name_alias` - (Optional) Name Alias of object Port Tracking. Type: String.
* `description` - (Optional) Description of object Port Tracking. Type: String.
* `admin_st` - (Optional) Port Tracking State.The administrative state of the object or policy. Allowed values are "off", "on". Type: String.
* `delay` - (Optional) Delay Timeout.The administrative port delay. Allowed range is "1"-"300". Type: String.
* `include_apic_ports` - (Optional) Include APIC Ports when port tracking is triggered. Allowed values are "no", "yes". Type: String. (Note: attribute include_apic_ports is supported for version 5 and above of APIC)
* `minlinks` - (Optional) Minimum links left up before trigger. Allowed range is "0"-"48". Type: String.


## Importing ##

An existing PortTracking can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_port_tracking.example <Dn>
```