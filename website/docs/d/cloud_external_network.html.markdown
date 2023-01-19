---
subcategory: "Cloud"
layout: "aci"
page_title: "ACI: aci_cloud_external_network"
sidebar_current: "docs-aci-data-source-aci_cloud_external_network"
description: |-
  Data source for ACI Cloud External Network
---

# aci_cloud_external_network #

Data source for ACI Cloud External Network
<b>Note: This resource is supported in Cloud APIC version > 25.0 only.</b>


## API Information ##

* `Class` - cloudtemplateExtNetwork
* `Distinguished Name` - uni/tn-{tenant_name}/infranetwork-{infra_network_name}/extnetwork-{external_network_name}

## GUI Information ##

* `Location` - Cloud APIC -> Application Management -> External Networks



## Example Usage ##

```hcl
data "aci_cloud_external_network" "example" {
  name  = "example"
}
```

## Argument Reference ##

* `name` - (Required) Name of the Cloud External Network.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the Cloud External Network.
* `annotation` - (Optional) Annotation of the Cloud External Network.
* `name_alias` - (Optional) Name Alias of the Cloud External Network.
* `vrf_dn` - (Optional) Distinguished name of the VRF. Note that the VRF has to be created under the infra tenant.
* `hub_network_name` - (Optional) Hub Network name of the Cloud External Network.
* `vpn_router_name` - (Optional) VPN Router name of the Cloud External Network. 
* `host_router_name` - (Optional) Host Router name of the Cloud External Network.
* `all_regions` - (Optional) Selects all regions available to the Cloud External Network. This option is always set to "yes" for Azure and AWS cAPICs and "no" in GCP cAPIC.
* `regions` - (Optional) Manually adds the regions to the Cloud External Network. This option is only available in GCP cAPICs.
* `router_type` - (Optional) Router type. Allowed values are "c8kv", "tgw". (Available only for AWS cAPIC).
* `cloud_vendor` - (Optional) Name of the vendor. Allowed values are "aws", "azure", "gcp".