---
subcategory: "Generic"
layout: "aci"
page_title: "ACI: annotation"
sidebar_current: "docs-aci-data-source-annotation"
description: |-
  Data source for Annotation
---

# aci_annotation #

Data source for Annotation

## API Information ##

* `Class` - `tagAnnotation`

* `Distinguished Name Formats`
  - `TaskDeployCsr/annotationKey-[{key}]`
  - `TaskDeployCtx/annotationKey-[{key}]`
  - `TaskDeploySg/annotationKey-[{key}]`
  - `TaskDeploySgRule/annotationKey-[{key}]`
  - `TaskPolicyUpdate/annotationKey-[{key}]`
  - `TaskResourceMap/annotationKey-[{key}]`
  - `acct-[{name}]/acctoper/fault-{code}/annotationKey-[{key}]`
  - `acct-[{name}]/adds/annotationKey-[{key}]`
  - `acct-[{name}]/apigw/annotationKey-[{key}]`
  - `acct-[{name}]/certstore-{store}/cert-[{name}]/certificateoper/fault-{code}/annotationKey-[{key}]`
  - `acct-[{name}]/cloudsvc-[{name}]-typ-[{type}]-rg-{rg}/adds/annotationKey-[{key}]`
  - `acct-[{name}]/cloudsvc-[{name}]-typ-[{type}]-rg-{rg}/apigw/annotationKey-[{key}]`
  - `acct-[{name}]/cloudsvc-[{name}]-typ-[{type}]-rg-{rg}/cntreg/annotationKey-[{key}]`
  - `acct-[{name}]/cloudsvc-[{name}]-typ-[{type}]-rg-{rg}/cosmosdb/annotationKey-[{key}]`
  - `acct-[{name}]/cloudsvc-[{name}]-typ-[{type}]-rg-{rg}/dbricks/annotationKey-[{key}]`
  - `acct-[{name}]/cloudsvc-[{name}]-typ-[{type}]-rg-{rg}/dbricks/rsbricksToPubSubnet/annotationKey-[{key}]`
  - `acct-[{name}]/cloudsvc-[{name}]-typ-[{type}]-rg-{rg}/dbricks/rsbricksToPubSubnet/fault-{code}/annotationKey-[{key}]`
  - `acct-[{name}]/cloudsvc-[{name}]-typ-[{type}]-rg-{rg}/dbricks/rsbricksToPvtSubnet/annotationKey-[{key}]`
  - `acct-[{name}]/cloudsvc-[{name}]-typ-[{type}]-rg-{rg}/dbricks/rsbricksToPvtSubnet/fault-{code}/annotationKey-[{key}]`
  - `acct-[{name}]/cloudsvc-[{name}]-typ-[{type}]-rg-{rg}/keyvault/annotationKey-[{key}]`
  - `Too many DN formats to display, see model documentation for all possible parents.`

## GUI Information ##

* `Location` - `Generic`

## Example Usage ##

```hcl

data "aci_annotation" "example" {
  parent_dn = aci_application_epg.example.id
  key       = "test_key"
}

```

## Schema

### Required

* `parent_dn` - (string) The distinquised name (DN) of the parent object, possible resources:
  - [aci_application_epg](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/application_epg) (`fvAEPg`)
  - [aci_contract_interface](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/contract_interface) (`fvRsConsIf`)
  - [aci_tenant](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/tenant) (`fvTenant`)
  - [aci_l3out_consumer_label](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/l3out_consumer_label) (`l3extConsLbl`)
  - [aci_l3_outside](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/l3_outside) (`l3extOut`)
  - [aci_l3out_redistribute_policy](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/l3out_redistribute_policy) (`l3extRsRedistributePol`)
  - [aci_l3out_management_network_instance_profile](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/l3out_management_network_instance_profile) (`mgmtInstP`)
  - [aci_l3out_management_network_oob_contract](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/l3out_management_network_oob_contract) (`mgmtRsOoBCons`)
  - [aci_l3out_management_network_subnet](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/l3out_management_network_subnet) (`mgmtSubnet`)
  - [aci_pim_route_map_entry](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/pim_route_map_entry) (`pimRouteMapEntry`)
  - [aci_pim_route_map_policy](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/pim_route_map_policy) (`pimRouteMapPol`)
  - [aci_route_control_profile](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/route_control_profile) (`rtctrlProfile`)
  - [aci_contract_interface](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/contract_interface) (`vzCPIf`)
  - [aci_out_of_band_contract](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/out_of_band_contract) (`vzOOBBrCP`)
  - Too many classes to display, see model documentation for all possible classes.
* `key` - (string) The key or password used to uniquely identify this configuration object.

### Read-Only

* `id` - (string) The distinquised name (DN) of the Annotation object.
* `value` - (string) The value of the property.
