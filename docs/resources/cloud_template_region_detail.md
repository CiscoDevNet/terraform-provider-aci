---
subcategory: "Cloud"
layout: "aci"
page_title: "ACI: aci_cloud_template_region_detail"
sidebar_current: "docs-aci-resource-cloud_template_region_detail"
description: |-
  Manages ACI Cloud Template Region Detail
---

# aci_cloud_template_region_detail #

Manages ACI Cloud Template Region Detail

## API Information ##

* `Class` - cloudtemplateRegionDetail
* `Distinguished Name` - uni/tn-{tenant_name}/infranetwork-{infra_name}/stats-{stats_name}/provider-{provider}-region-{region}/regiondetail

## GUI Information ##

* `Location` - Application Management -> Cloud Context Profiles -> Hub Network Peering

## Example Usage ##

```hcl
resource "aci_cloud_template_region_detail" "hub_network" {
  parent_dn      = "uni/tn-infra/infranetwork-default/intnetwork-default/provider-azure-region-westus"
  hub_networking = "disabled"
}
```

## Argument Reference ##

* `parent_dn` - (Required) Distinguished name of the parent object.
* `hub_networking` - (Optional) Disabling `hub_networking` blocks the traffic between VNets in the given region. In order to add the cloud subnets to the cloud context profile associated with the infra tenant, `hub_networking` needs to be explicitly "disabled". After the cloud subnets are added, `hub_networking` needs to be explicitly "enabled" again. Allowed values are "disabled", "enabled". Type: String.


## Importing ##

An existing Cloud Template Region Detail can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_cloud_template_region_detail.hub_network <Dn>
```