package aci

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ACI_USERNAME", nil),
				Description: "Username for the APIC Account",
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ACI_PASSWORD", nil),
				Description: "Password for the APIC Account",
			},
			"url": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ACI_URL", nil),
				Description: "URL of the Cisco ACI web interface",
			},
			"insecure": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Allow insecure HTTPS client",
			},
			"private_key": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ACI_PRIVATE_KEY", nil),
				Description: "Private key path for signature calculation",
			},
			"cert_name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ACI_CERT_NAME", nil),
				Description: "Certificate name for the User in Cisco ACI.",
			},
			"proxy_url": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ACI_PROXY_URL", nil),
				Description: "Proxy Server URL with port number",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"aci_tenant":                                   resourceAciTenant(),
			"aci_application_profile":                      resourceAciApplicationProfile(),
			"aci_bridge_domain":                            resourceAciBridgeDomain(),
			"aci_contract":                                 resourceAciContract(),
			"aci_application_epg":                          resourceAciApplicationEPG(),
			"aci_contract_subject":                         resourceAciContractSubject(),
			"aci_subnet":                                   resourceAciSubnet(),
			"aci_filter":                                   resourceAciFilter(),
			"aci_filter_entry":                             resourceAciFilterEntry(),
			"aci_vmm_domain":                               resourceAciVMMDomain(),
			"aci_vrf":                                      resourceAciVRF(),
			"aci_rest":                                     resourceAciRest(),
			"aci_external_network_instance_profile":        resourceAciExternalNetworkInstanceProfile(),
			"aci_l3_outside":                               resourceAciL3Outside(),
			"aci_leaf_interface_profile":                   resourceAciLeafInterfaceProfile(),
			"aci_interface_fc_policy":                      resourceAciInterfaceFCPolicy(),
			"aci_l2_interface_policy":                      resourceAciL2InterfacePolicy(),
			"aci_leaf_access_bundle_policy_group":          resourceAciPCVPCInterfacePolicyGroup(),
			"aci_leaf_access_port_policy_group":            resourceAciLeafAccessPortPolicyGroup(),
			"aci_lldp_interface_policy":                    resourceAciLLDPInterfacePolicy(),
			"aci_miscabling_protocol_interface_policy":     resourceAciMiscablingProtocolInterfacePolicy(),
			"aci_ospf_interface_policy":                    resourceAciOSPFInterfacePolicy(),
			"aci_access_port_selector":                     resourceAciAccessPortSelector(),
			"aci_access_port_block":                        resourceAciAccessPortBlock(),
			"aci_lacp_policy":                              resourceAciLACPPolicy(),
			"aci_port_security_policy":                     resourceAciPortSecurityPolicy(),
			"aci_leaf_profile":                             resourceAciLeafProfile(),
			"aci_end_point_retention_policy":               resourceAciEndPointRetentionPolicy(),
			"aci_attachable_access_entity_profile":         resourceAciAttachableAccessEntityProfile(),
			"aci_vlan_encapsulationfor_vxlan_traffic":      resourceAciVlanEncapsulationforVxlanTraffic(),
			"aci_logical_node_profile":                     resourceAciLogicalNodeProfile(),
			"aci_logical_interface_profile":                resourceAciLogicalInterfaceProfile(),
			"aci_l3_ext_subnet":                            resourceAciL3ExtSubnet(),
			"aci_cloud_applicationcontainer":               resourceAciCloudApplicationcontainer(),
			"aci_cloud_aws_provider":                       resourceAciCloudAWSProvider(),
			"aci_autonomous_system_profile":                resourceAciAutonomousSystemProfile(),
			"aci_cloud_cidr_pool":                          resourceAciCloudCIDRPool(),
			"aci_cloud_domain_profile":                     resourceAciCloudDomainProfile(),
			"aci_cloud_context_profile":                    resourceAciCloudContextProfile(),
			"aci_cloud_epg":                                resourceAciCloudEPg(),
			"aci_cloud_endpoint_selectorfor_external_epgs": resourceAciCloudEndpointSelectorforExternalEPgs(),
			"aci_cloud_endpoint_selector":                  resourceAciCloudEndpointSelector(),
			"aci_cloud_external_epg":                       resourceAciCloudExternalEPg(),
			"aci_cloud_provider_profile":                   resourceAciCloudProviderProfile(),
			"aci_cloud_providers_region":                   resourceAciCloudProvidersRegion(),
			"aci_cloud_subnet":                             resourceAciCloudSubnet(),
			"aci_cloud_availability_zone":                  resourceAciCloudAvailabilityZone(),
			"aci_local_user":                               resourceAciLocalUser(),
			"aci_pod_maintenance_group":                    resourceAciPODMaintenanceGroup(),
			"aci_maintenance_policy":                       resourceAciMaintenancePolicy(),
			"aci_monitoring_policy":                        resourceAciMonitoringPolicy(),
			"aci_physical_domain":                          resourceAciPhysicalDomain(),
			"aci_action_rule_profile":                      resourceAciActionRuleProfile(),
			"aci_trigger_scheduler":                        resourceAciTriggerScheduler(),
			"aci_taboo_contract":                           resourceAciTabooContract(),
			"aci_leaf_selector":                            resourceAciSwitchAssociation(),
			"aci_span_destination_group":                   resourceAciSPANDestinationGroup(),
			"aci_span_source_group":                        resourceAciSPANSourceGroup(),
			"aci_span_sourcedestination_group_match_label": resourceAciSPANSourcedestinationGroupMatchLabel(),
			"aci_vlan_pool":                                resourceAciVLANPool(),
			"aci_vxlan_pool":                               resourceAciVXLANPool(),
			"aci_vsan_pool":                                resourceAciVSANPool(),
			"aci_firmware_group":                           resourceAciFirmwareGroup(),
			"aci_firmware_policy":                          resourceAciFirmwarePolicy(),
			"aci_firmware_download_task":                   resourceAciFirmwareDownloadTask(),
			"aci_fc_domain":                                resourceAciFCDomain(),
			"aci_configuration_export_policy":              resourceAciConfigurationExportPolicy(),
			"aci_cdp_interface_policy":                     resourceAciCDPInterfacePolicy(),
			"aci_access_sub_port_block":                    resourceAciAccessSubPortBlock(),
			"aci_node_block_maintgrp":                      resourceAciNodeBlockMG(),
			"aci_node_block_firmware":                      resourceAciNodeBlockFW(),
			"aci_vpc_explicit_protection_group":            resourceAciVPCExplicitProtectionGroup(),
			"aci_configuration_import_policy":              resourceAciConfigurationImportPolicy(),
			"aci_fabric_node_member":                       resourceAciFabricNodeMember(),
			"aci_ranges":                                   resourceAciRanges(),
			"aci_l3_domain_profile":                        resourceAciL3DomainProfile(),
			"aci_x509_certificate":                         resourceAciX509Certificate(),
			"aci_epg_to_static_path":                       resourceAciStaticPath(),
			"aci_logical_node_to_fabric_node":              resourceAciFabricNode(),
			"aci_any":                                      resourceAciAny(),
			"aci_imported_contract":                        resourceAciImportedContract(),
			"aci_epg_to_contract":                          resourceAciContractProvider(),
			"aci_node_block":                               resourceAciNodeBlock(),
			"aci_epg_to_domain":                            resourceAciDomain(),
			"aci_access_generic":                           resourceAciAccessGeneric(),
			"aci_epgs_using_function":                      resourceAciEPGsUsingFunction(),
			"aci_service_redirect_policy":                  resourceAciServiceRedirectPolicy(),
			"aci_destination_of_redirected_traffic":        resourceAciDestinationofredirectedtraffic(),
			"aci_fex_profile":                              resourceAciFEXProfile(),
			"aci_fex_bundle_group":                         resourceAciFexBundleGroup(),
			"aci_access_group":                             resourceAciAccessGroup(),
			"aci_spine_profile":                            resourceAciSpineProfile(),
			"aci_spine_switch_association":                 resourceAciSwitchSpineAssociation(),
			"aci_spine_port_selector":                      resourceAciInterfaceProfile(),
			"aci_spine_interface_profile":                  resourceAciSpineInterfaceProfile(),
			"aci_spine_port_policy_group":                  resourceAciSpineAccessPortPolicyGroup(),
			"aci_fabric_if_pol":                            resourceAciLinkLevelPolicy(),
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
			"aci_l3out_hsrp_secondary_vip":                 resourceAciL3outHSRPSecondaryVIP(),
			"aci_l3out_bfd_interface_profile":              resourceAciBFDInterfaceProfile(),
			"aci_l3out_bgp_protocol_profile":               resourceAciL3outBGPProtocolProfile(),
			"aci_l3out_route_tag_policy":                   resourceAciL3outRouteTagPolicy(),
			"aci_l3out_static_route":                       resourceAciL3outStaticRoute(),
			"aci_l3out_static_route_next_hop":              resourceAciL3outStaticRouteNextHop(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"aci_tenant":                                   dataSourceAciTenant(),
			"aci_application_profile":                      dataSourceAciApplicationProfile(),
			"aci_bridge_domain":                            dataSourceAciBridgeDomain(),
			"aci_contract":                                 dataSourceAciContract(),
			"aci_application_epg":                          dataSourceAciApplicationEPG(),
			"aci_contract_subject":                         dataSourceAciContractSubject(),
			"aci_subnet":                                   dataSourceAciSubnet(),
			"aci_filter":                                   dataSourceAciFilter(),
			"aci_filter_entry":                             dataSourceAciFilterEntry(),
			"aci_vmm_domain":                               dataSourceAciVMMDomain(),
			"aci_vrf":                                      dataSourceAciVRF(),
			"aci_external_network_instance_profile":        dataSourceAciExternalNetworkInstanceProfile(),
			"aci_l3_outside":                               dataSourceAciL3Outside(),
			"aci_leaf_interface_profile":                   dataSourceAciLeafInterfaceProfile(),
			"aci_interface_fc_policy":                      dataSourceAciInterfaceFCPolicy(),
			"aci_l2_interface_policy":                      dataSourceAciL2InterfacePolicy(),
			"aci_leaf_access_bundle_policy_group":          dataSourceAciPCVPCInterfacePolicyGroup(),
			"aci_leaf_access_port_policy_group":            dataSourceAciLeafAccessPortPolicyGroup(),
			"aci_lldp_interface_policy":                    dataSourceAciLLDPInterfacePolicy(),
			"aci_miscabling_protocol_interface_policy":     dataSourceAciMiscablingProtocolInterfacePolicy(),
			"aci_ospf_interface_policy":                    dataSourceAciOSPFInterfacePolicy(),
			"aci_access_port_selector":                     dataSourceAciAccessPortSelector(),
			"aci_access_port_block":                        dataSourceAciAccessPortBlock(),
			"aci_lacp_policy":                              dataSourceAciLACPPolicy(),
			"aci_port_security_policy":                     dataSourceAciPortSecurityPolicy(),
			"aci_leaf_profile":                             dataSourceAciLeafProfile(),
			"aci_end_point_retention_policy":               dataSourceAciEndPointRetentionPolicy(),
			"aci_attachable_access_entity_profile":         dataSourceAciAttachableAccessEntityProfile(),
			"aci_vlan_encapsulationfor_vxlan_traffic":      dataSourceAciVlanEncapsulationforVxlanTraffic(),
			"aci_logical_node_profile":                     dataSourceAciLogicalNodeProfile(),
			"aci_logical_interface_profile":                dataSourceAciLogicalInterfaceProfile(),
			"aci_l3_ext_subnet":                            dataSourceAciL3ExtSubnet(),
			"aci_cloud_applicationcontainer":               dataSourceAciCloudApplicationcontainer(),
			"aci_cloud_aws_provider":                       dataSourceAciCloudAWSProvider(),
			"aci_autonomous_system_profile":                dataSourceAciAutonomousSystemProfile(),
			"aci_cloud_cidr_pool":                          dataSourceAciCloudCIDRPool(),
			"aci_cloud_domain_profile":                     dataSourceAciCloudDomainProfile(),
			"aci_cloud_context_profile":                    dataSourceAciCloudContextProfile(),
			"aci_cloud_epg":                                dataSourceAciCloudEPg(),
			"aci_cloud_endpoint_selectorfor_external_epgs": dataSourceAciCloudEndpointSelectorforExternalEPgs(),
			"aci_cloud_endpoint_selector":                  dataSourceAciCloudEndpointSelector(),
			"aci_cloud_external_epg":                       dataSourceAciCloudExternalEPg(),
			"aci_cloud_provider_profile":                   dataSourceAciCloudProviderProfile(),
			"aci_cloud_providers_region":                   dataSourceAciCloudProvidersRegion(),
			"aci_cloud_subnet":                             dataSourceAciCloudSubnet(),
			"aci_cloud_availability_zone":                  dataSourceAciCloudAvailabilityZone(),
			"aci_local_user":                               dataSourceAciLocalUser(),
			"aci_pod_maintenance_group":                    dataSourceAciPODMaintenanceGroup(),
			"aci_maintenance_policy":                       dataSourceAciMaintenancePolicy(),
			"aci_monitoring_policy":                        dataSourceAciMonitoringPolicy(),
			"aci_physical_domain":                          dataSourceAciPhysicalDomain(),
			"aci_action_rule_profile":                      dataSourceAciActionRuleProfile(),
			"aci_trigger_scheduler":                        dataSourceAciTriggerScheduler(),
			"aci_taboo_contract":                           dataSourceAciTabooContract(),
			"aci_leaf_selector":                            dataSourceAciSwitchAssociation(),
			"aci_span_destination_group":                   dataSourceAciSPANDestinationGroup(),
			"aci_span_source_group":                        dataSourceAciSPANSourceGroup(),
			"aci_span_sourcedestination_group_match_label": dataSourceAciSPANSourcedestinationGroupMatchLabel(),
			"aci_vlan_pool":                                dataSourceAciVLANPool(),
			"aci_vxlan_pool":                               dataSourceAciVXLANPool(),
			"aci_vsan_pool":                                dataSourceAciVSANPool(),
			"aci_firmware_group":                           dataSourceAciFirmwareGroup(),
			"aci_firmware_policy":                          dataSourceAciFirmwarePolicy(),
			"aci_firmware_download_task":                   dataSourceAciFirmwareDownloadTask(),
			"aci_fc_domain":                                dataSourceAciFCDomain(),
			"aci_configuration_export_policy":              dataSourceAciConfigurationExportPolicy(),
			"aci_cdp_interface_policy":                     dataSourceAciCDPInterfacePolicy(),
			"aci_access_sub_port_block":                    dataSourceAciAccessSubPortBlock(),
			"aci_node_block_maintgrp":                      dataSourceAciNodeBlockMG(),
			"aci_node_block_firmware":                      dataSourceAciNodeBlockFW(),
			"aci_vpc_explicit_protection_group":            dataSourceAciVPCExplicitProtectionGroup(),
			"aci_configuration_import_policy":              dataSourceAciConfigurationImportPolicy(),
			"aci_fabric_node_member":                       dataSourceAciFabricNodeMember(),
			"aci_ranges":                                   dataSourceAciRanges(),
			"aci_l3_domain_profile":                        dataSourceAciL3DomainProfile(),
			"aci_x509_certificate":                         dataSourceAciX509Certificate(),
			"aci_epg_to_static_path":                       dataSourceAciStaticPath(),
			"aci_logical_node_to_fabric_node":              dataSourceAciFabricNode(),
			"aci_any":                                      dataSourceAciAny(),
			"aci_imported_contract":                        dataSourceAciImportedContract(),
			"aci_epg_to_contract":                          dataSourceAciContractProvider(),
			"aci_node_block":                               dataSourceAciNodeBlock(),
			"aci_epg_to_domain":                            dataSourceAciDomain(),
			"aci_access_generic":                           dataSourceAciAccessGeneric(),
			"aci_epgs_using_function":                      dataSourceAciEPGsUsingFunction(),
			"aci_service_redirect_policy":                  dataSourceAciServiceRedirectPolicy(),
			"aci_destination_of_redirected_traffic":        dataSourceAciDestinationofredirectedtraffic(),
			"aci_fex_profile":                              dataSourceAciFEXProfile(),
			"aci_fex_bundle_group":                         dataSourceAciFexBundleGroup(),
			"aci_access_group":                             dataSourceAciAccessGroup(),
			"aci_spine_profile":                            dataSourceAciSpineProfile(),
			"aci_spine_switch_association":                 dataSourceAciSwitchSpineAssociation(),
			"aci_spine_port_selector":                      dataSourceAciInterfaceProfile(),
			"aci_spine_interface_profile":                  dataSourceAciSpineInterfaceProfile(),
			"aci_spine_port_policy_group":                  dataSourceAciSpineAccessPortPolicyGroup(),
			"aci_fabric_path_ep":                           dataSourceAciFabricPathEndpoint(),
			"aci_fabric_if_pol":                            dataSourceAciLinkLevelPolicy(),
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
			"aci_l3out_hsrp_secondary_vip":                 dataSourceAciL3outHSRPSecondaryVIP(),
			"aci_l3out_bfd_interface_profile":              dataSourceAciBFDInterfaceProfile(),
			"aci_l3out_bgp_protocol_profile":               dataSourceAciL3outBGPProtocolProfile(),
			"aci_l3out_route_tag_policy":                   dataSourceAciL3outRouteTagPolicy(),
			"aci_l3out_static_route":                       dataSourceAciL3outStaticRoute(),
			"aci_l3out_static_route_next_hop":              dataSourceAciL3outStaticRouteNextHop(),
		},

		ConfigureFunc: configureClient,
	}
}

func configureClient(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		Username:   d.Get("username").(string),
		Password:   d.Get("password").(string),
		URL:        d.Get("url").(string),
		IsInsecure: d.Get("insecure").(bool),
		PrivateKey: d.Get("private_key").(string),
		Certname:   d.Get("cert_name").(string),
		ProxyUrl:   d.Get("proxy_url").(string),
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
	}

	return nil
}

func (c Config) getClient() interface{} {
	if c.Password != "" {

		return client.GetClient(c.URL, c.Username, client.Password(c.Password), client.Insecure(c.IsInsecure), client.ProxyUrl(c.ProxyUrl))

	} else {

		return client.GetClient(c.URL, c.Username, client.PrivateKey(c.PrivateKey), client.AdminCert(c.Certname), client.Insecure(c.IsInsecure), client.ProxyUrl(c.ProxyUrl))
	}
}

// Config
type Config struct {
	Username   string
	Password   string
	URL        string
	IsInsecure bool
	PrivateKey string
	Certname   string
	ProxyUrl   string
}
