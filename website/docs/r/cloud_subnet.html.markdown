---
subcategory: "Cloud"
layout: "aci"
page_title: "ACI: aci_cloud_subnet"
sidebar_current: "docs-aci-resource-cloud_subnet"
description: |-
  Manages ACI Cloud Subnet
---

# aci_cloud_subnet #
Manages ACI Cloud Subnet
<b>Note: This resource is supported in Cloud APIC only.</b>
## Example Usage ##

```hcl

	resource "aci_cloud_subnet" "foocloud_subnet" {
		cloud_cidr_pool_dn = aci_cloud_cidr_pool.example.id
		description        = "sample cloud subnet"
		name			   = "subnet1"	
		ip                 = "14.12.0.0/28"
		annotation         = "tag_subnet"
		name_alias         = "alias_subnet"
		scope              = "public"
		usage              = "user"
		zone               = data.aci_cloud_availability_zone.az_us_east_1_aws.id
	}

```


## Argument Reference ##
* `cloud_cidr_pool_dn` - (Required) Distinguished name of parent CloudCIDRPool object.
* `ip` - (Required) CIDR block of Object cloud subnet.
* `name` - (Optional) Name for object cloud subnet.
* `description` - (Optional) Description for object cloud subnet.
* `annotation` - (Optional) Annotation for object cloud subnet.
* `name_alias` - (Optional) Name alias for object cloud subnet.
* `scope` - (Optional) The domain applicable to the capability. Allowed values are "public", "private" and "shared". Default is "private".
* `usage` - (Optional) The usage of the port. This property shows how the port is used. Allowed values are "user", "gateway" and "infra-router". Default is "user". To make any subnet a Gateway subnet use `usage` = "gateway".	
* `zone` - (Optional) [AWS Only] Availability zone where the subnet must be deployed. This property can carry both the actual zone or the ACI logical zone name. In the former case, the driver directly uses the value of this property. In the latter case, the Connector has to first resolve the mapping from ACI logical zone to the actual AWS zone. This parameter is required in APIC v5.0 or higher
                
* `relation_cloud_rs_subnet_to_flow_log` - (Optional) Relation to class cloudAwsFlowLogPol. Cardinality - N_TO_ONE. Type - String.
                


## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Cloud Subnet.

## Importing ##

An existing Cloud Subnet can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_cloud_subnet.example <Dn>
```