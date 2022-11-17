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
<b>Note: This resource is supported in Cloud APIC only. </b>
## API Information ##

* `Class` - cloudCtxProfile
* `Distinguished Name` - uni/tn-{tenant_name}/ctxprofile-{cloud_context_profile_name}

## GUI Information ##

* `Location` - Cloud Resources -> Cloud Context Profiles
## Example Usage ##

```hcl
resource "aci_cloud_context_profile" "foocloud_context_profile" {
  name                     = "cloud_ctx_prof"
  description              = "cloud_context_profile created while acceptance testing"
  tenant_dn                = aci_tenant.footenant.id
  primary_cidr             = "10.230.231.1/16"
  region                   = "us-west-1"
  cloud_vendor             = "aws"
  relation_cloud_rs_to_ctx = aci_vrf.example.id
  hub_network              = "uni/tn-infra/gwrouterp-default"
  annotation               = "context_profile"
  name_alias               = "alias_context_profile"
  type                     = "regular"
}
```


## Argument Reference ##
* `tenant_dn` - (Required) Distinguished name of the parent Tenant object.
* `description` -(Optional) Description of the Cloud Context Profile object.
* `name` - (Required) Name of the Cloud Context Profile object.
* `primary_cidr` - (Required) Primary CIDR block of the Cloud Context Profile object.
* `region` - (Required) Region of the Cloud Context Profile object.
* `cloud_vendor` - (Required) Name of the vendor. Allowed values are "aws", "azure", "gcp".
* `annotation` - (Optional) Annotation of the Cloud Context Profile object.
* `name_alias` - (Optional) Name alias of the Cloud Context Profile object.
* `type` - (Optional) Type of the Cloud Context Profile object. Allowed values are "regular", "shadow", "hosted" and "container-overlay". Default is "regular".
* `hub_network` - (Optional) Hub Network Dn which enables Transit Gateway.
* `relation_cloud_rs_ctx_to_flow_log` - (Optional) Relation to a AWS Flow Log Policy (class cloudAwsFlowLogPol). Cardinality - N TO ONE. Type - String.
* `relation_cloud_rs_to_ctx` - (Required) Relation to a VRF (class fvCtx). Cardinality - N TO ONE. Type - String.