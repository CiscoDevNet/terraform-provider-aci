---
layout: "aci"
page_title: "ACI: aci_error_disable_recovery"
sidebar_current: "docs-aci-resource-error_disable_recovery"
description: |-
  Manages ACI Error Disable Recovery
---

# aci_error_disable_recovery #

Manages ACI Error Disable Recovery

## API Information ##

* `Class` - edrErrDisRecoverPol and edrEventP
* `Distinguished Named` - uni/infra/edrErrDisRecoverPol-{name}

## GUI Information ##

* `Location` - Fabric -> Access Policies -> Policies -> Global -> Error Disabled Recovery Policy


## Example Usage ##

```hcl
resource "aci_error_disable_recovery" "example" {
  annotation          = "orchestrator:terraform"
  err_dis_recov_intvl = "300"
  name_alias          = "error_disable_recovery_alias"
  description         = "From Terraform"
  edr_event {
    event             = "event-mcp-loop"
    recover           = "yes"
    description       = "From Terraform"
    name_alias        = "event_alias"
    name              = "example"
    annotation        = "orchestrator:terraform"
  }
}
```

## Note ##
Users can use the resource of type aci_error_disable_recovery to change the configuration of the object Error Disable Recovery. Users cannot create more than one instance of object Error Disable Recovery.
## Argument Reference ##


* `annotation` - (Optional) Annotation of object Error Disable Recovery. Type String.
* `name_alias` - (Optional) Name Alias of object Error Disable Recovery. Type: String.
* `description` - (Optional) Description of object Error Disable Recovery. Type: String.
* `err_dis_recov_intvl` - (Optional) Error Disable Recovery Interval.Sets the error disable recovery interval, which specifies the time to recover from an error-disabled state. Allowed range is "30" - "65535". Type: String.
* `edr_event` - (Optional) To manage Error Disable Recovery Event from the Error Disable Recovery Policy resource. 
* `edr_event.event` - (Required) Event of object Error Disabled Recovery. The error disable recovery event type. Allowed values are "event-arp-inspection", "event-bpduguard", "event-debug-1", "event-debug-2", "event-debug-3", "event-debug-4", "event-debug-5", "event-dhcp-rate-lim", "event-ep-move", "event-ethpm", "event-ip-addr-conflict", "event-ipqos-dcbxp-compat-failure", "event-ipqos-mgr-error", "event-link-flap", "event-loopback", "event-mcp-loop", "event-psec-violation", "event-sec-violation", "event-set-port-state-failed", "event-storm-ctrl", "event-stp-inconsist-vpc-peerlink", "event-syserr-based", "event-udld", "unknown".  Type: String.
* `edr_event.recover` - (Optional) Enables or disables Error Disable Recovery. Allowed values are "no", "yes". Type: String.
* `edr_event.name` - (Optional) Name of object Error Disable Recovery Event. Type: String.
* `edr_event.name_alias` - (Optional) Name Alias of object Error Disable Recovery Event. Type: String.
* `edr_event.description` - (Optional) Description of object Error Disable Recovery Event. Type: String. 
* `edr_event.annotation` - (Optional) Annotation of object Error Disable Recovery Event. Type String.

## Importing ##

An existing Error Disable Recovery can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_error_disable_recovery.example <Dn>
```