package aci

import (
	"context"
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAciCloudTemplateforExternalNetwork() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciCloudTemplateforExternalNetworkCreate,
		UpdateContext: resourceAciCloudTemplateforExternalNetworkUpdate,
		ReadContext:   resourceAciCloudTemplateforExternalNetworkRead,
		DeleteContext: resourceAciCloudTemplateforExternalNetworkDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciCloudTemplateforExternalNetworkImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"hub_network_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vrf_dn": {
				Type:     schema.TypeString,
				Required: true,
			},
		})),
	}
}

func getRemoteCloudTemplateforExternalNetwork(client *client.Client, dn string) (*models.CloudTemplateforExternalNetwork, error) {
	cloudtemplateExtNetworkCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	cloudtemplateExtNetwork := models.CloudTemplateforExternalNetworkFromContainer(cloudtemplateExtNetworkCont)
	if cloudtemplateExtNetwork.DistinguishedName == "" {
		return nil, fmt.Errorf("CloudTemplateforExternalNetwork %s not found", cloudtemplateExtNetwork.DistinguishedName)
	}
	return cloudtemplateExtNetwork, nil
}

func setCloudTemplateforExternalNetworkAttributes(cloudtemplateExtNetwork *models.CloudTemplateforExternalNetwork, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(cloudtemplateExtNetwork.DistinguishedName)

	cloudtemplateExtNetworkMap, err := cloudtemplateExtNetwork.ToMap()
	if err != nil {
		return d, err
	}

	d.Set("annotation", cloudtemplateExtNetworkMap["annotation"])
	d.Set("hub_network_name", cloudtemplateExtNetworkMap["hubNetworkName"])
	d.Set("name", cloudtemplateExtNetworkMap["name"])
	d.Set("name_alias", cloudtemplateExtNetworkMap["nameAlias"])
	if cloudtemplateExtNetworkMap["vrfName"] != "" {
		d.Set("vrf_dn", fmt.Sprintf("uni/tn-infra/ctx-%s", cloudtemplateExtNetworkMap["vrfName"]))
	} else {
		d.Set("vrf_dn", "")
	}
	return d, nil
}

func resourceAciCloudTemplateforExternalNetworkImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	cloudtemplateExtNetwork, err := getRemoteCloudTemplateforExternalNetwork(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setCloudTemplateforExternalNetworkAttributes(cloudtemplateExtNetwork, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciCloudTemplateforExternalNetworkCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] CloudTemplateforExternalNetwork: Beginning Creation")
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)
	CloudInfraNetworkTemplateDn := "uni/tn-infra/infranetwork-default"

	cloudtemplateExtNetworkAttr := models.CloudTemplateforExternalNetworkAttributes{}

	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		cloudtemplateExtNetworkAttr.Annotation = Annotation.(string)
	} else {
		cloudtemplateExtNetworkAttr.Annotation = "{}"
	}

	if HubNetworkName, ok := d.GetOk("hub_network_name"); ok {
		cloudtemplateExtNetworkAttr.HubNetworkName = HubNetworkName.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		cloudtemplateExtNetworkAttr.Name = Name.(string)
	}

	if _, ok := d.GetOk("vrf_dn"); ok {
		cloudtemplateExtNetworkAttr.VrfName = GetMOName(d.Get("vrf_dn").(string))
	}
	cloudtemplateExtNetwork := models.NewCloudTemplateforExternalNetwork(fmt.Sprintf(models.RncloudtemplateExtNetwork, name), CloudInfraNetworkTemplateDn, nameAlias, cloudtemplateExtNetworkAttr)

	err := aciClient.Save(cloudtemplateExtNetwork)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(cloudtemplateExtNetwork.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciCloudTemplateforExternalNetworkRead(ctx, d, m)
}

func resourceAciCloudTemplateforExternalNetworkUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] CloudTemplateforExternalNetwork: Beginning Update")
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)
	CloudInfraNetworkTemplateDn := "uni/tn-infra/infranetwork-default"

	cloudtemplateExtNetworkAttr := models.CloudTemplateforExternalNetworkAttributes{}

	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		cloudtemplateExtNetworkAttr.Annotation = Annotation.(string)
	} else {
		cloudtemplateExtNetworkAttr.Annotation = "{}"
	}

	if HubNetworkName, ok := d.GetOk("hub_network_name"); ok {
		cloudtemplateExtNetworkAttr.HubNetworkName = HubNetworkName.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		cloudtemplateExtNetworkAttr.Name = Name.(string)
	}

	if _, ok := d.GetOk("vrf_dn"); ok {
		cloudtemplateExtNetworkAttr.VrfName = GetMOName(d.Get("vrf_dn").(string))
	}
	cloudtemplateExtNetwork := models.NewCloudTemplateforExternalNetwork(fmt.Sprintf(models.RncloudtemplateExtNetwork, name), CloudInfraNetworkTemplateDn, nameAlias, cloudtemplateExtNetworkAttr)

	cloudtemplateExtNetwork.Status = "modified"

	err := aciClient.Save(cloudtemplateExtNetwork)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(cloudtemplateExtNetwork.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciCloudTemplateforExternalNetworkRead(ctx, d, m)
}

func resourceAciCloudTemplateforExternalNetworkRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	cloudtemplateExtNetwork, err := getRemoteCloudTemplateforExternalNetwork(aciClient, dn)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	_, err = setCloudTemplateforExternalNetworkAttributes(cloudtemplateExtNetwork, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciCloudTemplateforExternalNetworkDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	err := aciClient.DeleteByDn(dn, "cloudtemplateExtNetwork")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}
