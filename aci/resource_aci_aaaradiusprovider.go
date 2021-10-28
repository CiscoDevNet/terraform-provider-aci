package aci

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAciRADIUSProvider() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciRADIUSProviderCreate,
		UpdateContext: resourceAciRADIUSProviderUpdate,
		ReadContext:   resourceAciRADIUSProviderRead,
		DeleteContext: resourceAciRADIUSProviderDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciRADIUSProviderImport,
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
			"type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"duo", "radius"}, false),
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

func getRemoteRADIUSProvider(client *client.Client, dn string) (*models.RADIUSProvider, error) {
	aaaRadiusProviderCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	aaaRadiusProvider := models.RADIUSProviderFromContainer(aaaRadiusProviderCont)
	if aaaRadiusProvider.DistinguishedName == "" {
		return nil, fmt.Errorf("RADIUSProvider %s not found", aaaRadiusProvider.DistinguishedName)
	}
	return aaaRadiusProvider, nil
}

func setRADIUSProviderAttributes(radius_type string, aaaRadiusProvider *models.RADIUSProvider, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(aaaRadiusProvider.DistinguishedName)
	d.Set("description", aaaRadiusProvider.Description)
	aaaRadiusProviderMap, err := aaaRadiusProvider.ToMap()
	if err != nil {
		return nil, err
	}
	d.Set("type", radius_type)
	d.Set("auth_port", aaaRadiusProviderMap["authPort"])
	d.Set("auth_protocol", aaaRadiusProviderMap["authProtocol"])
	d.Set("monitor_server", aaaRadiusProviderMap["monitorServer"])
	d.Set("monitoring_user", aaaRadiusProviderMap["monitoringUser"])
	d.Set("name", aaaRadiusProviderMap["name"])
	d.Set("retries", aaaRadiusProviderMap["retries"])
	d.Set("timeout", aaaRadiusProviderMap["timeout"])
	d.Set("name_alias", aaaRadiusProviderMap["nameAlias"])
	return d, nil
}

func resourceAciRADIUSProviderImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	aaaRadiusProvider, err := getRemoteRADIUSProvider(aciClient, dn)
	if err != nil {
		return nil, err
	}

	aaa_radius_provider := strings.Split(dn, "/")
	aaa_radius_provider_type := aaa_radius_provider[2]
	var aaa_radius_provider_set string

	if aaa_radius_provider_type == "duoext" {
		aaa_radius_provider_set = "duo"
	} else {
		aaa_radius_provider_set = "radius"
	}

	schemaFilled, err := setRADIUSProviderAttributes(aaa_radius_provider_set, aaaRadiusProvider, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciRADIUSProviderCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] RADIUSProvider: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	radius_provider_type := d.Get("type").(string)
	aaaRadiusProviderAttr := models.RADIUSProviderAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		aaaRadiusProviderAttr.Annotation = Annotation.(string)
	} else {
		aaaRadiusProviderAttr.Annotation = "{}"
	}

	if AuthPort, ok := d.GetOk("auth_port"); ok {
		aaaRadiusProviderAttr.AuthPort = AuthPort.(string)
	}

	if AuthProtocol, ok := d.GetOk("auth_protocol"); ok {
		aaaRadiusProviderAttr.AuthProtocol = AuthProtocol.(string)
	}

	if Key, ok := d.GetOk("key"); ok {
		aaaRadiusProviderAttr.Key = Key.(string)
	}

	if MonitorServer, ok := d.GetOk("monitor_server"); ok {
		aaaRadiusProviderAttr.MonitorServer = MonitorServer.(string)
	}

	if MonitoringPassword, ok := d.GetOk("monitoring_password"); ok {
		aaaRadiusProviderAttr.MonitoringPassword = MonitoringPassword.(string)
	}

	if MonitoringUser, ok := d.GetOk("monitoring_user"); ok {
		aaaRadiusProviderAttr.MonitoringUser = MonitoringUser.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		aaaRadiusProviderAttr.Name = Name.(string)
	}

	if Retries, ok := d.GetOk("retries"); ok {
		aaaRadiusProviderAttr.Retries = Retries.(string)
	}

	if Timeout, ok := d.GetOk("timeout"); ok {
		aaaRadiusProviderAttr.Timeout = Timeout.(string)
	}

	var aaaRadiusProvider *models.RADIUSProvider

	if radius_provider_type == "duo" {
		aaaRadiusProvider = models.NewRADIUSProvider(fmt.Sprintf("userext/duoext/radiusprovider-%s", name), "uni", desc, nameAlias, aaaRadiusProviderAttr)
	} else {
		aaaRadiusProvider = models.NewRADIUSProvider(fmt.Sprintf("userext/radiusext/radiusprovider-%s", name), "uni", desc, nameAlias, aaaRadiusProviderAttr)
	}
	err := aciClient.Save(aaaRadiusProvider)
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

	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}

	if relationToaaaRsProvToEpp, ok := d.GetOk("relation_aaa_rs_prov_to_epp"); ok {
		relationParam := relationToaaaRsProvToEpp.(string)
		err = aciClient.CreateRelationaaaRsProvToEpp(aaaRadiusProvider.DistinguishedName, aaaRadiusProviderAttr.Annotation, relationParam)

		if err != nil {
			return diag.FromErr(err)
		}

	}

	if relationToaaaRsSecProvToEpg, ok := d.GetOk("relation_aaa_rs_sec_prov_to_epg"); ok {
		relationParam := relationToaaaRsSecProvToEpg.(string)
		err = aciClient.CreateRelationaaaRsSecProvToEpg(aaaRadiusProvider.DistinguishedName, aaaRadiusProviderAttr.Annotation, relationParam)

		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(aaaRadiusProvider.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciRADIUSProviderRead(ctx, d, m)
}

func resourceAciRADIUSProviderUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] RADIUSProvider: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	radius_provider_type := d.Get("type").(string)
	aaaRadiusProviderAttr := models.RADIUSProviderAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		aaaRadiusProviderAttr.Annotation = Annotation.(string)
	} else {
		aaaRadiusProviderAttr.Annotation = "{}"
	}

	if AuthPort, ok := d.GetOk("auth_port"); ok {
		aaaRadiusProviderAttr.AuthPort = AuthPort.(string)
	}

	if AuthProtocol, ok := d.GetOk("auth_protocol"); ok {
		aaaRadiusProviderAttr.AuthProtocol = AuthProtocol.(string)
	}

	if Key, ok := d.GetOk("key"); ok {
		aaaRadiusProviderAttr.Key = Key.(string)
	}

	if MonitorServer, ok := d.GetOk("monitor_server"); ok {
		aaaRadiusProviderAttr.MonitorServer = MonitorServer.(string)
	}

	if MonitoringPassword, ok := d.GetOk("monitoring_password"); ok {
		aaaRadiusProviderAttr.MonitoringPassword = MonitoringPassword.(string)
	}

	if MonitoringUser, ok := d.GetOk("monitoring_user"); ok {
		aaaRadiusProviderAttr.MonitoringUser = MonitoringUser.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		aaaRadiusProviderAttr.Name = Name.(string)
	}

	if Retries, ok := d.GetOk("retries"); ok {
		aaaRadiusProviderAttr.Retries = Retries.(string)
	}

	if Timeout, ok := d.GetOk("timeout"); ok {
		aaaRadiusProviderAttr.Timeout = Timeout.(string)
	}

	var aaaRadiusProvider *models.RADIUSProvider
	if radius_provider_type == "duo" {
		aaaRadiusProvider = models.NewRADIUSProvider(fmt.Sprintf("userext/duoext/radiusprovider-%s", name), "uni", desc, nameAlias, aaaRadiusProviderAttr)
	} else {
		aaaRadiusProvider = models.NewRADIUSProvider(fmt.Sprintf("userext/radiusext/radiusprovider-%s", name), "uni", desc, nameAlias, aaaRadiusProviderAttr)
	}

	aaaRadiusProvider.Status = "modified"
	err := aciClient.Save(aaaRadiusProvider)
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

	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("relation_aaa_rs_prov_to_epp") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_aaa_rs_prov_to_epp")
		err = aciClient.DeleteRelationaaaRsProvToEpp(aaaRadiusProvider.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationaaaRsProvToEpp(aaaRadiusProvider.DistinguishedName, aaaRadiusProviderAttr.Annotation, newRelParam.(string))

		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_aaa_rs_sec_prov_to_epg") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_aaa_rs_sec_prov_to_epg")
		err = aciClient.DeleteRelationaaaRsSecProvToEpg(aaaRadiusProvider.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationaaaRsSecProvToEpg(aaaRadiusProvider.DistinguishedName, aaaRadiusProviderAttr.Annotation, newRelParam.(string))

		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(aaaRadiusProvider.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciRADIUSProviderRead(ctx, d, m)
}

func resourceAciRADIUSProviderRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	aaaRadiusProvider, err := getRemoteRADIUSProvider(aciClient, dn)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	radius_provider_group := strings.Split(dn, "/")
	var radius_group_type string
	if radius_provider_group[2] == "duoext" {
		radius_group_type = "duo"
	} else {
		radius_group_type = "radius"
	}
	_, err = setRADIUSProviderAttributes(radius_group_type, aaaRadiusProvider, d)
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

func resourceAciRADIUSProviderDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "aaaRadiusProvider")
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}
