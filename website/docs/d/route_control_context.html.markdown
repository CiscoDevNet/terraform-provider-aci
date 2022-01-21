---
subcategory: "Tenant Policies"
layout: "aci"
page_title: "ACI: aci_route_control_context"
sidebar_current: "docs-aci-data-source-route_control_context"
description: |-
  Data source for ACI Route Control Context
---

# aci_route_control_context #

Data source for ACI Route Control Context


## API Information ##

* `Class` - rtctrlCtxP
* `Distinguished Named` - uni/tn-{name}/prof-{name}/ctx-{name}

## GUI Information ##

* `Location` - Tenant -> Policies -> Protocol -> Route Maps for Route Control



## Example Usage ##

```hcl
data "aci_route_control_context" "control" {
  route_control_profile_dn  = aci_route_control_profile.bgp.id
  name  = "control"
}
```

## Argument Reference ##

* `route_control_profile_dn` - (Required) Distinguished name of parent Route Control Profile object.
* `name` - (Required) name of object Route Control Context.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the Route Control Context.
* `annotation` - (Optional) Annotation of object Route Control Context.
* `name_alias` - (Optional) Name Alias of object Route Control Context.
* `action` - (Optional) Action. The action required when the condition is met.
* `order` - (Optional) Local Order. The order of the policy context.
