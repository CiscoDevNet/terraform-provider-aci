---
subcategory: "Multicast"
layout: "aci"
page_title: "ACI: pim_route_map_entry"
sidebar_current: "docs-aci-resource-pim_route_map_entry"
description: |-
  Manages ACI Pim Route Map Entry
---

# aci_pim_route_map_entry #

Manages ACI Pim Route Map Entry

## API Information ##

* `Class` - [pimRouteMapEntry](https://pubhub.devnetcloud.com/media/model-doc-521/docs/app/index.html#/objects/pimRouteMapEntry/overview)

* `Distinguished Name Formats`
  - `uni/tn-{name}/rtmap-{name}/rtmapentry-{order}`

## GUI Information ##

* `Location` - `Tenants -> Policies -> Protocol -> Route Maps for Multicast -> Route Maps`

## Example Usage ##

The configuration snippet below creates a Pim Route Map Entry with only required attributes.

```hcl

resource "aci_pim_route_map_entry" "example" {
  parent_dn = aci_pim_route_map_policy.example.id
  order     = "1"
}
  ```

The configuration snippet below below shows all possible attributes of the Pim Route Map Entry.

!> This example might not be valid configuration and is only used to show all possible attributes.

```hcl

resource "aci_pim_route_map_entry" "example" {
  parent_dn   = aci_pim_route_map_policy.example.id
  action      = "deny"
  annotation  = "annotation"
  description = "description"
  grp         = "0.0.0.0"
  name        = "name"
  name_alias  = "name_alias"
  order       = "1"
  rp          = "0.0.0.0"
  src         = "1.1.1.1/30"
  annotations = [
    {
      key   = "annotations_1"
      value = "value_1"
    }
  ]
}

```

All examples for the Pim Route Map Entry resource can be found in the [examples](https://github.com/CiscoDevNet/terraform-provider-aci/examples/resources/aci_pim_route_map_entry) folder.

## Schema

### Required

* `parent_dn` - (string) The distinguished name (DN) of the parent object, possible resources:
  - [aci_pim_route_map_policy](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/pim_route_map_policy) ([pimRouteMapPol](https://pubhub.devnetcloud.com/media/model-doc-521/docs/app/index.html#/objects/pimRouteMapPol/overview))
* `order` - (string) PIM route map entry order.

### Read-Only

* `id` - (string) The distinguished name (DN) of the Pim Route Map Entry object.

### Optional
  
* `action` - (string) route action.
  - Default: `permit`
  - Valid Values: `deny`, `permit`.
* `annotation` - (string) The annotation of the Pim Route Map Entry object.
  - Default: `orchestrator:terraform`
* `description` - (string) The description of the Pim Route Map Entry object.
* `grp` - (string) Multicast group ip/prefix.
* `name` - (string) The name of the Pim Route Map Entry object.
* `name_alias` - (string) The name alias of the Pim Route Map Entry object.
* `rp` - (string) Multicast RP Ip.
* `src` - (string) Multicast Source Ip.

* `annotations` - (list) A list of Annotations objects ([tagAnnotation](https://pubhub.devnetcloud.com/media/model-doc-521/docs/app/index.html#/objects/tagAnnotation/overview)) which can be configured using the [aci_annotation](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/annotation) resource.
  
  #### Required
  
  * `key` - (string) The key or password used to uniquely identify this configuration object.
  * `value` - (string) The value of the property.

## Importing

An existing Pim Route Map Entry can be [imported](https://www.terraform.io/docs/import/index.html) into this resource via its distinguished name (DN), via the following command:

```
terraform import aci_pim_route_map_entry.example uni/tn-{name}/rtmap-{name}/rtmapentry-{order}
```

Starting in Terraform version 1.5, an existing Pim Route Map Entry can be imported 
using [import blocks](https://developer.hashicorp.com/terraform/language/import) via the following configuration:

```
import {
  id = "uni/tn-{name}/rtmap-{name}/rtmapentry-{order}"
  to = aci_pim_route_map_entry.example
}
```
