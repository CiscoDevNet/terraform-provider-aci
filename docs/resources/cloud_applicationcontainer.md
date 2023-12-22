---
subcategory: "Cloud"
layout: "aci"
page_title: "ACI: aci_cloud_applicationcontainer"
sidebar_current: "docs-aci-resource-aci_cloud_applicationcontainer"
description: |-
  Manages Cloud Network Controller Cloud Application container
---

# aci_cloud_applicationcontainer #
Manages Cloud Network Controller Cloud Application container
Note: This resource is supported in Cloud Network Controller only.
## Example Usage ##

```hcl
resource "aci_cloud_applicationcontainer" "foo_clou_app" {
  tenant_dn  = aci_tenant.dev_tenant.id
  name       = "demo_cloud_app"
  description = "From terraform"
  annotation = "tag_cloud_app"
  name_alias = "alias_app"
}

```
## Argument Reference ##
* `tenant_dn` - (Required) Distinguished name of parent Tenant object.
* `name` - (Required) Name of Object cloud applicationcontainer.
* `description` - (Optional) Description for object cloud applicationcontainer.
* `annotation` - (Optional) Annotation for object cloud applicationcontainer.
* `name_alias` - (Optional) Name alias for object cloud applicationcontainer.



## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Cloud Application container.

## Importing ##

An existing Cloud Application container can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_cloud_applicationcontainer.example <Dn>
```