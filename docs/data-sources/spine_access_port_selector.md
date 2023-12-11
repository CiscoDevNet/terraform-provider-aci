---
subcategory: "Access Policies"
layout: "aci"
page_title: "ACI: aci_spine_access_port_selector"
sidebar_current: "docs-aci-data-source-aci_spine_access_port_selector"
description: |-
  Data source for ACI Spine Access Port Selector
---

# aci_spine_access_port_selector #

Data source for ACI Spine Access Port Selector


## API Information ##

* `Class` - infraSHPortS
* `Distinguished Named` - uni/infra/spaccportprof-{name}/shports-{name}-typ-{type}

## GUI Information ##

* `Location` - Fabric > Access Policies > Interfaces > Spine Interfaces > Profiles > {interface_profile}:{interface_selector}



## Example Usage ##

```hcl
data "aci_spine_access_port_selector" "example" {
  spine_interface_profile_dn  = aci_spine_interface_profile.example.id
  name  = "example"
  spine_access_port_selector_type  = "ALL"
}
```

## Argument Reference ##

* `spine_interface_profile_dn` - (Required) Distinguished name of the parent Spine Interface Profile object.
* `name` - (Required) Name of the Spine Access Port Selector.
* `spine_access_port_selector_type` - (Required) The type of Spine Access Port Selector. Allowed values are "ALL" and "range". Default is "ALL".

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the Spine Access Port Selector.
* `annotation` - (Optional) Annotation of the Spine Access Port Selector.
* `name_alias` - (Optional) Name Alias of the Spine Access Port Selector.