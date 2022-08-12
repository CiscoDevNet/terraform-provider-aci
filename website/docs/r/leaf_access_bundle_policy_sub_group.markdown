---
subcategory: "Access Policies"
layout: "aci"
page_title: "ACI: aci_leaf_access_bundle_policy_sub_group"
sidebar_current: "docs-aci-resource-leaf_access_bundle_policy_sub_group"
description: |-
  Manages ACI Override Policy Group
---

# aci_leaf_access_bundle_policy_sub_group #

Manages ACI Override Policy Group

## API Information ##

* `Class` - infraAccBndlSubgrp
* `Distinguished Name` - uni/infra/funcprof/accbundle-{name}/accsubbndl-{name}

## GUI Information ##

* `Location` - Fabric - Access Policies - Interfaces - Leaf Interfaces - Policy Groups - [ PC | VPC ] Interface - Advanced Policies - Override Access Policy Groups

## Example Usage ##

```hcl
resource "aci_leaf_access_bundle_policy_sub_group" "example" {
  leaf_access_bundle_policy_group_dn  = aci_leaf_access_bundle_policy_group.example.id
  name  = "example"
  port_channel_member = aci_resource.example.id
}
```

## Argument Reference ##

* `leaf_access_bundle_policy_group_dn` - (Required) Distinguished name of the parent infraAccBndlGrp object.
* `name` - (Required) Name of the object Override Policy Group.
* `annotation` - (Optional) Annotation of the object Override Policy Group.
* `description` - (Optional) Description of the object Override Policy Group.
* `name_alias` - (Optional) Name alias.
* `port_channel_member` - (Optional) Represents the relation to a Relation to LACP Interface Policy (class lacpIfPol). The PortChannel member policy configured parameters. Type: String.

## Importing ##

An existing infraAccBndlSubgrp can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_leaf_access_bundle_policy_sub_group.example <Dn>
```