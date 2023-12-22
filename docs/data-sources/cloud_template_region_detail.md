---
subcategory: "Cloud"
layout: "aci"
page_title: "ACI: aci_cloud_template_region_detail"
sidebar_current: "docs-aci-data-source-aci_cloud_template_region_detail"
description: |-
  Data source for ACI Cloud Template Region Detail
---

# aci_cloud_template_region_detail #

Data source for ACI Cloud Template Region Detail

## API Information ##

* `Class` - cloudtemplateRegionDetail
* `Distinguished Name` - uni/tn-{tenant_name}/infranetwork-{infra_name}/stats-{stats_name}/provider-{provider}-region-{region}/regiondetail

## GUI Information ##

* `Location` - Application Management -> Cloud Context Profiles -> Hub Network Peering


## Example Usage ##

```hcl
data "aci_cloud_template_region_detail" "example" {
  parent_dn  = "uni/tn-infra/infranetwork-default/intnetwork-default/provider-azure-region-westus"
}
```

## Argument Reference ##

* `parent_dn` - (Required) Distinguished name of the parent object.

## Attribute Reference ##

* `id` - Attribute id set to the Dn of the Cloud template region detail.
* `hub_networking` - (Read-Only) Indicates if hub networking is "enabled" or "disabled" for a given region.
