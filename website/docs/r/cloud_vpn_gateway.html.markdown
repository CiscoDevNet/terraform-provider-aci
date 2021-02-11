---
 layout: "aci"
# page_title: "ACI: aci_cloud_router_profile"
# sidebar_current: "docs-aci-resource-cloud_router_profile"
# description: |-
#   Manages ACI Cloud Router Profile

page_title: "ACI: aci_cloud_vpn_gateway"
sidebar_current: "docs-aci-resource-cloud_vpn_gateway"
description: |-
  Manages ACI Cloud Vpn Gateway
---

<!-- # aci_cloud_router_profile #
Manages ACI Cloud Router Profile
Note: This resource is supported in Cloud APIC only.
## Example Usage ## -->

# aci_cloud_vpn_gateway #
Manages ACI Cloud Vpn Gateway
Note: This resource is supported in Cloud APIC only.
## Example Usage ## 

```hcl
 resource "aci_cloud_vpn_gateway" "example" {

  cloud_context_profile_dn  = "${aci_cloud_context_profile.example.id}"

  name  = "example"
  annotation  = "example"
  name_alias  = "example"
  num_instances  = "example"
  cloud_router_profile_type  = "example"
} 

```
## Argument Reference ##
* `cloud_context_profile_dn` - (Required) Distinguished name of parent CloudContextProfile object.
* `name` - (Required) name of Object cloud_router_profile.
* `annotation` - (Optional) annotation for object cloud_router_profile.
* `name_alias` - (Optional) name_alias for object cloud_router_profile.
* `num_instances` - (Optional) num_instances for object cloud_router_profile.
* `cloud_router_profile_type` - (Optional) component type

* `relation_cloud_rs_to_vpn_gw_pol` - (Optional) Relation to class cloudVpnGwPol. Cardinality - N_TO_ONE. Type - String.
                
* `relation_cloud_rs_to_direct_conn_pol` - (Optional) Relation to class cloudDirectConnPol. Cardinality - N_TO_ONE. Type - String.
                
* `relation_cloud_rs_to_host_router_pol` - (Optional) Relation to class cloudHostRouterPol. Cardinality - N_TO_ONE. Type - String.
                


## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Cloud Router Profile.

## Importing ##

An existing Cloud Router Profile can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html


```
terraform import aci_cloud_vpn_gateway.example <Dn>
```