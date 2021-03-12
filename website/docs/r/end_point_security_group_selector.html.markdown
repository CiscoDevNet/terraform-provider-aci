---
layout: "aci"
page_title: "ACI: aci_endpoint_security_group_selector"
sidebar_current: "docs-aci-resource-endpoint_security_group_selector"
description: |-
  Manages ACI Endpoint Security Group Selector
---

# aci_endpoint_security_group_selector #
Manages ACI Endpoint Security Group Selector

## Example Usage ##

```hcl
resource "aci_endpoint_security_group_selector" "example" {

  endpoint_security_group_dn  = "${aci_endpoint_security_group.example.id}"
  matchExpression  = "example"
  annotation  = "example"
  match_expression  = "example"
  name_alias  = "example"
  userdom  = "example"
}
```
## Argument Reference ##
* `endpoint_security_group_dn` - (Required) Distinguished name of parent EndpointSecurityGroup object.
* `annotation` - (Optional) annotation for object endpoint_security_group_selector.
* `match_expression` - (Required) match_expression for object endpoint_security_group_selector.
* `name_alias` - (Optional) name_alias for object endpoint_security_group_selector.
* `userdom` - (Optional) userdom for object endpoint_security_group_selector.
