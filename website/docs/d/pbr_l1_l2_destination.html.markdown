---
subcategory: "L4-L7 Services"
layout: "aci"
page_title: "ACI: aci_pbr_l1_l2_destination"
sidebar_current: "docs-aci-data-source-pbr_l1_l2_destination"
description: |-
  Data source for ACI Destination of L1/L2 redirected traffic
---

# aci_pbr_l1_l2_destination #

Data source for ACI Destination of L1/L2 redirected traffic


## API Information ##

* `Class` - vnsL1L2RedirectDest
* `Distinguished Name` - uni/tn-{tenant_name}/svcCont/svcRedirectPol-{pbr_name}/L1L2RedirectDest-{destName}
* `Distinguished Name` - uni/tn-{tenant_name}/svcCont/backupPol-{pbr_name}/L1L2RedirectDest-{destName}

## GUI Information ##

* `Location` - Tenants -> Policies -> Protocol -> L4-L7 Policy-Based Redirect -> L1/L2 Destinations
* `Location` - Tenants -> Policies -> Protocol -> L4-L7 Policy-Based Redirect Backup -> L1/L2 Destinations

## Example Usage ##

```hcl

# L4-L7 Policy-Based Redirect
data "aci_pbr_l1_l2_destination" "l1_l2_destination" {
  policy_based_redirect_dn = aci_service_redirect_policy.pbr.id
  dest_name                = "l1_l2_destination"
}

# L4-L7 Policy-Based Redirect Backup
data "aci_pbr_l1_l2_destination" "l1_l2_destination" {
  policy_based_redirect_dn = aci_service_redirect_backup_policy.pbr_backup.id
  dest_name                = "l1_l2_destination"
}

```

## Argument Reference ##

* `policy_based_redirect_dn` - (Required) Distinguished name of the parent Policy-Based Redirect object.
* `dest_name` - (Required) Destination Name of the L1/L2 redirected traffic.

## Attribute Reference ##
* `id` - Attribute id set to the Dn of the L1/L2 redirected traffic.
* `annotation` - (Optional) Annotation of the L1/L2 redirected traffic.
* `name_alias` - (Optional) Name Alias of the L1/L2 redirected traffic.
* `name` - (Optional) Name of the L1/L2 redirected traffic.
* `mac` - (Optional) MAC Address of the L1/L2 redirected traffic.
* `pod_id` - (Optional) Pod Id of the L1/L2 redirected traffic.
* `relation_vns_rs_l1_l2_redirect_health_group` - (Optional) Represents the relation to a L4-L7 Redirect Health Group (class vnsRedirectHealthGroup).
* `relation_vns_rs_to_c_if` - (Optional) Represents the relation to a Concrete Interface (class vnsCIf).
