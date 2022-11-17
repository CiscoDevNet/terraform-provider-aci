## 2.5.2 (August 2, 2022)
BUG FIX::
- Fix aci_bulk_epg_to_static_path idempotency and default values when optional attributes are not provided

## 2.5.1 (July 29, 2022)
BUG FIX::
- Add documentation for aci_bulk_epg_to_static_path resource

## 2.5.0 (July 29, 2022)
IMPROVEMENTS:
- Add aci_vrf_leak_epg_bd_subnet resource and data source (leakRoutes, leakInternalSubnet and leakTo) (#900)
- Add resource aci_bulk_epg_to_static_path for bulk static path creation (#896)

## 2.4.0 (July 21, 2022)
IMPROVEMENTS:
- Allow nil return option for datasource aci_client_end_point (#893)
- Add next_hop_addr, msnlb and anycast_mac attributes to resource aci_subnet (#895)

BUG FIXES:
- Fix aci_imported_contract relation_vz_rs_if to properly set the relationship tDn (#894)
- Fix idempotency issues in aci_l3out_bgp_protocol_profile with relation_bgp_rs_best_path_ctrl_pol attribute (#904)
- Add documentation for relation_bgp_rs_best_path_ctrl_pol attribute of the aci_l3out_bgp_protocol_profile resource
- Improve unreachable error messages from aci-go-client

## 2.3.0 (June 11, 2022)
IMPROVEMENTS:
- Add datasource aci_l4_l7_deployed_graph_connector_vlan
- Add resource and datasource for aci_ip_sla_monitoring_policy (#881)
- Add resources and datasources aci_contract_subject_filter and aci_contract_subject_one_way_filter and support for one-way contracts in aci_contract_subject (#839).
- Add resources and datasources aci_l4_l7_redirect_health_group (vnsRedirectHealthGroup), aci_l4_l7_logical_interface (vnsLIf), aci_l4_l7_device (vnsLDevVip), aci_l4_l7_concrete_interface (vnsCIf) and aci_l4_l7_concrete_device (vnsCDev) (#861, #865, #866, #873, #877)
- Add set_dampening block attribute to the aci_action_rule_profile resource and datasource (#857)
- Add enable_vm_folder attribute to aci_vmm_domain (#888)

BUG FIXES:
- Add example for aci_user_security_domain and aci_security_domain_role
- Set filter_ids and filter_entry_ids to computed in aci_contract to fix idempotency issue (#883)

## 2.2.1 (May 13, 2022)
BUG FIXES:
- Fix 71 resources to not fail if object does not exist when refreshing state.

## 2.2.0 (May 9, 2022)
IMPROVEMENTS:
- Add Set As Path, Multipath, Next Hop Propagation, Set Communities, Set Next Hop, Set Metric Type, Set Metric, Set Preference, Set Weight and Set Route Tag options to aci_action_rule_profile (#851, #843)
- Add aci_action_rule_additional_communities resource and datasource (#840)
- Add aci_match_regex_community_term and aci_match_community_factor resources and datasources (#835)
- Add aci_aaep_to_domain resource and data source (#824)
- Add aci_epg_to_contract_interface resource and datasource (#833)
- Add deprecation message to attribute relation_fv_rs_path_att
- Add custom_epg_name attribute to resource aci_epg_to_domain (#841)
- Add relation_bgp_rs_best_path_ctrl_pol to aci_l3out_bgp_protocol_profile
- Add enhanced_lag_policy attribute to aci_epg_to_domain (#852)
- Add support for M1 Mac

BUG FIXES:
- Update docs for aci_application_epg (#842)
- Change documentation for 'managed' mode in aci_function_node
- Fix for "encap" in aci_epg_to_static_path should be "Required : true" (#845)
- Fix aci_route_control_context read and import function to retrieve set_rule and relation_rtctrl_rs_ctx_p_to_subj_p properly
- Add relation_vns_rs_l_dev_ctx_to_l_dev as Required attribute in aci_logical_device_context
- Fix idempotency by removing dhcp_option_ids from aci_dhcp_option_policy (#831)
- Fix to avoid vmmSecP object mapping with not supported domains in aci_epg_to_domain (#830)
- Fix for import of relation_infra_rs_sp_acc_grp in aci_spine_access_port_selector resource (#829)
- Fix aci_bgp_peer_connectivity_profile read changes issue

## 2.1.0 (March 16, 2022)
IMPROVEMENTS:
- Allow user to enter value between 0-255 for "prot" attribute in aci_filter_entry (#820)
- Add support for aci_vrf_to_bgp_address_family_context resource and data source
- Add aci_aaa_domain_relationship resource and data source to map AAA domain relationship for the parent object

BUG FIXES:
- Add option none to aci_bgp_address_family_context and aci_bfd_interface_policy ctrl attribute (#813)
- Update ctrl attribute definition in aci_ospf_interface_policy documentation and example (#816)
- Add capability to accept IPv4 and IPv6 addresses in aci_dhcp_relay_policy (#823)
- Fix aci_l3out_bfd_interface_profile relationship to bfd policy not created issue
- Add example of relation_l3ext_rs_dyn_path_att for resource aci_l3out_floating_svi
- Fix forged_transmit, mac_change and promiscuous_mode default values in aci_l3out_floating_svi
- Add input validation for af sub-attribute in aci_vrf relation_fv_rs_ctx_to_bgp_ctx_af_pol attribute and fix documentation and examples.

## 2.0.1 (February 27, 2022)
BUG FIXES:
- Fix some documentation examples identation
- Fix some import examples
- Update aci-go-client to v1.23.2 to improve retry mechanism
- Removed unused module attributes for aci_external_network_instance_profile

## 2.0.0 (January 27, 2022)
BREAKING CHANGE:
- Remove aci_application_epg unused relationship relation_fv_rs_graph_def to avoid idempotency issues.
- aci_bgp_peer_connectivity_profile attributes addr_t_ctrl, ctrl, peer_ctrl and private_a_sctrl changed from string to list of strings.
- aci_hsrp_interface_policy attribute ctrl changed from string to list of strings.
- aci_l3out_ospf_external_policy attribute area_ctrl changed from string to list of strings.
- aci_ospf_timers attribute ctrl changed from string to list of strings.
- aci_cloud_subnet attribute scope changed froms tring to list of strings.

Most of those changes will require changes to your Terraform plan and your state file.
At your own risk you can either manually modify your state file or use the following commands:
```
terraform state rm the_resource_type.name_of_your_resource
terraform import the_resource_type.name_of_your_resource dn_of_your_object
```

BEHAVIOR CHANGE:
- Add support for the new aci-go-client retries mechanism when connection fails or server errors in provider and set default to 2 retries

IMPROVEMENTS:
- Improve aci_l3out_vpc_member example
- ParentDn and relation updates to prepare for Terraformer support
- Rename aci_bpg_route_control_policy to aci_route_control_policy and deprecate aci_bpg_route_control_policy
- Rename aci_spine_port_selector to aci_spine_interface_profile_selector and deprecate aci_spine_port_selector
- Add support for relationship import in aci_bridge_domain
- Add inline block support for relation_l3ext_rs_dyn_path_att in aci_l3out_floating_svi
- Add new resources and data sources: aci_rest_managed, aci_spine_access_port_selector, aci_snmp_community
- Deprecation of aci_vrf_snmp_context_community
- Add spine_selector block to aci_spine_profile

BUG FIXES:
- Fix idempotency issue with area_id backbone/0.0.0.0 in aci_l3out_ospf_external_policy
- Fix idempotency issue with different IPv6 syntax in aci_subnet
- Fix ep_move_detect code location making aci_bridge_domain crash if bridge domain was not present on APIC
- Fix various empty relationship value when relationship was not set in resources creating idempotency issues
- Fix incorrect field extraction in the aci_client_end_point datasource
- Various documentation fixes

## 1.2.0 (December 13, 2021)
IMPROVEMENTS:
- Add new resources and data sources: aci_tag and aci_annotation

## 1.1.0 (December 10, 2021)
IMPROVEMENTS:
- Add new resources and data sources: aci_access_switch_policy_group, aci_authentication_properties, aci_bfd_interface_policy, aci_console_authentication, aci_coop_policy, aci_default_authentication, aci_duo_provider_group, aci_encryption_key, aci_endpoint_controls, aci_endpoint_ip_aging_profile, aci_endpoint_loop_protection, aci_error_disable_recovery, aci_fabric_node_control, aci_fabric_wide_settings, aci_file_remote_path, aci_global_security, aci_interface_blacklist, aci_isis_domain_policy, aci_l3_interface_policy, aci_ldap_group_map, aci_ldap_group_map_rule, aci_ldap_group_map_rule_to_group_map, aci_ldap_provider, aci_login_domain, aci_login_domain_provider, aci_managed_node_connectivity_group, aci_mcp_instance_policy, aci_mgmt_preference, aci_mgmt_zone, aci_port_tracking, aci_qos_instance_policy, aci_radius_provider, aci_radius_provider_group, aci_recurring_window, aci_rsa_provider, aci_saml_provider, aci_saml_provider_group, aci_spine_switch_policy_group, aci_tacacs_accounting, aci_tacacs_accounting_destination, aci_tacacs_provider, aci_tacacs_provider_group, aci_tacacs_source, aci_user_security_domain, aci_user_security_domain_role, aci_vpc_domain_policy, aci_vrf_snmp_context, aci_vrf_snmp_context_community, aci_match_rule, aci_match_route_destination_rule, aci_route_control_context

## 1.0.1 (November 09, 2021)
BUG FIXES:
- Fix aci_cloud_vpn_gateway documentation subcategory issue

## 1.0.0 (November 09, 2021)
BREAKING CHANGE:
- Migration to Terraform Provider SDK v2. Remove support for Terraform v0.11.x or below
- Fix and update netflow monitor relation in aci_leaf_access_port_policy_group and aci_leaf_access_bundle_policy_group
- Fix tcp_rules from string to list in aci_filter_entry

IMPROVEMENTS:
- Add ESG Tag Selector and ESG EPG Selector resources and data sources
- Add support for admin_state attribute and relation_bgp_rs_peer_to_profile relation
- Add support for aci_bgp_peer_connectivity_profile to be defined at interface level and node level
- Add ability to disable endpoint learning (garp) in aci_bridge_domain
- Add support for level4 - level6 to aci_application_epg prio attribute
- Deprecate tn_rtctrl_profile_name and add replacement tn_rtctrl_profile_dn in aci_subnet
- Add references for provider_profile_dn in vmm_domain
- Update dependancy versions

BUG FIXES:
- Fix multiple idempotency issues across resources
- Fix issues found during testing of resources with TF provider SDK v2
- Add forced replacement of resource if path or class_name is changed in aci_rest
- Fix ESG Selector required parameter and documentation
- Fix VMM Controller descr argument not supported issue
- Fix ASN and Local ASN update function
- Fix multiple documentation issues
- Fix examples formating in examples directory and add examples for resources without examples

## 0.7.1 (June 25, 2021)
BREAKING CHANGE:
- Change aci_dhcp_relay_policy relation_dhcp_rs_prov argument from list of string to block definition to accomodate the addr argument.

BUG FIXES:
- Fix a regression introduced in aci_rest creating issues when use in parallel.
- Make management_profile_dn an optional parameter with "uni/tn-mgmt/mgmtp-default" as default value in aci_node_mgmt_epg.
- Deprecate use of filter argument in aci_contract and removal from documentation.
- Fix documentation of region argument in aci_cloud_aws_provider.
- Fix aci_bgp_peer_connectivity_profile documentation for as_number and local_asn.
- Fix aci_application_epg examples and documentation to make it clearer.
- Fix cert_name usage examples in documentation and README.md.
- Remove application_epg_dn argument from aci_client_end_point documentation as use case is not implemented yet (follow #513 for use case development).
- Add required_provider and provider definition in all examples to conform to new Terraform provider usage definitions.
- Fix aci_contract example to showcase how to create contract, subject, filter and filter entries.

## 0.7.0 (May 26, 2021)
BREAKING CHANGE:
- Fix "ctrl" attribute issues with list of items in OSPF Interface Policy resource/datasource.
- Fix "enforce_rtctrl" attribute issues with list of items in L3 Outside resource/datasource.
- Change aci_stp_if_pol resource name to aci_spanning_tree_interface_policy name

BUG FIXES:
- Update aci_spanning_tree_interface_policy documentation to add description attribute.
- Stop control(ctrl) from being added repeatedly when set to "unspecified" in OSPF Interface Policy, Subnet and STP Interface Policy resources.
- Add alloc_mode in documentation of aci_vlan_pool datasource and update resource example.
- Fix aci_l3out_path_attachment to accept custom MTU values.
- Fix relation_infra_rs_spine_acc_node_p_grp issue in aci_spine_switch_association resource.

## 0.6.0 (May 11, 2021)
IMPROVEMENTS:
- Updated documentation and examples for new terraform required_provider syntax.
- Add new resources for Cloud ACI VGW, L3Outs, L2Outs, routing, Service Graphs, ESGs, STP Interface Policy, DHCP options, DHCP relay, DHCP labels, breakout, OOB/inband EPG and VMM domain policies.
- Add vPC support for aci_fabric_path_ep

BUG FIXES:
- Fixed a few documentation issues.
- Ignore REST errors on destroy for object that cannot be deleted.
- Diverse fixes for issues.

## 0.5.4 (January 13, 2021)

BUG FIXES:
- Added Missing documentation for aci_monitoring_policy resource.

## 0.5.3 (December 22, 2020)

IMPROVEMENTS:
- Added New attribute named endpoint_path to fvcep data-source.
- Added More levels for priorities to the application_profile resource. (Supported in latest version of APIC)

BUG FIXES:
- Renamed `_from` attribute to `from` for aci_ranges resource.

BREAKING CHANGES:
- scope attribute for aci_l3_ext_subnet resource is now list of string rather than a single string. This change will break your infrastructure if you have l3extsubnet created with terraform. Consider removing the l3extsubnet resource from your terraform state file using `terraform state rm` and than run the `terraform apply` to make your configuration inline with the new changes. This will not affect the l3extsubnet which is already there.

## 0.5.2 (November 20, 2020)

BUG FIXES:
- Fixed an issue with aci_subnet ctrl attribute to have list value.
- Fixed an issue with aci_any relations being not created.
- Fixed an issue with aci_cloud_subnet to have name attribute.

## 0.5.1 (November 05, 2020)

IMPROVEMENTS:
- Added new data-source for fvCEP resource..

BUG FIXES:
- Fixed an issue with aci_physical_domain and aci_l3_domain_profile about unknown attribute error.


## 0.5.0 (October 23, 2020)

IMPROVEMENTS:
- Added new resources Spine Switch profiles and interfaces, L4-L7 interfaces.
- access_port_block have default name attribute with auto incrementor.
- Added resources to manage FEX profiles.

BUG FIXES:
- Fixed an issue with docs being not rendered via name in Hashicorp registry.
- Fixed an issue with subnet scope attribute to have list value.
- Fixed all the bugs reported.

## 0.4.1 (September 23, 2020)

IMPROVEMENTS:
- First Terraform Registry release.

## 0.4.0 (September 16, 2020)

IMPROVEMENTS:
- Improved checks in the parameters.
- Added resources to manage FEX profiles.

BUG FIXES:
- Fixed an issue with parameters not getting updated on first run.
- Fixed typo errors in documentations.

## 0.3.4 (July 20, 2020)

IMPROVEMENTS:
- Parameter `relation_cloud_rs_to_ctx` works on id now for Cloud Context Profile resource.

BREAKING CHANGES:
- Renamed all the t_dn attributes to tdn.

## 0.3.3 (July 16, 2020)

IMPROVEMENTS:
- Added zone parameter to cloud_subnet resource for APIC v5.0 or higher.

BREAKING CHANGES:
- Renamed all the e_pg attributes to epg.

## 0.3.2 (July 06, 2020)

IMPROVEMENTS:
- Updated objet model payload for l3out and vmmdomain relations.

BUG FIXES:
- Fixed the issue with vzany not updated in first run.
- FIxed the issue with switch id replaced while creating multiple switches.
## 0.3.1 (June 24, 2020)

IMPROVEMENTS:
- Updated object model for all the relation attributes compatible with new APIC versions.

## 0.3.0 (June 17, 2020)

IMPROVEMENTS:
- Added support for inline creation of filter and filter entry with contract.
- Added new resource to manage relations from epg to domain and contract with more control.
- aci_rest now supports more generic YAML/JSON payload.
- All the relation supports id only.

BUG FIXES:
- Fixed issues with domain and leaf attachment.
## 0.2.3 (May 19, 2020)

IMPROVEMENTS:
- Added new resource to manage imported_contracts.

## 0.2.2 (May 11, 2020)

BREAKING CHANGES:
- Renamed the aci_cloud_epg, aci_cloud_external_epg, aci_cloud_endpoint_selectorfor_external_epgs resources, removed an extra `_` in epg. New names for these resources will be aci_cloud_epg, aci_cloud_external_epg, aci_cloud_endpoint_selectorfor_external_epgs respectively.

IMPROVEMENTS:
- Removed the implicit status insertion for aci_rest resource.

BUG FIXES:
- Fixed the issue with l3extRsL3DomAtt not attaching properly.
## 0.2.1 (April 15, 2020)

IMPROVEMENTS:
- Added new resources for static leaf attachment, l3out profile, aci_any.
- Added support for inline private key for authentication.
## 0.2.0 (April 07, 2020)

BUG FIXES:

- Added singleton implementation for authentication endpoint.
## 0.1.8 (April 02, 2020)

IMPROVEMENTS:
- Added new modules for managing fabric and APIC management objects.
## 0.1.7 (January 27, 2020)
BUG FIXES:

- Fixed the issue with new Rn format for CloudExtEpgSelector class.
## 0.1.6 (January 25, 2020)

IMPROVEMENTS:
- Added support for new cipher suites and TLS version for the new release of cloud APIC.
## 0.1.5 (January 22, 2020)

IMPROVEMENTS:

- Added logic to handle panics and show proper error messages.
## 0.1.4 (December 20, 2019)
BUG FIXES:

- Fixed crashing of Terraform while using cert based authentication.

IMPROVEMENTS:

- Switched to terraform-plugin-sdk instead of legacy terraform package. 
## 0.1.3 (December 18, 2019)
BUG FIXES:

- Fixed issue of having 405 errors from APIC nginx.

## 0.1.2 (November 04, 2019)

BUG FIXES:

- Fixed issue of hanging sessions with Terraform 0.12.
## 0.1.1 (September 19, 2019)

IMPROVEMENTS:

- Added Docs for aci_rest resource.
- Markdown improvements.

BUG FIXES:

- Fixed issue of Terraform crashing while creating L3 Subnet.
## 0.1.0 (July 22, 2019)

- Initial Release
