---
subcategory: "Cloud"
layout: "aci"
page_title: "ACI: aci_cloud_template_external_network"
sidebar_current: "docs-aci-resource-aci_cloud_template_external_network"
description: |-
  Manages ACI Template for cloud External Network
---

# aci_cloud_template_external_network #

Manages ACI Template for cloud External Network

## API Information ##

* `Class` - cloudtemplateExtNetwork
* `Distinguished Name` - uni/tn-{tenant_name}/infranetwork-{name}/extnetwork-{name}

## GUI Information ##

* `Location` - Tenants -> {tenant_name} -> Application Management -> External Networks


## Example Usage ##
<b>Note: This resource is supported in Cloud APIC version > 25.0 only.</b>
```hcl
resource "aci_cloud_template_external_network" "example" {
  name  = "example"
  annotation = "orchestrator:terraform"
  hub_network_name = 
  vrf_dn = aci_vrf.vrf.id
}
```

## Argument Reference ##

* `name` - (Required) Name of the Template for cloud External Network object.
* `annotation` - (Optional) Annotation of the Template for cloud External Network object.
* `hub_network_name` - (Optional) Hub Network Name.
* `vrf_dn` - (Required) Distinguished name of the VRF. Note that the VRF has to be created under infra tenant.


## Importing ##

An existing CloudTemplateforExternalNetwork can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_cloud_template_external_network.example "<Dn>"
```