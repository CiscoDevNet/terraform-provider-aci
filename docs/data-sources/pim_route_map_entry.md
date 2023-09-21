---
subcategory: "Multicast"
layout: "aci"
page_title: "ACI: pim_route_map_entry"
sidebar_current: "docs-aci-data-source-pim_route_map_entry"
description: |-
  Data source for Pim Route Map Entry
---

# aci_pim_route_map_entry #

Data source for Pim Route Map Entry

## API Information ##

* `Class` - pimRouteMapEntry

* `Distinguished Name Formats`
  - uni/tn-{name}/rtmap-{name}/rtmapentry-{order}

## GUI Information ##

* `Location` - Tenants -> Policies -> Protocol -> Route Maps for Multicast -> Route Maps

## Example Usage ##

```hcl

data "aci_pim_route_map_entry" "example" {
  parent_dn = aci_pim_route_map_policy.example.id
  order     = "1"
}

```

## Schema

### Required

* `parent_dn` - (string) The distinquised name (DN) of the parent object, possible resources:
  - [aci_pim_route_map_policy](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/pim_route_map_policy) (`pimRouteMapPol`)
* `order` - (string) PIM route map entry order.

### Read-Only

* `id` - (string) The distinquised name (DN) of the Pim Route Map Entry object.
* `action` - (string) route action.
* `annotation` - (string) The annotation of the Pim Route Map Entry object.
* `description` - (string) The description of the Pim Route Map Entry object.
* `grp` - (string) Multicast group ip/prefix.
* `name` - (string) The name of the Pim Route Map Entry object.
* `name_alias` - (string) The name alias of the Pim Route Map Entry object.
* `rp` - (string) Multicast RP Ip.
* `src` - (string) Multicast Source Ip.

* `annotations` - (list) A list of Annotations objects `tagAnnotation`.
  * `key` - (string) The key or password used to uniquely identify this configuration object.
  * `value` - (string) The value of the property.