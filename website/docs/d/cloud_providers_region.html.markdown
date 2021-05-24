---
layout: "aci"
page_title: "ACI: aci_cloud_providers_region"
sidebar_current: "docs-aci-data-source-cloud_providers_region"
description: |-
  Data source for ACI Cloud Providers Region
---

# aci_cloud_providers_region #
Data source for ACI Cloud Providers Region  
<b>Note: This resource is supported in Cloud APIC only.</b>
## Example Usage ##

```hcl
data "aci_cloud_providers_region" "region_aws" {
  cloud_provider_profile_dn  = aci_cloud_provider_profile.aws_prov.id
  name                       = "us-east-1"
}
```
## Argument Reference ##
* `cloud_provider_profile_dn` - (Required) Distinguished name of parent CloudProviderProfile object.
* `name` - (Required) Name of Object cloud providers region.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Cloud Providers Region.
* `admin_st` - (Optional) Administrative state of the object or policy.
* `description` - (Optional) Description for object cloud providers region.
* `annotation` - (Optional) Annotation for object cloud providers region.
* `name_alias` - (Optional) Name alias for object cloud providers region.
