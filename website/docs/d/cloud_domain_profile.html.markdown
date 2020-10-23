---
layout: "aci"
page_title: "ACI: aci_cloud_domain_profile"
sidebar_current: "docs-aci-data-source-cloud_domain_profile"
description: |-
  Data source for ACI Cloud Domain Profile
---

# aci_cloud_domain_profile #
Data source for ACI Cloud Domain Profile  
<b>Note: This resource is supported in Cloud APIC only.</b>
## Example Usage ##

```hcl
data "aci_cloud_domain_profile" "default_domp" {

}
```
## Argument Reference ##
This data source doesn't require any arguments.

## Attribute Reference

* `id` - Attribute id set to the Dn of the Cloud Domain Profile.
* `annotation` - (Optional) annotation for object cloud_domain_profile.
* `name_alias` - (Optional) name_alias for object cloud_domain_profile.
* `site_id` - (Optional) site_id for object cloud_domain_profile.
