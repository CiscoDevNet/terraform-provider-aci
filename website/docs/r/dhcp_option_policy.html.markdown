---
layout: "aci"
page_title: "ACI: aci_dhcp_option_policy"
sidebar_current: "docs-aci-resource-dhcp_option_policy"
description: |-
  Manages ACI DHCP Option Policy
---

# aci_dhcp_option_policy #
Manages ACI DHCP Option Policy

## Example Usage ##

```hcl
resource "aci_dhcp_option_policy" "example" {

  tenant_dn  = "${aci_tenant.example.id}"
  name  = "example"
  annotation  = "example"
  name_alias  = "example"
}
```
## Argument Reference ##
* `tenant_dn` - (Required) Distinguished name of parent Tenant object.
* `name` - (Required) Name of Object  DHCP Option Policy.
* `annotation` - (Optional) Annotation for object  DHCP Option Policy.
* `name_alias` - (Optional) Name alias for object  DHCP Option Policy.



## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the DHCP Option Policy.

## Importing ##

An existing DHCP Option Policy can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_dhcp_option_policy.example <Dn>
```