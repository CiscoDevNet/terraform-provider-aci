---
subcategory: "Application Management"
layout: "aci"
page_title: "ACI: aci_application_epg"
sidebar_current: "docs-aci-resource-aci_application_epg"
description: |-
  Manages ACI Application EPG
---

# aci_application_epg

Manages ACI Application EPG

## Example Usage

```hcl
resource "aci_application_epg" "fooapplication_epg" {
    application_profile_dn  = aci_application_profile.app_profile_for_epg.id
    name                    = "demo_epg"
    description             = "from terraform"
    annotation              = "tag_epg"
    exception_tag           = "0"
    flood_on_encap          = "disabled"
    fwd_ctrl                = "none"
    has_mcast_source        = "no"
    is_attr_based_epg       = "no"
    match_t                 = "AtleastOne"
    name_alias              = "alias_epg"
    pc_enf_pref             = "unenforced"
    pref_gr_memb            = "exclude"
    prio                    = "unspecified"
    shutdown                = "no"
    relation_fv_rs_bd       = aci_bridge_domain.example.id
}
```

## Argument Reference ##
* `application_profile_dn` - (Required) Distinguished Name of the parent application profile. Type - String.
* `name` - (Required) Name of Object application epg. Type - String.
* `pc_tag` - (Read-Only) A numeric ID to represent a policy enforcement group.
* `annotation` - (Optional) Annotation for object application epg. Type - String.
* `description` - (Optional) Description for object application epg. Type - String.
* `exception_tag` - (Optional) Exception tag for object application epg. Range: "0" - "512" . Type - String.
* `flood_on_encap` - (Optional) Control at EPG level if the traffic L2 Multicast/Broadcast and Link Local Layer should be flooded only on ENCAP or based on bridg-domain settings. Allowed values are "disabled" and "enabled". Default is "disabled". Type - String.
* `fwd_ctrl` - (Optional) Forwarding control at EPG level. Allowed values are "none" and "proxy-arp". Default is "none". Type - String.
* `has_mcast_source` - (Optional) If the source for the EPG is multicast or not. Allowed values are "yes" and "no". Default values is "no". Type - String.
* `is_attr_based_epg` - (Optional) If the EPG is attribute based or not. Allowed values are "yes" and "no". Default is "no". Type - String.
* `match_t` - (Optional) The provider label match criteria for EPG. Allowed values are "All", "AtleastOne", "AtmostOne", "None". Default is "AtleastOne". Type - String.
* `name_alias` - (Optional) Name alias for object application epg. Type - String.
* `pc_enf_pref` - (Optional) The preferred policy control. Allowed values are "unenforced" and "enforced". Default is "unenforced". Type - String.
* `pref_gr_memb` - (Optional) Represents parameter used to determine if EPg is part of a group that does not a contract for communication. Allowed values are "exclude" and "include". Default is "exclude". Type - String.
* `prio` - (Optional) QoS priority class id. Allowed values are "unspecified", "level1", "level2", "level3", "level4","level5" and "level6". By default the value is inherited from the parent application profile. Type - String.
* `shutdown` - (Optional) Shutdown for object application epg. Allowed values are "yes" and "no". Default is "no". Type - String.

* `relation_fv_rs_bd` - (Optional) Relation to the Bridge domain associated with EPG (Point to class fvBD). This attribute is optional because the ACI API does not mandate it **but it is necessary for correct function of the resource.** Cardinality - N_TO_ONE. Type - String.

* `relation_fv_rs_cust_qos_pol` - (Optional) Relation to Custom Quality of Service traffic policy name (Point to class qosCustomPol). Cardinality - N_TO_ONE. Type - String.
<!-- tenant -> policies -> protocol -> Custom QoS -->

* `relation_fv_rs_fc_path_att` - (Optional) Relation to Fibre Channel (Paths) (Point to class fabricPathEp). Cardinality - N_TO_M. Type - Set of String.

* `relation_fv_rs_prov` - (Optional) Relation to Provided Contract (Point to class vzBrCP). Cardinality - N_TO_M. Type - Set of String.

* `relation_fv_rs_cons_if` - (Optional) Relation to Imported Contract (Point to class vzCPIf). Cardinality - N_TO_M. Type - Set of String.

* `relation_fv_rs_sec_inherited` - (Optional) Relation represents that the EPG is inheriting security configuration from other EPGs (Point to class fvEPg). Cardinality - N_TO_M. Type - Set of String.

* `relation_fv_rs_node_att` - (Optional) A block representing the relation to a Static Leaf binding (Point to class fabricNode). Cardinality - N_TO_M. Type: Block.

  - `node_dn` - (Required) The Distinguished Name of the Node object. Type: String.
  - `encap` - (Required) The port encapsulation of the Node Object. Type: String.
  - `description` - (Optional) The description for relation_fv_rs_node_att of the Node Object. Type: String.
  - `deployment_immediacy` - (Optional) The deployment immediacy of the Static Path of the Node Object. Allowed values: "immediate", "lazy". Default value: "lazy". Type: String.
  - `mode` - (Optional) The Application EPG mode of the Node Object. Allowed values: "regular", "native", "untagged". Default value: "regular". Type: String.
<!-- tenant -> Application Profile -> EPG ->Static Leaf -->

* `relation_fv_rs_dpp_pol` - (Optional) Relation to define a Data Plane Policing policy (Point to class qosDppPol). Cardinality - N_TO_ONE. Type - String.
<!-- tenant -> policies -> protocol -> Data Plane Policing -->

* `relation_fv_rs_cons` - (Optional) Relation to Consumed Contract (Point to class vzBrCP). Cardinality - N_TO_M. Type - Set of String.

* `relation_fv_rs_trust_ctrl` - (Optional) Relation to First Hop Security trust control (Point to class fhsTrustCtrlPol). Cardinality - N_TO_ONE. Type - String.
<!-- tenant -> policies -> protocol -> First Hop Security -->

* `relation_fv_rs_prot_by` - (Optional) Relation to Taboo Contract (Point to class vzTaboo). Cardinality - N_TO_M. Type - Set of String.

* `relation_fv_rs_aepg_mon_pol` - (Optional) Relation to create a container for monitoring policies associated with the tenant. This allows you to apply tenant-specific policies (Point to class monEPGPol). Cardinality - N_TO_ONE. Type - String.
<!-- tenant -> policies -> Monitoring -->

* `relation_fv_rs_intra_epg` - (Optional) Relation to Intra EPG Contract (Point to class vzBrCP). Cardinality - N_TO_M. Type - Set of String.


## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the Application EPG.

## Importing

An existing Application EPG can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_application_epg.example "<Dn>"
```
Starting in Terraform version 1.5, an existing EPG can be imported using [import blocks](https://developer.hashicorp.com/terraform/language/import) via the following configuration:

 ```hcl
 import {
    id = "<Dn>"
    to = aci_aci_application_epg.example
 }
 ```