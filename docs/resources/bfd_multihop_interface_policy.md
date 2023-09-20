---
subcategory: "Tenant Policies"
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
  detection_multiplier  = "3"
  min_transmit_interval = "250"
  min_receive_interval  = "250"
}
```

## Argument Reference ##

* `tenant_dn` - (Required) Distinguished name of the parent Tenant object. Type: String.
* `name` - (Required) Name of the BFD Multihop Interface Policy object. Type: String.
* `name_alias` - (Optional) Name Alias of the BFD Multihop Interface Policy object. Type: String.
* `annotation` - (Optional) Annotation of the BFD Multihop Interface Policy object. Type: String.
* `admin_state` - (Optional) The administrative state of the object or policy. Allowed values are "disabled", "enabled", and default value is "enabled". Type: String.
* `description` - (Optional) Description for the BFD Multihop Interface Policy object. Type: String.
* `detection_multiplier` - (Optional) Detection Multiplier. Allowed range is 1-50 and default value is "3".  Type: String.
* `min_receive_interval` - (Optional) Required Minimum Rx Interval. Allowed range is 250-999 and default value is "250".  Type: String.
* `min_transmit_interval` - (Optional) Desired Minimum Tx Interval. Allowed range is 250-999 and default value is "250".  Type: String.


## Importing ##

An existing BFD Multihop Interface Policy can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_bfd_multihop_interface_policy.example <Dn>
```