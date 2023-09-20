---
subcategory: "Networking"
layout: "aci"
page_title: "ACI: aci_vrf_snmp_context"
sidebar_current: "docs-aci-data-source-vrf_snmp_context"
description: |-
  Data source for ACI VRF SNMP Context
---

# aci_vrf_snmp_context #
Data source for ACI VRF SNMP Context


## API Information ##
* `Class` - snmpCtxP
* `Distinguished Name` - uni/tn-{name}/ctx-{name}/snmpctx

## GUI Information ##
* `Location` - Tenants -> Networking -> VRF -> Policy -> Create SNMP Context



## Example Usage ##
```hcl
data "aci_vrf_snmp_context" "example" {
  vrf_dn  = aci_vrf.example.id
}
```

## Argument Reference ##
* `vrf_dn` - (Required) Distinguished name of parent VRF object.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the VRF SNMP Context.
* `name` - (Optional) Name of object VRF SNMP Context
* `annotation` - (Optional) Annotation of object VRF SNMP Context.
* `name_alias` - (Optional) Name Alias of object VRF SNMP Context.
