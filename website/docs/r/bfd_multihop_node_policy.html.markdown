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
* `Distinguished Name` - uni/tn-{name}/bfdMhNodePol-{name}

## GUI Information ##

* `Location` - 


## Example Usage ##

```hcl
resource "aci_bfdmultihop_node_policy" "example" {
  tenant_dn  = aci_tenant.example.id
  name  = "example"
  admin_st = "enabled"
  annotation = "orchestrator:terraform"
  detect_mult = "3"
  min_rx_intvl = "250"
  min_tx_intvl = "250"

  name_alias = 
}
```

## Argument Reference ##

* `tenant_dn` - (Required) Distinguished name of the parent Tenant object. Type: String.
* `name` - (Required) Name of the BFD Multihop Node Policy object. Type: String.
* `annotation` - (Optional) Annotation of the BFD Multihop Node Policy object. Type: String.
* `name_alias` - (Optional) Name Alias of the BFD Multihop Node Policy object. Type: String.
* `admin_st` - (Optional) Enable Disable sessions.The administrative state of the object or policy. Allowed values are "disabled", "enabled", and default value is "enabled". Type: String.
* `detect_mult` - (Optional) Detection Multiplier. Allowed range is 1-50 and default value is "3". Type: String.
* `min_rx_intvl` - (Optional) Required Minimum RX Interval. Allowed range is 250-999 and default value is "250". Type: String.
* `min_tx_intvl` - (Optional) Desired Minimum TX Interval. Allowed range is 250-999 and default value is "250". Type: String.



## Importing ##

An existing BFDMultihopNodePolicy can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_bfdmultihop_node_policy.example <Dn>
```