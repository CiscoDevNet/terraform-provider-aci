---
subcategory: "Tenant Policies"
layout: "aci"
page_title: "ACI: aci_hsrp_interface_policy"
sidebar_current: "docs-aci-resource-hsrp_interface_policy"
description: |-
  Manages ACI HSRP Interface Policy
---

# aci_hsrp_interface_policy #
Manages ACI HSRP Interface Policy

## Example Usage ##

```hcl
resource "aci_hsrp_interface_policy" "example" {
  tenant_dn    = aci_tenant.tenentcheck.id
  name         = "one"
  annotation   = "example"
  description  = "from terraform"
  ctrl         = ["bia", "bfd"]
  delay        = "10"
  name_alias   = "example"
  reload_delay = "10"
}
```


## Argument Reference ##

* `tenant_dn` - (Required) Distinguished name of parent tenant object.
* `name` - (Required) Name of HSRP interface policy object.
* `annotation` - (Optional) Annotation for HSRP interface policy object.
* `description` - (Optional) Description for HSRP interface policy object.
* `ctrl` - (Optional) Control state for HSRP interface policy object. It is in the form of comma separated string and allowed values are "bia" and "bfd".
* `delay` - (Optional) Administrative port delay for HSRP interface policy object.Range: "0" to "10000". Default value is "0".
* `name_alias` - (Optional) Name alias for HSRP interface policy object.
* `reload_delay` - (Optional) Reload delay for HSRP interface policy object.Range: "0" to "10000". Default value is "0".



## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the HSRP Interface Policy.

## Importing ##

An existing HSRP Interface Policy can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_hsrp_interface_policy.example <Dn>
```