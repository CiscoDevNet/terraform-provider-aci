---
subcategory: "Multicast"
layout: "aci"
page_title: "ACI: pim_route_map_policy"
sidebar_current: "docs-aci-resource-pim_route_map_policy"
description: |-
  Manages ACI Pim Route Map Policy
---

# aci_pim_route_map_policy #

Manages ACI Pim Route Map Policy

## API Information ##

* `Class` - pimRouteMapPol

* `Distinguished Name Formats`
  - uni/tn-{name}/rtmap-{name}

## GUI Information ##

* `Location` - Tenants -> Policies -> Protocol -> Route Maps for Multicast

## Example Usage ##

```hcl

resource "aci_pim_route_map_policy" "example" {
  parent_dn = aci_tenant.example.id
  name      = "test_name"
  annotations = [
    {
      key   = "test_key"
      value = "test_value"
    },
  ]
}

```

## Schema

### Required

* `parent_dn` - (string) The distinquised name (DN) of the parent object, possible resources:
  - `aci_tenant` (class: fvTenant)
* `name` - (string) The name of the Pim Route Map Policy object.

### Read-Only

* `id` - (string) The distinquised name (DN) of the Pim Route Map Policy object.

### Optional
  
* `annotation` - (string) The annotation of the Pim Route Map Policy object.
  - Default: `orchestrator:terraform`
* `description` - (string) The description of the Pim Route Map Policy object.
* `name_alias` - (string) The name alias of the Pim Route Map Policy object.
* `owner_key` - (string) The key for enabling clients to own their data for entity correlation.
* `owner_tag` - (string) A tag for enabling clients to add their own data. For example, to indicate who created this object.

* `annotations` - (list) A list of Annotation objects (tagAnnotation) which can also be configured using the `aci_annotation` resource.
    
  #### Required
  
  * `key` - (string) The key or password used to uniquely identify this configuration object.
  * `value` - (string) The value of the property.

## Importing ##

An existing Pim Route Map Policy can be [imported](https://www.terraform.io/docs/import/index.html) into this resource via its distinquised name (DN), via the following command:

```
terraform import aci_pim_route_map_policy.example uni/tn-{name}/rtmap-{name}
```

Starting in Terraform version 1.5, an existing BFD Multihop Node Policy can be imported 
using [import blocks](https://developer.hashicorp.com/terraform/language/import) via the following configuration:

```
import {
  id = "uni/tn-{name}/rtmap-{name}"
  to = aci_pim_route_map_policy.example
}
```