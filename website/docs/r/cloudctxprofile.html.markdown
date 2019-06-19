---
layout: "aci"
page_title: "ACI: aci_cloud_context_profile"
sidebar_current: "docs-aci-resource-cloud_context_profile"
description: |-
  Manages ACI Cloud Context Profile
---

# aci_cloud_context_profile #
Manages ACI Cloud Context Profile

## Example Usage ##

```hcl
resource "aci_cloud_context_profile" "example" {

  tenant_dn  = "${aci_tenant.example.id}"

    name  = "example"

  annotation  = "example"
  name_alias  = "example"
  type  = "example"
}
```
## Argument Reference ##
* `tenant_dn` - (Required) Distinguished name of parent Tenant object.
* `name` - (Required) name of Object cloud_context_profile.
* `annotation` - (Optional) annotation for object cloud_context_profile.
* `name_alias` - (Optional) name_alias for object cloud_context_profile.
* `type` - (Optional) component type

* `relation_cloud_rs_ctx_to_flow_log` - (Optional) Relation to class cloudAwsFlowLogPol. Cardinality - N_TO_ONE. Type - String.
                
* `relation_cloud_rs_to_ctx` - (Optional) Relation to class fvCtx. Cardinality - N_TO_ONE. Type - String.
                
* `relation_cloud_rs_ctx_profile_to_region` - (Optional) Relation to class cloudRegion. Cardinality - N_TO_ONE. Type - String.
                


