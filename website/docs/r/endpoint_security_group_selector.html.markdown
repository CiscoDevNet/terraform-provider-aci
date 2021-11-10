---
subcategory: "Application Management"
layout: "aci"
page_title: "ACI: aci_endpoint_security_group_selector"
sidebar_current: "docs-aci-resource-endpoint_security_group_selector"
description: |-
  Manages ACI Endpoint Security Group Selector
---

# aci_endpoint_security_group_selector #

Manages ACI Endpoint Security Group Selector

## API Information ##

* `Class` - fvEPSelector
* `Distinguished Named` - uni/tn-{name}/ap-{name}/esg-{name}/epselector-{[matchExpression]}

## GUI Information ##

* `Location` - Tenants > {tenant_name} > Application Profiles > Endpoint Security Groups > Selectors

## Example Usage ##

```hcl
resource "aci_endpoint_security_group_selector" "example" {
  endpoint_security_group_dn  = aci_endpoint_security_group.example.id
  annotation = "orchestrator:terraform"
  description = "from terraform"
  name = "example"
  name_alias = "example"
  match_expression = "ip=='10.10.10.0/24'"
}
```

## Argument Reference ##

* `endpoint_security_group_dn` - (Required) Distinguished name of parent Endpoint Security Group object.
* `annotation` - (Optional) Annotation of object Endpoint Security Group Selector.
* `description` - (Optional) Description of object Endpoint Security Group Selector.
* `name` - (Optional) Name of object Endpoint Security Group Selector.
* `name_alias` - (Optional) Name Alias of object Endpoint Security Group Selector.
* `match_expression` - (Optional) Expression used to define matching tags.  

## Importing ##

An existing EndpointSecurityGroupSelector can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import endpoint_security_group_selector.example <Dn>
```