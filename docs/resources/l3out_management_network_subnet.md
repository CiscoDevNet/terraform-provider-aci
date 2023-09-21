---
subcategory: "L3Out"
layout: "aci"
page_title: "ACI: l3out_management_network_subnet"
sidebar_current: "docs-aci-resource-l3out_management_network_subnet"
description: |-
  Manages ACI L3out Management Network Subnet
---

# aci_l3out_management_network_subnet #

Manages ACI L3out Management Network Subnet

## API Information ##

* `Class` - mgmtSubnet

* `Distinguished Name Formats`
  - uni/tn-mgmt/extmgmt-default/instp-{name}/subnet-[{ip}]

## GUI Information ##

* `Location` - Tenants (mgmt) -> External Management Network Instance Profiles -> Subnets

## Example Usage ##

```hcl

resource "aci_l3out_management_network_subnet" "example" {
  parent_dn = aci_l3out_management_network_instance_profile.example.id
  ip        = "1.1.1.0/24"
}

```

## Schema

### Required

* `parent_dn` - (string) The distinquised name (DN) of the parent object, possible resources:
  - [aci_l3out_management_network_instance_profile](https://registry.terraform.io/providers/CiscoDevNet/aci/latest/docs/resources/l3out_management_network_instance_profile) (`mgmtInstP`)
* `ip` - (string) The external subnet IP address and subnet mask. This IP address is used for creating an external management entity. The subnet mask for the IP address to be imported from the outside into the fabric. The contracts associated with its parent instance profile (l3ext:InstP) are applied to the subnet.

### Read-Only

* `id` - (string) The distinquised name (DN) of the L3out Management Network Subnet object.

### Optional
  
* `annotation` - (string) The annotation of the L3out Management Network Subnet object.
  - Default: `orchestrator:terraform`
* `description` - (string) The description of the L3out Management Network Subnet object.
* `name` - (string) The name of the L3out Management Network Subnet object.
* `name_alias` - (string) The name alias of the L3out Management Network Subnet object.

## Importing ##

An existing L3out Management Network Subnet can be [imported](https://www.terraform.io/docs/import/index.html) into this resource via its distinquised name (DN), via the following command:

```
terraform import aci_l3out_management_network_subnet.example uni/tn-mgmt/extmgmt-default/instp-{name}/subnet-[{ip}]
```

Starting in Terraform version 1.5, an existing BFD Multihop Node Policy can be imported 
using [import blocks](https://developer.hashicorp.com/terraform/language/import) via the following configuration:

```
import {
  id = "uni/tn-mgmt/extmgmt-default/instp-{name}/subnet-[{ip}]"
  to = aci_l3out_management_network_subnet.example
}
```