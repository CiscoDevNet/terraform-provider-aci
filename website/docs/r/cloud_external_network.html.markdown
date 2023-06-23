---
subcategory: "Cloud"
layout: "aci"
page_title: "ACI: aci_cloud_external_network"
sidebar_current: "docs-aci-resource-aci_cloud_external_network"
description: |-
  Manages Cloud Network Controller Cloud External Network
---

# aci_cloud_external_network #

Manages Cloud Network Controller Cloud External Network.
<b>Note: This resource is supported in Cloud Network Controller version > 25.0 only.</b>

## API Information ##

* `Class` - cloudtemplateExtNetwork
* `Distinguished Name` - uni/tn-{tenant_name}/infranetwork-{infra_network_name}/extnetwork-{external_network_name}

## GUI Information ##

* `Location` - Cloud Network Controller -> Application Management -> External Networks


## Example Usage ##

```hcl
resource "aci_cloud_external_network" "example" {
  name       = "example"
  annotation = "orchestrator:terraform"
  vrf_dn     = aci_vrf.vrf.id
}

# GCP cloud - all_regions is set to "no" and regions can be set only in GCP Cloud
resource "aci_cloud_external_network" "external_network" {
  name         = "cloud_external_network"
  vrf_dn       = aci_vrf.vrf.id
  cloud_vendor = "gcp"
  regions      = ["europe-west3", "europe-west4"]
}

# Azure Cloud - all_regions is set to "yes" only in Azure Cloud
resource "aci_cloud_external_network" "external_network" {
  name        = "cloud_external_network"
  vrf_dn      = aci_vrf.vrf.id
  all_regions = "yes"
}
```

## Argument Reference ##

* `name` - (Required) Name of the Cloud External Network.
* `name_alias` - (Optional) Name Alias of the Cloud External Network.
* `annotation` - (Optional) Annotation of the Cloud External Network.
* `vrf_dn` - (Required) Distinguished name of the VRF. Note that the VRF has to be created under the infra tenant.
* `hub_network_name` - (Optional) Hub Network name of the Cloud External Network.
* `vpn_router_name` - (Optional) VPN Router name of the Cloud External Network. 
* `host_router_name` - (Optional) Host Router name of the Cloud External Network.
* `all_regions` - (Optional) Selects all regions available to the Cloud External Network. This option is always set to "yes" for Azure cAPICs.
* `regions` - (Optional) Manually adds the regions to the Cloud External Network. This option is only available in GCP cAPICs.
* `router_type` - (Optional) Router type. Allowed values are "c8kv", "tgw". (Available only for AWS cAPIC).
* `cloud_vendor` - (Optional) Name of the vendor. Allowed values are "aws", "azure", "gcp".



## Importing ##

An existing Cloud External Network can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_cloud_external_network.example "<Dn>"
```