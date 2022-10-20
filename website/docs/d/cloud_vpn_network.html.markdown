---
subcategory: "Cloud"
layout: "aci"
page_title: "ACI: aci_cloud_external_network_vpn_network"
sidebar_current: "docs-aci-data-source-aci_cloud_external_network_vpn_network"
description: |-
  Data source for ACI Template for VPN Network
---

# aci_cloud_external_network_vpn_network #

Data source for ACI Template for VPN Network


## API Information ##

* `Class` - cloudtemplateVpnNetwork
* `Distinguished Name` - uni/tn-{tenant_name}/infranetwork-{name}/extnetwork-{name}/vpnnetwork-{name}

## GUI Information ##

* `Location` - Tenants -> {tenant_name} -> Application Management -> External Networks -> VPN Networks



## Example Usage ##

```hcl
data "aci_cloud_external_network_vpn_network" "example" {
  aci_cloud_external_network_dn  = aci_cloud_external_network.example.id
  name  = "example"
}
```

## Argument Reference ##

* `aci_cloud_external_network_dn` - (Required) Distinguished name of parent TemplateforExternalNetwork object.
* `name` - (Required) Name of the Cloud VPN Network object.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the Cloud VPN Network.
* `annotation` - (Optional) Annotation of the Cloud VPN Network object.
* `name_alias` - (Optional) Name Alias of the Cloud VPN Network object.
* `remote_site_id` - (Optional) Remote Site ID. 
* `remote_site_name` - (Optional) Remote Site Name. 
