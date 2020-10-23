---
layout: "aci"
page_title: "ACI: aci_cloud_provider_profile"
sidebar_current: "docs-aci-data-source-cloud_provider_profile"
description: |-
  Data source for ACI Cloud Provider Profile
---

# aci_cloud_provider_profile #
Data source for ACI Cloud Provider Profile  
<b>Note: This resource is supported in Cloud APIC only.</b>
## Example Usage ##

```hcl
data "aci_cloud_provider_profile" "aws_prof" {
  vendor  = "aws"
}
```
## Argument Reference ##
* `vendor` - (Required) vendor of Object cloud_provider_profile.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Cloud Provider Profile.
* `annotation` - (Optional) annotation for object cloud_provider_profile.
