---
layout: "aci"
page_title: "ACI: aci_cloud_subnet"
sidebar_current: "docs-aci-data-source-cloud_subnet"
description: |-
  Data source for ACI Cloud Subnet
---

# aci_cloud_subnet #
Data source for ACI Cloud Subnet  
<b>Note: This resource is supported in Cloud APIC only.</b>
## Example Usage ##

```hcl

data "aci_cloud_subnet" "dev_subnet" {
  cloud_cidr_pool_dn  = aci_cloud_cidr_pool.dev_cidr_pool.id
  ip                  = "14.12.0.0/28"
}

```


## Argument Reference ##
* `cloud_cidr_pool_dn` - (Required) Distinguished name of parent CloudCIDRPool object.
* `ip` - (Required) CIDR block of Object cloud subnet.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Cloud Subnet.
* `name` - (Optional) Name for object cloud subnet.
* `description` - (Optional) Description for object cloud subnet.
* `annotation` - (Optional) Annotation for object cloud subnet.
* `name_alias` - (Optional) Name alias for object cloud subnet.
* `scope` - (Optional) The domain applicable to the capability.
* `usage` - (Optional) The usage of the port. This property shows how the port is used.
