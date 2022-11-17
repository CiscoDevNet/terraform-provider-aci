---
subcategory: "Cloud"
layout: "aci"
page_title: "ACI: aci_cloud_context_profile"
sidebar_current: "docs-aci-data-source-cloud_context_profile"
description: |-
  Data source for ACI Cloud Context Profile
---

# aci_cloud_context_profile

Data source for ACI Cloud Context Profile
<b>Note: This resource is supported in Cloud APIC only. </b>

## API Information ##

* `Class` - cloudCtxProfile
* `Distinguished Name` - uni/tn-{tenant_name}/ctxprofile-{cloud_context_profile_name}

## GUI Information ##

* `Location` - Cloud Resources -> Cloud Context Profiles

## Example Usage

```hcl

data "aci_cloud_context_profile" "sample_prof" {
  tenant_dn  = aci_tenant.dev_tenant.id
  name       = "demo_cloud_ctx_prof"
}

```

## Argument Reference

* `tenant_dn` - (Required) Distinguished name of the parent Tenant object.
* `name` - (Required) Name of the Cloud Context Profile object.

## Attribute Reference

* `id` - Dn of the Cloud Context Profile object.
* `description` - Description of the Cloud Context Profile object.
* `annotation` - Annotation of the Cloud Context Profile object.
* `primary_cidr` - Primary CIDR block of the Cloud Context Profile object.
* `region` - Region of the Cloud Context Profile object.
* `cloud_vendor` - Name of the vendor.
* `name_alias` - Name alias of the Cloud Context Profile object.
* `type` - Type of the Cloud Context Profile object.
* `hub_network` - Hub Network Dn which enables Transit Gateway.
* `relation_cloud_rs_ctx_to_flow_log` - Relation to a AWS Flow Log Policy (class cloudAwsFlowLogPol).
* `relation_cloud_rs_to_ctx` - Relation to a VRF (class fvCtx).