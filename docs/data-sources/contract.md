---
subcategory: "Contract"
layout: "aci"
page_title: "ACI: aci_contract"
sidebar_current: "docs-aci-data-source-aci_contract"
description: |-
  Data source for ACI Contract
---

# aci_contract #
Data source for ACI Contract

## Example Usage ##

```hcl
data "aci_contract" "example" {
  tenant_dn  = aci_tenant.dev_tenant.id
  name       = "contract_name"
}
```
## Argument Reference ##
* `tenant_dn` - (Required) Distinguished name of parent Tenant object.
* `name` - (Required) name of Object contract.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Contract.
* `description` - (Optional) Description for object contract.
* `annotation` - (Optional) Annotation for object contract.
* `name_alias` - (Optional) Name alias for object contract.
* `prio` - (Optional) Priority level of the service contract.
* `scope` - (Optional) Represents the scope of this contract. If the scope is set as application-profile, the epg can only communicate with epgs in the same application-profile.
* `target_dscp` - (Optional) The target differentiated services code point (DSCP) of the path attached to the layer 3 outside profile.
