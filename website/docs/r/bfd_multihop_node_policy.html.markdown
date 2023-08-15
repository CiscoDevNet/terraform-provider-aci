---
subcategory: -
layout: "aci"
page_title: "ACI: aci_bfdmultihop_node_policy"
sidebar_current: "docs-aci-resource-bfdmultihop_node_policy"
description: |-
  Manages ACI BFD Multihop Node Policy
---

# aci_bfdmultihop_node_policy #

Manages ACI BFD Multihop Node Policy

## API Information ##

* `Class` - bfdMhNodePol
* `Distinguished Name` - uni/tn-{tenant_name}/bfdMhNodePol-{name}

## GUI Information ##

* `Location` - Tenant -> Policies -> Protocol -> BFD Multihop -> Node Policies

## Example Usage ##

```hcl
resource "aci_bfd_multihop_node_policy" "example" {
  tenant_dn            = aci_tenant.example.id
  name                 = "example"
  admin_st             = "enabled"
  annotation           = "orchestrator:terraform"
  detection_multiplier = "3"
  min_rx_interval      = "250"
  min_tx_interval      = "250"
}
```

## Argument Reference ##

* `tenant_dn` - (Required) Distinguished name of the parent Tenant object. Type: String.
* `name` - (Required) Name of the BFD Multihop Node Policy object. Type: String.
* `annotation` - (Optional) Annotation of the BFD Multihop Node Policy object. Type: String.
* `name_alias` - (Optional) Name Alias of the BFD Multihop Node Policy object. Type: String.
* `admin_state` - (Optional) Administrative state of the object or policy. Allowed values are "disabled", "enabled", and default value is "enabled". Type: String.
* `detection_multiplier` - (Optional) Detection Multiplier. Allowed range is 1-50 and default value is "3". Type: String.
* `min_rx_interval` - (Optional) Required Minimum Rx Interval. Allowed range is 250-999 and default value is "250". Type: String.
* `min_tx_interval` - (Optional) Desired Minimum Tx Interval. Allowed range is 250-999 and default value is "250". Type: String.

## Importing ##

An existing BFD Multihop Node Policy can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_bfd_multihop_node_policy.example <Dn>
```

Starting in Terraform version 1.5, an existing BFD Multihop Node Policy can be imported 
using [import blocks](https://developer.hashicorp.com/terraform/language/import) via the following configuration:

```
import {
  id = "<Dn>"
  to = aci_bfd_multihop_node_policy.example
}
```