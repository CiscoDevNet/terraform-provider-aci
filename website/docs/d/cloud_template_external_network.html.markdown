---
layout: "aci"
page_title: "ACI: aci_templatefor_external_network"
sidebar_current: "docs-aci-data-source-templatefor_external_network"
description: |-
  Data source for ACI Template for External Network
---

# aci_templatefor_external_network #

Data source for ACI Template for External Network


## API Information ##

* `Class` - cloudtemplateExtNetwork
* `Distinguished Name` - uni/tn-{name}/infranetwork-{name}/extnetwork-{name}

## GUI Information ##

* `Location` - 



## Example Usage ##

```hcl
data "aci_templatefor_external_network" "example" {
  infra_network_template_dn  = aci_infra_network_template.example.id
  name  = "example"
}
```

## Argument Reference ##

* `infra_network_template_dn` - (Required) Distinguished name of parent InfraNetworkTemplate object.
* `name` - (Required) Name of object Template for External Network.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the Template for External Network.
* `annotation` - (Optional) Annotation of object Template for External Network.
* `name_alias` - (Optional) Name Alias of object Template for External Network.
* `hub_network_name` - (Optional) Hub Network Name. 
* `vrf_name` - (Optional) External Network VRF Name. The VRF name. This name can be up to 64 alphanumeric characters.
