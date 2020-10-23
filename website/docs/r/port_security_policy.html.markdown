---
layout: "aci"
page_title: "ACI: aci_port_security_policy"
sidebar_current: "docs-aci-resource-port_security_policy"
description: |-
  Manages ACI Port Security Policy
---

# aci_port_security_policy #
Manages ACI Port Security Policy

## Example Usage ##

```hcl
	resource "aci_port_security_policy" "fooport_security_policy" {
		description = "%s"
		name        = "demo_port_pol"
		annotation  = "tag_port_pol"
		maximum     = "12"
		name_alias  = "alias_port_pol"
		timeout     = "60"
		violation   = "protect"
	}
```
## Argument Reference ##
* `name` - (Required) name of Object port_security_policy.
* `annotation` - (Optional) annotation for object port_security_policy.
* `maximum` - (Optional) Port Security Maximum. Allowed value range is "0" - "12000". Default is "0".
* `mode` - (Optional) bgp domain mode
* `name_alias` - (Optional) name_alias for object port_security_policy.
* `timeout` - (Optional) amount of time between authentication attempts. Allowed value range is "60" - "3600". Default is "60".
* `violation` - (Optional) Port Security Violation. default value is "protect".
Allowed value: "protect"




## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Port Security Policy.

## Importing ##

An existing Port Security Policy can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_port_security_policy.example <Dn>
```