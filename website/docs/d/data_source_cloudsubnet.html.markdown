---
layout: "aci"
page_title: "ACI: aci_cloud_subnet"
sidebar_current: "docs-aci-data-source-cloud_subnet"
description: |-
  Data source for ACI Cloud Subnet
---

# aci_cloud_subnet #
Data source for ACI Cloud Subnet
Note: This resource is supported in Cloud APIC only.
## Example Usage ##

```hcl
data "aci_cloud_subnet" "example" {

  cloud_cidr_pool_dn  = "${aci_cloud_cidr_pool.example.id}"

  ip  = "example"
}
```
## Argument Reference ##
* `cloud_cidr_pool_dn` - (Required) Distinguished name of parent CloudCIDRPool object.
* `ip` - (Required) ip of Object cloud_subnet.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Cloud Subnet.
* `annotation` - (Optional) annotation for object cloud_subnet.
* `ip` - (Optional) ip address
* `name_alias` - (Optional) name_alias for object cloud_subnet.
* `scope` - (Optional) capability domain
* `usage` - (Optional) usage of the port
