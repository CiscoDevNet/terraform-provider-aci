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

func resourceAciLDAPProvider() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciLDAPProviderCreate,
		UpdateContext: resourceAciLDAPProviderUpdate,
		ReadContext:   resourceAciLDAPProviderRead,
		DeleteContext: resourceAciLDAPProviderDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciLDAPProviderImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{

			"ssl_validation_level": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"permissive",
					"strict",
				}, false),
			},
			"attribute": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"basedn": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enable_ssl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
			},
			"filter": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
				ValidateFunc: validation.StringInSlice([]string{"duo", "ldap"}, false),
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
			"rootdn": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"timeout": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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

func getRemoteLDAPProvider(client *client.Client, dn string) (*models.LDAPProvider, error) {
	aaaLdapProviderCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	aaaLdapProvider := models.LDAPProviderFromContainer(aaaLdapProviderCont)
	if aaaLdapProvider.DistinguishedName == "" {
		return nil, fmt.Errorf("LDAPProvider %s not found", aaaLdapProvider.DistinguishedName)
	}
	return aaaLdapProvider, nil
}

func setLDAPProviderAttributes(ldap_provider_group_type string, aaaLdapProvider *models.LDAPProvider, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(aaaLdapProvider.DistinguishedName)
	d.Set("description", aaaLdapProvider.Description)
	aaaLdapProviderMap, err := aaaLdapProvider.ToMap()
	if err != nil {
		return nil, err
	}
	d.Set("type", ldap_provider_group_type)
	d.Set("ssl_validation_level", aaaLdapProviderMap["SSLValidationLevel"])
	d.Set("annotation", aaaLdapProviderMap["annotation"])
	d.Set("attribute", aaaLdapProviderMap["attribute"])
	d.Set("basedn", aaaLdapProviderMap["basedn"])
	d.Set("enable_ssl", aaaLdapProviderMap["enableSSL"])
	d.Set("filter", aaaLdapProviderMap["filter"])
	d.Set("monitor_server", aaaLdapProviderMap["monitorServer"])
	d.Set("monitoring_user", aaaLdapProviderMap["monitoringUser"])
	d.Set("name", aaaLdapProviderMap["name"])
	d.Set("port", aaaLdapProviderMap["port"])
	d.Set("retries", aaaLdapProviderMap["retries"])
	d.Set("rootdn", aaaLdapProviderMap["rootdn"])
	d.Set("timeout", aaaLdapProviderMap["timeout"])
	d.Set("name_alias", aaaLdapProviderMap["nameAlias"])
	return d, nil
}

func resourceAciLDAPProviderImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	aaaLdapProvider, err := getRemoteLDAPProvider(aciClient, dn)
	if err != nil {
		return nil, err
	}

	ldap_provider_group := strings.Split(dn, "/")

	var ldap_group_type string
	if ldap_provider_group[2] == "ldapext" {
		ldap_group_type = "ldap"
	} else {
		ldap_group_type = "duo"
	}

	schemaFilled, err := setLDAPProviderAttributes(ldap_group_type, aaaLdapProvider, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciLDAPProviderCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] LDAPProvider: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	lsdp_type := d.Get("type").(string)
	aaaLdapProviderAttr := models.LDAPProviderAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		aaaLdapProviderAttr.Annotation = Annotation.(string)
	} else {
		aaaLdapProviderAttr.Annotation = "{}"
	}

	if SSLValidationLevel, ok := d.GetOk("ssl_validation_level"); ok {
		aaaLdapProviderAttr.SSLValidationLevel = SSLValidationLevel.(string)
	}

	if Attribute, ok := d.GetOk("attribute"); ok {
		aaaLdapProviderAttr.Attribute = Attribute.(string)
	}

	if Basedn, ok := d.GetOk("basedn"); ok {
		aaaLdapProviderAttr.Basedn = Basedn.(string)
	}

	if EnableSSL, ok := d.GetOk("enable_ssl"); ok {
		aaaLdapProviderAttr.EnableSSL = EnableSSL.(string)
	}

	if Filter, ok := d.GetOk("filter"); ok {
		aaaLdapProviderAttr.Filter = Filter.(string)
	}

	if Key, ok := d.GetOk("key"); ok {
		aaaLdapProviderAttr.Key = Key.(string)
	}

	if MonitorServer, ok := d.GetOk("monitor_server"); ok {
		aaaLdapProviderAttr.MonitorServer = MonitorServer.(string)
	}

	if MonitoringPassword, ok := d.GetOk("monitoring_password"); ok {
		aaaLdapProviderAttr.MonitoringPassword = MonitoringPassword.(string)
	}

	if MonitoringUser, ok := d.GetOk("monitoring_user"); ok {
		aaaLdapProviderAttr.MonitoringUser = MonitoringUser.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		aaaLdapProviderAttr.Name = Name.(string)
	}

	if Port, ok := d.GetOk("port"); ok {
		aaaLdapProviderAttr.Port = Port.(string)
	}

	if Retries, ok := d.GetOk("retries"); ok {
		aaaLdapProviderAttr.Retries = Retries.(string)
	}

	if Rootdn, ok := d.GetOk("rootdn"); ok {
		aaaLdapProviderAttr.Rootdn = Rootdn.(string)
	}

	if Timeout, ok := d.GetOk("timeout"); ok {
		aaaLdapProviderAttr.Timeout = Timeout.(string)
	}

	var aaaLdapProvider *models.LDAPProvider

	if lsdp_type == "duo" {
		aaaLdapProvider = models.NewLDAPProvider(fmt.Sprintf("userext/duoext/ldapprovider-%s", name), "uni", desc, nameAlias, aaaLdapProviderAttr)
	} else {
		aaaLdapProvider = models.NewLDAPProvider(fmt.Sprintf("userext/ldapext/ldapprovider-%s", name), "uni", desc, nameAlias, aaaLdapProviderAttr)
	}

	err := aciClient.Save(aaaLdapProvider)
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
		err = aciClient.CreateRelationaaaRsProvToEpp(aaaLdapProvider.DistinguishedName, aaaLdapProviderAttr.Annotation, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if relationToaaaRsSecProvToEpg, ok := d.GetOk("relation_aaa_rs_sec_prov_to_epg"); ok {
		relationParam := relationToaaaRsSecProvToEpg.(string)
		err = aciClient.CreateRelationaaaRsSecProvToEpg(aaaLdapProvider.DistinguishedName, aaaLdapProviderAttr.Annotation, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(aaaLdapProvider.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciLDAPProviderRead(ctx, d, m)
}

func resourceAciLDAPProviderUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] LDAPProvider: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	lsdp_type := d.Get("type").(string)
	aaaLdapProviderAttr := models.LDAPProviderAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		aaaLdapProviderAttr.Annotation = Annotation.(string)
	} else {
		aaaLdapProviderAttr.Annotation = "{}"
	}

	if SSLValidationLevel, ok := d.GetOk("ssl_validation_level"); ok {
		aaaLdapProviderAttr.SSLValidationLevel = SSLValidationLevel.(string)
	}

	if Attribute, ok := d.GetOk("attribute"); ok {
		aaaLdapProviderAttr.Attribute = Attribute.(string)
	}

	if Basedn, ok := d.GetOk("basedn"); ok {
		aaaLdapProviderAttr.Basedn = Basedn.(string)
	}

	if EnableSSL, ok := d.GetOk("enable_ssl"); ok {
		aaaLdapProviderAttr.EnableSSL = EnableSSL.(string)
	}

	if Filter, ok := d.GetOk("filter"); ok {
		aaaLdapProviderAttr.Filter = Filter.(string)
	}

	if Key, ok := d.GetOk("key"); ok {
		aaaLdapProviderAttr.Key = Key.(string)
	}

	if MonitorServer, ok := d.GetOk("monitor_server"); ok {
		aaaLdapProviderAttr.MonitorServer = MonitorServer.(string)
	}

	if MonitoringPassword, ok := d.GetOk("monitoring_password"); ok {
		aaaLdapProviderAttr.MonitoringPassword = MonitoringPassword.(string)
	}

	if MonitoringUser, ok := d.GetOk("monitoring_user"); ok {
		aaaLdapProviderAttr.MonitoringUser = MonitoringUser.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		aaaLdapProviderAttr.Name = Name.(string)
	}

	if Port, ok := d.GetOk("port"); ok {
		aaaLdapProviderAttr.Port = Port.(string)
	}

	if Retries, ok := d.GetOk("retries"); ok {
		aaaLdapProviderAttr.Retries = Retries.(string)
	}

	if Rootdn, ok := d.GetOk("rootdn"); ok {
		aaaLdapProviderAttr.Rootdn = Rootdn.(string)
	}

	if Timeout, ok := d.GetOk("timeout"); ok {
		aaaLdapProviderAttr.Timeout = Timeout.(string)
	}

	var aaaLdapProvider *models.LDAPProvider

	if lsdp_type == "duo" {
		aaaLdapProvider = models.NewLDAPProvider(fmt.Sprintf("userext/duoext/ldapprovider-%s", name), "uni", desc, nameAlias, aaaLdapProviderAttr)
	} else {
		aaaLdapProvider = models.NewLDAPProvider(fmt.Sprintf("userext/ldapext/ldapprovider-%s", name), "uni", desc, nameAlias, aaaLdapProviderAttr)
	}

	aaaLdapProvider.Status = "modified"
	err := aciClient.Save(aaaLdapProvider)
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
		err = aciClient.DeleteRelationaaaRsProvToEpp(aaaLdapProvider.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationaaaRsProvToEpp(aaaLdapProvider.DistinguishedName, aaaLdapProviderAttr.Annotation, newRelParam.(string))

		if err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("relation_aaa_rs_sec_prov_to_epg") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_aaa_rs_sec_prov_to_epg")
		err = aciClient.DeleteRelationaaaRsSecProvToEpg(aaaLdapProvider.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationaaaRsSecProvToEpg(aaaLdapProvider.DistinguishedName, aaaLdapProviderAttr.Annotation, newRelParam.(string))

		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(aaaLdapProvider.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciLDAPProviderRead(ctx, d, m)
}

func resourceAciLDAPProviderRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	aaaLdapProvider, err := getRemoteLDAPProvider(aciClient, dn)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	ldap_provider_group := strings.Split(dn, "/")

	var ldap_group_type string
	if ldap_provider_group[2] == "ldapext" {
		ldap_group_type = "ldap"
	} else {
		ldap_group_type = "duo"
	}
	_, err = setLDAPProviderAttributes(ldap_group_type, aaaLdapProvider, d)
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

func resourceAciLDAPProviderDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "aaaLdapProvider")
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}
