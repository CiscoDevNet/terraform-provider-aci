---
subcategory: "Cloud"
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
		name 		         = "cloud_ctx_prof"
		description              = "cloud_context_profile created while acceptance testing"
		tenant_dn                = aci_tenant.footenant.id
		primary_cidr             = "10.230.231.1/16"
		region                   = "us-west-1"
		cloud_vendor	         = "aws"
		relation_cloud_rs_to_ctx = aci_vrf.example.id
		hub_network  		 	 = "uni/tn-infra/gwrouterp-default"
		annotation			     = "context_profile"
		name_alias				 = "alias_context_profile"
		type					 = "regular"
	}

```


## Argument Reference ##
* `tenant_dn` - (Required) Distinguished name of parent Tenant object.
* `description` -(Optional) Description of object Cloud Context profile.
* `name` - (Required) Name of Object Cloud Context profile.
* `primary_cidr` - (Required) Primary CIDR block of Cloud Context profile. 
* `region` - (Required) AWS region in which profile is created.
* `cloud_vendor` - (Required) Name of the vendor. e.g. "aws", "azure".
* `annotation` - (Optional) Annotation for object Cloud Context profile.
* `name_alias` - (Optional) Name alias for object Cloud Context profile.
* `type` - (Optional) The specific type of the object or component. Allowed values are "regular", "shadow", "hosted" and "container-overlay". Default is "regular".

* `hub_network` - (Optional) Hub network Dn which enables Transit Gateway.

* `relation_cloud_rs_ctx_to_flow_log` - (Optional) Relation to class cloudAwsFlowLogPol. Cardinality - N TO ONE. Type - String.
                
* `relation_cloud_rs_to_ctx` - (Required) Relation to class fvCtx. Cardinality - N TO ONE. Type - String.
                
* `relation_cloud_rs_ctx_profile_to_region` - (Optional) Relation to class cloudRegion. Cardinality - N TO ONE. Type - String.
                


