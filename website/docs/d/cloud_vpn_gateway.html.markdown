---
layout: "aci"
page_title: "ACI: aci_cloud_vpn_gateway"
sidebar_current: "docs-aci-data-source-cloud_vpn_gateway"
description: |-
  Data source for ACI Cloud Vpn Gateway
---

# aci_cloud_vpn_gateway #
Data source for ACI Cloud Vpn Gateway.
Note: This resource is supported in Cloud APIC only.
## Example Usage ##

```hcl
 data "aci_cloud_vpn_gateway" "example" {

  cloud_context_profile_dn  = "${aci_cloud_context_profile.example.id}"

  name  = "example"
}

```
## Argument Reference ##
* `cloud_context_profile_dn` - (Required) Distinguished name of parent CloudContextProfile object.
* `name` - (Required) Name of Object cloud_router_profile.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Cloud Router Profile.
* `annotation` - (Optional) Annotation for object cloud_router_profile.
* `name_alias` - (Optional) name_alias for object cloud_router_profile.
* `num_instances` - (Optional) num_instances for object cloud_router_profile.
* `cloud_router_profile_type` - (Optional) Component type
