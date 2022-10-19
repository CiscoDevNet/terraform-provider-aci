---
layout: "aci"
page_title: "ACI: aci_templatefor_vpn_network"
sidebar_current: "docs-aci-resource-templatefor_vpn_network"
description: |-
  Manages ACI Template for VPN Network
---

# aci_templatefor_vpn_network #

Manages ACI Template for VPN Network

## API Information ##

* `Class` - cloudtemplateVpnNetwork
* `Distinguished Name` - uni/tn-{name}/infranetwork-{name}/extnetwork-{name}/vpnnetwork-{name}

## GUI Information ##

* `Location` - 


## Example Usage ##

```hcl
resource "aci_templatefor_vpn_network" "example" {
  aci_cloud_external_network_dn  = aci_templatefor_external_network.example.id
  name  = "example"
  annotation = "orchestrator:terraform"

  remote_site_id = "0"
  remote_site_name = 
}
```

## Argument Reference ##

* `aci_cloud_external_network_dn` - (Required) Distinguished name of the parent TemplateforExternalNetwork object.
* `name` - (Required) Name of the object Template for VPN Network.
* `annotation` - (Optional) Annotation of the object Template for VPN Network.

* `remote_site_id` - (Optional) Remote Site ID. Allowed range is 0-1000 and default value is "0".
* `remote_site_name` - (Optional) Remote Site Name.


## Importing ##

An existing TemplateforVPNNetwork can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_templatefor_vpn_network.example <Dn>
```