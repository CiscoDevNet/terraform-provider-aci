package aci

import (
	"context"
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAciCloudTemplateforVPNNetwork() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciCloudTemplateforVPNNetworkCreate,
		UpdateContext: resourceAciCloudTemplateforVPNNetworkUpdate,
		ReadContext:   resourceAciCloudTemplateforVPNNetworkRead,
		DeleteContext: resourceAciCloudTemplateforVPNNetworkDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciCloudTemplateforVPNNetworkImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"aci_cloud_external_network_dn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"remote_site_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"remote_site_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cloud_ipsec_tunnel": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ike_version": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{
								"ikev1",
								"ikev2",
							}, false),
						},
						"public_ip_address": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"subnet_pool_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"pre_shared_key": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		})),
	}
}

func getRemoteTemplateforVPNNetwork(client *client.Client, dn string) (*models.TemplateforVPNNetwork, error) {
	cloudtemplateVpnNetworkCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	cloudtemplateVpnNetwork := models.TemplateforVPNNetworkFromContainer(cloudtemplateVpnNetworkCont)
	if cloudtemplateVpnNetwork.DistinguishedName == "" {
		return nil, fmt.Errorf("TemplateforVPNNetwork %s not found", cloudtemplateVpnNetwork.DistinguishedName)
	}
	return cloudtemplateVpnNetwork, nil
}

func setTemplateforVPNNetworkAttributes(cloudtemplateVpnNetwork *models.TemplateforVPNNetwork, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()
	d.SetId(cloudtemplateVpnNetwork.DistinguishedName)
	d.Set("description", cloudtemplateVpnNetwork.Description)
	if dn != cloudtemplateVpnNetwork.DistinguishedName {
		d.Set("aci_cloud_external_network_dn", "")
	}
	cloudtemplateVpnNetworkMap, err := cloudtemplateVpnNetwork.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("annotation", cloudtemplateVpnNetworkMap["annotation"])
	d.Set("name", cloudtemplateVpnNetworkMap["name"])
	d.Set("remote_site_id", cloudtemplateVpnNetworkMap["remoteSiteId"])
	d.Set("remote_site_name", cloudtemplateVpnNetworkMap["remoteSiteName"])
	d.Set("name_alias", cloudtemplateVpnNetworkMap["nameAlias"])
	return d, nil
}

func getRemoteCloudTemplateforIpSectunnel(client *client.Client, dn string) (*models.CloudTemplateforIpSectunnel, error) {
	cloudtemplateIpSecTunnelCont, err := client.Get(dn + "/" + models.RncloudtemplateExtNetwork)
	if err != nil {
		return nil, err
	}
	cloudtemplateIpSecTunnel := models.CloudTemplateforIpSectunnelFromContainer(cloudtemplateIpSecTunnelCont)
	if cloudtemplateIpSecTunnel.DistinguishedName == "" {
		return nil, fmt.Errorf("CloudTemplateforIpSectunnel %s not found", cloudtemplateIpSecTunnel.DistinguishedName)
	}
	return cloudtemplateIpSecTunnel, nil
}

func setCloudTemplateforIpSectunnelAttributes(cloudtemplateIpSecTunnel *models.CloudTemplateforIpSectunnel, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()
	d.SetId(cloudtemplateIpSecTunnel.DistinguishedName)
	d.Set("description", cloudtemplateIpSecTunnel.Description)
	if dn != cloudtemplateIpSecTunnel.DistinguishedName {
		d.Set("templatefor_vpn_network_dn", "")
	}
	cloudtemplateIpSecTunnelMap, err := cloudtemplateIpSecTunnel.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("annotation", cloudtemplateIpSecTunnelMap["annotation"])
	d.Set("ike_version", cloudtemplateIpSecTunnelMap["ikeVersion"])
	d.Set("public_ip_address", cloudtemplateIpSecTunnelMap["peeraddr"])
	d.Set("subnet_pool_name", cloudtemplateIpSecTunnelMap["poolname"])
	d.Set("pre_shared_key", cloudtemplateIpSecTunnelMap["preSharedKey"])
	d.Set("name_alias", cloudtemplateIpSecTunnelMap["nameAlias"])
	return d, nil
}

func resourceAciCloudTemplateforVPNNetworkImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	cloudtemplateVpnNetwork, err := getRemoteTemplateforVPNNetwork(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setTemplateforVPNNetworkAttributes(cloudtemplateVpnNetwork, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciCloudTemplateforVPNNetworkCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] TemplateforVPNNetwork: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	TemplateforExternalNetworkDn := d.Get("aci_cloud_external_network_dn").(string)

	cloudtemplateVpnNetworkAttr := models.TemplateforVPNNetworkAttributes{}

	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		cloudtemplateVpnNetworkAttr.Annotation = Annotation.(string)
	} else {
		cloudtemplateVpnNetworkAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		cloudtemplateVpnNetworkAttr.Name = Name.(string)
	}

	if RemoteSiteId, ok := d.GetOk("remote_site_id"); ok {
		cloudtemplateVpnNetworkAttr.RemoteSiteId = RemoteSiteId.(string)
	}

	if RemoteSiteName, ok := d.GetOk("remote_site_name"); ok {
		cloudtemplateVpnNetworkAttr.RemoteSiteName = RemoteSiteName.(string)
	}
	cloudtemplateVpnNetwork := models.NewTemplateforVPNNetwork(fmt.Sprintf(models.RncloudtemplateVpnNetwork, name), TemplateforExternalNetworkDn, desc, nameAlias, cloudtemplateVpnNetworkAttr)

	err := aciClient.Save(cloudtemplateVpnNetwork)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("LOGs CREATE cloudtemplateVpnNetwork : %v ", cloudtemplateVpnNetwork)

	cloudIpSecTunnelPeers := make([]string, 0, 1)
	log.Printf("LOGs CREATE IF cloud_ipsec_tunnel : %v ", d.Get("cloud_ipsec_tunnel"))
	if ipSecTunnelPeers, ok := d.GetOk("cloud_ipsec_tunnel"); ok {
		clopudIpSecTunnels := ipSecTunnelPeers.([]interface{})
		cloudVPNNetworkDn := cloudtemplateVpnNetwork.DistinguishedName
		log.Printf("LOGs CREATE IF cloudVPNNetworkDn : %v ", cloudVPNNetworkDn)

		for _, val := range clopudIpSecTunnels {
			log.Printf("LOGs INSIDE FOR")
			ipSecTunnels := val.(map[string]interface{})
			log.Printf("LOGs INSIDE FOR222222 ipSecTunnels  :: %v ", ipSecTunnels)
			log.Printf("LOGs INSIDE FOR3333 ipSecTunnels  :: %v ", ipSecTunnels["subnet_pool_name"].(string))

			peeraddr := ipSecTunnels["public_ip_address"].(string)
			ikeVersion := ipSecTunnels["ike_version"].(string)
			poolname := ipSecTunnels["subnet_pool_name"].(string)

			log.Printf("LOGs CREATE FOR RncloudtemplateIpSecTunnel : %v ", fmt.Sprintf(models.RncloudtemplateIpSecTunnel, peeraddr))
			log.Printf("LOGs INSIDE FOR2 poolname  :: %s ", poolname)
			preSharedKey := ipSecTunnels["pre_shared_key"].(string)
			log.Printf("LOGs INSIDE FOR2 preSharedKey  :: %s ", preSharedKey)

			cloudtemplateIpSecTunnelAttr := models.CloudTemplateforIpSectunnelAttributes{}
			cloudtemplateIpSecTunnelAttr.Annotation = "{}"
			cloudtemplateIpSecTunnelAttr.IkeVersion = ikeVersion
			cloudtemplateIpSecTunnelAttr.Poolname = poolname
			cloudtemplateIpSecTunnelAttr.PreSharedKey = preSharedKey

			log.Printf("LOGs INSIDE FOR3 cloudtemplateIpSecTunnelAttr.Poolname :: %v ", cloudtemplateIpSecTunnelAttr.Poolname)

			// IkeVersion
			// Peeraddr
			// Poolname
			// PreSharedKey
			log.Printf("LOGs CREATE FOR RncloudtemplateIpSecTunnel : %v ", fmt.Sprintf(models.RncloudtemplateIpSecTunnel, peeraddr))
			log.Printf("LOGs CREATE FOR cloudVPNNetworkDn : %v ", cloudVPNNetworkDn)
			log.Printf("LOGs CREATE FOR cloudtemplateIpSecTunnelAttr : %v ", cloudtemplateIpSecTunnelAttr)

			cloudtemplateIpSecTunnel := models.NewCloudTemplateforIpSectunnel(fmt.Sprintf(models.RncloudtemplateIpSecTunnel, peeraddr), cloudVPNNetworkDn, cloudtemplateIpSecTunnelAttr)
			log.Printf("LOGs CREATE cloudtemplateIpSecTunnel : %v ", cloudtemplateIpSecTunnel)
			err := aciClient.Save(cloudtemplateIpSecTunnel)
			if err != nil {
				return diag.FromErr(err)
			}
			cloudIpSecTunnelPeers = append(cloudIpSecTunnelPeers, cloudtemplateIpSecTunnel.DistinguishedName)
		}
	}

	d.SetId(cloudtemplateVpnNetwork.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciCloudTemplateforVPNNetworkRead(ctx, d, m)
}

func resourceAciCloudTemplateforVPNNetworkUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] TemplateforVPNNetwork: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	TemplateforExternalNetworkDn := d.Get("aci_cloud_external_network_dn").(string)

	cloudtemplateVpnNetworkAttr := models.TemplateforVPNNetworkAttributes{}

	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		cloudtemplateVpnNetworkAttr.Annotation = Annotation.(string)
	} else {
		cloudtemplateVpnNetworkAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		cloudtemplateVpnNetworkAttr.Name = Name.(string)
	}

	if RemoteSiteId, ok := d.GetOk("remote_site_id"); ok {
		cloudtemplateVpnNetworkAttr.RemoteSiteId = RemoteSiteId.(string)
	}

	if RemoteSiteName, ok := d.GetOk("remote_site_name"); ok {
		cloudtemplateVpnNetworkAttr.RemoteSiteName = RemoteSiteName.(string)
	}
	cloudtemplateVpnNetwork := models.NewTemplateforVPNNetwork(fmt.Sprintf("vpnnetwork-%s", name), TemplateforExternalNetworkDn, desc, nameAlias, cloudtemplateVpnNetworkAttr)

	cloudtemplateVpnNetwork.Status = "modified"

	err := aciClient.Save(cloudtemplateVpnNetwork)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(cloudtemplateVpnNetwork.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciCloudTemplateforVPNNetworkRead(ctx, d, m)
}

func resourceAciCloudTemplateforVPNNetworkRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	cloudtemplateVpnNetwork, err := getRemoteTemplateforVPNNetwork(aciClient, dn)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	_, err = setTemplateforVPNNetworkAttributes(cloudtemplateVpnNetwork, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciCloudTemplateforVPNNetworkDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	err := aciClient.DeleteByDn(dn, "cloudtemplateVpnNetwork")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}
