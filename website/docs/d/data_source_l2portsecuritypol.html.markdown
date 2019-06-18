---
layout: "aci"
page_title: "ACI: aci_port_security_policy"
sidebar_current: "docs-aci-data-source-port_security_policy"
description: |-
  Data source for ACI Port Security Policy
---

# aci_port_security_policy #
Data source for ACI Port Security Policy

## Example Usage ##

```hcl
data "aci_port_security_policy" "example" {


  name  = "example"
}
```
## Argument Reference ##
* `name` - (Required) name of Object port_security_policy.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Port Security Policy.
* `annotation` - (Optional) annotation for object port_security_policy.
* `maximum` - (Optional) Port Security Maximum.
* `mode` - (Optional) bgp domain mode
* `name_alias` - (Optional) name_alias for object port_security_policy.
* `timeout` - (Optional) amount of time between authentication attempts
* `violation` - (Optional) Port security violation.
