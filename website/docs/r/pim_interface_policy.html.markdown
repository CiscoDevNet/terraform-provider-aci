---
subcategory: "Tenant Policies"
layout: "aci"
page_title: "ACI: aci_pim_interface_policy"
sidebar_current: "docs-aci-resource-pim_interface_policy"
description: |-
  Manages ACI PIM Interface Policy
---

# aci_pim_interface_policy #

Manages ACI PIM Interface Policy

## API Information ##

* `Class` - pimIfPol
* `Distinguished Name` - uni/tn-{tenant_name}/pimifpol-{name}

## GUI Information ##

* `Location` - Tenants -> Policies -> Protocol -> PIM

## Example Usage ##

```hcl
resource "aci_pim_interface_policy" "example_ip" {
  tenant_dn                  = aci_tenant.example.id
  name                       = "example_ip"
  designated_router_delay    = "3"
  designated_router_priority = "1"
  hello_interval             = "30000"
  join_prune_interval        = "60"
  control_state              = ["border"]
}
```

## Argument Reference ##

* `tenant_dn` - (Required) Distinguished name of the parent Tenant object.
* `name` - (Required) Name of the PIM Interface Policy object.
* `annotation` - (Optional) Annotation of the PIM Interface Policy object.
* `name_alias` - (Optional) Name Alias of the PIM Interface Policy object.
* `auth_key` - (Optional) Secure authentication key.
* `auth_type` - (Optional) Authentication type. Allowed values are "ah-md5", "none" and the default value is "none". Type: String.
* `control_state` - (Optional) Interface controls. Allowed values are "border", "passive" and "strict-rfc-compliant". Type: List.
* `designated_router_delay` - (Optional) Designated router delay. Allowed range is "1-65535" and the default value is "3".
* `designated_router_priority` - (Optional) Designated router priority. Allowed range is "1-4294967295" and the default value is "1".
* `hello_interval` - (Optional) Hello traffic policy. Allowed range is "1-18724286" and the default value is "30000".
* `join_prune_interval` - (Optional) Join Prune Traffic Policy. Allowed range is "60-65520" and the default value is "60".
* `inbound_join_prune_filter_policy` - (Optional) Inbound join prune filter policy which represents the relation to a PIM Route Map Filter (class rtdmcARtMapPol). Type: String.
* `outbound_join_prune_filter_policy` - (Optional) Outbound join prune filter policy which represents the relation to a PIM Route Map Filter (class rtdmcARtMapPol). Type: String.
* `neighbor_filter_policy` - (Optional) Neighbor filter policy which represents the relation to a PIM Route Map Filter (class rtdmcARtMapPol). Type: String.

## Importing ##

An existing PIM Interface Policy can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_pim_interface_policy.example <Dn>
```