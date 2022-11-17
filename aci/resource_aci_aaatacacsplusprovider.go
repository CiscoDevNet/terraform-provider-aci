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

func resourceAciTACACSProvider() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciTACACSProviderCreate,
		UpdateContext: resourceAciTACACSProviderUpdate,
		ReadContext:   resourceAciTACACSProviderRead,
		DeleteContext: resourceAciTACACSProviderDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciTACACSProviderImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{

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
			"port": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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

func getRemoteTACACSProvider(client *client.Client, dn string) (*models.TACACSProvider, error) {
	aaaTacacsPlusProviderCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	aaaTacacsPlusProvider := models.TACACSProviderFromContainer(aaaTacacsPlusProviderCont)
	if aaaTacacsPlusProvider.DistinguishedName == "" {
		return nil, fmt.Errorf("TACACSProvider %s not found", aaaTacacsPlusProvider.DistinguishedName)
	}
	return aaaTacacsPlusProvider, nil
}

func setTACACSProviderAttributes(aaaTacacsPlusProvider *models.TACACSProvider, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(aaaTacacsPlusProvider.DistinguishedName)
	d.Set("description", aaaTacacsPlusProvider.Description)
	aaaTacacsPlusProviderMap, err := aaaTacacsPlusProvider.ToMap()
	if err != nil {
		return nil, err
	}
	d.Set("annotation", aaaTacacsPlusProviderMap["annotation"])
	d.Set("auth_protocol", aaaTacacsPlusProviderMap["authProtocol"])
	d.Set("monitor_server", aaaTacacsPlusProviderMap["monitorServer"])
	d.Set("monitoring_user", aaaTacacsPlusProviderMap["monitoringUser"])
	d.Set("name", aaaTacacsPlusProviderMap["name"])
	d.Set("port", aaaTacacsPlusProviderMap["port"])
	d.Set("retries", aaaTacacsPlusProviderMap["retries"])
	d.Set("timeout", aaaTacacsPlusProviderMap["timeout"])
	d.Set("name_alias", aaaTacacsPlusProviderMap["nameAlias"])
	return d, nil
}

func resourceAciTACACSProviderImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	aaaTacacsPlusProvider, err := getRemoteTACACSProvider(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setTACACSProviderAttributes(aaaTacacsPlusProvider, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciTACACSProviderCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] TACACSProvider: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	aaaTacacsPlusProviderAttr := models.TACACSProviderAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		aaaTacacsPlusProviderAttr.Annotation = Annotation.(string)
	} else {
		aaaTacacsPlusProviderAttr.Annotation = "{}"
	}

	if AuthProtocol, ok := d.GetOk("auth_protocol"); ok {
		aaaTacacsPlusProviderAttr.AuthProtocol = AuthProtocol.(string)
	}

	if Key, ok := d.GetOk("key"); ok {
		aaaTacacsPlusProviderAttr.Key = Key.(string)
	}

	if MonitorServer, ok := d.GetOk("monitor_server"); ok {
		aaaTacacsPlusProviderAttr.MonitorServer = MonitorServer.(string)
	}

	if MonitoringPassword, ok := d.GetOk("monitoring_password"); ok {
		aaaTacacsPlusProviderAttr.MonitoringPassword = MonitoringPassword.(string)
	}

	if MonitoringUser, ok := d.GetOk("monitoring_user"); ok {
		aaaTacacsPlusProviderAttr.MonitoringUser = MonitoringUser.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		aaaTacacsPlusProviderAttr.Name = Name.(string)
	}

	if Port, ok := d.GetOk("port"); ok {
		aaaTacacsPlusProviderAttr.Port = Port.(string)
	}

	if Retries, ok := d.GetOk("retries"); ok {
		aaaTacacsPlusProviderAttr.Retries = Retries.(string)
	}

	if Timeout, ok := d.GetOk("timeout"); ok {
		aaaTacacsPlusProviderAttr.Timeout = Timeout.(string)
	}
	aaaTacacsPlusProvider := models.NewTACACSProvider(fmt.Sprintf("userext/tacacsext/tacacsplusprovider-%s", name), "uni", desc, nameAlias, aaaTacacsPlusProviderAttr)
	err := aciClient.Save(aaaTacacsPlusProvider)
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
		err = aciClient.CreateRelationaaaRsProvToEpp(aaaTacacsPlusProvider.DistinguishedName, aaaTacacsPlusProviderAttr.Annotation, relationParam)

		if err != nil {
			return diag.FromErr(err)
		}

	}

	if relationToaaaRsSecProvToEpg, ok := d.GetOk("relation_aaa_rs_sec_prov_to_epg"); ok {
		relationParam := relationToaaaRsSecProvToEpg.(string)
		err = aciClient.CreateRelationaaaRsSecProvToEpg(aaaTacacsPlusProvider.DistinguishedName, aaaTacacsPlusProviderAttr.Annotation, relationParam)

		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(aaaTacacsPlusProvider.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciTACACSProviderRead(ctx, d, m)
}

func resourceAciTACACSProviderUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] TACACSProvider: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	aaaTacacsPlusProviderAttr := models.TACACSProviderAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		aaaTacacsPlusProviderAttr.Annotation = Annotation.(string)
	} else {
		aaaTacacsPlusProviderAttr.Annotation = "{}"
	}

	if AuthProtocol, ok := d.GetOk("auth_protocol"); ok {
		aaaTacacsPlusProviderAttr.AuthProtocol = AuthProtocol.(string)
	}

	if Key, ok := d.GetOk("key"); ok {
		aaaTacacsPlusProviderAttr.Key = Key.(string)
	}

	if MonitorServer, ok := d.GetOk("monitor_server"); ok {
		aaaTacacsPlusProviderAttr.MonitorServer = MonitorServer.(string)
	}

	if MonitoringPassword, ok := d.GetOk("monitoring_password"); ok {
		aaaTacacsPlusProviderAttr.MonitoringPassword = MonitoringPassword.(string)
	}

	if MonitoringUser, ok := d.GetOk("monitoring_user"); ok {
		aaaTacacsPlusProviderAttr.MonitoringUser = MonitoringUser.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		aaaTacacsPlusProviderAttr.Name = Name.(string)
	}

	if Port, ok := d.GetOk("port"); ok {
		aaaTacacsPlusProviderAttr.Port = Port.(string)
	}

	if Retries, ok := d.GetOk("retries"); ok {
		aaaTacacsPlusProviderAttr.Retries = Retries.(string)
	}

	if Timeout, ok := d.GetOk("timeout"); ok {
		aaaTacacsPlusProviderAttr.Timeout = Timeout.(string)
	}
	aaaTacacsPlusProvider := models.NewTACACSProvider(fmt.Sprintf("userext/tacacsext/tacacsplusprovider-%s", name), "uni", desc, nameAlias, aaaTacacsPlusProviderAttr)
	aaaTacacsPlusProvider.Status = "modified"
	err := aciClient.Save(aaaTacacsPlusProvider)
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
		err = aciClient.DeleteRelationaaaRsProvToEpp(aaaTacacsPlusProvider.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationaaaRsProvToEpp(aaaTacacsPlusProvider.DistinguishedName, aaaTacacsPlusProviderAttr.Annotation, newRelParam.(string))
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_aaa_rs_sec_prov_to_epg") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_aaa_rs_sec_prov_to_epg")
		err = aciClient.DeleteRelationaaaRsSecProvToEpg(aaaTacacsPlusProvider.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationaaaRsSecProvToEpg(aaaTacacsPlusProvider.DistinguishedName, aaaTacacsPlusProviderAttr.Annotation, newRelParam.(string))

		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(aaaTacacsPlusProvider.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciTACACSProviderRead(ctx, d, m)
}

func resourceAciTACACSProviderRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	aaaTacacsPlusProvider, err := getRemoteTACACSProvider(aciClient, dn)
	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setTACACSProviderAttributes(aaaTacacsPlusProvider, d)
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

func resourceAciTACACSProviderDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "aaaTacacsPlusProvider")
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}
