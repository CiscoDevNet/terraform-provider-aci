---
subcategory: "Access Policies"
layout: "aci"
page_title: "ACI: aci_leaf_access_bundle_policy_sub_group"
sidebar_current: "docs-aci-data-source-leaf_access_bundle_policy_sub_group"
description: |-
  Data source for ACI Override Policy Group
---

# aci_leaf_access_bundle_policy_sub_group #

Data source for ACI Override Policy Group

## API Information ##

* `Class` - infraAccBndlSubgrp
* `Distinguished Name` - uni/infra/funcprof/accbundle-{name}/accsubbndl-{name}

## GUI Information ##

* `Location` - Fabric - Access Policies - Interfaces - Leaf Interfaces - Policy Groups - [ PC | VPC ] Interface - Advanced Policies - Override Access Policy Groups

## Example Usage ##

```hcl
data "aci_leaf_access_bundle_policy_sub_group" "example" {
  leaf_access_bundle_policy_group_dn  = aci_leaf_access_bundle_policy_group.example.id
  name  = "example"
}
```

## Argument Reference ##

* `leaf_access_bundle_policy_group_dn` - (Required) Distinguished name of the parent infraAccBndlGrp object.
* `name` - (Required) Name of the object Override Policy Group.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the Override Policy Group.
* `annotation` - (Optional) Annotation of the object Override Policy Group.
* `name_alias` - (Optional) Name Alias of the object Override Policy Group.
