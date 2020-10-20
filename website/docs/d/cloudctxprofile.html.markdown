---
layout: "aci"
page_title: "ACI: aci_cloud_context_profile"
sidebar_current: "docs-aci-data-source-cloud_context_profile"
description: |-
  Data source for ACI Cloud Context Profile
---

# aci_cloud_context_profile #
Data source for ACI Cloud Context Profile
<b>Note: This resource is supported in Cloud APIC only. </b>
## Example Usage ##

```hcl

data "aci_cloud_context_profile" "sample_prof" {
  tenant_dn  = "${aci_tenant.dev_tenant.id}"
  name       = "demo_cloud_ctx_prof"
}

```


## Argument Reference ##
* `tenant_dn` - (Required) Distinguished name of parent Tenant object.
* `name` - (Required) name of Object cloud-ctx-profile.



## Attribute Reference

* `id` - Attribute id set to the Dn of the Cloud Context profile.
* `annotation` - annotation for object Cloud Context profile.
* `name_alias` - name_alias for object Cloud Context Profile.
* `type` - The specific type of the object or component. 
* `primary_cidr` - Primary CIDR block of Cloud Context profile. 
* `region` - AWS region in which profile is created.
* `hub_network` - hub network Dn which enables Transit Gateway.

