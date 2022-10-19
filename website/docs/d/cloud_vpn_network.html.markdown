---
layout: "aci"
page_title: "ACI: aci_templatefor_vpn_network"
sidebar_current: "docs-aci-data-source-templatefor_vpn_network"
description: |-
  Data source for ACI Template for VPN Network
---

# aci_templatefor_vpn_network #

Data source for ACI Template for VPN Network


## API Information ##

* `Class` - cloudtemplateVpnNetwork
* `Distinguished Name` - uni/tn-{name}/infranetwork-{name}/extnetwork-{name}/vpnnetwork-{name}

## GUI Information ##

* `Location` - 



## Example Usage ##

```hcl
data "aci_templatefor_vpn_network" "example" {
  aci_cloud_external_network_dn  = aci_templatefor_external_network.example.id
  name  = "example"
}
```

## Argument Reference ##

* `aci_cloud_external_network_dn` - (Required) Distinguished name of parent TemplateforExternalNetwork object.
* `name` - (Required) Name of object Template for VPN Network.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the Template for VPN Network.
* `annotation` - (Optional) Annotation of object Template for VPN Network.
* `name_alias` - (Optional) Name Alias of object Template for VPN Network.
* `remote_site_id` - (Optional) Remote Site ID. 
* `remote_site_name` - (Optional) Remote Site Name. 
