---
subcategory: "L3Out"
layout: "aci"
page_title: "ACI: aci_logical_interface_profile"
sidebar_current: "docs-aci-resource-logical_interface_profile"
description: |-
  Manages ACI Logical Interface Profile
---

# aci_logical_interface_profile

Manages ACI Logical Interface Profile

## API Information

- `Class` - l3extLIfP
- `Distinguished Name` - uni/tn-{tenant_name}/out-{l3out}/lnodep-{logical_node_profile}/lifp-{logical_interface_profile}

## GUI Information

- `Location` - Tenants -> Networking -> L3Outs -> Logical Node Profiles -> Logical Interface Profiles

## Example Usage

```hcl
resource "aci_logical_interface_profile" "foological_interface_profile" {
  logical_node_profile_dn               = aci_logical_node_profile.foological_node_profile.id
  description                           = "aci_logical_interface_profile from terraform"
  name                                  = "demo_int_prof"
  annotation                            = "tag_prof"
  name_alias                            = "alias_prof"
  prio                                  = "unspecified"
  tag                                   = "black"
  relation_l3ext_rs_egress_qos_dpp_pol  = "uni/tn-l3out_tf_tenant/qosdpppol-egress_data_plane"
  relation_l3ext_rs_ingress_qos_dpp_pol = "uni/tn-l3out_tf_tenant/qosdpppol-ingress_data_plane"
  relation_l3ext_rs_l_if_p_cust_qos_pol = "uni/tn-l3out_tf_tenant/qoscustom-qos"
  relation_l3ext_rs_nd_if_pol           = "uni/tn-l3out_tf_tenant/ndifpol-nd"
  relation_l3ext_rs_l_if_p_to_netflow_monitor_pol {
    tn_netflow_monitor_pol_dn = "uni/tn-l3out_tf_tenant/monitorpol-netflow"
    flt_type                  = "ipv4"
  }
}
```

## Argument Reference

- `logical_node_profile_dn` - (Required) Distinguished name of the parent Logical Node Profile object.
- `name` - (Required) Name of the logical interface profile object.
- `annotation` - (Optional) Annotation of the logical interface profile object.
- `description` - (Optional) Description of the logical interface profile object.
- `name_alias` - (Optional) Name alias of the logical interface profile object.
- `prio` - (Optional) QoS priority class id. Allowed values are "unspecified", "level1", "level2", "level3", "level4", "level5" and "level6". Default is "unspecified.
- `tag` - (Optional) Specifies the color of a policy label. Allowed values are "black", "navy", "dark-blue", "medium-blue", "blue", "dark-green", "green", "teal", "dark-cyan", "deep-sky-blue", "dark-turquoise", "medium-spring-green", "lime", "spring-green", "aqua", "cyan", "midnight-blue", "dodger-blue", "light-sea-green", "forest-green", "sea-green", "dark-slate-gray", "lime-green", "medium-sea-green", "turquoise", "royal-blue", "steel-blue", "dark-slate-blue", "medium-turquoise", "indigo", "dark-olive-green", "cadet-blue", "cornflower-blue", "medium-aquamarine", "dim-gray", "slate-blue", "olive-drab", "slate-gray", "light-slate-gray", "medium-slate-blue", "lawn-green", "chartreuse", "aquamarine", "maroon", "purple", "olive", "gray", "sky-blue", "light-sky-blue", "blue-violet", "dark-red", "dark-magenta", "saddle-brown", "dark-sea-green", "light-green", "medium-purple", "dark-violet", "pale-green", "dark-orchid", "yellow-green", "sienna", "brown", "dark-gray", "light-blue", "green-yellow", "pale-turquoise", "light-steel-blue", "powder-blue", "fire-brick", "dark-goldenrod", "medium-orchid", "rosy-brown", "dark-khaki", "silver", "medium-violet-red", "indian-red", "peru", "chocolate", "tan", "light-gray", "thistle", "orchid", "goldenrod", "pale-violet-red", "crimson", "gainsboro", "plum", "burlywood", "light-cyan", "lavender", "dark-salmon", "violet", "pale-goldenrod", "light-coral", "khaki", "alice-blue", "honeydew", "azure", "sandy-brown", "wheat", "beige", "white-smoke", "mint-cream", "ghost-white", "salmon", "antique-white", "linen", "light-goldenrod-yellow", "old-lace", "red", "fuchsia", "magenta", "deep-pink", "orange-red", "tomato", "hot-pink", "coral", "dark-orange", "light-salmon", "orange", "light-pink", "pink", "gold", "peachpuff", "navajo-white", "moccasin", "bisque", "misty-rose", "blanched-almond", "papaya-whip", "lavender-blush", "seashell", "cornsilk", "lemon-chiffon", "floral-white", "snow", "yellow", "light-yellow", "ivory" and "white".

- `relation_l3ext_rs_pim_ip_if_pol` - (Optional) Represents the relation to the PIM Interface Policy (class pimIfPol). Type: String.
- `relation_l3ext_rs_pim_ipv6_if_pol` - (Optional) Represents the relation to the PIM IPv6 Interface Policy (class pimIfPol). Type: String.
- `relation_l3ext_rs_igmp_if_pol` - (Optional) Represents the relation to the IGMP Interface Policy (class igmpIfPol). Type: String.
- `relation_l3ext_rs_l_if_p_to_netflow_monitor_pol` - (Optional) Relation to the Netflow Monitor Policy (class netflowMonitorPol). Cardinality - N_TO_M. Type - Block.
  - `tn_netflow_monitor_pol_name` - (Deprecated) Distinguished name of the target Netflow Monitor Policy.
	- `tn_netflow_monitor_pol_dn` - (Required) Distinguished name of the target Netflow Monitor Policy.
	- `flt_type` - (Required) Netflow IP filter type. Allowed values: "ce", "ipv4", "ipv6".
- `relation_l3ext_rs_egress_qos_dpp_pol` - (Optional) Relation to the Egress Data Plane Policing Policy (class qosDppPol). Cardinality - N_TO_ONE. Type - String.
- `relation_l3ext_rs_ingress_qos_dpp_pol` - (Optional) Relation to the Ingress Data Plane Policing Policy (class qosDppPol). Cardinality - N_TO_ONE. Type - String.
- `relation_l3ext_rs_l_if_p_cust_qos_pol` - (Optional) Relation to the Custom QoS Policy (class qosCustomPol). Cardinality - N_TO_ONE. Type - String.
- `relation_l3ext_rs_nd_if_pol` - (Optional) Relation to the IPv6 ND policy (class ndIfPol). Cardinality - N_TO_ONE. Type - String.

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Logical Interface Profile.

## Importing

An existing Logical Interface Profile can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_logical_interface_profile.example <Dn>
```
