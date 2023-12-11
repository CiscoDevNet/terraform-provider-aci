---
subcategory: "L4-L7 Services"
layout: "aci"
page_title: "ACI: aci_l4_l7_redirect_health_group"
sidebar_current: "docs-aci-data-source-aci_l4_l7_redirect_health_group"
description: |-
  Data source for ACI L4-L7 Redirect Health Group
---

# aci_l4_l7_redirect_health_group #

Data source for ACI L4-L7 Redirect Health Group

## API Information ##

* `Class` - vnsRedirectHealthGroup
* `Distinguished Name` - uni/tn-{name}/svcCont/redirectHealthGroup-{name}

## GUI Information ##

* `Location` - Tenant -> Policies -> Protocol -> L4-L7 Redirect Health Groups

## Example Usage ##

```hcl
data "aci_l4_l7_redirect_health_group" "example" {
  tenant_dn  = aci_tenant.example.id
  name       = "example"
}
```

## Argument Reference ##

* `tenant_dn` - (Required) Distinguished name of the parent Tenant object.
* `name` - (Required) Name of the L4-L7 Redirect Health Group object.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the L4-L7 Redirect Health Group.
* `annotation` - (Optional) Annotation of the L4-L7 Redirect Health Group object.
* `name_alias` - (Optional) Name Alias of the L4-L7 Redirect Health Group object.
* `description` - (Optional) Description of the L4-L7 Redirect Health Group object.
