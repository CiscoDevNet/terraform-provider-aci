---
subcategory: "Tenant Policies"
layout: "aci"
page_title: "ACI: aci_bgp_route_summarization"
sidebar_current: "docs-aci-resource-aci_bgp_route_summarization"
description: |-
  Manages ACI BGP Route Summarization
---

# aci_bgp_route_summarization

Manages ACI BGP Route Summarization

## API Information ##

* `Class` - bgpRtSummPol
* `Distinguished Name` - uni/tn-{tenant_name}/bgprtsum-{bgp_rt_summ_pol_name}

## GUI Information ##

* `Location` - Tenants -> Policies -> Protocol -> BGP -> BGP Route Summarization

## Example Usage

```hcl
resource "aci_bgp_route_summarization" "bgp_rt_summ_pol" {
  tenant_dn             = aci_tenant.tf_tenant.id
  name                  = "bgp_rt_summ_pol"
  description           = "from terraform"
  attrmap               = "sample attrmap"
  ctrl                  = ["summary-only", "as-set"]
  address_type_controls = ["af-ucast", "af-mcast", "af-label-ucast"]
}
```

## Argument Reference

- `tenant_dn` - (Required) Distinguished name of the parent Tenant object. Type: String.
- `name` - (Required) Name of the BGP Route Summarization object. Type: String.
- `annotation` - (Optional) Annotation of the BGP Route Summarization object. Type: String.
- `description` - (Optional) Description of the BGP Route Summarization object. Type: String.
- `attrmap` - (Optional) Route Map Summary of the BGP Route Summarization object. Type: String.
- `ctrl` - (Optional) Control State of the BGP Route Summarization object. Allowed values are "as-set", "summary-only". Type: List.
- `address_type_controls` - (Optional) Address Type Controls of the BGP Route Summarization object. Allowed values are "af-ucast", "af-mcast", "af-label-ucast". The APIC defaults to "af-ucast" when unset during creation. Type: List.
- `name_alias` - (Optional) Name alias of the BGP Route Summarization object. Type: String.

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the BGP Route Summarization.

## Importing

An existing BGP Route Summarization can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_bgp_route_summarization.example <Dn>
```
