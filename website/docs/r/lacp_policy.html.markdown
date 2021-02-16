---
layout: "aci"
page_title: "ACI: aci_lacp_policy"
sidebar_current: "docs-aci-resource-lacp_policy"
description: |-
  Manages ACI LACP Policy
---

# aci_lacp_policy #
Manages ACI LACP Policy

## Example Usage ##

```hcl

	resource "aci_lacp_policy" "foolacp_policy" {
		description = "%s"
		name        = "demo_lacp_pol"
		annotation  = "tag_lacp"
		ctrl        = ["susp-individual", "load-defer", "graceful-conv"]
		max_links   = "16"
		min_links   = "1"
		mode        = "%s"
		name_alias  = "alias_lacp"
	}

```


## Argument Reference ##
* `name` - (Required) name of Object lacp_policy.
* `annotation` - (Optional) annotation for object lacp_policy.
* `ctrl` - (Optional) List of LAG control properties. Allowed values are "symmetric-hash", "susp-individual", "graceful-conv", "load-defer" and "fast-sel-hot-stdby".
* `max_links` - (Optional) maximum number of links. Allowed value range is "1" - "16". Default is "16".
* `min_links` - (Optional) minimum number of links in port channel. Allowed value range is "1" - "16". Default is "1".
* `mode` - (Optional) policy mode. Allowed values are "off", "active", "passive", "mac-pin" and "mac-pin-nicload". Default is "off".
* `name_alias` - (Optional) name_alias for object lacp_policy.



## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the LACP Policy.

## Importing ##

An existing LACP Policy can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_lacp_policy.example <Dn>
```