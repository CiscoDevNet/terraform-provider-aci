---
layout: "aci"
page_title: "ACI: aci_rest"
sidebar_current: "docs-aci-data-source-rest"
description: |-
  Data source for ACI Rest
---

# aci_rest #
Data source for ACI Rest

## Example Usage ##

```hcl
resource "aci_tenant" "tenentcheck" {
  name       = "phase2"
  annotation = "atag"
  name_alias = "alias_tenant"
}

// To get information regarding the tenant object without children
data "aci_rest" "tenant_rest" {
  path = "/api/node/mo/${aci_tenant.tenentcheck.id}.json"
}

// To get information regarding the tenant object with children
data "aci_rest" "tenant_rest" {
  path = "/api/node/mo/${aci_tenant.tenentcheck.id}.json?rsp-subtree=children"
}

resource "aci_bgp_peer_prefix" "example" {
  tenant_dn    = "${aci_tenant.tenentcheck.id}"
  name         = "one"
  description  = "from terraform"
  action       = "shut"
  annotation   = "example"
  max_pfx      = "200"
  name_alias   = "example"
  restart_time = "200"
  thresh       = "85"
}

// To get information regarding the BGP peer prefix object without children
data "aci_rest" "bgp_peer_prefix_Rest" {
  path = "/api/node/mo/${aci_bgp_peer_prefix.example.id}.json"
}

// To get information regarding the BGP peer prefix object with children
data "aci_rest" "tenant_rest" {
  path = "/api/node/mo/${aci_bgp_peer_prefix.example.id}.json?rsp-subtree=children"
}
```

## Argument Reference ##

* `path` - (Required) ACI path for object which should should be get. Starting with api/node/mo/{parent-dn}(if applicable)/{rn of object}.json

<strong>Note</strong> : To extract children, use path format as "api/node/mo/{parent-dn}(if applicable)/{rn of object}.json?rsp-subtree=children"


## Attribute Reference

* `id` - Dishtiguished name of object being managed.
* `class_name` - Class name of object being managed.
* `content` - Map of key-value pairs which represents the attributes for the object being managed.
* `dn` - Distinguished name of object being managed.

* `children` - Set of children of the object being managed.
* `children.child_class_name` - Class name of the child of the object being managed.
* `children.child_content` - Map of key-value pairs which represents the attributes for child of the object being managed.
