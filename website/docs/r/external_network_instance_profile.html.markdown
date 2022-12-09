---
subcategory: "L3Out"
layout: "aci"
page_title: "ACI: aci_external_network_instance_profile"
sidebar_current: "docs-aci-resource-external_network_instance_profile"
description: |-
  Manages ACI External Network Instance Profile
---

# aci_external_network_instance_profile

Manages ACI External Network Instance Profile

## API Information ##

* `Class` - l3extOut
* `Distinguished Name` - uni/tn-{tenant}/out-{l3out}/instP-{external_epg}

## GUI Information ##

* `Location` - Tenant -> Networking -> L3Outs -> External EPGs

## Example Usage

```hcl
resource "aci_external_network_instance_profile" "external_epg" {
  l3_outside_dn = aci_l3_outside.l3_outside.id
  name          = "external_epg"
  exception_tag  = "2"
  flood_on_encap = "disabled"
  match_t        = "All"
  pref_gr_memb   = "exclude"
  prio           = "level1"
  target_dscp    = "unspecified"

  # Route Control Profile - Every direction(export/import) allows only one object
  relation_l3ext_rs_inst_p_to_profile {
    tn_rtctrl_profile_dn = aci_route_control_profile.route_control_profile_export.id
    direction            = "export"
  }
  relation_l3ext_rs_inst_p_to_profile {
    tn_rtctrl_profile_dn = aci_route_control_profile.route_control_profile_import.id
    direction            = "import"
  }

  relation_fv_rs_sec_inherited = [aci_external_network_instance_profile.external_epg_1.id]
  relation_fv_rs_cons_if = [aci_imported_contract.contract_interface.id]
  relation_fv_rs_prov = [aci_contract.web_contract.id]
  relation_fv_rs_cons = [aci_contract.web_contract.id]
  relation_fv_rs_prot_by = [aci_taboo_contract.taboo_contract.id]
}

```

## Argument Reference

* `l3_outside_dn` - (Required) Distinguished name of the parent L3Outside object.
* `name` - (Required) Name of the External Network Instance Profile object.
* `annotation` - (Optional) Annotation of the External Network Instance Profile object.
* `exception_tag` - (Optional) Exception tag of the External Network Instance Profile object. Allowed value range is "0" - "512".
* `flood_on_encap` - (Optional) Control at EPG level if the traffic L2 Multicast/Broadcast and Link Local Layer should be flooded only on ENCAP or based on bridg-domain settings. Allowed values are "disabled" and "enabled". Default value is "disabled".
* `match_t` - (Optional) The provider label match criteria of the External Network Instance Profile object. Allowed values are "All", "AtleastOne", "AtmostOne" and "None". Default is "AtleastOne".
* `name_alias` - (Optional) Name alias of the External Network Instance Profile object.
* `pref_gr_memb` - (Optional) Represents parameter used to determine if EPg is part of a group that does not a contract for communication. Allowed values are "include" and "exclude". Default is "exclude".
* `prio` - (Optional) The QoS priority class identifier of the External Network Instance Profile object. Allowed values are "unspecified", "level1", "level2", "level3", "level4", "level5" and "level6". Default is "unspecified".
* `target_dscp` - (Optional) The target differentiated services code point (DSCP) of the path attached to the layer 3 outside profile. Allowed values are "CS0", "CS1", "AF11", "AF12", "AF13", "CS2", "AF21", "AF22", "AF23", "CS3", "AF31", "AF32", "AF33", "CS4", "AF41", "AF42", "AF43", "CS5", "VA", "EF", "CS6", "CS7" and "unspecified". Default is "unspecified".
* `relation_fv_rs_sec_inherited` - (Optional) Relation to class fvEPg. Cardinality - N_TO_M. Type - Set of String.
* `relation_fv_rs_prov` - (Optional) Relation to class vzBrCP. Cardinality - N_TO_M. Type - Set of String.
* `relation_fv_rs_cons_if` - (Optional) Relation to class vzCPIf. Cardinality - N_TO_M. Type - Set of String.
* `relation_l3ext_rs_inst_p_to_profile` - (Optional) Relation to class rtctrlProfile. Every direction(export/import) allows only one object. Cardinality - N_TO_M. Type: Block.
  * tn_rtctrl_profile_name - (Deprecated) Name of the Route map for import and export route control.
  * tn_rtctrl_profile_dn - (Optional) Distinguished name of the Route map for import and export route control.
  * direction - (Optional) Direction of the Route Control Profile. Allowed values are "export" and "import".
* `relation_fv_rs_cons` - (Optional) Relation to class vzBrCP. Cardinality - N_TO_M. Type - Set of String.
* `relation_fv_rs_prot_by` - (Optional) Relation to class vzTaboo. Cardinality - N_TO_M. Type - Set of String.

## Attribute Reference

The only attribute that this resource exports is the `id`, which is set to the
Dn of the External Network Instance Profile.

## Importing

An existing External Network Instance Profile can be [imported][docs-import] into this resource via its Dn, via the following command:
[docs-import]: https://www.terraform.io/docs/import/index.html

```
terraform import aci_external_network_instance_profile.example <Dn>
```
