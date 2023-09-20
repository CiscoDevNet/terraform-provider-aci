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
* `name` - (Read-Only) Name of the Cloud Subnet object.
* `description` - (Read-Only) Description of the Cloud Subnet object.
* `annotation` - (Read-Only) Annotation of the Cloud Subnet object.
* `name_alias` - (Read-Only) Name alias of the Cloud Subnet object.
* `scope` - (Read-Only) List of domain applicable to the capability. Allowed values are "public", "private" and "shared".
* `usage` - (Read-Only) The usage of the port. This property shows how the port is used.
* `zone` - (Read-Only) Relation to a Cloud Resource Zone (class cloudRsZoneAttach).
* `relation_cloud_rs_subnet_to_flow_log` - (Read-Only) Relation to the AWS Flow Log Policy (class cloudAwsFlowLogPol).
* `relation_cloud_rs_subnet_to_ctx` - (Read-Only) Relation to associate the subnet with a secondary VRF (class cloudRsSubnetToCtx).
* `subnet_group_label` - (Read-Only) Subnet Group Label of the Cloud Subnet object. It is only applicable to the GCP vendor.