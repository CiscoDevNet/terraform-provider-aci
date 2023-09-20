---
subcategory: "Networking"
layout: "aci"
page_title: "ACI: aci_vrf_snmp_context"
sidebar_current: "docs-aci-resource-vrf_snmp_context"
description: |-
  Manages ACI VRF SNMP Context
---

# aci_vrf_snmp_context #
Manages ACI VRF SNMP Context

## API Information ##
* `Class` - snmpCtxP
* `Distinguished Name` - uni/tn-{name}/ctx-{name}/snmpctx

## GUI Information ##
* `Location` - Tenants -> Networking -> VRF -> Policy -> Create SNMP Context


## Example Usage ##

```hcl
resource "aci_vrf_snmp_context" "example" {
  vrf_dn = aci_vrf.example.id
  name = "example"
  annotation = "orchestrator:terraform"
  name_alias = "example_name_alias"
}
```

## Argument Reference ##
* `vrf_dn` - (Required) Distinguished name of parent VRF object.
* `name` - (Required) Name of object VRF SNMP Context
* `annotation` - (Optional) Annotation of object VRF SNMP Context.
* `name_alias` - (Optional) Name Alias of object VRF SNMP Context.

## Importing ##
An existing VRF SNMP Context can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_vrf_snmp_context.example <Dn>
```

## NOTE ##
User can create only one VRF SNMP Context under one VRF.