---
layout: "aci"
page_title: "ACI: aci_templatefor_external_network"
sidebar_current: "docs-aci-resource-templatefor_external_network"
description: |-
  Manages ACI Template for External Network
---

# aci_templatefor_external_network #

Manages ACI Template for External Network

## API Information ##

* `Class` - cloudtemplateExtNetwork
* `Distinguished Name` - uni/tn-{name}/infranetwork-{name}/extnetwork-{name}

## GUI Information ##

* `Location` - 


## Example Usage ##

```hcl
resource "aci_templatefor_external_network" "example" {
  infra_network_template_dn  = aci_infra_network_template.example.id
  name  = "example"
  annotation = "orchestrator:terraform"
  hub_network_name = 

  vrf_name = "overlay-1"
}
```

## Argument Reference ##

* `infra_network_template_dn` - (Required) Distinguished name of the parent InfraNetworkTemplate object.
* `name` - (Required) Name of the object Template for External Network.
* `annotation` - (Optional) Annotation of the object Template for External Network.

* `hub_network_name` - (Optional) Hub Network Name.
* `vrf_name` - (Optional) External Network VRF Name.The VRF name. This name can be up to 64 alphanumeric characters. Allowed values are and default value is "overlay-1".


## Importing ##

An existing CloudTemplateforExternalNetwork can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_templatefor_external_network.example <Dn>
```