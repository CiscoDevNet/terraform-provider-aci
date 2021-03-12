---
layout: "aci"
page_title: "ACI: aci_endpoint_security_group_selector"
sidebar_current: "docs-aci-data-source-endpoint_security_group_selector"
description: |-
  Data source for ACI Endpoint Security Group Selector
---

# aci_endpoint_security_group_selector #
Data source for ACI Endpoint Security Group Selector

## Example Usage ##

```hcl
data "aci_endpoint_security_group_selector" "example" {

  endpoint_security_group_dn  = "${aci_endpoint_security_group.example.id}"

  match_expression  = "example"
}
```
## Argument Reference ##
* `endpoint_security_group_dn` - (Required) Distinguished name of parent EndpointSecurityGroup object.
* `matchExpression` - (Required) matchExpression of Object endpoint_security_group_selector.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Endpoint Security Group Selector.
* `annotation` - (Optional) annotation for object endpoint_security_group_selector.
* `match_expression` - (Optional) match_expression for object endpoint_security_group_selector.
* `name_alias` - (Optional) name_alias for object endpoint_security_group_selector.
* `userdom` - (Optional) userdom for object endpoint_security_group_selector.
