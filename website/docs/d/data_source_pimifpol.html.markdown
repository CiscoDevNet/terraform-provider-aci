---
subcategory: "Tenant Policies"
layout: "aci"
page_title: "ACI: aci_pim_interface_policy"
sidebar_current: "docs-aci-data-source-pim_interface_policy"
description: |-
  Data source for ACI PIM Interface Policy
---

# aci_pim_interface_policy #

Data source for ACI PIM Interface Policy

## API Information ##

* `Class` - pimIfPol
* `Distinguished Name` - uni/tn-{tenant_name}/pimifpol-{name}

## GUI Information ##

* `Location` - Tenants -> Policies -> Protocol -> PIM

## Example Usage ##

```hcl
data "aci_pim_interface_policy" "example" {
  tenant_dn  = aci_tenant.example.id
  name       = "example"
}
```

## Argument Reference ##

* `tenant_dn` - (Required) Distinguished name of the parent Tenant object.
* `name` - (Required) Name of the PIM Interface Policy object.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the PIM Interface Policy.
* `annotation` - (Read-Only) Annotation of the PIM Interface Policy object.
* `name_alias` - (Read-Only) Name Alias of the PIM Interface Policy object.
* `auth_t` - (Read-Only) Authentication Type. 
* `ctrl` - (Read-Only) Interface Controls. 
* `dr_delay` - (Read-Only) Designated Router Delay. 
* `dr_prio` - (Read-Only) Designated Router Priority. 
* `hello_itvl` - (Read-Only) Hello Traffic Policy
* `jp_interval` - (Read-Only) JP Traffic Policy
