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
  cloud_context_profile_dn  = aci_cloud_context_profile.example.id
  name  = "example"
  description = "from terraform"
}

```
## Argument Reference ##
* `cloud_context_profile_dn` - (Required) Distinguished name of parent CloudContextProfile object.
* `name` - (Required) Name of Object Cloud Router Profile.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Cloud Router Profile.
* `description` - (Optional) Description for object Cloud Router Profile.
* `annotation` - (Optional) Annotation for object Cloud Router Profile.
* `name_alias` - (Optional) Name alias for object Cloud Router Profile.
* `num_instances` - (Optional) Num instances for object Cloud Router Profile.
* `cloud_router_profile_type` - (Optional) Component type
