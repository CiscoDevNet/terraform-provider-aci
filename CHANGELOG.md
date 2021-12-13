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
