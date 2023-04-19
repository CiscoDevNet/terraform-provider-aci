package aci

import (
	"context"
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAciSAMLProvider() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciSAMLProviderCreate,
		UpdateContext: resourceAciSAMLProviderUpdate,
		ReadContext:   resourceAciSAMLProviderRead,
		DeleteContext: resourceAciSAMLProviderDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciSAMLProviderImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{

			"entity_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"gui_banner_message": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"https_proxy": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"id_p": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"adfs",
					"okta",
					"ping identity",
				}, false),
			},
			"key": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"metadata_url": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"monitor_server": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"disabled",
					"enabled",
				}, false),
			},
			"monitoring_password": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"monitoring_user": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"retries": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: checkAtleastOne(),
			},
			"sig_alg": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"SIG_RSA_SHA1",
					"SIG_RSA_SHA224",
					"SIG_RSA_SHA256",
					"SIG_RSA_SHA384",
					"SIG_RSA_SHA512",
				}, false),
			},
			"timeout": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tp": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"want_assertions_encrypted": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
			},
			"want_assertions_signed": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
			},
			"want_requests_signed": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
			},
			"want_response_signed": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
			},

			"relation_aaa_rs_prov_to_epp": &schema.Schema{
				Type: schema.TypeString,

				Optional:    true,
				Description: "Create relation to fv:AREpP",
			},
			"relation_aaa_rs_sec_prov_to_epg": &schema.Schema{
				Type: schema.TypeString,

				Optional:    true,
				Description: "Create relation to fv:ATg",
			}})),
	}
}

func getRemoteSAMLProvider(client *client.Client, dn string) (*models.SAMLProvider, error) {
	aaaSamlProviderCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	aaaSamlProvider := models.SAMLProviderFromContainer(aaaSamlProviderCont)
	if aaaSamlProvider.DistinguishedName == "" {
		return nil, fmt.Errorf("SAML Provider %s not found", dn)
	}
	return aaaSamlProvider, nil
}

func setSAMLProviderAttributes(aaaSamlProvider *models.SAMLProvider, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(aaaSamlProvider.DistinguishedName)
	d.Set("description", aaaSamlProvider.Description)
	aaaSamlProviderMap, err := aaaSamlProvider.ToMap()
	if err != nil {
		return nil, err
	}
	d.Set("annotation", aaaSamlProviderMap["annotation"])
	d.Set("entity_id", aaaSamlProviderMap["entityId"])
	d.Set("gui_banner_message", aaaSamlProviderMap["guiBannerMessage"])
	d.Set("https_proxy", aaaSamlProviderMap["httpsProxy"])
	d.Set("id_p", aaaSamlProviderMap["idP"])
	d.Set("metadata_url", aaaSamlProviderMap["metadataUrl"])
	d.Set("monitor_server", aaaSamlProviderMap["monitorServer"])
	d.Set("monitoring_user", aaaSamlProviderMap["monitoringUser"])
	d.Set("name", aaaSamlProviderMap["name"])
	d.Set("retries", aaaSamlProviderMap["retries"])
	d.Set("sig_alg", aaaSamlProviderMap["sigAlg"])
	d.Set("timeout", aaaSamlProviderMap["timeout"])
	d.Set("tp", aaaSamlProviderMap["tp"])
	d.Set("want_assertions_encrypted", aaaSamlProviderMap["wantAssertionsEncrypted"])
	d.Set("want_assertions_signed", aaaSamlProviderMap["wantAssertionsSigned"])
	d.Set("want_requests_signed", aaaSamlProviderMap["wantRequestsSigned"])
	d.Set("want_response_signed", aaaSamlProviderMap["wantResponseSigned"])
	d.Set("name_alias", aaaSamlProviderMap["nameAlias"])
	return d, nil
}

func resourceAciSAMLProviderImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	aaaSamlProvider, err := getRemoteSAMLProvider(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setSAMLProviderAttributes(aaaSamlProvider, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciSAMLProviderCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] SAMLProvider: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	aaaSamlProviderAttr := models.SAMLProviderAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		aaaSamlProviderAttr.Annotation = Annotation.(string)
	} else {
		aaaSamlProviderAttr.Annotation = "{}"
	}

	if EntityId, ok := d.GetOk("entity_id"); ok {
		aaaSamlProviderAttr.EntityId = EntityId.(string)
	}

	if GuiBannerMessage, ok := d.GetOk("gui_banner_message"); ok {
		aaaSamlProviderAttr.GuiBannerMessage = GuiBannerMessage.(string)
	}

	if HttpsProxy, ok := d.GetOk("https_proxy"); ok {
		aaaSamlProviderAttr.HttpsProxy = HttpsProxy.(string)
	}

	if IdP, ok := d.GetOk("id_p"); ok {
		aaaSamlProviderAttr.IdP = IdP.(string)
	}

	if Key, ok := d.GetOk("key"); ok {
		aaaSamlProviderAttr.Key = Key.(string)
	}

	if MetadataUrl, ok := d.GetOk("metadata_url"); ok {
		aaaSamlProviderAttr.MetadataUrl = MetadataUrl.(string)
	}

	if MonitorServer, ok := d.GetOk("monitor_server"); ok {
		aaaSamlProviderAttr.MonitorServer = MonitorServer.(string)
	}

	if MonitoringPassword, ok := d.GetOk("monitoring_password"); ok {
		aaaSamlProviderAttr.MonitoringPassword = MonitoringPassword.(string)
	}

	if MonitoringUser, ok := d.GetOk("monitoring_user"); ok {
		aaaSamlProviderAttr.MonitoringUser = MonitoringUser.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		aaaSamlProviderAttr.Name = Name.(string)
	}

	if Retries, ok := d.GetOk("retries"); ok {
		aaaSamlProviderAttr.Retries = Retries.(string)
	}

	if SigAlg, ok := d.GetOk("sig_alg"); ok {
		aaaSamlProviderAttr.SigAlg = SigAlg.(string)
	}

	if Timeout, ok := d.GetOk("timeout"); ok {
		aaaSamlProviderAttr.Timeout = Timeout.(string)
	}

	if Tp, ok := d.GetOk("tp"); ok {
		aaaSamlProviderAttr.Tp = Tp.(string)
	}

	if WantAssertionsEncrypted, ok := d.GetOk("want_assertions_encrypted"); ok {
		aaaSamlProviderAttr.WantAssertionsEncrypted = WantAssertionsEncrypted.(string)
	}

	if WantAssertionsSigned, ok := d.GetOk("want_assertions_signed"); ok {
		aaaSamlProviderAttr.WantAssertionsSigned = WantAssertionsSigned.(string)
	}

	if WantRequestsSigned, ok := d.GetOk("want_requests_signed"); ok {
		aaaSamlProviderAttr.WantRequestsSigned = WantRequestsSigned.(string)
	}

	if WantResponseSigned, ok := d.GetOk("want_response_signed"); ok {
		aaaSamlProviderAttr.WantResponseSigned = WantResponseSigned.(string)
	}
	aaaSamlProvider := models.NewSAMLProvider(fmt.Sprintf("userext/samlext/samlprovider-%s", name), "uni", desc, nameAlias, aaaSamlProviderAttr)
	err := aciClient.Save(aaaSamlProvider)
	if err != nil {
		return diag.FromErr(err)
	}
	checkDns := make([]string, 0, 1)

	if relationToaaaRsProvToEpp, ok := d.GetOk("relation_aaa_rs_prov_to_epp"); ok {
		relationParam := relationToaaaRsProvToEpp.(string)
		checkDns = append(checkDns, relationParam)

	}

	if relationToaaaRsSecProvToEpg, ok := d.GetOk("relation_aaa_rs_sec_prov_to_epg"); ok {
		relationParam := relationToaaaRsSecProvToEpg.(string)
		checkDns = append(checkDns, relationParam)

	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if relationToaaaRsProvToEpp, ok := d.GetOk("relation_aaa_rs_prov_to_epp"); ok {
		relationParam := relationToaaaRsProvToEpp.(string)
		err = aciClient.CreateRelationaaaRsProvToEpp(aaaSamlProvider.DistinguishedName, aaaSamlProviderAttr.Annotation, relationParam)

		if err != nil {
			return diag.FromErr(err)
		}

	}

	if relationToaaaRsSecProvToEpg, ok := d.GetOk("relation_aaa_rs_sec_prov_to_epg"); ok {
		relationParam := relationToaaaRsSecProvToEpg.(string)
		err = aciClient.CreateRelationaaaRsSecProvToEpg(aaaSamlProvider.DistinguishedName, aaaSamlProviderAttr.Annotation, relationParam)

		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(aaaSamlProvider.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciSAMLProviderRead(ctx, d, m)
}

func resourceAciSAMLProviderUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] SAMLProvider: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	aaaSamlProviderAttr := models.SAMLProviderAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		aaaSamlProviderAttr.Annotation = Annotation.(string)
	} else {
		aaaSamlProviderAttr.Annotation = "{}"
	}

	if EntityId, ok := d.GetOk("entity_id"); ok {
		aaaSamlProviderAttr.EntityId = EntityId.(string)
	}

	if GuiBannerMessage, ok := d.GetOk("gui_banner_message"); ok {
		aaaSamlProviderAttr.GuiBannerMessage = GuiBannerMessage.(string)
	}

	if HttpsProxy, ok := d.GetOk("https_proxy"); ok {
		aaaSamlProviderAttr.HttpsProxy = HttpsProxy.(string)
	}

	if IdP, ok := d.GetOk("id_p"); ok {
		aaaSamlProviderAttr.IdP = IdP.(string)
	}

	if Key, ok := d.GetOk("key"); ok {
		aaaSamlProviderAttr.Key = Key.(string)
	}

	if MetadataUrl, ok := d.GetOk("metadata_url"); ok {
		aaaSamlProviderAttr.MetadataUrl = MetadataUrl.(string)
	}

	if MonitorServer, ok := d.GetOk("monitor_server"); ok {
		aaaSamlProviderAttr.MonitorServer = MonitorServer.(string)
	}

	if MonitoringPassword, ok := d.GetOk("monitoring_password"); ok {
		aaaSamlProviderAttr.MonitoringPassword = MonitoringPassword.(string)
	}

	if MonitoringUser, ok := d.GetOk("monitoring_user"); ok {
		aaaSamlProviderAttr.MonitoringUser = MonitoringUser.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		aaaSamlProviderAttr.Name = Name.(string)
	}

	if Retries, ok := d.GetOk("retries"); ok {
		aaaSamlProviderAttr.Retries = Retries.(string)
	}

	if SigAlg, ok := d.GetOk("sig_alg"); ok {
		aaaSamlProviderAttr.SigAlg = SigAlg.(string)
	}

	if Timeout, ok := d.GetOk("timeout"); ok {
		aaaSamlProviderAttr.Timeout = Timeout.(string)
	}

	if Tp, ok := d.GetOk("tp"); ok {
		aaaSamlProviderAttr.Tp = Tp.(string)
	}

	if WantAssertionsEncrypted, ok := d.GetOk("want_assertions_encrypted"); ok {
		aaaSamlProviderAttr.WantAssertionsEncrypted = WantAssertionsEncrypted.(string)
	}

	if WantAssertionsSigned, ok := d.GetOk("want_assertions_signed"); ok {
		aaaSamlProviderAttr.WantAssertionsSigned = WantAssertionsSigned.(string)
	}

	if WantRequestsSigned, ok := d.GetOk("want_requests_signed"); ok {
		aaaSamlProviderAttr.WantRequestsSigned = WantRequestsSigned.(string)
	}

	if WantResponseSigned, ok := d.GetOk("want_response_signed"); ok {
		aaaSamlProviderAttr.WantResponseSigned = WantResponseSigned.(string)
	}
	aaaSamlProvider := models.NewSAMLProvider(fmt.Sprintf("userext/samlext/samlprovider-%s", name), "uni", desc, nameAlias, aaaSamlProviderAttr)
	aaaSamlProvider.Status = "modified"
	err := aciClient.Save(aaaSamlProvider)
	if err != nil {
		return diag.FromErr(err)
	}
	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_aaa_rs_prov_to_epp") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_aaa_rs_prov_to_epp")
		checkDns = append(checkDns, newRelParam.(string))

	}

	if d.HasChange("relation_aaa_rs_sec_prov_to_epg") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_aaa_rs_sec_prov_to_epg")
		checkDns = append(checkDns, newRelParam.(string))

	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if d.HasChange("relation_aaa_rs_prov_to_epp") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_aaa_rs_prov_to_epp")
		err = aciClient.DeleteRelationaaaRsProvToEpp(aaaSamlProvider.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationaaaRsProvToEpp(aaaSamlProvider.DistinguishedName, aaaSamlProviderAttr.Annotation, newRelParam.(string))

		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_aaa_rs_sec_prov_to_epg") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_aaa_rs_sec_prov_to_epg")
		err = aciClient.DeleteRelationaaaRsSecProvToEpg(aaaSamlProvider.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationaaaRsSecProvToEpg(aaaSamlProvider.DistinguishedName, aaaSamlProviderAttr.Annotation, newRelParam.(string))

		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(aaaSamlProvider.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciSAMLProviderRead(ctx, d, m)
}

func resourceAciSAMLProviderRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	aaaSamlProvider, err := getRemoteSAMLProvider(aciClient, dn)
	if err != nil {
		return errorForObjectNotFound(err, dn, d)
	}
	_, err = setSAMLProviderAttributes(aaaSamlProvider, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	aaaRsProvToEppData, err := aciClient.ReadRelationaaaRsProvToEpp(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation aaaRsProvToEpp %v", err)
		d.Set("relation_aaa_rs_prov_to_epp", "")
	} else {
		setRelationAttribute(d, "relation_aaa_rs_prov_to_epp", aaaRsProvToEppData)
	}

	aaaRsSecProvToEpgData, err := aciClient.ReadRelationaaaRsSecProvToEpg(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation aaaRsSecProvToEpg %v", err)
		d.Set("relation_aaa_rs_sec_prov_to_epg", "")
	} else {
		setRelationAttribute(d, "relation_aaa_rs_sec_prov_to_epg", aaaRsSecProvToEpgData)
	}
	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciSAMLProviderDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "aaaSamlProvider")
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}
