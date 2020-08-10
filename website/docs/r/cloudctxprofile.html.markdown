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
	resource "aci_cloud_context_profile" "foocloud_context_profile" {
		name 		                 = "%s"
		description              = "cloud_context_profile created while acceptance testing"
		tenant_dn                = "${aci_tenant.footenant.id}"
		primary_cidr             = "10.230.231.1/16"
		region                   = "us-west-1"
		cloud_vendor			 = "aws"
		relation_cloud_rs_to_ctx = "${aci_vrf.vrf1.id}"
	}
```
## Argument Reference ##
* `tenant_dn` - (Required) Distinguished name of parent Tenant object.
* `name` - (Required) name of Object cloud_context_profile.
* `primary_cidr` - (Required) Primary CIDR block of Cloud Context profile. 
* `region` - (Required) AWS region in which profile is created.
* `cloud_vendor` - (Required) name of the vendor. e.g. "aws", "azure".
* `annotation` - (Optional) annotation for object cloud_context_profile.
* `name_alias` - (Optional) name_alias for object cloud_context_profile.
* `type` - (Optional) The specific type of the object or component. Allowed values are "regular" and "shadow". Default is "regular".

* `relation_cloud_rs_ctx_to_flow_log` - (Optional) Relation to class cloudAwsFlowLogPol. Cardinality - N_TO_ONE. Type - String.
                
* `relation_cloud_rs_to_ctx` - (Required) Relation to class fvCtx. Cardinality - N_TO_ONE. Type - String.
                
* `relation_cloud_rs_ctx_profile_to_region` - (Optional) Relation to class cloudRegion. Cardinality - N_TO_ONE. Type - String.
                


