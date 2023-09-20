---
subcategory: "Tenant Policies"
layout: "aci"
page_title: "ACI: aci_bgp_route_summarization"
sidebar_current: "docs-aci-data-source-bgp_route_summarization"
description: |-
  Data source for ACI BGP Route Summarization
---

# aci_bgp_route_summarization

Data source for ACI BGP Route Summarization

## API Information ##

* `Class` - bgpRtSummPol
* `Distinguished Name` - uni/tn-{tenant_name}/bgprtsum-{bgp_rt_summ_pol_name}

## GUI Information ##

* `Location` - Tenants -> Policies -> Protocol -> BGP -> BGP Route Summarization

## Example Usage

```hcl
data "aci_bgp_route_summarization" "bgp_rt_summ_pol" {
  tenant_dn = aci_tenant.tf_tenant.id
  name      = "bgp_rt_summ_pol"
}
```

## Argument Reference

- `tenant_dn` - (Required) Distinguished name of the parent Tenant object. Type: String.
- `name` - (Required) Name of the BGP Route Summarization object. Type: String.

## Attribute Reference

- `id` - Attribute ID set to the Dn of the BGP Route Summarization object. Type: String.
- `annotation` - (Optional) Annotation of the BGP Route Summarization object. Type: String.
- `description` - (Optional) Description of the BGP Route Summarization object. Type: String.
- `attrmap` - (Optional) Route Map Summary of the BGP Route Summarization object. Type: String.
- `ctrl` - (Optional) Control State of the BGP Route Summarization object. Type: List.
- `address_type_controls` - (Optional) Address Type Controls of the BGP Route Summarization object. Type: List.
- `name_alias` - (Optional) Name alias of the BGP Route Summarization object. Type: String.
