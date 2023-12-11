---
subcategory: "L4-L7 Services"
layout: "aci"
page_title: "ACI: aci_l4_l7_redirect_health_group"
sidebar_current: "docs-aci-resource-aci_l4_l7_redirect_health_group"
description: |-
  Manages ACI L4-L7 Redirect Health Group
---

# aci_l4_l7_redirect_health_group #

Manages ACI L4-L7 Redirect Health Group

## API Information ##

* `Class` - vnsRedirectHealthGroup
* `Distinguished Name` - uni/tn-{name}/svcCont/redirectHealthGroup-{name}

## GUI Information ##

* `Location` - Tenant -> Policies -> Protocol -> L4-L7 Redirect Health Groups

## Example Usage ##

```hcl
resource "aci_l4_l7_redirect_health_group" "example" {
  tenant_dn  = aci_tenant.example.id
  name       = "example"
}
```

## Argument Reference ##

* `tenant_dn` - (Required) Distinguished name of the parent Tenant object.
* `name` - (Required) Name of the L4-L7 Redirect Health Group object.
* `name_alias` - (Optional) Name Alias of the L4-L7 Redirect Health Group object.
* `annotation` - (Optional) Annotation of the L4-L7 Redirect Health Group object.
* `description` - (Optional) Description of the L4-L7 Redirect Health Group object.

## Importing ##

An existing L4-L7 Redirect Health Group can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_l4_l7_redirect_health_group.example <Dn>
```