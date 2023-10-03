---
subcategory: "L3Out"
layout: "aci"
page_title: "ACI: l3out_management_network_oob_contract"
sidebar_current: "docs-aci-resource-l3out_management_network_oob_contract"
description: |-
  Manages ACI L3out Management Network Oob Contract
---

# aci_l3out_management_network_oob_contract #

Manages ACI L3out Management Network Oob Contract

## API Information ##

* `Class` - `mgmtRsOoBCons`

* `Distinguished Name Formats`
  - `uni/tn-{name}/extmgmt-{name}/instp-{name}/rsooBCons-{tnVzOOBBrCPName}`

## GUI Information ##

* `Location` - `Tenants (mgmt) -> External Management Network Instance Profiles -> Contracts`

## Example Usage ##

```hcl

resource "aci_l3out_management_network_oob_contract" "example" {
  parent_dn                 = aci_l3out_management_network_instance_profile.example.id
  out_of_band_contract_name = "test_tn_vz_oob_br_cp_name"
}

```

## Schema

### Required

* `parent_dn` - (string) The distinguished name (DN) of the parent object, possible resources:
  - [aci_l3out_management_network_instance_profile](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/l3out_management_network_instance_profile) (`mgmtInstP`)
* `out_of_band_contract_name` - (string) An out-of-band management endpoint group contract consists of switches (leaves/spines) and APICs that are part of the associated out-of-band management zone. Each node in the group is assigned an IP address that is dynamically allocated from the address pool associated with the corresponding out-of-band management zone.

### Read-Only

* `id` - (string) The distinguished name (DN) of the L3out Management Network Oob Contract object.

### Optional
  
* `annotation` - (string) The annotation of the L3out Management Network Oob Contract object.
  - Default: `orchestrator:terraform`
* `priority` - (string) The Quality of service (QoS) priority class ID. QoS refers to the capability of a network to provide better service to selected network traffic over various technologies. The primary goal of QoS is to provide priority including dedicated bandwidth, controlled jitter and latency (required by some real-time and interactive traffic), and improved loss characteristics. You can configure the bandwidth of each QoS level using QoS profiles.
  - Default: `unspecified`
  - Valid Values: `level1`, `level2`, `level3`, `level4`, `level5`, `level6`, `unspecified`.

## Importing

An existing L3out Management Network Oob Contract can be [imported](https://www.terraform.io/docs/import/index.html) into this resource via its distinguished name (DN), via the following command:

```
terraform import aci_l3out_management_network_oob_contract.example uni/tn-{name}/extmgmt-{name}/instp-{name}/rsooBCons-{tnVzOOBBrCPName}
```

Starting in Terraform version 1.5, an existing L3out Management Network Oob Contract can be imported 
using [import blocks](https://developer.hashicorp.com/terraform/language/import) via the following configuration:

```
import {
  id = "uni/tn-{name}/extmgmt-{name}/instp-{name}/rsooBCons-{tnVzOOBBrCPName}"
  to = aci_l3out_management_network_oob_contract.example
}
```

## Child Resources
  
  - [aci_annotation](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/annotation)