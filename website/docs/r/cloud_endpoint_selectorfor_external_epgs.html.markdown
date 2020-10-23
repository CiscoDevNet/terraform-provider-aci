---
layout: "aci"
page_title: "ACI: aci_cloud_endpoint_selectorfor_external_epgs"
sidebar_current: "docs-aci-resource-cloud_endpoint_selectorfor_external_epgs"
description: |-
  Manages ACI Cloud Endpoint Selector for External EPgs
---

# aci_cloud_endpoint_selectorfor_external_epgs #
Manages ACI Cloud Endpoint Selector for External EPgs
<b>Note: This resource is supported in Cloud APIC only.</b>
## Example Usage ##

```hcl
	resource "aci_cloud_endpoint_selectorfor_external_epgs" "foocloud_endpoint_selectorfor_external_epgs" {
		cloud_external_epg_dn = "${aci_cloud_external_epg.foocloud_external_epg.id}"
		description            = "sample external ep selector"
		name                   = "ext_ep_selector"
		annotation             = "tag_ext_selector"
		is_shared              = "yes"
		name_alias             = "alias_select"
		subnet                 = "0.0.0.0/0"
	}
```
## Argument Reference ##
* `cloud_external_epg_dn` - (Required) Distinguished name of parent CloudExternalEPg object.
* `name` - (Required) name of Object cloud_endpoint_selectorfor_external_epgs.
* `annotation` - (Optional) annotation for object cloud_endpoint_selectorfor_external_epgs.
* `is_shared` - (Optional) For Selectors set the shared route control. Allowed values are "yes" and "no". Default value is "yes".
* `name_alias` - (Optional) name_alias for object cloud_endpoint_selectorfor_external_epgs.
* `subnet` - (Optional) Subnet from which EP to select. Any valid CIDR block is allowed here.



## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Cloud Endpoint Selector for External EPgs.

## Importing ##

An existing Cloud Endpoint Selector for External EPgs can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_cloud_endpoint_selectorfor_external_epgs.example <Dn>
```