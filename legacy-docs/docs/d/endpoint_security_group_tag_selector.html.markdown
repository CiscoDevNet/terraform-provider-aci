---
subcategory: "Application Management"
layout: "aci"
page_title: "ACI: aci_endpoint_security_group_tag_selector"
sidebar_current: "docs-aci-data-source-endpoint_security_group_tag_selector"
description: |-
  Data source for ACI Endpoint Security Group Tag Selector
---

# aci_endpoint_security_group_tag_selector #

Data source for ACI Endpoint Security Group Tag Selector


## API Information ##

* `Class` - fvTagSelector
* `Distinguished Name` - uni/tn-{name}/ap-{name}/esg-{name}/tagselectorkey-[{matchKey}]-value-[{matchValue}]

## GUI Information ##

* `Location` - Tenants > {tenant_name} > Application Profiles > Endpoint Security Groups > Selectors > Tag Selectors



## Example Usage ##

```hcl
data "aci_endpoint_security_group_tag_selector" "example" {
  endpoint_security_group_dn  = aci_endpoint_security_group.example.id
  match_key  = "example"
  match_value  = "example"
}
```

## Argument Reference ##

* `endpoint_security_group_dn` - (Required) Distinguished name of parent EndpointSecurityGroup object.
* `match_key` - (Required) Match key of object Endpoint Security Group Tag Selector.
* `match_value` - (Required) Match value of object Endpoint Security Group Tag Selector.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the Endpoint Security Group Tag Selector.
* `annotation` - (Optional) Annotation of object Endpoint Security Group Tag Selector.
* `match_key` - (Optional) Key of Tag to be associated with. 
* `match_value` - (Optional) Value of Tag to be associated with. 
* `value_operator` - (Optional) Match Value Operator. 
