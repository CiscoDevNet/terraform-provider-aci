---
subcategory: "Cloud"
layout: "aci"
page_title: "ACI: aci_cloud_cidr_pool"
sidebar_current: "docs-aci-data-source-aci_cloud_cidr_pool"
description: |-
  Data source for Cloud Network Controller Cloud CIDR Pool
---

# aci_cloud_cidr_pool #
Data source for Cloud Network Controller Cloud CIDR Pool.  
<b>Note: This resource is supported in Cloud Network Controller only.</b>
## Example Usage ##

```hcl
data "aci_cloud_cidr_pool" "dev_cloud_cidr" {

  cloud_context_profile_dn  = aci_cloud_context_profile.dev_ctx_prof.id
  addr  = "10.0.1.10/28"
}
```
## Argument Reference ##
* `cloud_context_profile_dn` - (Required) Distinguished name of parent CloudContextProfile object.
* `addr` - (Required) CIDR IPv4 block.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Cloud CIDR Pool.
* `description` - (Optional) Description for object Cloud CIDR Pool.
* `annotation` - (Optional) Annotation for object Cloud CIDR Pool.
* `name_alias` - (Optional) Name alias for object Cloud CIDR Pool.
* `primary` - (Optional) This will represent whether CIDR is primary CIDR or not.
