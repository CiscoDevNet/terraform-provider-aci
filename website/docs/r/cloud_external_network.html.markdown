---
subcategory: "Cloud"
layout: "aci"
page_title: "ACI: aci_cloud_external_network"
sidebar_current: "docs-aci-resource-aci_cloud_external_network"
description: |-
  Manages ACI Cloud External Network
---

# aci_cloud_external_network #

Manages ACI Cloud External Network.
<b>Note: This resource is supported in Cloud APIC version > 25.0 only.</b>

## API Information ##

* `Class` - cloudtemplateExtNetwork
* `Distinguished Name` - uni/tn-{tenant_name}/infranetwork-{name}/extnetwork-{name}

## GUI Information ##

* `Location` - Tenants -> {tenant_name} -> Application Management -> External Networks


## Example Usage ##

```hcl
resource "aci_cloud_external_network" "example" {
  name  = "example"
  annotation = "orchestrator:terraform"
  vrf_dn = aci_vrf.vrf.id
}
```

## Argument Reference ##

* `name` - (Required) Name of the Cloud External Network.
* `annotation` - (Optional) Annotation of the Cloud External Network.
* `hub_network_name` - (Optional) Hub Network name of the Cloud External Network.
* `vrf_dn` - (Required) Distinguished name of the VRF. Note that the VRF has to be created under the infra tenant.


## Importing ##

An existing Cloud External Network can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_cloud_external_network.example "<Dn>"
```