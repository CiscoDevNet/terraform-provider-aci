---
layout: "aci"
page_title: "ACI: aci_span_destination_group"
sidebar_current: "docs-aci-resource-span_destination_group"
description: |-
  Manages ACI SPAN Destination Group
---

# aci_span_destination_group #
Manages ACI SPAN Destination Group

## Example Usage ##

```hcl
resource "aci_span_destination_group" "example" {

  tenant_dn  = "${aci_tenant.example.id}"

  name  = "example"
  annotation  = "example"
  name_alias  = "example"
}
```
## Argument Reference ##
* `tenant_dn` - (Required) Distinguished name of parent Tenant object.
* `name` - (Required) name of Object span_destination_group.
* `annotation` - (Optional) 
* `name_alias` - (Optional) 



## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the SPAN Destination Group.

## Importing ##

An existing SPAN Destination Group can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_span_destination_group.example <Dn>
```