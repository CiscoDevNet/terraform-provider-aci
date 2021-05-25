---
layout: "aci"
page_title: "ACI: aci_cloud_context_profile"
sidebar_current: "docs-aci-data-source-cloud_context_profile"
description: |-
  Data source for ACI Cloud Context Profile
---

# aci_cloud_context_profile

Data source for ACI Cloud Context Profile
<b>Note: This resource is supported in Cloud APIC only. </b>

## Example Usage

```hcl

data "aci_cloud_context_profile" "sample_prof" {
  tenant_dn  = "${aci_tenant.dev_tenant.id}"
  name       = "demo_cloud_ctx_prof"
}

```

## Argument Reference

- `tenant_dn` - (Required) Distinguished name of parent Tenant object.
- `name` - (Required) Name of Object Cloud Context profile.

## Attribute Reference

- `id` - Attribute id set to the Dn of the Cloud Context profile.
- `description` - Description of object Cloud Context profile.
- `annotation` - Annotation for object Cloud Context profile.
- `name_alias` - Name alias for object Cloud Context Profile.
- `type` - The specific type of the object or component.
