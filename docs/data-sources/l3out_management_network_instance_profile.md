---
subcategory: "L3Out"
layout: "aci"
page_title: "ACI: l3out_management_network_instance_profile"
sidebar_current: "docs-aci-data-source-l3out_management_network_instance_profile"
description: |-
  Data source for L3out Management Network Instance Profile
---

# aci_l3out_management_network_instance_profile #

Data source for L3out Management Network Instance Profile

## API Information ##

* `Class` - `mgmtInstP`

* `Distinguished Name Formats`
  - `uni/tn-mgmt/extmgmt-default/instp-{name}`

## GUI Information ##

* `Location` - `Tenants (mgmt) -> External Management Network Instance Profiles`

## Example Usage ##

```hcl
data "aci_l3out_management_network_instance_profile" "example" {
  name = "test_name"
}
```

## Schema

### Required

* `name` - (string) The name of the L3out Management Network Instance Profile object.

### Read-Only

* `id` - (string) The distinguished name (DN) of the L3out Management Network Instance Profile object.
* `annotation` - (string) The annotation of the L3out Management Network Instance Profile object.
* `description` - (string) The description of the L3out Management Network Instance Profile object.
* `name_alias` - (string) The name alias of the L3out Management Network Instance Profile object.
* `priority` - (string) The QoS priority class identifier.

* `l3out_management_network_oob_contracts` - (list) A list of L3out Management Network Oob Contracts relationship objects `mgmtRsOoBCons` pointing to the Out Of Band Contract `vzOOBBrCP` object.
  * `annotation` - (string) The annotation of the L3out Management Network Oob Contract object.
  * `priority` - (string) The Quality of service (QoS) priority class ID. QoS refers to the capability of a network to provide better service to selected network traffic over various technologies. The primary goal of QoS is to provide priority including dedicated bandwidth, controlled jitter and latency (required by some real-time and interactive traffic), and improved loss characteristics. You can configure the bandwidth of each QoS level using QoS profiles.
  * `out_of_band_contract_name` - (string) An out-of-band management endpoint group contract consists of switches (leaves/spines) and APICs that are part of the associated out-of-band management zone. Each node in the group is assigned an IP address that is dynamically allocated from the address pool associated with the corresponding out-of-band management zone.