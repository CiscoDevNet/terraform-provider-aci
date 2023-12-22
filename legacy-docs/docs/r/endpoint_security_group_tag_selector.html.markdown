---
subcategory: "Application Management"
layout: "aci"
page_title: "ACI: aci_endpoint_security_group_tag_selector"
sidebar_current: "docs-aci-resource-aci_endpoint_security_group_tag_selector"
description: |-
  Manages ACI Endpoint Security Group Tag Selector
---

# aci_endpoint_security_group_tag_selector #

Manages ACI Endpoint Security Group Tag Selector

## API Information ##

* `Class` - fvTagSelector
* `Distinguished Name` - uni/tn-{name}/ap-{name}/esg-{name}/tagselectorkey-[{matchKey}]-value-[{matchValue}]

## GUI Information ##

* `Location` - Tenants > {tenant_name} > Application Profiles > Endpoint Security Groups > Selectors > Tag Selectors


## Example Usage ##

```hcl
resource "aci_endpoint_security_group_tag_selector" "example" {
  endpoint_security_group_dn  = aci_endpoint_security_group.example.id
  annotation = "orchestrator:terraform"
  match_key = "example-Key"
  match_value = "example-Value"
  value_operator = "equals"
}
```

## Argument Reference ##

* `endpoint_security_group_dn` - (Required) Distinguished name of parent EndpointSecurityGroup object.
* `match_key` - (Required) Match key of object Endpoint Security Group Tag Selector.
* `match_value` - (Required) Match value of object Endpoint Security Group Tag Selector.
* `annotation` - (Optional) Annotation of object Endpoint Security Group Tag Selector.
* `value_operator` - (Optional) Match Value Operator. Allowed values are "contains", "equals", "regex", and default value is "equals". Type: String.


## Importing ##

An existing EndpointSecurityGroupTagSelector can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_endpoint_security_group_tag_selector.example <Dn>
```