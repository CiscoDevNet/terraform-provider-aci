---
layout: "aci"
page_title: "ACI: aci_service_redirect_policy"
sidebar_current: "docs-aci-data-source-service_redirect_policy"
description: |-
  Data source for ACI Service Redirect Policy
---

# aci_service_redirect_policy #
Data source for ACI Service Redirect Policy

## Example Usage ##

```hcl

data "aci_service_redirect_policy" "example" {
  tenant_dn   = "${aci_tenant.example.id}"
  name        = "example"
}

```


## Argument Reference ##
* `tenant_dn` - (Required) Distinguished name of parent Tenant object.
* `name` - (Required) name of Object service_redirect_policy.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Service Redirect Policy.
* `anycast_enabled` - (Optional) anycast_enabled for object service_redirect_policy.
* `annotation` - (Optional) annotation for object service_redirect_policy.
* `dest_type` - (Optional) dest_type for object service_redirect_policy.
* `hashing_algorithm` - (Optional) hashing_algorithm for object service_redirect_policy.
* `max_threshold_percent` - (Optional) max_threshold_percent for object service_redirect_policy.
* `min_threshold_percent` - (Optional) min_threshold_percent for object service_redirect_policy.
* `name_alias` - (Optional) name_alias for object service_redirect_policy.
* `program_local_pod_only` - (Optional) program_local_pod_only for object service_redirect_policy.
* `resilient_hash_enabled` - (Optional) resilient_hash_enabled for object service_redirect_policy.
* `threshold_down_action` - (Optional) threshold_down_action for object service_redirect_policy.
* `threshold_enable` - (Optional) threshold_enable for object service_redirect_policy.
