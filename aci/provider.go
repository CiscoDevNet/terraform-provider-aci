package aci

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Username for the APIC Account. This can also be set as the ACI_USERNAME environment variable.",
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Password for the APIC Account. This can also be set as the ACI_PASSWORD environment variable.",
			},
			"url": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "URL of the Cisco ACI web interface. This can also be set as the ACI_URL environment variable.",
			},
			"insecure": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Allow insecure HTTPS client. This can also be set as the ACI_INSECURE environment variable. Defaults to `true`.",
			},
			"private_key": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Private key path for signature calculation. This can also be set as the ACI_PRIVATE_KEY environment variable.",
			},
			"cert_name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Certificate name for the User in Cisco ACI. This can also be set as the ACI_CERT_NAME environment variable.",
			},
			"proxy_url": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Proxy Server URL with port number. This can also be set as the ACI_PROXY_URL environment variable.",
			},
			"proxy_creds": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Proxy server credentials in the form of username:password. This can also be set as the ACI_PROXY_CREDS environment variable.",
			},
			"validate_relation_dn": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Flag to validate if a object with entered relation Dn exists in the APIC. Defaults to `true`.",
			},
			"retries": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Number of retries for REST API calls. This can also be set as the ACI_RETRIES environment variable. Defaults to `2`.",
			},
			"annotation": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Global annotation for the provider. This can also be set as the ACI_ANNOTATION environment variable.",
			},
			"allow_existing_on_create": &schema.Schema{
				Type:        schema.TypeBool,
				Description: "Allow existing objects to be managed. This can also be set as the ACI_ALLOW_EXISTING_ON_CREATE environment variable.",
				Optional:    true,
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"aci_contract":                                 resourceAciContract(),
			"aci_contract_subject":                         resourceAciContractSubject(),
			"aci_contract_subject_filter":                  resourceAciSubjectFilter(),
			"aci_contract_subject_one_way_filter":          resourceAciFilterRelationship(),
			"aci_subnet":                                   resourceAciSubnet(),
			"aci_filter":                                   resourceAciFilter(),
			"aci_filter_entry":                             resourceAciFilterEntry(),
			"aci_vmm_domain":                               resourceAciVMMDomain(),
			"aci_vmm_controller":                           resourceAciVMMController(),
			"aci_vswitch_policy":                           resourceAciVSwitchPolicyGroup(),
			"aci_vrf_to_bgp_address_family_context":        resourceAciBGPAddressFamilyContextPolicyRelationship(),
			"aci_rest":                                     resourceAciRest(),
			"aci_external_network_instance_profile":        resourceAciExternalNetworkInstanceProfile(),
			"aci_l3_outside":                               resourceAciL3Outside(),
			"aci_bfd_multihop_interface_profile":           resourceAciBfdMultihopInterfaceProfile(),
			"aci_bfd_multihop_interface_policy":            resourceAciBfdMultihopInterfacePolicy(),
			"aci_bfd_multihop_node_policy":                 resourceAciBFDMultihopNodePolicy(),
			"aci_interface_fc_policy":                      resourceAciInterfaceFCPolicy(),
			"aci_leaf_access_bundle_policy_group":          resourceAciPCVPCInterfacePolicyGroup(),
			"aci_leaf_access_bundle_policy_sub_group":      resourceAciOverridePCVPCPolicyGroup(),
			"aci_leaf_access_port_policy_group":            resourceAciLeafAccessPortPolicyGroup(),
			"aci_lldp_interface_policy":                    resourceAciLLDPInterfacePolicy(),
			"aci_miscabling_protocol_interface_policy":     resourceAciMiscablingProtocolInterfacePolicy(),
			"aci_ospf_interface_policy":                    resourceAciOSPFInterfacePolicy(),
			"aci_lacp_policy":                              resourceAciLACPPolicy(),
			"aci_lacp_member_policy":                       resourceAciLACPMemberPolicy(),
			"aci_port_security_policy":                     resourceAciPortSecurityPolicy(),
			"aci_leaf_profile":                             resourceAciLeafProfile(),
			"aci_end_point_retention_policy":               resourceAciEndPointRetentionPolicy(),
			"aci_vlan_encapsulationfor_vxlan_traffic":      resourceAciVlanEncapsulationforVxlanTraffic(),
			"aci_logical_node_profile":                     resourceAciLogicalNodeProfile(),
			"aci_logical_interface_profile":                resourceAciLogicalInterfaceProfile(),
			"aci_l3_ext_subnet":                            resourceAciL3ExtSubnet(),
			"aci_cloud_applicationcontainer":               resourceAciCloudApplicationcontainer(),
			"aci_cloud_ipsec_tunnel_subnet_pool":           resourceAciSubnetPoolforIpSecTunnels(),
			"aci_cloud_external_network":                   resourceAciCloudTemplateforExternalNetwork(),
			"aci_cloud_external_network_vpn_network":       resourceAciCloudTemplateforVPNNetwork(),
			"aci_cloud_aws_provider":                       resourceAciCloudAWSProvider(),
			"aci_cloud_cidr_pool":                          resourceAciCloudCIDRPool(),
			"aci_cloud_domain_profile":                     resourceAciCloudDomainProfile(),
			"aci_cloud_context_profile":                    resourceAciCloudContextProfile(),
			"aci_cloud_epg":                                resourceAciCloudEPg(),
			"aci_cloud_endpoint_selectorfor_external_epgs": resourceAciCloudEndpointSelectorforExternalEPgs(),
			"aci_cloud_endpoint_selector":                  resourceAciCloudEndpointSelector(),
			"aci_cloud_external_epg":                       resourceAciCloudExternalEPg(),
			"aci_cloud_service_epg":                        resourceAciCloudServiceEPg(),
			"aci_cloud_service_endpoint_selector":          resourceAciCloudServiceEndpointSelector(),
			"aci_cloud_private_link_label":                 resourceAciCloudPrivateLinkLabel(),
			"aci_cloud_subnet":                             resourceAciCloudSubnet(),
			"aci_cloud_account":                            resourceAciCloudAccount(),
			"aci_tenant_to_cloud_account":                  resourceAciTenantToCloudAccountAssociation(),
			"aci_cloud_ad":                                 resourceAciCloudActiveDirectory(),
			"aci_cloud_credentials":                        resourceAciCloudCredentials(),
			"aci_local_user":                               resourceAciLocalUser(),
			"aci_pod_maintenance_group":                    resourceAciPODMaintenanceGroup(),
			"aci_maintenance_policy":                       resourceAciMaintenancePolicy(),
			"aci_monitoring_policy":                        resourceAciMonitoringPolicy(),
			"aci_physical_domain":                          resourceAciPhysicalDomain(),
			"aci_action_rule_profile":                      resourceAciActionRuleProfile(),
			"aci_trigger_scheduler":                        resourceAciTriggerScheduler(),
			"aci_leaf_selector":                            resourceAciSwitchAssociation(),
			"aci_span_destination_group":                   resourceAciSPANDestinationGroup(),
			"aci_span_source_group":                        resourceAciSPANSourceGroup(),
			"aci_span_sourcedestination_group_match_label": resourceAciSPANSourcedestinationGroupMatchLabel(),
			"aci_vlan_pool":                                resourceAciVLANPool(),
			"aci_vxlan_pool":                               resourceAciVXLANPool(),
			"aci_vsan_pool":                                resourceAciVSANPool(),
			"aci_multicast_pool":                           resourceAciMulticastAddressPool(),
			"aci_multicast_pool_block":                     resourceAciMulticastAddressBlock(),
			"aci_firmware_group":                           resourceAciFirmwareGroup(),
			"aci_firmware_policy":                          resourceAciFirmwarePolicy(),
			"aci_firmware_download_task":                   resourceAciFirmwareDownloadTask(),
			"aci_fc_domain":                                resourceAciFCDomain(),
			"aci_configuration_export_policy":              resourceAciConfigurationExportPolicy(),
			"aci_cdp_interface_policy":                     resourceAciCDPInterfacePolicy(),
			"aci_access_sub_port_block":                    resourceAciAccessSubPortBlock(),
			"aci_maintenance_group_node":                   resourceAciNodeBlockMG(),
			"aci_node_block_firmware":                      resourceAciNodeBlockFW(),
			"aci_vpc_explicit_protection_group":            resourceAciVPCExplicitProtectionGroup(),
			"aci_configuration_import_policy":              resourceAciConfigurationImportPolicy(),
			"aci_fabric_node_member":                       resourceAciFabricNodeMember(),
			"aci_ranges":                                   resourceAciRanges(),
			"aci_l3_domain_profile":                        resourceAciL3DomainProfile(),
			"aci_x509_certificate":                         resourceAciX509Certificate(),
			"aci_epg_to_static_path":                       resourceAciStaticPath(),
			"aci_logical_node_to_fabric_node":              resourceAciFabricNode(),
			"aci_imported_contract":                        resourceAciImportedContract(),
			"aci_epg_to_contract":                          resourceAciContractProvider(),
			"aci_epg_to_contract_interface":                resourceAciContractInterfaceRelationship(),
			"aci_node_block":                               resourceAciNodeBlock(),
			"aci_epg_to_domain":                            resourceAciDomain(),
			"aci_access_generic":                           resourceAciAccessGeneric(),
			"aci_epgs_using_function":                      resourceAciEPGsUsingFunction(),
			"aci_service_redirect_policy":                  resourceAciServiceRedirectPolicy(),
			"aci_destination_of_redirected_traffic":        resourceAciDestinationofredirectedtraffic(),
			"aci_fex_bundle_group":                         resourceAciFexBundleGroup(),
			"aci_access_group":                             resourceAciAccessGroup(),
			"aci_spine_profile":                            resourceAciSpineProfile(),
			"aci_spine_switch_association":                 resourceAciSwitchSpineAssociation(),
			"aci_spine_interface_profile_selector":         resourceAciInterfaceProfile(),
			"aci_spine_port_selector":                      resourceAciInterfaceProfileDeprecated(),
			"aci_spine_port_policy_group":                  resourceAciSpineAccessPortPolicyGroup(),
			"aci_fabric_if_pol":                            resourceAciLinkLevelPolicy(),
			"aci_spanning_tree_interface_policy":           resourceAciSpanningTreeInterfacePolicy(),
			"aci_aaa_domain":                               resourceAciSecurityDomain(),
			"aci_l4_l7_service_graph_template":             resourceAciL4L7ServiceGraphTemplate(),
			"aci_logical_device_context":                   resourceAciLogicalDeviceContext(),
			"aci_function_node":                            resourceAciFunctionNode(),
			"aci_cloud_vpn_gateway":                        resourceAciCloudVpnGateway(),
			"aci_logical_interface_context":                resourceAciLogicalInterfaceContext(),
			"aci_dhcp_option_policy":                       resourceAciDHCPOptionPolicy(),
			"aci_bd_dhcp_label":                            resourceAciBDDHCPLabel(),
			"aci_dhcp_relay_policy":                        resourceAciDHCPRelayPolicy(),
			"aci_leaf_breakout_port_group":                 resourceAciLeafBreakoutPortGroup(),
			"aci_l2_domain":                                resourceAciL2Domain(),
			"aci_l2out_extepg":                             resourceAciL2outExternalEpg(),
			"aci_l2_outside":                               resourceAciL2Outside(),
			"aci_node_mgmt_epg":                            resourceAciNodeManagementEPg(),
			"aci_connection":                               resourceAciConnection(),
			"aci_l3out_bgp_external_policy":                resourceAciL3outBgpExternalPolicy(),
			"aci_l3out_ospf_external_policy":               resourceAciL3outOspfExternalPolicy(),
			"aci_l3out_path_attachment":                    resourceAciL3outPathAttachment(),
			"aci_l3out_path_attachment_secondary_ip":       resourceAciL3outPathAttachmentSecondaryIp(),
			"aci_bgp_route_summarization":                  resourceAciBgpRouteSummarization(),
			"aci_static_node_mgmt_address":                 resourceAciMgmtStaticNode(),
			"aci_l3out_ospf_interface_profile":             resourceAciOSPFInterfaceProfile(),
			"aci_l3out_loopback_interface_profile":         resourceAciLoopBackInterfaceProfile(),
			"aci_bgp_peer_prefix":                          resourceAciBGPPeerPrefixPolicy(),
			"aci_bgp_peer_connectivity_profile":            resourceAciBgpPeerConnectivityProfile(),
			"aci_bgp_best_path_policy":                     resourceAciBgpBestPathPolicy(),
			"aci_bgp_timers":                               resourceAciBGPTimersPolicy(),
			"aci_ospf_route_summarization":                 resourceAciOspfRouteSummarization(),
			"aci_bgp_address_family_context":               resourceAciBGPAddressFamilyContextPolicy(),
			"aci_hsrp_group_policy":                        resourceAciHSRPGroupPolicy(),
			"aci_l3out_hsrp_interface_profile":             resourceAciL3outHSRPInterfaceProfile(),
			"aci_ospf_timers":                              resourceAciOSPFTimersPolicy(),
			"aci_hsrp_interface_policy":                    resourceAciHSRPInterfacePolicy(),
			"aci_bgp_route_control_profile":                resourceAciBgpRouteControlProfile(),
			"aci_l3out_hsrp_interface_group":               resourceAciHSRPGroupProfile(),
			"aci_l3out_floating_svi":                       resourceAciVirtualLogicalInterfaceProfile(),
			"aci_l3out_hsrp_secondary_vip":                 resourceAciL3outHSRPSecondaryVIP(),
			"aci_l3out_bfd_interface_profile":              resourceAciBFDInterfaceProfile(),
			"aci_l3out_bgp_protocol_profile":               resourceAciL3outBGPProtocolProfile(),
			"aci_l3out_route_tag_policy":                   resourceAciL3outRouteTagPolicy(),
			"aci_l3out_static_route":                       resourceAciL3outStaticRoute(),
			"aci_l3out_static_route_next_hop":              resourceAciL3outStaticRouteNextHop(),
			"aci_l3out_vpc_member":                         resourceAciL3outVPCMember(),
			"aci_endpoint_security_group_selector":         resourceAciEndpointSecurityGroupSelector(),
			"aci_endpoint_security_group_epg_selector":     resourceAciEndpointSecurityGroupEPgSelector(),
			"aci_endpoint_security_group_tag_selector":     resourceAciEndpointSecurityGroupTagSelector(),
			"aci_bfd_interface_policy":                     resourceAciBFDInterfacePolicy(),
			"aci_l3_interface_policy":                      resourceAciL3InterfacePolicy(),
			"aci_access_switch_policy_group":               resourceAciAccessSwitchPolicyGroup(),
			"aci_managed_node_connectivity_group":          resourceAciManagedNodeConnectivityGroup(),
			"aci_vpc_domain_policy":                        resourceAciVPCDomainPolicy(),
			"aci_spine_switch_policy_group":                resourceAciSpineSwitchPolicyGroup(),
			"aci_recurring_window":                         resourceAciRecurringWindow(),
			"aci_file_remote_path":                         resourceAciRemotePathofaFile(),
			"aci_snmp_user":                                resourceAciSnmpUserProfile(),
			"aci_vrf_snmp_context_community":               resourceAciSNMPCommunityDeprecated(),
			"aci_snmp_community":                           resourceAciSNMPCommunity(),
			"aci_mgmt_zone":                                resourceAciManagedNodesZone(),
			"aci_vrf_snmp_context":                         resourceAciSNMPContextProfile(),
			"aci_endpoint_ip_aging_profile":                resourceAciIPAgingPolicy(),
			"aci_mgmt_preference":                          resourceAciMgmtconnectivitypreference(),
			"aci_endpoint_controls":                        resourceAciEndpointControlPolicy(),
			"aci_fabric_node_control":                      resourceAciFabricNodeControl(),
			"aci_coop_policy":                              resourceAciCOOPGroupPolicy(),
			"aci_endpoint_loop_protection":                 resourceAciEPLoopProtectionPolicy(),
			"aci_port_tracking":                            resourceAciPortTracking(),
			"aci_user_security_domain":                     resourceAciUserDomain(),
			"aci_encryption_key":                           resourceAciAESEncryptionPassphraseandKeysforConfigExportImport(),
			"aci_mcp_instance_policy":                      resourceAciMiscablingProtocolInstancePolicy(),
			"aci_qos_instance_policy":                      resourceAciQOSInstancePolicy(),
			"aci_user_security_domain_role":                resourceAciUserRole(),
			"aci_console_authentication":                   resourceAciConsoleAuthenticationMethod(),
			"aci_error_disable_recovery":                   resourceAciErrorDisabledRecoveryPolicy(),
			"aci_fabric_wide_settings":                     resourceAciFabricWideSettingsPolicy(),
			"aci_authentication_properties":                resourceAciAAAAuthentication(),
			"aci_duo_provider_group":                       resourceAciDuoProviderGroup(),
			"aci_ldap_provider":                            resourceAciLDAPProvider(),
			"aci_radius_provider_group":                    resourceAciRADIUSProviderGroup(),
			"aci_ldap_group_map_rule":                      resourceAciLDAPGroupMapRule(),
			"aci_tacacs_accounting_destination":            resourceAciTACACSDestination(),
			"aci_ldap_group_map_rule_to_group_map":         resourceAciLDAPGroupMapruleref(),
			"aci_tacacs_accounting":                        resourceAciTACACSMonitoringDestinationGroup(),
			"aci_rsa_provider":                             resourceAciRSAProvider(),
			"aci_saml_provider":                            resourceAciSAMLProvider(),
			"aci_login_domain":                             resourceAciLoginDomain(),
			"aci_default_authentication":                   resourceAciDefaultAuthenticationMethodforallLogins(),
			"aci_tacacs_provider_group":                    resourceAciTACACSPlusProviderGroup(),
			"aci_tacacs_provider":                          resourceAciTACACSProvider(),
			"aci_saml_provider_group":                      resourceAciSAMLProviderGroup(),
			"aci_ldap_group_map":                           resourceAciLDAPGroupMap(),
			"aci_global_security":                          resourceAciUserManagement(),
			"aci_login_domain_provider":                    resourceAciProviderGroupMember(),
			"aci_tacacs_source":                            resourceAciTACACSSource(),
			"aci_isis_domain_policy":                       resourceAciISISDomainPolicy(),
			"aci_radius_provider":                          resourceAciRADIUSProvider(),
			"aci_interface_blacklist":                      resourceAciOutofServiceFabricPath(),
			"aci_route_control_context":                    resourceAciRouteControlContext(),
			"aci_match_rule":                               resourceAciMatchRule(),
			"aci_match_community_terms":                    resourceAciMatchCommunityTerm(),
			"aci_match_regex_community_terms":              resourceAciMatchRuleBasedonCommunityRegularExpression(),
			"aci_match_route_destination_rule":             resourceAciMatchRouteDestinationRule(),
			"aci_aaa_domain_relationship":                  resourceAciDomainRelationship(),
			"aci_aaep_to_domain":                           resourceAciInfraRsDomP(),
			"aci_action_rule_additional_communities":       resourceAciRtctrlSetAddComm(),
			"aci_l4_l7_device":                             resourceAciL4ToL7Devices(),
			"aci_concrete_device":                          resourceAciConcreteDevice(),
			"aci_concrete_interface":                       resourceAciConcreteInterface(),
			"aci_l4_l7_logical_interface":                  resourceAciLogicalInterface(),
			"aci_l4_l7_redirect_health_group":              resourceAciL4L7RedirectHealthGroup(),
			"aci_ip_sla_monitoring_policy":                 resourceAciIPSLAMonitoringPolicy(),
			"aci_bulk_epg_to_static_path":                  resourceAciBulkStaticPath(),
			"aci_vrf_leak_epg_bd_subnet":                   resourceAciLeakInternalSubnet(),
			"aci_cloud_vrf_leak_routes":                    resourceAciLeakInternalPrefix(),
			"aci_service_redirect_backup_policy":           resourceAciPBRBackupPolicy(),
			"aci_pbr_l1_l2_destination":                    resourceAciL1L2RedirectDestTraffic(),
			"aci_interface_config":                         resourceAciInterfaceConfiguration(),
			"aci_cloud_template_region_detail":             resourceAciCloudTemplateRegion(),
			"aci_power_supply_redundancy_policy":           resourceAciPowerSupplyRedundancyPolicy(),
			"aci_pim_interface_policy":                     resourceAciPIMInterfacePolicy(),
			"aci_igmp_interface_policy":                    resourceAciIGMPInterfacePolicy(),
			"aci_cloud_l4_l7_native_load_balancer":         resourceAciCloudL4L7LoadBalancer(),
			"aci_cloud_l4_l7_third_party_device":           resourceAciCloudL4L7Device(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"aci_contract":                                 dataSourceAciContract(),
			"aci_contract_subject":                         dataSourceAciContractSubject(),
			"aci_contract_subject_filter":                  dataSourceAciSubjectFilter(),
			"aci_contract_subject_one_way_filter":          dataSourceAciFilterRelationship(),
			"aci_subnet":                                   dataSourceAciSubnet(),
			"aci_filter":                                   dataSourceAciFilter(),
			"aci_filter_entry":                             dataSourceAciFilterEntry(),
			"aci_vmm_domain":                               dataSourceAciVMMDomain(),
			"aci_vmm_controller":                           dataSourceAciVMMController(),
			"aci_vswitch_policy":                           dataSourceAciVSwitchPolicyGroup(),
			"aci_vrf_to_bgp_address_family_context":        dataSourceAciBGPAddressFamilyContextPolicyRelationship(),
			"aci_external_network_instance_profile":        dataSourceAciExternalNetworkInstanceProfile(),
			"aci_l3_outside":                               dataSourceAciL3Outside(),
			"aci_bfd_multihop_interface_profile":           dataSourceAciBfdMultihopInterfaceProfile(),
			"aci_bfd_multihop_interface_policy":            dataSourceAciBfdMultihopInterfacePolicy(),
			"aci_bfd_multihop_node_policy":                 dataSourceAciBFDMultihopNodePolicy(),
			"aci_interface_fc_policy":                      dataSourceAciInterfaceFCPolicy(),
			"aci_leaf_access_bundle_policy_group":          dataSourceAciPCVPCInterfacePolicyGroup(),
			"aci_leaf_access_bundle_policy_sub_group":      dataSourceAciOverridePCVPCPolicyGroup(),
			"aci_leaf_access_port_policy_group":            dataSourceAciLeafAccessPortPolicyGroup(),
			"aci_lldp_interface_policy":                    dataSourceAciLLDPInterfacePolicy(),
			"aci_miscabling_protocol_interface_policy":     dataSourceAciMiscablingProtocolInterfacePolicy(),
			"aci_ospf_interface_policy":                    dataSourceAciOSPFInterfacePolicy(),
			"aci_lacp_policy":                              dataSourceAciLACPPolicy(),
			"aci_lacp_member_policy":                       dataSourceAciLACPMemberPolicy(),
			"aci_port_security_policy":                     dataSourceAciPortSecurityPolicy(),
			"aci_leaf_profile":                             dataSourceAciLeafProfile(),
			"aci_end_point_retention_policy":               dataSourceAciEndPointRetentionPolicy(),
			"aci_vlan_encapsulationfor_vxlan_traffic":      dataSourceAciVlanEncapsulationforVxlanTraffic(),
			"aci_logical_node_profile":                     dataSourceAciLogicalNodeProfile(),
			"aci_logical_interface_profile":                dataSourceAciLogicalInterfaceProfile(),
			"aci_l3_ext_subnet":                            dataSourceAciL3ExtSubnet(),
			"aci_cloud_applicationcontainer":               dataSourceAciCloudApplicationcontainer(),
			"aci_cloud_ipsec_tunnel_subnet_pool":           dataSourceAciSubnetPoolforIpSecTunnels(),
			"aci_cloud_external_network":                   dataSourceAciCloudTemplateforExternalNetwork(),
			"aci_cloud_external_network_vpn_network":       dataSourceAciCloudTemplateforVPNNetwork(),
			"aci_cloud_aws_provider":                       dataSourceAciCloudAWSProvider(),
			"aci_autonomous_system_profile":                dataSourceAciAutonomousSystemProfile(),
			"aci_cloud_cidr_pool":                          dataSourceAciCloudCIDRPool(),
			"aci_cloud_domain_profile":                     dataSourceAciCloudDomainProfile(),
			"aci_cloud_context_profile":                    dataSourceAciCloudContextProfile(),
			"aci_cloud_epg":                                dataSourceAciCloudEPg(),
			"aci_cloud_endpoint_selectorfor_external_epgs": dataSourceAciCloudEndpointSelectorforExternalEPgs(),
			"aci_cloud_endpoint_selector":                  dataSourceAciCloudEndpointSelector(),
			"aci_cloud_external_epg":                       dataSourceAciCloudExternalEPg(),
			"aci_cloud_private_link_label":                 dataSourceAciCloudPrivateLinkLabel(),
			"aci_cloud_provider_profile":                   dataSourceAciCloudProviderProfile(),
			"aci_cloud_providers_region":                   dataSourceAciCloudProvidersRegion(),
			"aci_cloud_service_epg":                        dataSourceAciCloudServiceEPg(),
			"aci_cloud_service_endpoint_selector":          dataSourceAciCloudServiceEndpointSelector(),
			"aci_cloud_subnet":                             dataSourceAciCloudSubnet(),
			"aci_cloud_availability_zone":                  dataSourceAciCloudAvailabilityZone(),
			"aci_cloud_account":                            dataSourceAciCloudAccount(),
			"aci_tenant_to_cloud_account":                  dataSourceAciTenantToCloudAccountAssociation(),
			"aci_cloud_ad":                                 dataSourceAciCloudActiveDirectory(),
			"aci_cloud_credentials":                        dataSourceAciCloudCredentials(),
			"aci_local_user":                               dataSourceAciLocalUser(),
			"aci_pod_maintenance_group":                    dataSourceAciPODMaintenanceGroup(),
			"aci_maintenance_policy":                       dataSourceAciMaintenancePolicy(),
			"aci_monitoring_policy":                        dataSourceAciMonitoringPolicy(),
			"aci_physical_domain":                          dataSourceAciPhysicalDomain(),
			"aci_action_rule_profile":                      dataSourceAciActionRuleProfile(),
			"aci_trigger_scheduler":                        dataSourceAciTriggerScheduler(),
			"aci_leaf_selector":                            dataSourceAciSwitchAssociation(),
			"aci_span_destination_group":                   dataSourceAciSPANDestinationGroup(),
			"aci_span_source_group":                        dataSourceAciSPANSourceGroup(),
			"aci_span_sourcedestination_group_match_label": dataSourceAciSPANSourcedestinationGroupMatchLabel(),
			"aci_vlan_pool":                                dataSourceAciVLANPool(),
			"aci_vxlan_pool":                               dataSourceAciVXLANPool(),
			"aci_vsan_pool":                                dataSourceAciVSANPool(),
			"aci_multicast_pool":                           dataSourceAciMulticastAddressPool(),
			"aci_multicast_pool_block":                     dataSourceAciMulticastAddressBlock(),
			"aci_firmware_group":                           dataSourceAciFirmwareGroup(),
			"aci_firmware_policy":                          dataSourceAciFirmwarePolicy(),
			"aci_firmware_download_task":                   dataSourceAciFirmwareDownloadTask(),
			"aci_fc_domain":                                dataSourceAciFCDomain(),
			"aci_configuration_export_policy":              dataSourceAciConfigurationExportPolicy(),
			"aci_cdp_interface_policy":                     dataSourceAciCDPInterfacePolicy(),
			"aci_access_sub_port_block":                    dataSourceAciAccessSubPortBlock(),
			"aci_maintenance_group_node":                   dataSourceAciNodeBlockMG(),
			"aci_node_block_firmware":                      dataSourceAciNodeBlockFW(),
			"aci_vpc_explicit_protection_group":            dataSourceAciVPCExplicitProtectionGroup(),
			"aci_configuration_import_policy":              dataSourceAciConfigurationImportPolicy(),
			"aci_fabric_node_member":                       dataSourceAciFabricNodeMember(),
			"aci_ranges":                                   dataSourceAciRanges(),
			"aci_l3_domain_profile":                        dataSourceAciL3DomainProfile(),
			"aci_x509_certificate":                         dataSourceAciX509Certificate(),
			"aci_epg_to_static_path":                       dataSourceAciStaticPath(),
			"aci_logical_node_to_fabric_node":              dataSourceAciFabricNode(),
			"aci_imported_contract":                        dataSourceAciImportedContract(),
			"aci_epg_to_contract":                          dataSourceAciContractProvider(),
			"aci_epg_to_contract_interface":                dataSourceAciContractInterfaceRelationship(),
			"aci_node_block":                               dataSourceAciNodeBlock(),
			"aci_epg_to_domain":                            dataSourceAciDomain(),
			"aci_access_generic":                           dataSourceAciAccessGeneric(),
			"aci_epgs_using_function":                      dataSourceAciEPGsUsingFunction(),
			"aci_service_redirect_policy":                  dataSourceAciServiceRedirectPolicy(),
			"aci_destination_of_redirected_traffic":        dataSourceAciDestinationofredirectedtraffic(),
			"aci_fex_bundle_group":                         dataSourceAciFexBundleGroup(),
			"aci_access_group":                             dataSourceAciAccessGroup(),
			"aci_spine_profile":                            dataSourceAciSpineProfile(),
			"aci_spine_switch_association":                 dataSourceAciSwitchSpineAssociation(),
			"aci_spine_interface_profile_selector":         dataSourceAciInterfaceProfile(),
			"aci_spine_port_selector":                      dataSourceAciInterfaceProfileDeprecated(),
			"aci_spine_port_policy_group":                  dataSourceAciSpineAccessPortPolicyGroup(),
			"aci_fabric_path_ep":                           dataSourceAciFabricPathEndpoint(),
			"aci_fabric_if_pol":                            dataSourceAciLinkLevelPolicy(),
			"aci_spanning_tree_interface_policy":           dataSourceAciSpanningTreeInterfacePolicy(),
			"aci_aaa_domain":                               dataSourceAciSecurityDomain(),
			"aci_client_end_point":                         dataSourceAciClientEndPoint(),
			"aci_l4_l7_service_graph_template":             dataSourceAciL4L7ServiceGraphTemplate(),
			"aci_logical_device_context":                   dataSourceAciLogicalDeviceContext(),
			"aci_function_node":                            dataSourceAciFunctionNode(),
			"aci_cloud_vpn_gateway":                        dataSourceAciCloudVpnGateway(),
			"aci_logical_interface_context":                dataSourceAciLogicalInterfaceContext(),
			"aci_dhcp_option_policy":                       dataSourceAciDHCPOptionPolicy(),
			"aci_dhcp_option":                              dataSourceAciDHCPOption(),
			"aci_bd_dhcp_label":                            dataSourceAciBDDHCPLabel(),
			"aci_dhcp_relay_policy":                        dataSourceAciDHCPRelayPolicy(),
			"aci_leaf_breakout_port_group":                 dataSourceAciLeafBreakoutPortGroup(),
			"aci_l2_domain":                                dataSourceAciL2Domain(),
			"aci_l2out_extepg":                             dataSourceAciL2outExternalEpg(),
			"aci_l2_outside":                               dataSourceAciL2Outside(),
			"aci_node_mgmt_epg":                            dataSourceAciNodeManagementEPg(),
			"aci_connection":                               dataSourceAciConnection(),
			"aci_l3out_bgp_external_policy":                dataSourceAciL3outBgpExternalPolicy(),
			"aci_l3out_ospf_external_policy":               dataSourceAciL3outOspfExternalPolicy(),
			"aci_l3out_path_attachment":                    dataSourceAciL3outPathAttachment(),
			"aci_l3out_path_attachment_secondary_ip":       dataSourceAciL3outPathAttachmentSecondaryIp(),
			"aci_bgp_route_summarization":                  dataSourceAciBgpRouteSummarization(),
			"aci_static_node_mgmt_address":                 dataSourceAciMgmtStaticNode(),
			"aci_l3out_ospf_interface_profile":             dataSourceAciOSPFInterfaceProfile(),
			"aci_l3out_loopback_interface_profile":         dataSourceAciLoopBackInterfaceProfile(),
			"aci_bgp_peer_prefix":                          dataSourceAciBGPPeerPrefixPolicy(),
			"aci_rest":                                     dataSourceAciRest(),
			"aci_bgp_peer_connectivity_profile":            dataSourceAciBgpPeerConnectivityProfile(),
			"aci_bgp_best_path_policy":                     dataSourceAciBgpBestPathPolicy(),
			"aci_bgp_timers":                               dataSourceAciBGPTimersPolicy(),
			"aci_ospf_route_summarization":                 dataSourceAciOspfRouteSummarization(),
			"aci_bgp_address_family_context":               dataSourceAciBGPAddressFamilyContextPolicy(),
			"aci_hsrp_group_policy":                        dataSourceAciHSRPGroupPolicy(),
			"aci_l3out_hsrp_interface_profile":             dataSourceAciL3outHSRPInterfaceProfile(),
			"aci_ospf_timers":                              dataSourceAciOSPFTimersPolicy(),
			"aci_hsrp_interface_policy":                    dataSourceAciHSRPInterfacePolicy(),
			"aci_bgp_route_control_profile":                dataSourceAciBgpRouteControlProfile(),
			"aci_l3out_hsrp_interface_group":               dataSourceAciHSRPGroupProfile(),
			"aci_l3out_floating_svi":                       dataSourceAciVirtualLogicalInterfaceProfile(),
			"aci_l3out_hsrp_secondary_vip":                 dataSourceAciL3outHSRPSecondaryVIP(),
			"aci_l3out_bfd_interface_profile":              dataSourceAciBFDInterfaceProfile(),
			"aci_l3out_bgp_protocol_profile":               dataSourceAciL3outBGPProtocolProfile(),
			"aci_l3out_route_tag_policy":                   dataSourceAciL3outRouteTagPolicy(),
			"aci_l3out_static_route":                       dataSourceAciL3outStaticRoute(),
			"aci_l3out_static_route_next_hop":              dataSourceAciL3outStaticRouteNextHop(),
			"aci_l3out_vpc_member":                         dataSourceAciL3outVPCMember(),
			"aci_endpoint_security_group_selector":         dataSourceAciEndpointSecurityGroupSelector(),
			"aci_endpoint_security_group_epg_selector":     dataSourceAciEndpointSecurityGroupEPgSelector(),
			"aci_endpoint_security_group_tag_selector":     dataSourceAciEndpointSecurityGroupTagSelector(),
			"aci_bfd_interface_policy":                     dataSourceAciBFDInterfacePolicy(),
			"aci_l3_interface_policy":                      dataSourceAciL3InterfacePolicy(),
			"aci_fabric_node":                              dataSourceAciFabricNodeOrg(),
			"aci_access_switch_policy_group":               dataSourceAciAccessSwitchPolicyGroup(),
			"aci_managed_node_connectivity_group":          dataSourceAciManagedNodeConnectivityGroup(),
			"aci_vpc_domain_policy":                        dataSourceAciVPCDomainPolicy(),
			"aci_spine_switch_policy_group":                dataSourceAciSpineSwitchPolicyGroup(),
			"aci_recurring_window":                         dataSourceAciRecurringWindow(),
			"aci_file_remote_path":                         dataSourceAciRemotePathofaFile(),
			"aci_snmp_user":                                dataSourceAciSnmpUserProfile(),
			"aci_vrf_snmp_context_community":               dataSourceAciSNMPCommunityDeprecated(),
			"aci_snmp_community":                           dataSourceAciSNMPCommunity(),
			"aci_mgmt_zone":                                dataSourceAciManagedNodesZone(),
			"aci_vrf_snmp_context":                         dataSourceAciSNMPContextProfile(),
			"aci_endpoint_ip_aging_profile":                dataSourceAciIPAgingPolicy(),
			"aci_mgmt_preference":                          dataSourceAciMgmtconnectivitypreference(),
			"aci_endpoint_controls":                        dataSourceAciEndpointControlPolicy(),
			"aci_fabric_node_control":                      dataSourceAciFabricNodeControl(),
			"aci_coop_policy":                              dataSourceAciCOOPGroupPolicy(),
			"aci_endpoint_loop_protection":                 dataSourceAciEPLoopProtectionPolicy(),
			"aci_port_tracking":                            dataSourceAciPortTracking(),
			"aci_user_security_domain":                     dataSourceAciUserDomain(),
			"aci_encryption_key":                           dataSourceAciAESEncryptionPassphraseandKeysforConfigExportImport(),
			"aci_mcp_instance_policy":                      dataSourceAciMiscablingProtocolInstancePolicy(),
			"aci_qos_instance_policy":                      dataSourceAciQOSInstancePolicy(),
			"aci_user_security_domain_role":                dataSourceAciUserRole(),
			"aci_console_authentication":                   dataSourceAciConsoleAuthenticationMethod(),
			"aci_error_disable_recovery":                   dataSourceAciErrorDisabledRecoveryPolicy(),
			"aci_fabric_wide_settings":                     dataSourceAciFabricWideSettingsPolicy(),
			"aci_authentication_properties":                dataSourceAciAAAAuthentication(),
			"aci_duo_provider_group":                       dataSourceAciDuoProviderGroup(),
			"aci_ldap_provider":                            dataSourceAciLDAPProvider(),
			"aci_saml_certificate":                         dataSourceAciKeypairforSAMLEncryption(),
			"aci_radius_provider_group":                    dataSourceAciRADIUSProviderGroup(),
			"aci_ldap_group_map_rule":                      dataSourceAciLDAPGroupMapRule(),
			"aci_tacacs_accounting_destination":            dataSourceAciTACACSDestination(),
			"aci_ldap_group_map_rule_to_group_map":         dataSourceAciLDAPGroupMapruleref(),
			"aci_tacacs_accounting":                        dataSourceAciTACACSMonitoringDestinationGroup(),
			"aci_rsa_provider":                             dataSourceAciRSAProvider(),
			"aci_saml_provider":                            dataSourceAciSAMLProvider(),
			"aci_login_domain":                             dataSourceAciLoginDomain(),
			"aci_default_authentication":                   dataSourceAciDefaultAuthenticationMethodforallLogins(),
			"aci_tacacs_provider_group":                    dataSourceAciTACACSPlusProviderGroup(),
			"aci_tacacs_provider":                          dataSourceAciTACACSProvider(),
			"aci_saml_provider_group":                      dataSourceAciSAMLProviderGroup(),
			"aci_ldap_group_map":                           dataSourceAciLDAPGroupMap(),
			"aci_global_security":                          dataSourceAciUserManagement(),
			"aci_login_domain_provider":                    dataSourceAciProviderGroupMember(),
			"aci_tacacs_source":                            dataSourceAciTACACSSource(),
			"aci_isis_domain_policy":                       dataSourceAciISISDomainPolicy(),
			"aci_radius_provider":                          dataSourceAciRADIUSProvider(),
			"aci_interface_blacklist":                      dataSourceAciOutofServiceFabricPath(),
			"aci_route_control_context":                    dataSourceAciRouteControlContext(),
			"aci_match_rule":                               dataSourceAciMatchRule(),
			"aci_match_community_terms":                    dataSourceAciMatchCommunityTerm(),
			"aci_match_regex_community_terms":              dataSourceAciMatchRuleBasedonCommunityRegularExpression(),
			"aci_match_route_destination_rule":             dataSourceAciMatchRouteDestinationRule(),
			"aci_aaa_domain_relationship":                  dataSourceAciDomainRelationship(),
			"aci_aaep_to_domain":                           dataSourceAciInfraRsDomP(),
			"aci_action_rule_additional_communities":       dataSourceAciRtctrlSetAddComm(),
			"aci_l4_l7_device":                             dataSourceAciL4ToL7Devices(),
			"aci_concrete_device":                          dataSourceAciConcreteDevice(),
			"aci_concrete_interface":                       dataSourceAciConcreteInterface(),
			"aci_l4_l7_logical_interface":                  dataSourceAciLogicalInterface(),
			"aci_l4_l7_redirect_health_group":              dataSourceAciL4L7RedirectHealthGroup(),
			"aci_ip_sla_monitoring_policy":                 dataSourceAciIPSLAMonitoringPolicy(),
			"aci_l4_l7_deployed_graph_connector_vlan":      dataSourceAciEPgDef(),
			"aci_vrf_leak_epg_bd_subnet":                   dataSourceAciLeakInternalSubnet(),
			"aci_cloud_vrf_leak_routes":                    dataSourceAciLeakInternalPrefix(),
			"aci_service_redirect_backup_policy":           dataSourceAciPBRBackupPolicy(),
			"aci_pbr_l1_l2_destination":                    dataSourceAciL1L2RedirectDestTraffic(),
			"aci_interface_config":                         dataSourceAciInterfaceConfiguration(),
			"aci_cloud_template_region_detail":             dataSourceAciCloudTemplateRegion(),
			"aci_power_supply_redundancy_policy":           dataSourceAciPowerSupplyRedundancyPolicy(),
			"aci_pim_interface_policy":                     dataSourceAciPIMInterfacePolicy(),
			"aci_igmp_interface_policy":                    dataSourceAciIGMPInterfacePolicy(),
			"aci_cloud_l4_l7_native_load_balancer":         dataSourceAciCloudL4L7LoadBalancer(),
			"aci_cloud_l4_l7_third_party_device":           dataSourceAciCloudL4L7Device(),
		},

		ConfigureFunc: configureClient,
	}
}

func configureClient(d *schema.ResourceData) (interface{}, error) {

	config := Config{
		Username:           getStringAttribute(d, "username", "ACI_USERNAME"),
		Password:           getStringAttribute(d, "password", "ACI_PASSWORD"),
		URL:                getStringAttribute(d, "url", "ACI_URL"),
		IsInsecure:         getBoolAttribute(d, "insecure", "ACI_INSECURE", true),
		PrivateKey:         getStringAttribute(d, "private_key", "ACI_PRIVATE_KEY"),
		Certname:           getStringAttribute(d, "cert_name", "ACI_CERT_NAME"),
		ProxyUrl:           getStringAttribute(d, "proxy_url", "ACI_PROXY_URL"),
		ProxyCreds:         getStringAttribute(d, "proxy_creds", "ACI_PROXY_CREDS"),
		ValidateRelationDn: getBoolAttribute(d, "validate_relation_dn", "ACI_VAL_REL_DN", true),
		MaxRetries:         getIntAttribute(d, "retries", "ACI_RETRIES", 2),
	}

	if err := config.Valid(); err != nil {
		return nil, err
	}

	return config.getClient(), nil
}

func (c Config) Valid() error {

	if c.Username == "" {
		return fmt.Errorf("Username must be provided for the ACI provider")
	}

	if c.Password == "" {
		if c.PrivateKey == "" && c.Certname == "" {
			return fmt.Errorf("Either of private_key/cert_name or password is required")
		} else if c.PrivateKey == "" || c.Certname == "" {
			return fmt.Errorf("private_key and cert_name both must be provided")
		}
	}

	if c.URL == "" {
		return fmt.Errorf("The URL must be provided for the ACI provider")
	} else if !strings.HasPrefix(c.URL, "http://") && !strings.HasPrefix(c.URL, "https://") {
		return fmt.Errorf(fmt.Sprintf("The URL '%s' must start with 'http://' or 'https://'", c.URL))
	}

	if c.MaxRetries < 0 || c.MaxRetries > 9 {
		return fmt.Errorf("retries must be between 0 and 9 inclusive, got: %d", c.MaxRetries)
	}

	return nil
}

func (c Config) getClient() interface{} {
	if c.Password != "" {

		return client.GetClient(c.URL, c.Username, client.Password(c.Password), client.Insecure(c.IsInsecure), client.ProxyUrl(c.ProxyUrl), client.ProxyCreds(c.ProxyCreds), client.ValidateRelationDn(c.ValidateRelationDn), client.MaxRetries(c.MaxRetries))

	} else {

		return client.GetClient(c.URL, c.Username, client.PrivateKey(c.PrivateKey), client.AdminCert(c.Certname), client.Insecure(c.IsInsecure), client.ProxyUrl(c.ProxyUrl), client.ProxyCreds(c.ProxyCreds), client.ValidateRelationDn(c.ValidateRelationDn), client.MaxRetries(c.MaxRetries))
	}
}

// Config
type Config struct {
	Username           string
	Password           string
	URL                string
	IsInsecure         bool
	PrivateKey         string
	Certname           string
	ProxyUrl           string
	ProxyCreds         string
	ValidateRelationDn bool
	MaxRetries         int
}

func getStringAttribute(d *schema.ResourceData, attributeName, envKey string) string {

	if envValue := os.Getenv(envKey); envValue != "" {
		return envValue
	}
	return d.Get(attributeName).(string)
}

func getBoolAttribute(d *schema.ResourceData, attributeName, envKey string, defaultValue bool) bool {

	var err error
	var boolValue = defaultValue

	if tmpBoolValue, ok := d.GetOk(attributeName); ok {
		boolValue, err = strconv.ParseBool(tmpBoolValue.(string))
		if err != nil {
			log.Fatal(err)
		}
	} else if envValue := os.Getenv(envKey); envValue != "" {
		boolValue, err = strconv.ParseBool(envValue)
		if err != nil {
			log.Fatal(err)
		}
	}
	return boolValue
}

func getIntAttribute(d *schema.ResourceData, attributeName, envKey string, defaultValue int) int {

	var err error
	var intValue = defaultValue

	if tmpIntValue, ok := d.GetOk(attributeName); ok {
		intValue, err = strconv.Atoi(tmpIntValue.(string))
		if err != nil {
			log.Fatal(err)
		}
	} else if envValue := os.Getenv(envKey); envValue != "" {
		intValue, err = strconv.Atoi(envValue)
		if err != nil {
			log.Fatal(err)
		}
	}
	return intValue
}
