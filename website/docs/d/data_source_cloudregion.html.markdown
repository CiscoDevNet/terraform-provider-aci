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

  cloud_provider_profile_dn  = "${aci_cloud_provider_profile.aws_prov.id}"
  name                       = "us-east-1"
}
```
## Argument Reference ##
* `cloud_provider_profile_dn` - (Required) Distinguished name of parent CloudProviderProfile object.
* `name` - (Required) name of Object cloud_providers_region.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Cloud Providers Region.
* `admin_st` - (Optional) administrative state of the object or policy
* `annotation` - (Optional) annotation for object cloud_providers_region.
* `name_alias` - (Optional) name_alias for object cloud_providers_region.
