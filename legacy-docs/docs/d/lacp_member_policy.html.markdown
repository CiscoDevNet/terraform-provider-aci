---
subcategory: "Access Policies"
layout: "aci"
page_title: "ACI: aci_lacp_member_policy"
sidebar_current: "docs-aci-data-source-aci_lacp_member_policy"
description: |-
  Data source for ACI LACP Member Policy
---

# aci_lacpmember_policy #

Data source of the ACI LACP Member Policy

## API Information ##

* `Class` - lacpIfPol
* `Distinguished Name` - uni/infra/lacpifp-{name}

## GUI Information ##

* `Location` - Fabric - Access Policies - Policies - Interface - Port Channel Member

## Example Usage ##

```hcl
data "aci_lacp_member_policy" "example" {
  name  = "example"
}
```

## Argument Reference ##

* `name` - (Required) Name of LACP Member Policy object.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the LACP Member Policy.
* `annotation` - (Optional) Annotation of the LACP Member Policy object.
* `name_alias` - (Optional) Name Alias of the LACP Member Policy object.
* `priority` - (Optional) Port priority - LACP uses the port priority to decide which ports should be put in standby mode when there is a limitation that prevents all compatible ports from aggregating and which ports should be put into active mode. A higher port priority value means a lower priority for LACP.
* `transmit_rate` - (Optional) Transmission Rate. The configured transmit rate of the LACP packets.
