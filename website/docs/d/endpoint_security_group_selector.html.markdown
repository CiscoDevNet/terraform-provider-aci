---
subcategory: "Application Management"
layout: "aci"
page_title: "ACI: aci_endpoint_security_group_selector"
sidebar_current: "docs-aci-data-source-endpoint_security_group_selector"
description: |-
  Data source for ACI Endpoint Security Group Selector
---

# aci_endpoint_security_group_selector #

Data source for ACI Endpoint Security Group Selector

## API Information ##

* `Class` - fvEPSelector
* `Distinguished Name` - uni/tn-{name}/ap-{name}/esg-{name}/epselector-{[matchExpression]}

## GUI Information ##

* `Location` - Tenants > {tenant_name} > Application Profiles > Endpoint Security Groups > Selectors

## Example Usage ##

```hcl
data "aci_endpoint_security_group_selector" "example" {
  endpoint_security_group_dn  = aci_endpoint_security_group.example.id
  match_expression = "ip=='10.10.10.0/24'"
}
```

## Argument Reference ##

* `endpoint_security_group_dn` - (Required) Distinguished name of parent Endpoint Security Group object.
* `match_expression` - (Required) Expression used to define matching tags.

## Attribute Reference ##

* `id` - Attribute id set to the Dn of the Endpoint Security Group Selector.
* `annotation` - (Optional) Annotation of object Endpoint Security Group Selector.
* `description` - (Optional) Description of object Endpoint Security Group Selector.
* `name` - (Optional) Name of object Endpoint Security Group Selector.
* `name_alias` - (Optional) Name Alias of object Endpoint Security Group Selector.
