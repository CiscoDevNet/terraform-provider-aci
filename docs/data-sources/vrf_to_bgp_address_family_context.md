---
subcategory: "Networking"
layout: "aci"
page_title: "ACI: aci_vrf_to_bgp_address_family_context"
sidebar_current: "docs-aci-data-source-vrf_to_bgp_address_family_context"
description: |-
  Data source for the ACI Relationship object between VRF and the BGP Address Family Context Policy
---

# aci_vrf_to_bgp_address_family_context #

Data source for the ACI Relationship object between VRF and the BGP Address Family Context Policy

## API Information ##

* `Class` - fvRsCtxToBgpCtxAfPol
* `Distinguished Name` - uni/tn-{name}/ctx-{name}/rsctxToBgpCtxAfPol-[{tnBgpCtxAfPolName}]-{af}

## GUI Information ##

* `Location` - Tenant -> {tenant_name} -> Networking -> VRFs -> {vrf_name} -> Policy -> BGP Context per Address Family

## Example Usage ##

```hcl
data "aci_vrf_to_bgp_address_family_context" "example" {
  vrf_dn  = aci_vrf.example.id
  bgp_address_family_context_dn  = aci_bgp_address_family_context.example.id
  address_family  = "ipv4-ucast"
}
```

## Argument Reference ##

* `vrf_dn` - (Required) Distinguished name (DN) of the parent VRF object.
* `bgp_address_family_context_dn` - (Required) Distinguished name (DN) of the BGP  address family context policy object.
* `address_family` - (Required) The BGP address family. Allowed values are "ipv4-ucast", "ipv6-ucast", and default value is "ipv4-ucast". Type: String.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the BGP address family context policy relationship.
* `annotation` - Annotation of object BGP address family context policy relationship.
