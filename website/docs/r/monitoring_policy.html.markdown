---
layout: "aci"
page_title: "ACI: aci_monitoring_policy"
sidebar_current: "docs-aci-resource-monitoring_policy"
description: |-
  Manages ACI Monitoring Policy
---

# aci_monitoring_policy #
Manages ACI Monitoring Policy

## Example Usage ##

```hcl
resource "aci_monitoring_policy" "example" {
  			tenant_dn  = aci_tenant.example.id
			description = "From Terraform"
			name        = "example"
			annotation  = "example"
			name_alias  = "example"
}
```


## Argument Reference ##

* `name` - (Required) Name of object monitoring policy.
* `tenant_dn` - (Required) Tenant dn for monitoring policy.
* `name_alias` - (Optional) Name alias for monitoring policy.
* `description` - (Optional) Description for object monitoring policy.
* `annotation` - (Optional) Annotation for object monitoring policy.
                


## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the monitoring Policy.

## Importing ##

An existing monitoring Policy can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_monitoring_policy.example <Dn>
```