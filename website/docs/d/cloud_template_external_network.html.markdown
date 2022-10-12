---
subcategory: "Cloud"
layout: "aci"
page_title: "ACI: aci_cloud_template_external_network"
sidebar_current: "docs-aci-data-source-aci_cloud_template_external_network"
description: |-
  Data source for ACI Template for External Network
---

# aci_cloud_template_external_network #

Data source for ACI Template for cloud External Network


## API Information ##

* `Class` - cloudtemplateExtNetwork
* `Distinguished Name` - uni/tn-{tenant_name}/infranetwork-{name}/extnetwork-{name}

## GUI Information ##

* `Location` - Tenants -> {tenant_name} -> Application Management -> External Networks



## Example Usage ##
<b>Note: This resource is supported in Cloud APIC version > 25.0 only.</b>
```hcl
data "aci_cloud_template_external_network" "example" {
  name  = "example"
}
```

## Argument Reference ##

* `name` - (Required) Name of the Template for cloud External Network object.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the Template for cloud External Network object.
* `annotation` - (Optional) Annotation of the Template for cloud External Network object.
* `name_alias` - (Optional) Name Alias of the Template for cloud External Network object.
* `hub_network_name` - (Optional) Hub Network Name. 
* `vrf_dn` - (Optional) Distinguished name of the VRF.. Note that the VRF has to be created under infra tenant.
