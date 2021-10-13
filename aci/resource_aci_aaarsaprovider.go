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

func resourceAciRSAProvider() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciRSAProviderCreate,
		UpdateContext: resourceAciRSAProviderUpdate,
		ReadContext:   resourceAciRSAProviderRead,
		DeleteContext: resourceAciRSAProviderDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciRSAProviderImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{

			"auth_port": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"auth_protocol": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"chap",
					"mschap",
					"pap",
				}, false),
			},
			"key": &schema.Schema{
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
			"timeout": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: checkAtleastOne(),
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

func getRemoteRSAProvider(client *client.Client, dn string) (*models.RSAProvider, error) {
	aaaRsaProviderCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	aaaRsaProvider := models.RSAProviderFromContainer(aaaRsaProviderCont)
	if aaaRsaProvider.DistinguishedName == "" {
		return nil, fmt.Errorf("RSAProvider %s not found", aaaRsaProvider.DistinguishedName)
	}
	return aaaRsaProvider, nil
}

func setRSAProviderAttributes(aaaRsaProvider *models.RSAProvider, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(aaaRsaProvider.DistinguishedName)
	d.Set("description", aaaRsaProvider.Description)
	aaaRsaProviderMap, err := aaaRsaProvider.ToMap()
	if err != nil {
		return nil, err
	}
	d.Set("annotation", aaaRsaProviderMap["annotation"])
	d.Set("auth_port", aaaRsaProviderMap["authPort"])
	d.Set("auth_protocol", aaaRsaProviderMap["authProtocol"])
	d.Set("key", aaaRsaProviderMap["key"])
	d.Set("monitor_server", aaaRsaProviderMap["monitorServer"])
	d.Set("monitoring_password", aaaRsaProviderMap["monitoringPassword"])
	d.Set("monitoring_user", aaaRsaProviderMap["monitoringUser"])
	d.Set("name", aaaRsaProviderMap["name"])
	d.Set("retries", aaaRsaProviderMap["retries"])
	d.Set("timeout", aaaRsaProviderMap["timeout"])
	d.Set("name_alias", aaaRsaProviderMap["nameAlias"])
	return d, nil
}

func resourceAciRSAProviderImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	aaaRsaProvider, err := getRemoteRSAProvider(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setRSAProviderAttributes(aaaRsaProvider, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciRSAProviderCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] RSAProvider: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	aaaRsaProviderAttr := models.RSAProviderAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		aaaRsaProviderAttr.Annotation = Annotation.(string)
	} else {
		aaaRsaProviderAttr.Annotation = "{}"
	}

	if AuthPort, ok := d.GetOk("auth_port"); ok {
		aaaRsaProviderAttr.AuthPort = AuthPort.(string)
	}

	if AuthProtocol, ok := d.GetOk("auth_protocol"); ok {
		aaaRsaProviderAttr.AuthProtocol = AuthProtocol.(string)
	}

	if Key, ok := d.GetOk("key"); ok {
		aaaRsaProviderAttr.Key = Key.(string)
	}

	if MonitorServer, ok := d.GetOk("monitor_server"); ok {
		aaaRsaProviderAttr.MonitorServer = MonitorServer.(string)
	}

	if MonitoringPassword, ok := d.GetOk("monitoring_password"); ok {
		aaaRsaProviderAttr.MonitoringPassword = MonitoringPassword.(string)
	}

	if MonitoringUser, ok := d.GetOk("monitoring_user"); ok {
		aaaRsaProviderAttr.MonitoringUser = MonitoringUser.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		aaaRsaProviderAttr.Name = Name.(string)
	}

	if Retries, ok := d.GetOk("retries"); ok {
		aaaRsaProviderAttr.Retries = Retries.(string)
	}

	if Timeout, ok := d.GetOk("timeout"); ok {
		aaaRsaProviderAttr.Timeout = Timeout.(string)
	}
	aaaRsaProvider := models.NewRSAProvider(fmt.Sprintf("userext/rsaext/rsaprovider-%s", name), "uni", desc, nameAlias, aaaRsaProviderAttr)
	err := aciClient.Save(aaaRsaProvider)
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
		err = aciClient.CreateRelationaaaRsProvToEpp(aaaRsaProvider.DistinguishedName, aaaRsaProviderAttr.Annotation, relationParam)

		if err != nil {
			return diag.FromErr(err)
		}

	}

	if relationToaaaRsSecProvToEpg, ok := d.GetOk("relation_aaa_rs_sec_prov_to_epg"); ok {
		relationParam := relationToaaaRsSecProvToEpg.(string)
		err = aciClient.CreateRelationaaaRsSecProvToEpg(aaaRsaProvider.DistinguishedName, aaaRsaProviderAttr.Annotation, relationParam)

		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(aaaRsaProvider.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciRSAProviderRead(ctx, d, m)
}

func resourceAciRSAProviderUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] RSAProvider: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	aaaRsaProviderAttr := models.RSAProviderAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		aaaRsaProviderAttr.Annotation = Annotation.(string)
	} else {
		aaaRsaProviderAttr.Annotation = "{}"
	}

	if AuthPort, ok := d.GetOk("auth_port"); ok {
		aaaRsaProviderAttr.AuthPort = AuthPort.(string)
	}

	if AuthProtocol, ok := d.GetOk("auth_protocol"); ok {
		aaaRsaProviderAttr.AuthProtocol = AuthProtocol.(string)
	}

	if Key, ok := d.GetOk("key"); ok {
		aaaRsaProviderAttr.Key = Key.(string)
	}

	if MonitorServer, ok := d.GetOk("monitor_server"); ok {
		aaaRsaProviderAttr.MonitorServer = MonitorServer.(string)
	}

	if MonitoringPassword, ok := d.GetOk("monitoring_password"); ok {
		aaaRsaProviderAttr.MonitoringPassword = MonitoringPassword.(string)
	}

	if MonitoringUser, ok := d.GetOk("monitoring_user"); ok {
		aaaRsaProviderAttr.MonitoringUser = MonitoringUser.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		aaaRsaProviderAttr.Name = Name.(string)
	}

	if Retries, ok := d.GetOk("retries"); ok {
		aaaRsaProviderAttr.Retries = Retries.(string)
	}

	if Timeout, ok := d.GetOk("timeout"); ok {
		aaaRsaProviderAttr.Timeout = Timeout.(string)
	}
	aaaRsaProvider := models.NewRSAProvider(fmt.Sprintf("userext/rsaext/rsaprovider-%s", name), "uni", desc, nameAlias, aaaRsaProviderAttr)
	aaaRsaProvider.Status = "modified"
	err := aciClient.Save(aaaRsaProvider)
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
		err = aciClient.DeleteRelationaaaRsProvToEpp(aaaRsaProvider.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationaaaRsProvToEpp(aaaRsaProvider.DistinguishedName, aaaRsaProviderAttr.Annotation, newRelParam.(string))

		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_aaa_rs_sec_prov_to_epg") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_aaa_rs_sec_prov_to_epg")
		err = aciClient.DeleteRelationaaaRsSecProvToEpg(aaaRsaProvider.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationaaaRsSecProvToEpg(aaaRsaProvider.DistinguishedName, aaaRsaProviderAttr.Annotation, newRelParam.(string))

		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(aaaRsaProvider.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciRSAProviderRead(ctx, d, m)
}

func resourceAciRSAProviderRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	aaaRsaProvider, err := getRemoteRSAProvider(aciClient, dn)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	_, err = setRSAProviderAttributes(aaaRsaProvider, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	aaaRsProvToEppData, err := aciClient.ReadRelationaaaRsProvToEpp(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation aaaRsProvToEpp %v", err)
		d.Set("relation_aaa_rs_prov_to_epp", "")
	} else {
		if _, ok := d.GetOk("relation_aaa_rs_prov_to_epp"); ok {
			tfName := d.Get("relation_aaa_rs_prov_to_epp").(string)
			if tfName != aaaRsProvToEppData {
				d.Set("relation_aaa_rs_prov_to_epp", "")
			}
		}
	}

	aaaRsSecProvToEpgData, err := aciClient.ReadRelationaaaRsSecProvToEpg(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation aaaRsSecProvToEpg %v", err)
		d.Set("relation_aaa_rs_sec_prov_to_epg", "")
	} else {
		if _, ok := d.GetOk("relation_aaa_rs_sec_prov_to_epg"); ok {
			tfName := d.Get("relation_aaa_rs_sec_prov_to_epg").(string)
			if tfName != aaaRsSecProvToEpgData {
				d.Set("relation_aaa_rs_sec_prov_to_epg", "")
			}
		}
	}
	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciRSAProviderDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "aaaRsaProvider")
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}
