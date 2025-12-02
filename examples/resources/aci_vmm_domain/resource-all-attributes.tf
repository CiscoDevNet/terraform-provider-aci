
resource "aci_vmm_domain" "full_example" {
  parent_dn                        = "uni/vmmp-VMware"
  access_mode                      = "read-only"
  annotation                       = "annotation"
  arp_learning                     = ""
  ave_time_out                     = "30"
  configure_infra_port_groups      = "no"
  endpoint_data_path_verification  = "epDpVerify"
  custom_switch_name               = ""
  delimiter                        = "@"
  enable_ave_mode                  = "no"
  enable_tag_collection            = "no"
  enable_vm_folder_data_retrieval  = "no"
  encapsulation_mode               = "ivxlan"
  switching_enforcement_preference = "hw"
  endpoint_inventory_type          = "none"
  endpoint_retention_time          = "0"
  host_availability_assurance      = "no"
  multicast_address                = "224.0.1.0"
  switch_type                      = "cf"
  name                             = "test_name"
  name_alias                       = "name_alias_1"
  owner_key                        = "owner_key_1"
  owner_tag                        = "owner_tag_1"
  default_encapsulation_mode       = "unspecified"
  relation_to_ip_address_pool = {
    annotation = "annotation_1"
    target_dn  = aci_ip_address_pool.test_ip_address_pool_0.id
    annotations = [
      {
        key   = "key_0"
        value = "value_1"
      }
    ]
    tags = [
      {
        key   = "key_0"
        value = "value_1"
      }
    ]
  }
  relation_to_vlan_pool = {
    annotation = "annotation_1"
    target_dn  = aci_vlan_pool.test_vlan_pool_1.id
    annotations = [
      {
        key   = "key_0"
        value = "value_1"
      }
    ]
    tags = [
      {
        key   = "key_0"
        value = "value_1"
      }
    ]
  }
  annotations = [
    {
      key   = "key_0"
      value = "value_1"
    }
  ]
  tags = [
    {
      key   = "key_0"
      value = "value_1"
    }
  ]
  relation_to_multicast_pool = {
    annotation = "test_value_for_child"
    target_dn  = "target_dn_0"
    annotations = [
      {
        key   = "key_0"
        value = "value_1"
      }
    ]
    tags = [
      {
        key   = "key_0"
        value = "value_1"
      }
    ]
  }
  relation_to_lacp_enhanced_lag_policy = {
    annotation = "test_value_for_child"
    target_dn  = "target_dn_0"
    annotations = [
      {
        key   = "key_0"
        value = "value_1"
      }
    ]
    tags = [
      {
        key   = "key_0"
        value = "value_1"
      }
    ]
  }
  vmm_uplink_container = {
    annotation        = "annotation_1"
    name_alias        = "name_alias_1"
    number_of_uplinks = "1"
    annotations = [
      {
        key   = "key_0"
        value = "value_1"
      }
    ]
    tags = [
      {
        key   = "key_0"
        value = "value_1"
      }
    ]
    uplink_policies = [
      {
        annotation  = "annotation_1"
        name_alias  = "name_alias_1"
        uplink_id   = "2"
        uplink_name = "uplink_name_2"
        annotations = [
          {
            key   = "key_0"
            value = "value_1"
          }
        ]
        tags = [
          {
            key   = "key_0"
            value = "value_1"
          }
        ]
      }
    ]
  }
}
