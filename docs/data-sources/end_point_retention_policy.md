---
subcategory: "Tenant Policies"
layout: "aci"
page_title: "ACI: aci_end_point_retention_policy"
sidebar_current: "docs-aci-data-source-end_point_retention_policy"
description: |-
  Data source for ACI End Point Retention Policy
---

# aci_end_point_retention_policy #
Data source for ACI End Point Retention Policy

## Example Usage ##

```hcl
data "aci_end_point_retention_policy" "dev_ret_pol" {
  tenant_dn  = aci_tenant.dev_tenant.id
  name       = "foo_ret_pol"
}
```
## Argument Reference ##
* `tenant_dn` - (Required) Distinguished name of parent Tenant object.
* `name` - (Required) name of Object end_point_retention_policy.



## Attribute Reference

* `id` - Attribute id set to the Dn of the End Point Retention Policy.
* `descripton` - (Optional) Descripton for object end point retention policy.
* `annotation` - (Optional) Annotation for object end point retention policy.
* `bounce_age_intvl` - (Optional) The aging interval for a bounce entry. When an endpoint (VM) migrates to another switch, the endpoint is marked as bouncing for the specified aging interval and is deleted afterwards.   
* `bounce_trig` - (Optional) Specifies whether to install the bounce entry by RARP flood or by COOP protocol. Allowed values are "rarp-flood" and "protocol".  
* `hold_intvl` - (Optional) A time period during which new endpoint learn events will not be honored. This interval is triggered when the maximum endpoint move frequency is exceeded.  
* `local_ep_age_intvl` - (Optional) The aging interval for all local endpoints learned in this bridge domain. When 75% of the interval is reached, 3 ARP requests are sent to verify the existence of the endpoint. If no response is received, the endpoint is deleted.  
* `move_freq` - (Optional) A maximum allowed number of endpoint moves per second. If the move frequency is exceeded, the hold interval is triggered, and new endpoint learn events will not be honored until after the hold interval expires.  
* `name_alias` - (Optional) Name alias for object end point retention policy.
* `remote_ep_age_intvl` - (Optional) The aging interval for all remote endpoints learned in this bridge domain.  
