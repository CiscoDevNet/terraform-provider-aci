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

func resourceAciDefaultAuthenticationMethodforallLogins() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciDefaultAuthenticationMethodforallLoginsCreate,
		UpdateContext: resourceAciDefaultAuthenticationMethodforallLoginsUpdate,
		ReadContext:   resourceAciDefaultAuthenticationMethodforallLoginsRead,
		DeleteContext: resourceAciDefaultAuthenticationMethodforallLoginsDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciDefaultAuthenticationMethodforallLoginsImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{

			"fallback_check": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"false",
					"true",
				}, false),
			},
			"provider_group": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"realm": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ldap",
					"local",
					"radius",
					"rsa",
					"saml",
					"tacacs",
				}, false),
			},
			"realm_sub_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"default",
					"duo",
				}, false),
			},
		})),
	}
}

func getRemoteDefaultAuthenticationMethodforallLogins(client *client.Client, dn string) (*models.DefaultAuthenticationMethodforallLogins, error) {
	aaaDefaultAuthCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	aaaDefaultAuth := models.DefaultAuthenticationMethodforallLoginsFromContainer(aaaDefaultAuthCont)
	if aaaDefaultAuth.DistinguishedName == "" {
		return nil, fmt.Errorf("DefaultAuthenticationMethodforallLogins %s not found", aaaDefaultAuth.DistinguishedName)
	}
	return aaaDefaultAuth, nil
}

func setDefaultAuthenticationMethodforallLoginsAttributes(aaaDefaultAuth *models.DefaultAuthenticationMethodforallLogins, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(aaaDefaultAuth.DistinguishedName)
	d.Set("description", aaaDefaultAuth.Description)
	aaaDefaultAuthMap, err := aaaDefaultAuth.ToMap()
	if err != nil {
		return nil, err
	}
	d.Set("annotation", aaaDefaultAuthMap["annotation"])
	d.Set("fallback_check", aaaDefaultAuthMap["fallbackCheck"])
	d.Set("provider_group", aaaDefaultAuthMap["providerGroup"])
	d.Set("realm", aaaDefaultAuthMap["realm"])
	d.Set("realm_sub_type", aaaDefaultAuthMap["realmSubType"])
	d.Set("name_alias", aaaDefaultAuthMap["nameAlias"])
	return d, nil
}

func resourceAciDefaultAuthenticationMethodforallLoginsImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	aaaDefaultAuth, err := getRemoteDefaultAuthenticationMethodforallLogins(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setDefaultAuthenticationMethodforallLoginsAttributes(aaaDefaultAuth, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciDefaultAuthenticationMethodforallLoginsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] DefaultAuthenticationMethodforallLogins: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	aaaDefaultAuthAttr := models.DefaultAuthenticationMethodforallLoginsAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		aaaDefaultAuthAttr.Annotation = Annotation.(string)
	} else {
		aaaDefaultAuthAttr.Annotation = "{}"
	}

	if FallbackCheck, ok := d.GetOk("fallback_check"); ok {
		aaaDefaultAuthAttr.FallbackCheck = FallbackCheck.(string)
	}

	aaaDefaultAuthAttr.Name = ""

	if ProviderGroup, ok := d.GetOk("provider_group"); ok {
		aaaDefaultAuthAttr.ProviderGroup = ProviderGroup.(string)
	}

	if Realm, ok := d.GetOk("realm"); ok {
		aaaDefaultAuthAttr.Realm = Realm.(string)
	}

	if RealmSubType, ok := d.GetOk("realm_sub_type"); ok {
		aaaDefaultAuthAttr.RealmSubType = RealmSubType.(string)
	}
	aaaDefaultAuth := models.NewDefaultAuthenticationMethodforallLogins(fmt.Sprintf("userext/authrealm/defaultauth"), "uni", desc, nameAlias, aaaDefaultAuthAttr)
	aaaDefaultAuth.Status = "modified"
	err := aciClient.Save(aaaDefaultAuth)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(aaaDefaultAuth.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciDefaultAuthenticationMethodforallLoginsRead(ctx, d, m)
}

func resourceAciDefaultAuthenticationMethodforallLoginsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] DefaultAuthenticationMethodforallLogins: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	aaaDefaultAuthAttr := models.DefaultAuthenticationMethodforallLoginsAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		aaaDefaultAuthAttr.Annotation = Annotation.(string)
	} else {
		aaaDefaultAuthAttr.Annotation = "{}"
	}

	if FallbackCheck, ok := d.GetOk("fallback_check"); ok {
		aaaDefaultAuthAttr.FallbackCheck = FallbackCheck.(string)
	}

	aaaDefaultAuthAttr.Name = ""

	if ProviderGroup, ok := d.GetOk("provider_group"); ok {
		aaaDefaultAuthAttr.ProviderGroup = ProviderGroup.(string)
	}

	if Realm, ok := d.GetOk("realm"); ok {
		aaaDefaultAuthAttr.Realm = Realm.(string)
	}

	if RealmSubType, ok := d.GetOk("realm_sub_type"); ok {
		aaaDefaultAuthAttr.RealmSubType = RealmSubType.(string)
	}
	aaaDefaultAuth := models.NewDefaultAuthenticationMethodforallLogins(fmt.Sprintf("userext/authrealm/defaultauth"), "uni", desc, nameAlias, aaaDefaultAuthAttr)
	aaaDefaultAuth.Status = "modified"
	err := aciClient.Save(aaaDefaultAuth)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(aaaDefaultAuth.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciDefaultAuthenticationMethodforallLoginsRead(ctx, d, m)
}

func resourceAciDefaultAuthenticationMethodforallLoginsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	aaaDefaultAuth, err := getRemoteDefaultAuthenticationMethodforallLogins(aciClient, dn)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	_, err = setDefaultAuthenticationMethodforallLoginsAttributes(aaaDefaultAuth, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciDefaultAuthenticationMethodforallLoginsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	d.SetId("")
	var diags diag.Diagnostics
	diags = append(diags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "Resource with class name aaaDefaultAuth cannot be deleted",
	})
	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	return diags
}
