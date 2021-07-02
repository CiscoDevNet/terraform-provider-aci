---
layout: "aci"
page_title: "ACI: aci_port_security_policy"
sidebar_current: "docs-aci-data-source-port_security_policy"
description: |-
  Data source for ACI Port Security Policy
---

# aci_port_security_policy

Data source for ACI Port Security Policy

## Example Usage

```hcl
data "aci_port_security_policy" "dev_port_sec_pol" {
  name  = "foo_port_sec_pol"
}
```

## Argument Reference

- `name` - (Required) name of Object port security policy.

## Attribute Reference

- `id` - Attribute id set to the Dn of the port security policy.
- `description` - (Optional) Description for object port security policy.
- `annotation` - (Optional) Annotation for object port security policy.
- `maximum` - (Optional) Port Security Maximum.
- `name_alias` - (Optional) Name alias for object port security policy.
- `timeout` - (Optional) Amount of time between authentication attempts.
- `violation` - (Optional) Port Security Violation.
