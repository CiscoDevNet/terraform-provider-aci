---
layout: "aci"
page_title: "ACI: aci_vrf"
sidebar_current: "docs-aci-resource-vrf"
description: |-
  Manages ACI VRF
---

# aci_vrf #

Manages ACI VRF

## API Information ##

* `Class` - fvCtx
* `Distinguished Named` - uni/tn-{name}/ctx-{name}

## GUI Information ##

* `Location` - Tenant -> Networking -> VRFs

## Example Usage ##

```hcl
resource "aci_vrf" "foovrf" {
  tenant_dn              = aci_tenant.tenant_for_vrf.id
  name                   = "demo_vrf"
  annotation             = "tag_vrf"
  bd_enforced_enable     = "no"
  ip_data_plane_learning = "enabled"
  knw_mcast_act          = "permit"
  name_alias             = "alias_vrf"
  pc_enf_dir             = "egress"
  pc_enf_pref            = "unenforced"
  relation_fv_rs_ctx_to_bgp_ctx_af_pol {
    af                     = "ipv4-ucast"
    tn_bgp_ctx_af_pol_name = "test_bgp"
  }
}
```

## Argument Reference ##

* `tenant_dn` - (Required) Distinguished name of parent Tenant object.
* `name` - (Required) name of Object vrf.
* `annotation` - (Optional) annotation tags for object vrf.
* `bd_enforced_enable` - (Optional) Flag to enable/disable bd_enforced for VRF.Allowed values are "yes" and "no". Default is "no".  
* `ip_data_plane_learning` - (Optional) Flag to enable/disable ip-data-plane learning for VRF. Allowed values are "enabled" and disabled". Default is "enabled".  
* `knw_mcast_act` - (Optional) specifies if known multicast traffic is forwarded or not. Allowed values are "permit" and "deny". Default is "permit".
* `name_alias` - (Optional) name_alias for object vrf.
* `pc_enf_dir` - (Optional) Policy Control Enforcement Direction. It is used for defining policy enforcement direction for the traffic coming to or from an L3Out. Egress and Ingress directions are wrt L3Out. Default will be Ingress. But on the existing L3Outs during upgrade it will get set to Egress so that right after upgrade behavior doesn't change for them. This also means that there is no special upgrade sequence needed for upgrading to the release introducing this feature. After upgrade user would have to change the property value to Ingress. Once changed, system will reprogram the rules and prefix entry. Rules will get removed from the egress leaf and will get installed on the ingress leaf. Actrl prefix entry, if not already, will get installed on the ingress leaf. This feature will be ignored for the following cases: 1. Golf: Gets applied at Ingress by design. 2. Transit Rules get applied at Ingress by design. 4. vzAny 5. Taboo. Allowed values are "egress" and "ingress". Default is "ingress".
* `pc_enf_pref` - (Optional) Determines if the fabric should enforce contract policies to allow routing and packet forwarding. Allowed values are "enforced" and "unenforced". Default is "enforced".

* `relation_fv_rs_ospf_ctx_pol` - (Optional) Relation to class ospfCtxPol. Cardinality - N_TO_ONE. Type - String.
* `relation_fv_rs_vrf_validation_pol` - (Optional) Relation to class l3extVrfValidationPol. Cardinality - N_TO_ONE. Type - String.
* `relation_fv_rs_ctx_mcast_to` - (Optional) Relation to class vzFilter. Cardinality - N_TO_M. Type - List.
* `relation_fv_rs_ctx_to_eigrp_ctx_af_pol` - (Optional) Relation to class eigrpCtxAfPol. Cardinality - N_TO_M. Type - Block.
* `relation_fv_rs_ctx_to_ospf_ctx_pol` - (Optional) Relation to class ospfCtxPol. Cardinality - N_TO_M. Type - Block.
* `relation_fv_rs_ctx_to_ep_ret` - (Optional) Relation to class fvEpRetPol. Cardinality - N_TO_ONE. Type - String.
* `relation_fv_rs_bgp_ctx_pol` - (Optional) Relation to class bgpCtxPol. Cardinality - N_TO_ONE. Type - String.
* `relation_fv_rs_ctx_mon_pol` - (Optional) Relation to class monEPGPol. Cardinality - N_TO_ONE. Type - String.
* `relation_fv_rs_ctx_to_ext_route_tag_pol` - (Optional) Relation to class l3extRouteTagPol. Cardinality - N_TO_ONE. Type - String.
* `relation_fv_rs_ctx_to_bgp_ctx_af_pol` - (Optional) Relation to class bgpCtxAfPol. Cardinality - N_TO_M. Type - Block.

Note: In the APIC GUI,a VRF (fvCtx) was called a "Context"or "PrivateNetwork."

## Attribute Reference ##

The only attribute that this resource exports is the `id`, which is set to the
Dn of the VRF.

## Importing ##

An existing VRF can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: <https://www.terraform.io/docs/import/index.html>

```bash
terraform import aci_vrf.example <Dn>
```
