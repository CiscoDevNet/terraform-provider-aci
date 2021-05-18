---
layout: "aci"
page_title: "ACI: aci_cloud_cidr_pool"
sidebar_current: "docs-aci-data-source-cloud_cidr_pool"
description: |-
  Data source for ACI Cloud CIDR Pool
---

# aci_cloud_cidr_pool #
Data source for ACI Cloud CIDR Pool.  
<b>Note: This resource is supported in Cloud APIC only.</b>
## Example Usage ##

```hcl
data "aci_cloud_cidr_pool" "dev_cloud_cidr" {

  cloud_context_profile_dn  = "${aci_cloud_context_profile.dev_ctx_prof.id}"
  addr  = "10.0.1.10/28"
}
```
## Argument Reference ##
* `cloud_context_profile_dn` - (Required) Distinguished name of parent CloudContextProfile object.
* `addr` - (Required) CIDR IPv4 block.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Cloud CIDR Pool.
* `annotation` - (Optional) annotation for object cloud_cidr_pool.
* `name_alias` - (Optional) name_alias for object cloud_cidr_pool.
* `primary` - (Optional) This will represent whether CIDR is primary CIDR or not.
