package aci

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"sort"
	"strings"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAciDuoProviderGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciDuoProviderGroupCreate,
		UpdateContext: resourceAciDuoProviderGroupUpdate,
		ReadContext:   resourceAciDuoProviderGroupRead,
		DeleteContext: resourceAciDuoProviderGroupDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciDuoProviderGroupImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{

			"auth_choice": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"CiscoAVPair",
					"LdapGroupMap",
				}, false),
			},
			"ldap_group_map_ref": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"provider_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ldap",
					"radius",
				}, false),
			},
			"sec_fac_auth_methods": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"auto",
						"passcode",
						"phone",
						"push",
					}, false),
				},
			},
		})),
	}
}

func getRemoteDuoProviderGroup(client *client.Client, dn string) (*models.DuoProviderGroup, error) {
	aaaDuoProviderGroupCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	aaaDuoProviderGroup := models.DuoProviderGroupFromContainer(aaaDuoProviderGroupCont)
	if aaaDuoProviderGroup.DistinguishedName == "" {
		return nil, fmt.Errorf("DuoProviderGroup %s not found", aaaDuoProviderGroup.DistinguishedName)
	}
	return aaaDuoProviderGroup, nil
}

func setDuoProviderGroupAttributes(aaaDuoProviderGroup *models.DuoProviderGroup, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(aaaDuoProviderGroup.DistinguishedName)
	d.Set("description", aaaDuoProviderGroup.Description)
	aaaDuoProviderGroupMap, err := aaaDuoProviderGroup.ToMap()
	if err != nil {
		return nil, err
	}
	d.Set("annotation", aaaDuoProviderGroupMap["annotation"])
	d.Set("auth_choice", aaaDuoProviderGroupMap["authChoice"])
	d.Set("ldap_group_map_ref", aaaDuoProviderGroupMap["ldapGroupMapRef"])
	d.Set("name", aaaDuoProviderGroupMap["name"])
	d.Set("provider_type", aaaDuoProviderGroupMap["providerType"])
	secFacAuthMethodsGet := make([]string, 0, 1)
	if aaaDuoProviderGroupMap["secFacAuthMethods"] == "" {
		d.Set("sec_fac_auth_methods", secFacAuthMethodsGet)
	} else {
		for _, val := range strings.Split(aaaDuoProviderGroupMap["secFacAuthMethods"], ",") {
			secFacAuthMethodsGet = append(secFacAuthMethodsGet, strings.Trim(val, " "))
		}
		sort.Strings(secFacAuthMethodsGet)
		if secFacAuthMethodsIntr, ok := d.GetOk("sec_fac_auth_methods"); ok {
			secFacAuthMethodsAct := make([]string, 0, 1)
			for _, val := range secFacAuthMethodsIntr.([]interface{}) {
				secFacAuthMethodsAct = append(secFacAuthMethodsAct, val.(string))
			}
			sort.Strings(secFacAuthMethodsAct)
			if reflect.DeepEqual(secFacAuthMethodsAct, secFacAuthMethodsGet) {
				d.Set("sec_fac_auth_methods", d.Get("sec_fac_auth_methods").([]interface{}))
			} else {
				d.Set("sec_fac_auth_methods", secFacAuthMethodsGet)
			}
		} else {
			d.Set("sec_fac_auth_methods", secFacAuthMethodsGet)
		}
	}
	d.Set("name_alias", aaaDuoProviderGroupMap["nameAlias"])
	return d, nil
}

func resourceAciDuoProviderGroupImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	aaaDuoProviderGroup, err := getRemoteDuoProviderGroup(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setDuoProviderGroupAttributes(aaaDuoProviderGroup, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciDuoProviderGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] DuoProviderGroup: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	aaaDuoProviderGroupAttr := models.DuoProviderGroupAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		aaaDuoProviderGroupAttr.Annotation = Annotation.(string)
	} else {
		aaaDuoProviderGroupAttr.Annotation = "{}"
	}

	if AuthChoice, ok := d.GetOk("auth_choice"); ok {
		aaaDuoProviderGroupAttr.AuthChoice = AuthChoice.(string)
	}

	if LdapGroupMapRef, ok := d.GetOk("ldap_group_map_ref"); ok {
		aaaDuoProviderGroupAttr.LdapGroupMapRef = LdapGroupMapRef.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		aaaDuoProviderGroupAttr.Name = Name.(string)
	}

	if ProviderType, ok := d.GetOk("provider_type"); ok {
		aaaDuoProviderGroupAttr.ProviderType = ProviderType.(string)
	}

	if SecFacAuthMethods, ok := d.GetOk("sec_fac_auth_methods"); ok {
		secFacAuthMethodsList := make([]string, 0, 1)
		for _, val := range SecFacAuthMethods.([]interface{}) {
			secFacAuthMethodsList = append(secFacAuthMethodsList, val.(string))
		}
		err := checkDuplicate(secFacAuthMethodsList)
		if err != nil {
			return diag.FromErr(err)
		}
		SecFacAuthMethods := strings.Join(secFacAuthMethodsList, ",")
		aaaDuoProviderGroupAttr.SecFacAuthMethods = SecFacAuthMethods
	}
	aaaDuoProviderGroup := models.NewDuoProviderGroup(fmt.Sprintf("userext/duoext/duoprovidergroup-%s", name), "uni", desc, nameAlias, aaaDuoProviderGroupAttr)
	err := aciClient.Save(aaaDuoProviderGroup)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(aaaDuoProviderGroup.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciDuoProviderGroupRead(ctx, d, m)
}

func resourceAciDuoProviderGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] DuoProviderGroup: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	aaaDuoProviderGroupAttr := models.DuoProviderGroupAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		aaaDuoProviderGroupAttr.Annotation = Annotation.(string)
	} else {
		aaaDuoProviderGroupAttr.Annotation = "{}"
	}

	if AuthChoice, ok := d.GetOk("auth_choice"); ok {
		aaaDuoProviderGroupAttr.AuthChoice = AuthChoice.(string)
	}

	if LdapGroupMapRef, ok := d.GetOk("ldap_group_map_ref"); ok {
		aaaDuoProviderGroupAttr.LdapGroupMapRef = LdapGroupMapRef.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		aaaDuoProviderGroupAttr.Name = Name.(string)
	}

	if ProviderType, ok := d.GetOk("provider_type"); ok {
		aaaDuoProviderGroupAttr.ProviderType = ProviderType.(string)
	}
	if SecFacAuthMethods, ok := d.GetOk("sec_fac_auth_methods"); ok {
		secFacAuthMethodsList := make([]string, 0, 1)
		for _, val := range SecFacAuthMethods.([]interface{}) {
			secFacAuthMethodsList = append(secFacAuthMethodsList, val.(string))
		}
		err := checkDuplicate(secFacAuthMethodsList)
		if err != nil {
			return diag.FromErr(err)
		}
		SecFacAuthMethods := strings.Join(secFacAuthMethodsList, ",")
		aaaDuoProviderGroupAttr.SecFacAuthMethods = SecFacAuthMethods
	}
	aaaDuoProviderGroup := models.NewDuoProviderGroup(fmt.Sprintf("userext/duoext/duoprovidergroup-%s", name), "uni", desc, nameAlias, aaaDuoProviderGroupAttr)
	aaaDuoProviderGroup.Status = "modified"
	err := aciClient.Save(aaaDuoProviderGroup)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(aaaDuoProviderGroup.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciDuoProviderGroupRead(ctx, d, m)
}

func resourceAciDuoProviderGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	aaaDuoProviderGroup, err := getRemoteDuoProviderGroup(aciClient, dn)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	_, err = setDuoProviderGroupAttributes(aaaDuoProviderGroup, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciDuoProviderGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "aaaDuoProviderGroup")
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}
