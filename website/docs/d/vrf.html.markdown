---
layout: "aci"
page_title: "ACI: aci_vrf"
sidebar_current: "docs-aci-data-source-vrf"
description: |-
  Data source for ACI VRF
---

# aci_vrf

Data source for ACI VRF

## Example Usage

```hcl
data "aci_vrf" "dev_ctx" {
  tenant_dn  = aci_tenant.dev_tenant.id
  name       = "foo_ctx"
}
```

## Argument Reference

- `tenant_dn` - (Required) Distinguished name of parent Tenant object.
- `name` - (Required) name of Object vrf.

## Attribute Reference

- `id` - Attribute id set to the Dn of the VRF.
- `annotation` - (Optional) Annotation tags for object VRF.
- `description` - (Optional) Description tags for object VRF.
- `bd_enforced_enable` - (Optional) Flag to enable/disable enforced bridge domain for VRF.
- `ip_data_plane_learning` - (Optional) Flag to enable/disable IP-data-plane learning for VRF.
- `knw_mcast_act` - (Optional) Specifies if known multicast traffic is forwarded or not.
- `name_alias` - (Optional) Name alias for object VRF.
- `pc_enf_dir` - (Optional) Policy Control Enforcement Direction. It is used for defining policy enforcement direction for the traffic coming to or from an L3Out. Egress and Ingress directions are wrt L3Out. Default will be Ingress. But on the existing L3Outs during upgrade it will get set to Egress so that right after upgrade behavior doesn't change for them. This also means that there is no special upgrade sequence needed for upgrading to the release introducing this feature. After upgrade user would have to change the property value to Ingress. Once changed, system will reprogram the rules and prefix entry. Rules will get removed from the egress leaf and will get installed on the ingress leaf. Actrl prefix entry, if not already, will get installed on the ingress leaf. This feature will be ignored for the following cases: 1. Golf: Gets applied at Ingress by design. 2. Transit Rules get applied at Ingress by design. 4. vzAny 5. Taboo.
- `pc_enf_pref` - (Optional) Determines if the fabric should enforce contract policies to allow routing and packet forwarding.
