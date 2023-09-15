---
subcategory: "Generic"
layout: "aci"
page_title: "ACI: annotation"
sidebar_current: "docs-aci-resource-annotation"
description: |-
  Manages ACI Annotation
---

# aci_annotation #

Manages ACI Annotation

## API Information ##

* `Class` - tagAnnotation

* `Distinguished Name Formats`
  - TaskDeployCsr/annotationKey-[{key}]
  - TaskDeployCtx/annotationKey-[{key}]
  - TaskDeploySg/annotationKey-[{key}]
  - TaskDeploySgRule/annotationKey-[{key}]
  - TaskPolicyUpdate/annotationKey-[{key}]
  - TaskResourceMap/annotationKey-[{key}]
  - acct-[{name}]/acctoper/fault-{code}/annotationKey-[{key}]
  - acct-[{name}]/adds/annotationKey-[{key}]
  - acct-[{name}]/apigw/annotationKey-[{key}]
  - acct-[{name}]/certstore-{store}/cert-[{name}]/certificateoper/fault-{code}/annotationKey-[{key}]
  - acct-[{name}]/cloudsvc-[{name}]-typ-[{type}]-rg-{rg}/adds/annotationKey-[{key}]
  - acct-[{name}]/cloudsvc-[{name}]-typ-[{type}]-rg-{rg}/apigw/annotationKey-[{key}]
  - acct-[{name}]/cloudsvc-[{name}]-typ-[{type}]-rg-{rg}/cntreg/annotationKey-[{key}]
  - acct-[{name}]/cloudsvc-[{name}]-typ-[{type}]-rg-{rg}/cosmosdb/annotationKey-[{key}]
  - acct-[{name}]/cloudsvc-[{name}]-typ-[{type}]-rg-{rg}/dbricks/annotationKey-[{key}]
  - acct-[{name}]/cloudsvc-[{name}]-typ-[{type}]-rg-{rg}/dbricks/rsbricksToPubSubnet/annotationKey-[{key}]
  - acct-[{name}]/cloudsvc-[{name}]-typ-[{type}]-rg-{rg}/dbricks/rsbricksToPubSubnet/fault-{code}/annotationKey-[{key}]
  - acct-[{name}]/cloudsvc-[{name}]-typ-[{type}]-rg-{rg}/dbricks/rsbricksToPvtSubnet/annotationKey-[{key}]
  - acct-[{name}]/cloudsvc-[{name}]-typ-[{type}]-rg-{rg}/dbricks/rsbricksToPvtSubnet/fault-{code}/annotationKey-[{key}]
  - acct-[{name}]/cloudsvc-[{name}]-typ-[{type}]-rg-{rg}/keyvault/annotationKey-[{key}]
  - Too many DN formats to display, see model documentation for all possible parents.

## GUI Information ##

* `Location` - Generic

## Example Usage ##

```hcl

resource "aci_annotation" "example" {
  parent_dn = aci_application_epg.example.id
  key       = "test_annotation"
}

```

## Schema

### Required

* `parent_dn` - (string) The distinquised name (DN) of the parent object, possible resources:
  - `aci_application_epg` (class: fvAEPg)
  - `aci_tenant` (class: fvTenant)
  - `aci_l3out_consumer_label` (class: l3extConsLbl)
  - `aci_l3_outside` (class: l3extOut)
  - `aci_l3out_redistribute_policy` (class: l3extRsRedistributePol)
  - `aci_l3out_management_network_instance_profile` (class: mgmtInstP)
  - `aci_l3out_management_network_oob_contract` (class: mgmtRsOoBCons)
  - `aci_l3out_management_network_subnet` (class: mgmtSubnet)
  - `aci_pim_route_map_entry` (class: pimRouteMapEntry)
  - `aci_pim_route_map_policy` (class: pimRouteMapPol)
  - `aci_out_of_band_contract` (class: vzOOBBrCP)
  - Too many classes to display, see model documentation for all possible classes.
* `key` - (string) The key or password used to uniquely identify this configuration object.

### Read-Only

* `id` - (string) The distinquised name (DN) of the Annotation object.

### Optional

* `value` - (string) The value of the property.


## Importing ##

An existing Annotation can be [imported](https://www.terraform.io/docs/import/index.html) into this resource via its distinquised name (DN), via the following command:

```
terraform import aci_annotation.example uni/fabric/dcswitchconnprof/rsdcProfToEpg/annotationKey-[{key}]
```

Starting in Terraform version 1.5, an existing BFD Multihop Node Policy can be imported 
using [import blocks](https://developer.hashicorp.com/terraform/language/import) via the following configuration:

```
import {
  id = "uni/fabric/dcswitchconnprof/rsdcProfToEpg/annotationKey-[{key}]"
  to = aci_annotation.example
}
```