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
* `Distinguished Name` - uni/tn-{tenant_name}/infranetwork-{name}/extnetwork-{name}

## GUI Information ##

* `Location` - Tenants -> {tenant_name} -> Application Management -> External Networks



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
* `hub_network_name` - (Optional) Hub Network Name. 
* `vrf_dn` - (Optional) Distinguished name of the VRF. Note that the VRF has to be created under the infra tenant.
