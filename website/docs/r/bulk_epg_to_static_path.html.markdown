---
subcategory: "Application Management"
layout: "aci"
page_title: "ACI: aci_bulk_epg_to_static_path"
sidebar_current: "docs-aci-resource-bulk_epg_to_static_path"
description: |-
  Manages ACI EPG to Static Paths in bulk
---

# aci_bulk_epg_to_static_path #

Manages ACI EPG to Static Paths in bulk

## API Information ##
* `Class` - fvRsPathAtt
* `Distinguished Name` - uni/tn-{tenant_name}/ap-{anp_name}/epg-{epg_name}/rspathAtt-[{interface_dn}]

## GUI Information ##
* `Location` - Tenants -> {tenant_name} -> Application Profiles -> {anp_name} -> Application EPGs -> {epg_name} -> Static Ports

## Example Usage ##

```hcl
resource "aci_bulk_epg_to_static_path" "example" {
  application_epg_dn = aci_application_epg.example_epg.id
  static_path {
    interface_dn         = "topology/pod-1/paths-129/pathep-[eth1/5]"
    encap                = "vlan-1000"
  }
  static_path {
    interface_dn         = "topology/pod-1/paths-129/pathep-[eth1/6]"
    encap                = "vlan-1001"
    description          = "this is updated desc for another bulk static path"
    deployment_immediacy = "immediate"
    mode                 = "regular"
  }
  static_path {
    interface_dn         = "topology/pod-1/paths-129/pathep-[eth1/7]"
    encap                = "vlan-1002"
    description          = "this is desc for third bulk static path"
    deployment_immediacy = "lazy"
    mode                 = "untagged"
    primary_encap        = "vlan-900"
  }
  static_path {
    interface_dn         = "topology/pod-1/paths-129/pathep-[eth1/8]"
    encap                = "vlan-1003"
    description          = "this is desc for fourth bulk static path"
    deployment_immediacy = "lazy"
    mode                 = "native"
  }
}
```

## Argument Reference ##

* `application_epg_dn` - (Required) Distinguished name of the parent Application EPG object.
* `static_path` - (Optional) A block representing a Static Path object. Type: Block.
  * `interface_dn` - (Required) Distinguished name of the interface to assign to this EPG. Type: String.
  * `encap` - (Required) Encapsulation to use for the Static Path (for example: vlan-100). Type: String.
  * `deployment_immediacy` - (Optional) Deployment immediacy of the Static Path. Allowed values: "immediate", "lazy". Default value: "lazy". Type: String.
  * `mode` - (Optional) Mode of the static association of the interface. Allowed values: "regular", "native", "untagged". Default value: "regular". Type: String.
  * `primary_encap` - (Optional) Primary encapsulation for the Static Path object (used for micro-segmentation). Type: String.

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the Dn of the EPG for which the static paths are configured for.

## Importing ##

The existing Static Paths of an EPG can be [imported][docs-import] into this resource using the EPG Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_bulk_epg_to_static_path.example <Dn>
```
