---
subcategory: "Access Policies"
layout: "aci"
page_title: "ACI: aci_lacp_member_policy"
sidebar_current: "docs-aci-resource-lacp_member_policy"
description: |-
  Manages ACI LACP Member Policy
---

# aci_lacpmember_policy #

Manages ACI LACP Member Policy

## API Information ##

* `Class` - lacpIfPol
* `Distinguished Name` - uni/infra/lacpifp-{name}

## GUI Information ##

* `Location` - Fabric - Access Policies - Policies - Interface - Port Channel Member

## Example Usage ##

```hcl
resource "aci_lacp_member_policy" "example" {
  name  = "example"
  prio = "32768"
  tx_rate = "normal"
}
```

## Argument Reference ##

* `name` - (Required) Name of the object LACP Member Policy.
* `annotation` - (Optional) Annotation of the object LACP Member Policy.
* `description` - (Optional) Description of the object LACP Member Policy.
* `name_alias` - (Optional) Name alias.
* `prio` - (Optional) Priority.Port priority - LACP uses the port priority to decide which ports should be put in standby mode when there is a limitation that prevents all compatible ports from aggregating and which ports should be put into active mode. A higher port priority value means a lower priority for LACP Allowed range is 1-65535 and default value is "32768".
* `tx_rate` - (Optional) Transmission Rate.The configured transmit rate of the LACP packets. Allowed values are "fast", "normal", and default value is "normal". Type: String.

## Importing ##

An existing LACPMemberPolicy can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_lacpmember_policy.example <Dn>
```