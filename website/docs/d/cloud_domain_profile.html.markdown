---
subcategory: "Cloud"
layout: "aci"
page_title: "ACI: aci_cloud_domain_profile"
sidebar_current: "docs-aci-data-source-cloud_domain_profile"
description: |-
  Data source for ACI Cloud Domain Profile
---

# aci_cloud_domain_profile

Data source for ACI Cloud Domain Profile  
<b>Note: This resource is supported in Cloud APIC only.</b>

## Example Usage

```hcl
data "aci_cloud_domain_profile" "default_domp" {

}
```

## Argument Reference

This data source doesn't require any arguments.

## Attribute Reference

- `id` - Attribute id set to the Dn of the Cloud Domain Profile.
- `annotation` - (Optional) Specifies the Annotation of the cloud domain profile.
- `description` - (Optional) Specifies the Description of the cloud domain profile.
- `name_alias` - (Optional) Specifies the alias-name of the cloud domain profile.
- `site_id` - (Optional) Site-ID of the cloud domain profile.
