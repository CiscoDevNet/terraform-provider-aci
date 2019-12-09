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
			"aci_tenant":                                    resourceAciTenant(),
			"aci_application_profile":                       resourceAciApplicationProfile(),
			"aci_bridge_domain":                             resourceAciBridgeDomain(),
			"aci_contract":                                  resourceAciContract(),
			"aci_application_epg":                           resourceAciApplicationEPG(),
			"aci_contract_subject":                          resourceAciContractSubject(),
			"aci_subnet":                                    resourceAciSubnet(),
			"aci_filter":                                    resourceAciFilter(),
			"aci_filter_entry":                              resourceAciFilterEntry(),
			"aci_vmm_domain":                                resourceAciVMMDomain(),
			"aci_vrf":                                       resourceAciVRF(),
			"aci_rest":                                      resourceAciRest(),
			"aci_external_network_instance_profile":         resourceAciExternalNetworkInstanceProfile(),
			"aci_l3_outside":                                resourceAciL3Outside(),
			"aci_leaf_interface_profile":                    resourceAciLeafInterfaceProfile(),
			"aci_interface_fc_policy":                       resourceAciInterfaceFCPolicy(),
			"aci_l2_interface_policy":                       resourceAciL2InterfacePolicy(),
			"aci_pcvpc_interface_policy_group":              resourceAciPCVPCInterfacePolicyGroup(),
			"aci_leaf_access_port_policy_group":             resourceAciLeafAccessPortPolicyGroup(),
			"aci_lldp_interface_policy":                     resourceAciLLDPInterfacePolicy(),
			"aci_miscabling_protocol_interface_policy":      resourceAciMiscablingProtocolInterfacePolicy(),
			"aci_ospf_interface_policy":                     resourceAciOSPFInterfacePolicy(),
			"aci_access_port_selector":                      resourceAciAccessPortSelector(),
			"aci_access_port_block":                         resourceAciAccessPortBlock(),
			"aci_lacp_policy":                               resourceAciLACPPolicy(),
			"aci_port_security_policy":                      resourceAciPortSecurityPolicy(),
			"aci_leaf_profile":                              resourceAciLeafProfile(),
			"aci_end_point_retention_policy":                resourceAciEndPointRetentionPolicy(),
			"aci_attachable_access_entity_profile":          resourceAciAttachableAccessEntityProfile(),
			"aci_vlan_encapsulationfor_vxlan_traffic":       resourceAciVlanEncapsulationforVxlanTraffic(),
			"aci_logical_node_profile":                      resourceAciLogicalNodeProfile(),
			"aci_logical_interface_profile":                 resourceAciLogicalInterfaceProfile(),
			"aci_l3_ext_subnet":                             resourceAciL3ExtSubnet(),
			"aci_cloud_applicationcontainer":                resourceAciCloudApplicationcontainer(),
			"aci_cloud_aws_provider":                        resourceAciCloudAWSProvider(),
			"aci_autonomous_system_profile":                 resourceAciAutonomousSystemProfile(),
			"aci_cloud_cidr_pool":                           resourceAciCloudCIDRPool(),
			"aci_cloud_domain_profile":                      resourceAciCloudDomainProfile(),
			"aci_cloud_context_profile":                     resourceAciCloudContextProfile(),
			"aci_cloud_e_pg":                                resourceAciCloudEPg(),
			"aci_cloud_endpoint_selectorfor_external_e_pgs": resourceAciCloudEndpointSelectorforExternalEPgs(),
			"aci_cloud_endpoint_selector":                   resourceAciCloudEndpointSelector(),
			"aci_cloud_external_e_pg":                       resourceAciCloudExternalEPg(),
			"aci_cloud_provider_profile":                    resourceAciCloudProviderProfile(),
			"aci_cloud_providers_region":                    resourceAciCloudProvidersRegion(),
			"aci_cloud_subnet":                              resourceAciCloudSubnet(),
			"aci_cloud_availability_zone":                   resourceAciCloudAvailabilityZone(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"aci_tenant":                                    dataSourceAciTenant(),
			"aci_application_profile":                       dataSourceAciApplicationProfile(),
			"aci_bridge_domain":                             dataSourceAciBridgeDomain(),
			"aci_contract":                                  dataSourceAciContract(),
			"aci_application_epg":                           dataSourceAciApplicationEPG(),
			"aci_contract_subject":                          dataSourceAciContractSubject(),
			"aci_subnet":                                    dataSourceAciSubnet(),
			"aci_filter":                                    dataSourceAciFilter(),
			"aci_filter_entry":                              dataSourceAciFilterEntry(),
			"aci_vmm_domain":                                dataSourceAciVMMDomain(),
			"aci_vrf":                                       dataSourceAciVRF(),
			"aci_external_network_instance_profile":         dataSourceAciExternalNetworkInstanceProfile(),
			"aci_l3_outside":                                dataSourceAciL3Outside(),
			"aci_leaf_interface_profile":                    dataSourceAciLeafInterfaceProfile(),
			"aci_interface_fc_policy":                       dataSourceAciInterfaceFCPolicy(),
			"aci_l2_interface_policy":                       dataSourceAciL2InterfacePolicy(),
			"aci_pcvpc_interface_policy_group":              dataSourceAciPCVPCInterfacePolicyGroup(),
			"aci_leaf_access_port_policy_group":             dataSourceAciLeafAccessPortPolicyGroup(),
			"aci_lldp_interface_policy":                     dataSourceAciLLDPInterfacePolicy(),
			"aci_miscabling_protocol_interface_policy":      dataSourceAciMiscablingProtocolInterfacePolicy(),
			"aci_ospf_interface_policy":                     dataSourceAciOSPFInterfacePolicy(),
			"aci_access_port_selector":                      dataSourceAciAccessPortSelector(),
			"aci_access_port_block":                         dataSourceAciAccessPortBlock(),
			"aci_lacp_policy":                               dataSourceAciLACPPolicy(),
			"aci_port_security_policy":                      dataSourceAciPortSecurityPolicy(),
			"aci_leaf_profile":                              dataSourceAciLeafProfile(),
			"aci_end_point_retention_policy":                dataSourceAciEndPointRetentionPolicy(),
			"aci_attachable_access_entity_profile":          dataSourceAciAttachableAccessEntityProfile(),
			"aci_vlan_encapsulationfor_vxlan_traffic":       dataSourceAciVlanEncapsulationforVxlanTraffic(),
			"aci_logical_node_profile":                      dataSourceAciLogicalNodeProfile(),
			"aci_logical_interface_profile":                 dataSourceAciLogicalInterfaceProfile(),
			"aci_l3_ext_subnet":                             dataSourceAciL3ExtSubnet(),
			"aci_cloud_applicationcontainer":                dataSourceAciCloudApplicationcontainer(),
			"aci_cloud_aws_provider":                        dataSourceAciCloudAWSProvider(),
			"aci_autonomous_system_profile":                 dataSourceAciAutonomousSystemProfile(),
			"aci_cloud_cidr_pool":                           dataSourceAciCloudCIDRPool(),
			"aci_cloud_domain_profile":                      dataSourceAciCloudDomainProfile(),
			"aci_cloud_context_profile":                     dataSourceAciCloudContextProfile(),
			"aci_cloud_e_pg":                                dataSourceAciCloudEPg(),
			"aci_cloud_endpoint_selectorfor_external_e_pgs": dataSourceAciCloudEndpointSelectorforExternalEPgs(),
			"aci_cloud_endpoint_selector":                   dataSourceAciCloudEndpointSelector(),
			"aci_cloud_external_e_pg":                       dataSourceAciCloudExternalEPg(),
			"aci_cloud_provider_profile":                    dataSourceAciCloudProviderProfile(),
			"aci_cloud_providers_region":                    dataSourceAciCloudProvidersRegion(),
			"aci_cloud_subnet":                              dataSourceAciCloudSubnet(),
			"aci_cloud_availability_zone":                   dataSourceAciCloudAvailabilityZone(),
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
