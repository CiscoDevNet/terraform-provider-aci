---
subcategory: "L3Out"
layout: "aci"
page_title: "ACI: l3out_management_network_subnet"
sidebar_current: "docs-aci-data-source-l3out_management_network_subnet"
description: |-
  Data source for L3out Management Network Subnet
---

# aci_l3out_management_network_subnet #

Data source for L3out Management Network Subnet

## API Information ##

* `Class` - `mgmtSubnet`

* `Distinguished Name Formats`
  - `uni/tn-mgmt/extmgmt-default/instp-{name}/subnet-[{ip}]`

## GUI Information ##

* `Location` - `Tenants (mgmt) -> External Management Network Instance Profiles -> Subnets`

## Example Usage ##

```hcl

data "aci_l3out_management_network_subnet" "example" {
  parent_dn = aci_l3out_management_network_instance_profile.example.id
  ip        = "1.1.1.0/24"
}

```

## Schema

### Required

* `parent_dn` - (string) The distinguished name (DN) of the parent object, possible resources:
  - [aci_l3out_management_network_instance_profile](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/l3out_management_network_instance_profile) (`mgmtInstP`)
* `ip` - (string) The external subnet IP address and subnet mask. This IP address is used for creating an external management entity. The subnet mask for the IP address to be imported from the outside into the fabric. The contracts associated with its parent instance profile (l3ext:InstP) are applied to the subnet.

### Read-Only

* `id` - (string) The distinguished name (DN) of the L3out Management Network Subnet object.
* `annotation` - (string) The annotation of the L3out Management Network Subnet object.
* `description` - (string) The description of the L3out Management Network Subnet object.
* `name` - (string) The name of the L3out Management Network Subnet object.
* `name_alias` - (string) The name alias of the L3out Management Network Subnet object.

* `annotations` - (list) A list of Annotations objects `tagAnnotation`.
  * `key` - (string) The key or password used to uniquely identify this configuration object.
  * `value` - (string) The value of the property.