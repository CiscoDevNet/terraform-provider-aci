---
subcategory: "Cloud"
layout: "aci"
page_title: "ACI: aci_cloud_vpn_gateway"
sidebar_current: "docs-aci-resource-cloud_vpn_gateway"
description: |-
  Manages ACI Cloud Vpn Gateway
--- 

# aci_cloud_vpn_gateway #
Manages ACI Cloud Vpn Gateway.


Note: This resource is supported in Cloud APIC only.

## Example Usage ##

```hcl
 resource "aci_cloud_vpn_gateway" "example" {
  cloud_context_profile_dn  = aci_cloud_context_profile.example.id
  name  = "example_name"
  description = "from terraform"
  annotation  = "example_annotation"
  name_alias  = "Alias_name"
  num_instances  = "1"
  cloud_router_profile_type  = "vpn-gw"
}

```
## Argument Reference ##
* `cloud_context_profile_dn` - (Required) Distinguished name of parent CloudContextProfile object.
* `name` - (Required) Name of Object  Cloud Router Profile.
* `description` - (Optional) Description for object Cloud Router Profile.
* `annotation` - (Optional) Annotation for object  Cloud Router Profile.
* `name_alias` - (Optional) Name_alias for object  Cloud Router Profile.
* `num_instances` - (Optional) Num instances for object  Cloud Router Profile. Default value: "1"
* `cloud_router_profile_type` - (Optional) Component type Allowed values are "host-router" and "vpn-gw". Default value is "vpn-gw". 

* `relation_cloud_rs_to_vpn_gw_pol` - (Optional) Relation to class cloudVpnGwPol. Cardinality - N_TO_ONE. Type - String.
                
* `relation_cloud_rs_to_direct_conn_pol` - (Optional) Relation to class cloudDirectConnPol. Cardinality - N_TO_ONE. Type - String.
                
* `relation_cloud_rs_to_host_router_pol` - (Optional) Relation to class cloudHostRouterPol. Cardinality - N_TO_ONE. Type - String.

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Cloud Router Profile.

## Importing

An existing Cloud Router Profile can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_cloud_vpn_gateway.example <Dn>
```
