---
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
  tenant_dn               = "${aci_tenant.tenentcheck.id}"
  name                    = "first"
  dest_type               = "L3"
  max_threshold_percent   = "50"
  hashing_algorithm       = "sip"
  description             = "hello"
  anycast_enabled         = "no"
  resilient_hash_enabled  = "no"
  threshold_enable        = "no"
}

```


## Argument Reference ##
* `tenant_dn` - (Required) Distinguished name of parent Tenant object.
* `name` - (Required) name of Object service_redirect_policy.
* `anycast_enabled` - (Optional) anycast_enabled for object service_redirect_policy.
* `annotation` - (Optional) annotation for object service_redirect_policy.
* `dest_type` - (Optional) dest_type for object service_redirect_policy.
* `hashing_algorithm` - (Optional) hashing_algorithm for object service_redirect_policy.
* `max_threshold_percent` - (Optional) max_threshold_percent for object service_redirect_policy.
* `min_threshold_percent` - (Optional) min_threshold_percent for object service_redirect_policy.
* `name_alias` - (Optional) name_alias for object service_redirect_policy.
* `program_local_pod_only` - (Optional) program_local_pod_only for object service_redirect_policy.
* `resilient_hash_enabled` - (Optional) resilient_hash_enabled for object service_redirect_policy.
* `threshold_down_action` - (Optional) threshold_down_action for object service_redirect_policy.
* `threshold_enable` - (Optional) threshold_enable for object service_redirect_policy.

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