---
layout: "aci"
page_title: "ACI: aci_end_point_retention_policy"
sidebar_current: "docs-aci-resource-end_point_retention_policy"
description: |-
  Manages ACI End Point Retention Policy
---

# aci_end_point_retention_policy #
Manage End Point (EP) retention protocol policies


## Example Usage ##

```hcl
	resource "aci_end_point_retention_policy" "fooend_point_retention_policy" {
		tenant_dn   		= aci_tenant.tenant_for_ret_pol.id
		description 		= "From Terraform"
		name                = "demo_ret_pol"
		annotation          = "tag_ret_pol"
		bounce_age_intvl    = "630"
		bounce_trig         = "protocol"
		hold_intvl          = "6"
		local_ep_age_intvl  = "900"
		move_freq           = "256"
		name_alias          = "alias_demo"
		remote_ep_age_intvl = "300"
	} 
```
## Argument Reference ##
* `tenant_dn` - (Required) Distinguished name of parent Tenant object.
* `name` - (Required) Name of Object end point retention policy.
* `descripton` - (Optional) Descripton for object end point retention policy.
* `annotation` - (Optional) Annotation for object end point retention policy.
* `bounce_age_intvl` - (Optional)  The aging interval for a bounce entry. When an endpoint (VM) migrates to another switch, the endpoint is marked as bouncing for the specified aging interval and is deleted afterwards. Allowed value range is "0" - "0xffff". Default is "630".
* `bounce_trig` - (Optional) Specifies whether to install the bounce entry by RARP flood or by COOP protocol. Allowed values are "rarp-flood" and "protocol". Default is "protocol".
* `hold_intvl` - (Optional) A time period during which new endpoint learn events will not be honored. This interval is triggered when the maximum endpoint move frequency is exceeded. Allowed value range is "5" - "0xffff". Default is "300".  
* `local_ep_age_intvl` - (Optional) The aging interval for all local endpoints learned in this bridge domain. When 75% of the interval is reached, 3 ARP requests are sent to verify the existence of the endpoint. If no response is received, the endpoint is deleted. Allowed value range is "120" - "0xffff". Default is "900". "0" is treated as special value here. Providing interval as "0" is treated as infinite interval.
* `move_freq` - (Optional) A maximum allowed number of endpoint moves per second. If the move frequency is exceeded, the hold interval is triggered, and new endpoint learn events will not be honored until after the hold interval expires. Allowed value range is "0" - "0xffff". Default is "256".
* `name_alias` - (Optional) Name alias for object end point retention policy.
* `remote_ep_age_intvl` - (Optional) The aging interval for all remote endpoints learned in this bridge domain.Allowed value range is "120" - "0xffff". Default is "900". "0" is treated as special value here. Providing interval as "0" is treated as infinite interval.



## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the End Point Retention Policy.

## Importing ##

An existing End Point Retention Policy can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_end_point_retention_policy.example <Dn>
```