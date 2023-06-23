---
subcategory: "Cloud"
layout: "aci"
page_title: "ACI: aci_cloud_subnet"
sidebar_current: "docs-aci-data-source-cloud_subnet"
description: |-
  Data source for Cloud Network Controller Cloud Subnet
---

# aci_cloud_subnet #
Data source for Cloud Network Controller Cloud Subnet  
<b>Note: This resource is supported in Cloud Network Controller only.</b>

## API Information ##

* `Class` - cloudSubnet
* `Distinguished Name` - uni/tn-{tenant_name}/ctxprofile-{cloud_context_profile_name}/cidr-[{addr}]/subnet-[{ip}]

## GUI Information ##

* `Location` - Application Management -> Cloud Context Profile -> CIDR Block Range Subnets -> Subnet -> Subnet Group Label

## Example Usage ##

```hcl
data "aci_cloud_subnet" "dev_subnet" {
  cloud_cidr_pool_dn  = aci_cloud_cidr_pool.dev_cidr_pool.id
  ip                  = "14.12.0.0/28"
}
```

## Argument Reference ##
* `cloud_cidr_pool_dn` - (Required) Distinguished name of the Cloud CIDR Pool parent object.
* `ip` - (Required) CIDR block of Object cloud subnet.


## Attribute Reference

* `id` - Dn of the Cloud Subnet object.
* `name` - (Optional) Name of the Cloud Subnet object.
* `description` - (Optional) Description of the Cloud Subnet object.
* `annotation` - (Optional) Annotation of the Cloud Subnet object.
* `name_alias` - (Optional) Name alias of the Cloud Subnet object.
* `scope` - (Optional) List of domain applicable to the capability. Allowed values are "public", "private" and "shared". Default is ["private"].
* `usage` - (Optional) The usage of the port. This property shows how the port is used.
* `zone` - (Optional) Relation to a Cloud Resource Zone (class cloudRsZoneAttach). It is only applicable to the AWS vendor.
* `relation_cloud_rs_subnet_to_flow_log` - (Optional) Relation to the AWS Flow Log Policy (class cloudAwsFlowLogPol).
* `subnet_group_label` - (Optional) Subnet Group Label of the Cloud Subnet object. It is only applicable to the GCP vendor.