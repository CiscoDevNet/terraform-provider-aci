---
subcategory: "L3Out"
layout: "aci"
page_title: "ACI: l3out_management_network_instance_profile"
sidebar_current: "docs-aci-resource-l3out_management_network_instance_profile"
description: |-
  Manages ACI L3out Management Network Instance Profile
---

# aci_l3out_management_network_instance_profile #

Manages ACI L3out Management Network Instance Profile

## API Information ##

* `Class` - [mgmtInstP](https://pubhub.devnetcloud.com/media/model-doc-521/docs/app/index.html#/objects/mgmtInstP/overview)

* `Distinguished Name Formats`
  - `uni/tn-mgmt/extmgmt-default/instp-{name}`

## GUI Information ##

* `Location` - `Tenants (mgmt) -> External Management Network Instance Profiles`

## Example Usage ##

The configuration snippet below creates a L3out Management Network Instance Profile with only required attributes.

```hcl
resource "aci_l3out_management_network_instance_profile" "example" {
  name = "test_name"
}
```

The configuration snippet below below shows all possible attributes of the L3out Management Network Instance Profile.

!> This example might not be valid configuration and is only used to show all possible attributes.

```hcl
resource "aci_l3out_management_network_instance_profile" "example" {
  annotation  = "annotation"
  description = "description"
  name        = "test_name"
  name_alias  = "name_alias"
  priority    = "level1"
  l3out_management_network_oob_contracts = [
    {
      annotation                = "orchestrator:terraform"
      priority                  = "level1"
      out_of_band_contract_name = "l3out_management_network_oob_contracts_1"
    }
  ]
  annotations = [
    {
      key   = "annotations_1"
      value = "value_1"
    }
  ]
}
```

All examples for the L3out Management Network Instance Profile resource can be found in the [examples](https://github.com/CiscoDevNet/terraform-provider-aci/examples/resources/aci_l3out_management_network_instance_profile) folder.

## Schema

### Required

* `name` - (string) The name of the L3out Management Network Instance Profile object.

### Read-Only

* `id` - (string) The distinguished name (DN) of the L3out Management Network Instance Profile object.

### Optional
  
* `annotation` - (string) The annotation of the L3out Management Network Instance Profile object.
  - Default: `orchestrator:terraform`
* `description` - (string) The description of the L3out Management Network Instance Profile object.
* `name_alias` - (string) The name alias of the L3out Management Network Instance Profile object.
* `priority` - (string) The QoS priority class identifier.
  - Default: `unspecified`
  - Valid Values: `level1`, `level2`, `level3`, `level4`, `level5`, `level6`, `unspecified`.

* `l3out_management_network_oob_contracts` - (list) A list of L3out Management Network Oob Contracts relationship objects ([mgmtRsOoBCons](https://pubhub.devnetcloud.com/media/model-doc-521/docs/app/index.html#/objects/mgmtRsOoBCons/overview)) pointing to the Out Of Band Contract ([vzOOBBrCP](https://pubhub.devnetcloud.com/media/model-doc-521/docs/app/index.html#/objects/vzOOBBrCP/overview)) which can be configured using the [aci_out_of_band_contract](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/out_of_band_contract) resource.
  
  #### Required
  
  * `out_of_band_contract_name` - (string) An out-of-band management endpoint group contract consists of switches (leaves/spines) and APICs that are part of the associated out-of-band management zone. Each node in the group is assigned an IP address that is dynamically allocated from the address pool associated with the corresponding out-of-band management zone.

  #### Optional
    
  * `annotation` - (string) The annotation of the L3out Management Network Oob Contract object.
      - Default: `orchestrator:terraform`
  * `priority` - (string) The Quality of service (QoS) priority class ID. QoS refers to the capability of a network to provide better service to selected network traffic over various technologies. The primary goal of QoS is to provide priority including dedicated bandwidth, controlled jitter and latency (required by some real-time and interactive traffic), and improved loss characteristics. You can configure the bandwidth of each QoS level using QoS profiles.
      - Default: `unspecified`
      - Valid Values: `level1`, `level2`, `level3`, `level4`, `level5`, `level6`, `unspecified`.

* `annotations` - (list) A list of Annotations objects ([tagAnnotation](https://pubhub.devnetcloud.com/media/model-doc-521/docs/app/index.html#/objects/tagAnnotation/overview)) which can be configured using the [aci_annotation](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/annotation) resource.
  
  #### Required
  
  * `key` - (string) The key or password used to uniquely identify this configuration object.
  * `value` - (string) The value of the property.

## Importing

An existing L3out Management Network Instance Profile can be [imported](https://www.terraform.io/docs/import/index.html) into this resource via its distinguished name (DN), via the following command:

```
terraform import aci_l3out_management_network_instance_profile.example uni/tn-mgmt/extmgmt-default/instp-{name}
```

Starting in Terraform version 1.5, an existing L3out Management Network Instance Profile can be imported 
using [import blocks](https://developer.hashicorp.com/terraform/language/import) via the following configuration:

```
import {
  id = "uni/tn-mgmt/extmgmt-default/instp-{name}"
  to = aci_l3out_management_network_instance_profile.example
}
```

## Child Resources
  
  - [aci_l3out_management_network_subnet](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/l3out_management_network_subnet)