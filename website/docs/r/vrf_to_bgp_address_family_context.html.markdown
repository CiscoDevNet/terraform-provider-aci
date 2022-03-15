---
subcategory: "Networking"
layout: "aci"
page_title: "ACI: aci_vrf_to_bgp_address_family_context"
sidebar_current: "docs-aci-resource-vrf_to_bgp_address_family_context"
description: |-
  Manages the ACI Relationship object between VRF and the BGP Address Family Context Policy
---

# aci_vrf_to_bgp_address_family_context #

Manages the ACI Relationship object between VRF and the BGP Address Family Context Policy

## API Information ##

* `Class` - fvRsCtxToBgpCtxAfPol
* `Distinguished Name` - uni/tn-{name}/ctx-{name}/rsctxToBgpCtxAfPol-[{tnBgpCtxAfPolName}]-{af}

## GUI Information ##

* `Location` - Tenant -> {tenant_name} -> Networking -> VRFs -> {vrf_name} -> Policy -> BGP Context per Address Family

## Example Usage ##

```hcl
resource "aci_vrf_to_bgp_address_family_context" "example" {
  vrf_dn  = aci_vrf.example.id
  bgp_address_family_context_dn = aci_bgp_address_family_context.example.id
  address_family  = "ipv4-ucast"
}
```

## Argument Reference ##

* `vrf_dn` - (Required) Distinguished name (DN) of the parent VRF object.
* `bgp_address_family_context_dn` - (Required) Distinguished name (DN) of the BGP address family context policy object.
* `address_family` - (Required) The BGP address family. Allowed values are "ipv4-ucast", "ipv6-ucast", and default value is "ipv4-ucast". Type: String.
* `annotation` - (Optional) Annotation of the VRF to BGP address family context policy relationship object.

## Importing ##

An existing VRF to BGP address family context policy relationship object can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_vrf_to_bgp_address_family_context.example <Dn>
```