terraform {
  required_providers {
    aci = {
      source = "CiscoDevNet/aci"
    }
  }
}

# Configure the provider with your Cisco APIC credentials.
provider "aci" {
  # APIC Username
  username = " "
  # APIC Password
  password = " "
  # APIC URL
  url      = " "
  insecure = true
}

# Defines an Application EPG Resource.
/*
resource "aci_application_epg" "fooapplication_epg" {
    parent_dn  = "uni/tn-common/ap-default"
    name                    = "new_epg"
    description             = "from terraform"
    annotation              = "tag_epg"
    contract_exception_tag  = "0"
    flood_in_encapsulation  = "disabled"
    forwarding_control      = "none"
    
    relation_to_bridge_domain= [{
      annotation = "annotation1"
      bridge_domain_name = "default"
    }]
}
*/
/*
resource "aci_application_epg" "fooapplication_epg2" {
    parent_dn  = "uni/tn-common/ap-default222"
    name                    = "new_epg2"
    description             = "from terraform"
    annotation              = "tag_epg"
    contract_exception_tag  = "0"
    flood_in_encapsulation  = "disabled"
    forwarding_control      = "none"
    
    relation_to_bridge_domain= [{
      annotation = "annotation1"
      bridge_domain_name = "default"
    }]
}

*/
# Define an ACI Tenant Resource.

/*
resource "aci_endpoint_tag_ip" "full_example_tenant" {
  parent_dn    = "uni/tn-common"
  annotation   = "annotation"
  vrf_name     = "test_ctx_name"
  id_attribute = "1"
  ip           = "10.0.0.2"
  name         = "WOW"
  name_alias   = "name_alias"
  annotations = [
    {
      key   = "key_0"
      value = "vaule_1"
    }
  ]
  tags = [
    {
      key   = "key_0"
      value = "value_1"
    }
  ]
}
*/

//SUCCESS
/*
resource "aci_vrf_fallback_route_group" "full_example_vrf" {
  parent_dn   = "uni/tn-common"
  annotation  = "annotation"
  description = "description"
  name        = "fallback_route_group"
  name_alias  = "name_alias"
  vrf_fallback_route_group_members = [
    {
      annotation      = "annotation_1"
      description     = "description_1"
      name            = "name_1"
      name_alias      = "name_alias_1"
      fallback_member = "2.2.2.2"
    }
  ]
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

resource "aci_external_management_network_instance_profile" "full_example" {
  annotation  = "woidid"
  description = "description"
  name        = "test_name"
  name_alias  = "name_alias"
  priority    = "level1"
  relation_to_consumed_out_of_band_contracts = [
    {
      annotation                = "annotation_1"
      priority                  = "level1"
      out_of_band_contract_name = "aci_out_of_band_contract.example.name"
    }
  ]
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
*/
/* SUCCESS

resource "aci_external_management_network_instance_profile" "full_example333" {
  annotation  = "annotation"
  description = "description"
  name        = "test_name"
  name_alias  = "name_alias"
  priority    = "level1"
  relation_to_consumed_out_of_band_contracts = [
    {
      annotation                = "annotation_1"
      priority                  = "level1"
      out_of_band_contract_name = "aci_out_of_band_contract.example.name"
    }
  ]
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

*/

/*
resource "aci_netflow_monitor_policy" "full_example_tenant" {
  parent_dn   = "uni/tn-mgmt/extmgmt-default/instp-test_name"
  annotation  = "annotation"
  description = "description"
  name        = "netfow_monitor"
  name_alias  = "name_alias"
  owner_key   = "owner_key"
  owner_tag   = "owner_tag"
  relation_to_netflow_exporters = [
    {
      annotation                   = "annotation_1"
      netflow_exporter_policy_name = "aci_netflow_exporter_policy.example.name"
    }
  ]
  relation_to_netflow_record = [
    {
      annotation                 = "annotation_1"
      netflow_record_policy_name = "aci_netflow_record_policy.example.name"
    }
  ]
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
*/


/*
resource "aci_tag" "example_tenant" {
  parent_dn = "uni/tn-example_tenant"
  key       = "test_key"
  value     = "test_value"
}
*/


/*



*/

//DEFINES AN ACI ANNOTATION ------- TEST


resource "aci_annotation" "terraform_annotation" {
  parent_dn = "uni/tn-common/eptags/epiptag-[10.0.0.2]-test_ctx_name"
  key       = "test_key"
  value     = "test_value"
}


resource "aci_annotation" "terraform_annotation2" {
  parent_dn = "uni/tn-common/eptags/epiptag-[10.0.0.2]-test_ctx_name/annotationKey-[test_key]"
  key       = "test_keya"
  value     = "test_valuea"
}



resource "aci_pim_route_map_policy" "full_example_tenant" {
  parent_dn   = "uni/tn-demo_tenant"
  annotation  = "annotation"
  description = "description"
  name        = "test_name"
  name_alias  = "name_alias"
  owner_key   = "owner_key"
  owner_tag   = "owner_tag"
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


# Define an ACI Tenant VRF Resource.
# resource "aci_vrf" "terraform_vrf" {
#     tenant_dn   = aci_tenant.terraform_tenant.id
#     description = "VRF Created Using Terraform"
#     name        = var.vrf
# }

# # Define an ACI Tenant BD Resource.
# resource "aci_bridge_domain" "terraform_bd" {
#     tenant_dn          = aci_tenant.terraform_tenant.id
#     relation_fv_rs_ctx = aci_vrf.terraform_vrf.id
#     description        = "BD Created Using Terraform"
#     name               = var.bd
# }

# # Define an ACI Tenant BD Subnet Resource.
# resource "aci_subnet" "terraform_bd_subnet" {
#     parent_dn   = aci_bridge_domain.terraform_bd.id
#     description = "Subnet Created Using Terraform"
#     ip          = var.subnet
# }

# # Define an ACI Filter Resource.
# resource "aci_filter" "terraform_filter" {
#     for_each    = var.filters
#     tenant_dn   = aci_tenant.terraform_tenant.id
#     description = "This is filter ${each.key} created by terraform"
#     name        = each.value.filter
# }

# # Define an ACI Filter Entry Resource.
# resource "aci_filter_entry" "terraform_filter_entry" {
#     for_each      = var.filters
#     filter_dn     = aci_filter.terraform_filter[each.key].id
#     name          = each.value.entry
#     ether_t       = "ipv4"
#     prot          = each.value.protocol
#     d_from_port   = each.value.port
#     d_to_port     = each.value.port
# }

# # Define an ACI Contract Resource.
# resource "aci_contract" "terraform_contract" {
#     for_each      = var.contracts
#     tenant_dn     = aci_tenant.terraform_tenant.id
#     name          = each.value.contract
#     description   = "Contract created using Terraform"
# }

# # Define an ACI Contract Subject Resource.
# resource "aci_contract_subject" "terraform_contract_subject" {
#     for_each                      = var.contracts
#     contract_dn                   = aci_contract.terraform_contract[each.key].id
#     name                          = each.value.subject
#     relation_vz_rs_subj_filt_att  = [aci_filter.terraform_filter[each.value.filter].id]
# }

# # Define an ACI Application Profile Resource.
# resource "aci_application_profile" "terraform_ap" {
#     tenant_dn  = aci_tenant.terraform_tenant.id
#     name       = var.ap
#     description = "App Profile Created Using Terraform"
# }

 

# # Associate the EPG Resources with a VMM Domain.
# resource "aci_epg_to_domain" "terraform_epg_domain" {
#     for_each              = var.epgs
#     application_epg_dn    = aci_application_epg.terraform_epg[each.key].id
#     tdn   = "uni/vmmp-VMware/dom-aci_terraform_lab"
# }

# # Associate the EPGs with the contrats
# resource "aci_epg_to_contract" "terraform_epg_contract" {
#     for_each           = var.epg_contracts
#     application_epg_dn = aci_application_epg.terraform_epg[each.value.epg].id
#     contract_dn        = aci_contract.terraform_contract[each.value.contract].id
#     contract_type      = each.value.contract_type
# }

