---
layout: "aci"
page_title: "ACI: aci_service_redirect_policy"
sidebar_current: "docs-aci-data-source-service_redirect_policy"
description: |-
  Data source for ACI Service Redirect Policy
---

# aci_service_redirect_policy

Data source for ACI Service Redirect Policy

## Example Usage

```hcl

data "aci_service_redirect_policy" "example" {
  tenant_dn   = aci_tenant.example.id
  name        = "example"
}

```

## Argument Reference

- `tenant_dn` - (Required) Distinguished name of parent Tenant object.
- `name` - (Required) Name of Object Service Redirect Policy.

## Attribute Reference

- `id` - Attribute id set to the Dn of the Service Redirect Policy.
- `anycast_enabled` - (Optional) Anycast enabled for object Service Redirect Policy.
- `annotation` - (Optional) annotation for object Service Redirect Policy.
- `description` - (Optional) Description of object Service Redirect Policy.
- `dest_type` - (Optional) Dest type for object Service Redirect Policy.
- `hashing_algorithm` - (Optional) Hashing algorithm for object Service Redirect Policy.
- `max_threshold_percent` - (Optional) Max threshold percent for object Service Redirect Policy.
- `min_threshold_percent` - (Optional) Min threshold percent for object Service Redirect Policy.
- `name_alias` - (Optional) Name alias for object Service Redirect Policy.
- `program_local_pod_only` - (Optional) Program local pod only for object Service Redirect Policy.
- `resilient_hash_enabled` - (Optional) Resilient hash enabled for object Service Redirect Policy.
- `threshold_down_action` - (Optional) Threshold down the action for object Service Redirect Policy.
- `threshold_enable` - (Optional) Threshold enable for object Service Redirect Policy.
