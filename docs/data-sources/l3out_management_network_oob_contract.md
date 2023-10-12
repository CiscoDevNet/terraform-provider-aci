---
subcategory: "L3Out"
layout: "aci"
page_title: "ACI: l3out_management_network_oob_contract"
sidebar_current: "docs-aci-data-source-l3out_management_network_oob_contract"
description: |-
  Data source for L3out Management Network Oob Contract
---

# aci_l3out_management_network_oob_contract #

Data source for L3out Management Network Oob Contract

## API Information ##

* `Class` - [mgmtRsOoBCons](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/mgmtRsOoBCons/overview)

* `Distinguished Name Formats`
  - `uni/tn-{name}/extmgmt-{name}/instp-{name}/rsooBCons-{tnVzOOBBrCPName}`

## GUI Information ##

* `Location` - `Tenants (mgmt) -> External Management Network Instance Profiles -> Contracts`

## Example Usage ##

```hcl

data "aci_l3out_management_network_oob_contract" "example" {
  parent_dn                 = aci_l3out_management_network_instance_profile.example.id
  out_of_band_contract_name = "test_tn_vz_oob_br_cp_name"
}

```

## Schema

### Required

* `parent_dn` - (string) The distinguished name (DN) of the parent object, possible resources:
  - [aci_l3out_management_network_instance_profile](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/l3out_management_network_instance_profile) ([mgmtInstP](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/mgmtInstP/overview))
* `out_of_band_contract_name` - (string) An out-of-band management endpoint group contract consists of switches (leaves/spines) and APICs that are part of the associated out-of-band management zone. Each node in the group is assigned an IP address that is dynamically allocated from the address pool associated with the corresponding out-of-band management zone.

### Read-Only

* `id` - (string) The distinguished name (DN) of the L3out Management Network Oob Contract object.
* `annotation` - (string) The annotation of the L3out Management Network Oob Contract object.
* `priority` - (string) The Quality of service (QoS) priority class ID. QoS refers to the capability of a network to provide better service to selected network traffic over various technologies. The primary goal of QoS is to provide priority including dedicated bandwidth, controlled jitter and latency (required by some real-time and interactive traffic), and improved loss characteristics. You can configure the bandwidth of each QoS level using QoS profiles.

* `annotations` - (list) A list of Annotations objects ([tagAnnotation](https://pubhub.devnetcloud.com/media/model-doc-latest/docs/app/index.html#/objects/tagAnnotation/overview)).
  * `key` - (string) The key or password used to uniquely identify this configuration object.
  * `value` - (string) The value of the property.