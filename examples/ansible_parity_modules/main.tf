provider "aci" {
  username = "admin"
  password = "ins3965!ins3965!"
  url      = "https://173.36.219.26/"
  insecure = true
}


resource "aci_local_user" "foolocal_user" {
  name          = "user_demo1"
  pwd           = "123456786"
  description   = "This user is created by terraform"
}

resource "aci_pod_maintenance_group" "foopod_maintenance_group" {
  name = "maintGrp1"
}

resource "aci_maintenance_policy" "foomaintenance_policy" {
  name = "maintP1"
}

resource "aci_tenant" "test_tenant1" {
  name        = "tf_test_rel_tenant2"
  description = "This tenant is created by terraform"
}

resource "aci_monitoring_policy" "foomonitoring_policy" {
  tenant_dn = aci_tenant.test_tenant1.id
  name = "monepgpol1"
}

resource "aci_action_rule_profile" "fooaction_rule_profile" {
  tenant_dn = aci_tenant.test_tenant1.id
  name = "rtctrlAttrP1"
}

resource "aci_trigger_scheduler" "footrigger_scheduler" {
  name = "trigSchedP"
}

resource "aci_physical_domain" "foophysical_domain" {
  name = "physDomP1"
}

resource "aci_taboo_contract" "footaboo_contract" {
  tenant_dn = aci_tenant.test_tenant1.id
  name = "vztaboo1"
}

resource "aci_leaf_profile" "tf_leaf_prof" {
    name = "tf_leaf_prof1"
}

resource "aci_switch_association" "fooswitch_association" {
  leaf_profile_dn = aci_leaf_profile.tf_leaf_prof.id
  name = "infraLeafs1"
  switch_association_type = "ALL"
}

resource "aci_span_destination_group" "foospan_destination_group" {
  tenant_dn = aci_tenant.test_tenant1.id
  name = "spanDestGrp1"
}

resource "aci_span_source_group" "foospan_source_group" {
  tenant_dn = aci_tenant.test_tenant1.id
  name = "spanSrcGrp1"
}

resource "aci_span_sourcedestination_group_match_label" "foospan_sourcedestination_group_match_label" {
  span_source_group_dn = aci_span_source_group.foospan_source_group.id
  name = "spanLbl1"
}

resource "aci_vlan_pool" "foovlan_pool" {
  name       = "vlanInstP1"
  alloc_mode = "static"
}

resource "fcDomP" "fooranges" {
  vlan_pool_dn = aci_vlan_pool.foovlan_pool.id
  _from        = "vlan-25"
  to           = "vlan-35"
  alloc_mode   = "inherit"
}

resource "aci_vxlan_pool" "foovxlan_pool" {
  name = "vxlanInstP1"
}

resource "aci_vsan_pool" "foovsan_pool" {
  name = "vsanInstP1"
  alloc_mode = "static"
}

resource "aci_attachable_access_entity_profile" "tst1" {
    name = "infraAttEntity1"
}

resource "aci_firmware_group" "foofirmware_group" {
  name = "firmwareFwGrp1"
}

resource "aci_firmware_policy" "foofirmware_policy" {
  name = "firmwareFwp1"
}

resource "aci_firmware_download_task" "foofirmware_download_task" {
  name = "firmwareOSource1"
}

resource "aci_fc_domain" "foofc_domain" {
  name = "fcDomP1"
}

resource "aci_fabric_node_member" "foofabric_node_member" {
  name = "te1"
  serial = "127"
  node_id = "127"
}

resource "aci_configuration_export_policy" "fooconfiguration_export_policy" {
  name = "configExportP1"
}

resource "aci_cdp_interface_policy" "foocdp_interface_policy" {
  name = "cdpIfPol1"
}

resource "aci_leaf_interface_profile" "test_leaf_profile" {
    name = "tf_leaf"
}

resource "aci_access_port_selector" "test_selector" {
    leaf_interface_profile_dn = aci_leaf_interface_profile.test_leaf_profile.id
    name = "tf_test"
    access_port_selector_type = "default"
}

resource "aci_access_sub_port_block" "fooaccess_sub_port_block" {
  access_port_selector_dn = aci_access_port_selector.test_selector.id
  name = "infraSubportBlk1"
}

resource "aci_vpc_explicit_protection_group" "foovpc_explicit_protection_group" {
  name = "FabricExplicitGrp1"
  switch1 = "145"
  switch2 = "123"
  vpc_explicit_protection_group_id = "10"
}

resource "aci_node_block_maintgrp" "foonode_block_maintgrp" {
  pod_maintenance_group_dn = aci_pod_maintenance_group.foopod_maintenance_group.id
  name = "fabricNodeBlkMG"
}

resource "aci_node_block_firmware" "foonode_block_firmware" {
 firmware_group_dn = aci_firmware_group.foofirmware_group.id
 name = "fabricNodeBlkFW" 
}

resource "aci_configuration_import_policy" "fooconfiguration_import_policy" {
  name = "configImportP1"
}

resource "aci_l3_domain_profile" "fool3_domain_profile" {
  name = "l3extDomp1"
}
