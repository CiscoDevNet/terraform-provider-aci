---
layout: "aci"
page_title: "ACI: aci_port_tracking"
sidebar_current: "docs-aci-data-source-port_tracking"
description: |-
  Data source for ACI Port Tracking
---

# aci_port_tracking #

Data source for ACI Port Tracking


## API Information ##

* `Class` - infraPortTrackPol
* `Distinguished Named` - uni/infra/trackEqptFabP-{name}

## GUI Information ##

* `Location` - System -> System Settings -> Port Tracking



## Example Usage ##

```hcl
data "aci_port_tracking" "example" {}
```

## Argument Reference ##

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the Port Tracking.
* `annotation` - (Optional) Annotation of object Port Tracking.
* `name_alias` - (Optional) Name Alias of object Port Tracking.
* `admin_st` - (Optional) Port Tracking State. The administrative state of the object or policy.
* `delay` - (Optional) Delay Timeout. The administrative port delay.
* `include_apic_ports` - (Optional) Include APIC Ports when port tracking is triggered. 
* `minlinks` - (Optional) Minimum links left up before trigger. 
