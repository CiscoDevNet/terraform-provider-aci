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
  tenant_dn = "example"
  name  = "example"
  name_alias = "example"
}
```


## Argument Reference ##

* `name` - (Required) name of object monitoring policy.
* `tenant_dn` - (Required) tenant dn for monitoring policy.
* `name_alias` - (Optional) name alias for monitoring policy.
* `annotation` - (Optional) annotation for object monitoring policy.
                


## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the monitoring Policy.

## Importing ##

An existing monitoring Policy can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_monitoring_policy.example <Dn>
```