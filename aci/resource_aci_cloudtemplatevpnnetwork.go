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
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ike_version": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
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
							Required: true,
							ForceNew: true,
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

	cloudtemplateVpnNetworkMap, err := cloudtemplateVpnNetwork.ToMap()
	if err != nil {
		return d, err
	}

	if dn != cloudtemplateVpnNetwork.DistinguishedName {
		d.Set("aci_cloud_external_network_dn", "")
	} else {
		d.Set("aci_cloud_external_network_dn", GetParentDn(dn, "/"+fmt.Sprintf(models.RncloudtemplateVpnNetwork, cloudtemplateVpnNetworkMap["name"])))
	}
	log.Printf(" LOGs GET getRemoteCloudTemplateforIpSecTunnel : %v ", dn)
	log.Printf(" LOGs GET RN : %v ", models.RncloudtemplateIpSecTunnel)
	log.Printf(" LOGs GET  GetParentDn   %v ", GetParentDn(dn, "/"+fmt.Sprintf(models.RncloudtemplateVpnNetwork, cloudtemplateVpnNetworkMap["name"])))

	d.Set("annotation", cloudtemplateVpnNetworkMap["annotation"])
	d.Set("name", cloudtemplateVpnNetworkMap["name"])
	d.Set("remote_site_id", cloudtemplateVpnNetworkMap["remoteSiteId"])
	d.Set("remote_site_name", cloudtemplateVpnNetworkMap["remoteSiteName"])
	d.Set("name_alias", cloudtemplateVpnNetworkMap["nameAlias"])

	return d, nil
}

func getRemoteCloudTemplateforIpSecTunnel(client *client.Client, dn string) (*models.CloudTemplateforIpSectunnel, error) {
	cloudtemplateIpSecTunnelCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	log.Printf("LOGs GET getRemoteCloudTemplateforIpSecTunnel : %v ", cloudtemplateIpSecTunnelCont)
	cloudtemplateIpSecTunnel := models.CloudTemplateforIpSectunnelFromContainer(cloudtemplateIpSecTunnelCont)
	if cloudtemplateIpSecTunnel.DistinguishedName == "" {
		return nil, fmt.Errorf("CloudTemplateforIpSectunnel %s not found", cloudtemplateIpSecTunnel.DistinguishedName)
	}
	return cloudtemplateIpSecTunnel, nil
}

func setCloudTemplateforIpSecTunnelAttributes(cloudtemplateIpSecTunnel *models.CloudTemplateforIpSectunnel, d map[string]string) (map[string]string, error) {

	cloudtemplateIpSecTunnelMap, err := cloudtemplateIpSecTunnel.ToMap()
	if err != nil {
		return nil, err
	}

	d = map[string]string{
		"ike_version":       cloudtemplateIpSecTunnelMap["ikeVersion"],
		"public_ip_address": cloudtemplateIpSecTunnelMap["peeraddr"],
		"subnet_pool_name":  cloudtemplateIpSecTunnelMap["poolname"],
		"pre_shared_key":    cloudtemplateIpSecTunnelMap["preSharedKey"],
	}

	log.Printf("LOGs SET Dd : %v ", d)
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

	cloudtemplateIpSecTunnelData, err := aciClient.ListCloudTemplateforIpSectunnel(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading cloud IPSec Tunnel attributes %v", err)
	}
	cloudtemplateIpSecTunnelSet := make([]map[string]string, 0, 1)
	for _, cloudtemplateIpSecTunnel := range cloudtemplateIpSecTunnelData {

		cloudIpSecTunnelAttMap, err := setCloudTemplateforIpSecTunnelAttributes(cloudtemplateIpSecTunnel, make(map[string]string))
		if err != nil {
			d.SetId("")
			return nil, err
		}
		log.Printf("LOGs IN ELSEEE cloudIpSecTunnelAttMap : %v ", cloudIpSecTunnelAttMap)
		cloudtemplateIpSecTunnelSet = append(cloudtemplateIpSecTunnelSet, cloudIpSecTunnelAttMap)
		log.Printf("LOGs IN ELSEEE cloudtemplateIpSecTunnelSet : %v ", cloudtemplateIpSecTunnelSet)
	}
	d.Set("cloud_ipsec_tunnel", cloudtemplateIpSecTunnelSet)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciCloudTemplateforVPNNetworkCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] TemplateforVPNNetwork: Beginning Creation")
	aciClient := m.(*client.Client)
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
	cloudtemplateVpnNetwork := models.NewTemplateforVPNNetwork(fmt.Sprintf(models.RncloudtemplateVpnNetwork, name), TemplateforExternalNetworkDn, nameAlias, cloudtemplateVpnNetworkAttr)

	err := aciClient.Save(cloudtemplateVpnNetwork)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("LOGs CREATE cloudtemplateVpnNetwork : %v ", cloudtemplateVpnNetwork)
	log.Printf("LOGs CREATE d.GetOk.cloud_ipsec_tunnel : %v ", d.Get("cloud_ipsec_tunnel"))

	cloudIpSecTunnelPeers := make([]string, 0, 1)
	if ipSecTunnelPeers, ok := d.GetOk("cloud_ipsec_tunnel"); ok {
		clopudIpSecTunnels := ipSecTunnelPeers.(*schema.Set).List()
		for _, val := range clopudIpSecTunnels {
			ipSecTunnels := val.(map[string]interface{})

			cloudtemplateIpSecTunnelAttr := models.CloudTemplateforIpSectunnelAttributes{}
			cloudtemplateIpSecTunnelAttr.Annotation = "{}"
			cloudtemplateIpSecTunnelAttr.IkeVersion = ipSecTunnels["ike_version"].(string)
			cloudtemplateIpSecTunnelAttr.Poolname = ipSecTunnels["subnet_pool_name"].(string)
			cloudtemplateIpSecTunnelAttr.PreSharedKey = ipSecTunnels["pre_shared_key"].(string)
			cloudtemplateIpSecTunnelAttr.Peeraddr = ipSecTunnels["public_ip_address"].(string)

			cloudtemplateIpSecTunnel := models.NewCloudTemplateforIpSectunnel(fmt.Sprintf(models.RncloudtemplateIpSecTunnel, cloudtemplateIpSecTunnelAttr.Peeraddr), cloudtemplateVpnNetwork.DistinguishedName, cloudtemplateIpSecTunnelAttr)
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
	cloudtemplateVpnNetwork := models.NewTemplateforVPNNetwork(fmt.Sprintf(models.RncloudtemplateVpnNetwork, name), TemplateforExternalNetworkDn, nameAlias, cloudtemplateVpnNetworkAttr)

	cloudtemplateVpnNetwork.Status = "modified"

	err := aciClient.Save(cloudtemplateVpnNetwork)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("LOGs UPDATE cloudtemplateVpnNetwork : %v ", cloudtemplateVpnNetwork)
	cloudIpSecTunnelPeers := make([]string, 0, 1)
	if ipSecTunnelPeers, ok := d.GetOk("cloud_ipsec_tunnel"); ok {
		clopudIpSecTunnels := ipSecTunnelPeers.(*schema.Set).List()
		for _, val := range clopudIpSecTunnels {
			ipSecTunnels := val.(map[string]interface{})

			cloudtemplateIpSecTunnelAttr := models.CloudTemplateforIpSectunnelAttributes{}
			cloudtemplateIpSecTunnelAttr.Annotation = "{}"
			cloudtemplateIpSecTunnelAttr.IkeVersion = ipSecTunnels["ike_version"].(string)
			cloudtemplateIpSecTunnelAttr.Poolname = ipSecTunnels["subnet_pool_name"].(string)
			cloudtemplateIpSecTunnelAttr.PreSharedKey = ipSecTunnels["pre_shared_key"].(string)
			cloudtemplateIpSecTunnelAttr.Peeraddr = ipSecTunnels["public_ip_address"].(string)

			cloudtemplateIpSecTunnel := models.NewCloudTemplateforIpSectunnel(fmt.Sprintf(models.RncloudtemplateIpSecTunnel, cloudtemplateIpSecTunnelAttr.Peeraddr), cloudtemplateVpnNetwork.DistinguishedName, cloudtemplateIpSecTunnelAttr)
			err := aciClient.Save(cloudtemplateIpSecTunnel)
			if err != nil {
				return diag.FromErr(err)
			}
			cloudIpSecTunnelPeers = append(cloudIpSecTunnelPeers, cloudtemplateIpSecTunnel.DistinguishedName)
		}
	}

	d.SetId(cloudtemplateVpnNetwork.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciCloudTemplateforVPNNetworkRead(ctx, d, m)
}

func resourceAciCloudTemplateforVPNNetworkRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	log.Printf("LOGs I ENTER READ ")

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

	cloudtemplateIpSecTunnelData, err := aciClient.ListCloudTemplateforIpSectunnel(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading cloud IPSec Tunnel attributes %v", err)
	}
	cloudtemplateIpSecTunnelSet := make([]map[string]string, 0, 1)
	for _, cloudtemplateIpSecTunnel := range cloudtemplateIpSecTunnelData {

		cloudIpSecTunnelAttMap, err := setCloudTemplateforIpSecTunnelAttributes(cloudtemplateIpSecTunnel, make(map[string]string))
		if err != nil {
			d.SetId("")
			return nil
		}
		log.Printf("LOGs IN ELSEEE cloudIpSecTunnelAttMap : %v ", cloudIpSecTunnelAttMap)
		cloudtemplateIpSecTunnelSet = append(cloudtemplateIpSecTunnelSet, cloudIpSecTunnelAttMap)
		log.Printf("LOGs IN ELSEEE cloudtemplateIpSecTunnelSet : %v ", cloudtemplateIpSecTunnelSet)
	}
	d.Set("cloud_ipsec_tunnel", cloudtemplateIpSecTunnelSet)

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
