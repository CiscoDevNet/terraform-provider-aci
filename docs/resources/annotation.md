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

* `Class` - [tagAnnotation](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/tagAnnotation/overview)

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

The configuration snippet below creates a Annotation with only required attributes.

```hcl

resource "aci_annotation" "example" {
  parent_dn = aci_application_epg.example.id
  key       = "test_key"
  value     = "test_value"
}
  ```

The configuration snippet below below shows all possible attributes of the Annotation.

!> This example might not be valid configuration and is only used to show all possible attributes.

```hcl

resource "aci_annotation" "example" {
  parent_dn = aci_application_epg.example.id
  key       = "test_key"
  value     = "test_value"
}

```

All examples for the Annotation resource can be found in the [examples](https://github.com/CiscoDevNet/terraform-provider-aci/examples/resources/aci_annotation) folder.

## Schema

### Required

* `parent_dn` - (string) The distinguished name (DN) of the parent object, possible resources:
  - [aci_application_epg](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/application_epg) ([fvAEPg](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/fvAEPg/overview))
  - [aci_contract_interface](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/contract_interface) ([fvRsConsIf](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/fvRsConsIf/overview))
  - [aci_tenant](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/tenant) ([fvTenant](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/fvTenant/overview))
  - [aci_l3out_consumer_label](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/l3out_consumer_label) ([l3extConsLbl](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/l3extConsLbl/overview))
  - [aci_l3_outside](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/l3_outside) ([l3extOut](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/l3extOut/overview))
  - [aci_l3out_redistribute_policy](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/l3out_redistribute_policy) ([l3extRsRedistributePol](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/l3extRsRedistributePol/overview))
  - [aci_l3out_management_network_instance_profile](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/l3out_management_network_instance_profile) ([mgmtInstP](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/mgmtInstP/overview))
  - [aci_l3out_management_network_oob_contract](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/l3out_management_network_oob_contract) ([mgmtRsOoBCons](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/mgmtRsOoBCons/overview))
  - [aci_l3out_management_network_subnet](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/l3out_management_network_subnet) ([mgmtSubnet](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/mgmtSubnet/overview))
  - [aci_pim_route_map_entry](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/pim_route_map_entry) ([pimRouteMapEntry](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/pimRouteMapEntry/overview))
  - [aci_pim_route_map_policy](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/pim_route_map_policy) ([pimRouteMapPol](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/pimRouteMapPol/overview))
  - [aci_route_control_profile](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/route_control_profile) ([rtctrlProfile](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/rtctrlProfile/overview))
  - [aci_contract_interface](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/contract_interface) ([vzCPIf](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/vzCPIf/overview))
  - [aci_out_of_band_contract](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/out_of_band_contract) ([vzOOBBrCP](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/vzOOBBrCP/overview))
  - Too many classes to display, see model documentation for all possible classes.
* `key` - (string) The key or password used to uniquely identify this configuration object.
* `value` - (string) The value of the property.

### Read-Only

* `id` - (string) The distinguished name (DN) of the Annotation object.

## Importing

An existing Annotation can be [imported](https://www.terraform.io/docs/import/index.html) into this resource via its distinguished name (DN), via the following command:

```
terraform import aci_annotation.example uni/fabric/dcswitchconnprof/rsdcProfToEpg/annotationKey-[{key}]
```

Starting in Terraform version 1.5, an existing Annotation can be imported 
using [import blocks](https://developer.hashicorp.com/terraform/language/import) via the following configuration:

```
import {
  id = "uni/fabric/dcswitchconnprof/rsdcProfToEpg/annotationKey-[{key}]"
  to = aci_annotation.example
}
```
