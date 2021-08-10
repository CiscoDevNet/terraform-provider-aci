---
subcategory: "L4-L7 Services"
layout: "aci"
page_title: "ACI: aci_service_redirect_policy"
sidebar_current: "docs-aci-resource-service_redirect_policy"
description: |-
  Manages ACI Service Redirect Policy
---

# aci_service_redirect_policy #
Manages ACI Service Redirect Policy

## Example Usage ##

```hcl

resource "aci_service_redirect_policy" "example" {
  tenant_dn               = aci_tenant.tenentcheck.id
  name                    = "first"
  name_alias              = "name_alias"
  dest_type               = "L3"
  min_threshold_percent   = "30"
  max_threshold_percent   = "50"
  hashing_algorithm       = "sip"
  description             = "hello"
  anycast_enabled         = "no"
  resilient_hash_enabled  = "no"
  threshold_enable        = "no"
  program_local_pod_only  = "no"
  threshold_down_action   = "permit"

}

```


## Argument Reference ##
* `tenant_dn` - (Required) Distinguished name of parent Tenant object.
* `name` - (Required) Name of Object Service Redirect Policy.
* `description` - (Optional) Description of Object Service Redirect Policy.
* `anycast_enabled` - (Optional) Anycast enabled for object Service Redirect Policy. NOTE: `anycast_enabled` and `program_local_pod_only` cannot be "yes" simultaneously.
Allowed values: "yes", "no". Default value: "no".
* `annotation` - (Optional) Annotation for object Service Redirect Policy.
* `dest_type` - (Optional) Dest type for object Service Redirect Policy. Allowed values: "L1", "L2", "L3". Default value: "L3".
* `hashing_algorithm` - (Optional) Hashing algorithm for object Service Redirect Policy. Allowed values: "sip", "dip", "sip-dip-prototype", Default value: "sip-dip-prototype".
* `max_threshold_percent` - (Optional) Max threshold_percent for object Service Redirect Policy. Range : 1-100. Default Value: 0.
* `min_threshold_percent` - (Optional) Min threshold_percent for object Service Redirect Policy. Range : 1-100. Default Value: 0.
* `name_alias` - (Optional) Name alias for object Service Redirect Policy.
* `program_local_pod_only` - (Optional) Program local pod only for object Service Redirect Policy.
Allowed values: "yes", "no". Default value: "no".
* `resilient_hash_enabled` - (Optional) Resilient hash enabled for object Service Redirect Policy.
Allowed values: "yes", "no". Default value: "no".  
* `threshold_down_action` - (Optional) Threshold down the action for object Service Redirect Policy.
Allowed values: "bypass", "deny", "permit". Default value: "permit".
* `threshold_enable` - (Optional) Threshold enable for object Service Redirect Policy.
Allowed values: "yes", "no". Default value: "no".

* `relation_vns_rs_ipsla_monitoring_pol` - (Optional) Relation to class fvIPSLAMonitoringPol. Cardinality - N_TO_ONE. Type - String.
                


## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Service Redirect Policy.

## Importing ##

An existing Service Redirect Policy can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_service_redirect_policy.example <Dn>
```