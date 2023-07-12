---
layout: "aci"
page_title: "ACI: aci_bfd_multihop_interface_policy"
sidebar_current: "docs-aci-resource-bfd_multihop_interface_policy"
description: |-
  Manages ACI BFD Multihop Interface Policy
---

# aci_bfd_multihop_interface_policy #

Manages ACI BFD Multihop Interface Policy

## API Information ##

* `Class` - bfdMhIfPol
* `Distinguished Name` - uni/tn-{tn_name}/bfdMhIfPol-{policy_name}

## GUI Information ##

* `Location` - Tenant -> Policies -> Protocol -> BFD Multihop -> Interface Policies


## Example Usage ##

```hcl
resource "aci_bfd_multihop_interface_policy" "example" {
  tenant_dn             = aci_tenant.example.id
  name                  = "example"
  admin_state           = "enabled"
  annotation            = "orchestrator: terraform"
  detection_multiplier  = "3"
  min_transmit_interval = "250"
  min_receive_interval  = "250"
}
```

## Argument Reference ##

* `tenant_dn` - (Required) Distinguished name of the parent Tenant object.
* `name` - (Required) Name of the BFD Multihop Interface Policy object.
* `annotation` - (Optional) Annotation of the BFD Multihop Interface Policy object.
* `admin_state` - (Optional) Enable Disable sessions.The administrative state of the object or policy. Allowed values are "disabled", "enabled", and default value is "enabled". Type: String.

* `detection_multiplier` - (Optional) Detection Multiplier.Detection multiplier. Allowed range is 1-50 and default value is "3".
* `min_receive_interval` - (Optional) Required Minimum RX Interval.Required minimum rx interval. Allowed range is 250-999 and default value is "250".
* `min_transmit_interval` - (Optional) Desired Minimum TX Interval.Desired minimum tx interval. Allowed range is 250-999 and default value is "250".


## Importing ##

An existing BFD Multihop Interface Policy can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_bfd_multihop_interface_policy.example <Dn>
```