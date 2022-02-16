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

func resourceAciConsoleAuthenticationMethod() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciConsoleAuthenticationMethodCreate,
		UpdateContext: resourceAciConsoleAuthenticationMethodUpdate,
		ReadContext:   resourceAciConsoleAuthenticationMethodRead,
		DeleteContext: resourceAciConsoleAuthenticationMethodDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciConsoleAuthenticationMethodImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{

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

func GetRemoteConsoleAuthenticationMethod(client *client.Client, dn string) (*models.ConsoleAuthenticationMethod, error) {
	aaaConsoleAuthCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	aaaConsoleAuth := models.ConsoleAuthenticationMethodFromContainer(aaaConsoleAuthCont)
	if aaaConsoleAuth.DistinguishedName == "" {
		return nil, fmt.Errorf("ConsoleAuthenticationMethod %s not found", aaaConsoleAuth.DistinguishedName)
	}
	return aaaConsoleAuth, nil
}

func setConsoleAuthenticationMethodAttributes(aaaConsoleAuth *models.ConsoleAuthenticationMethod, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(aaaConsoleAuth.DistinguishedName)
	d.Set("description", aaaConsoleAuth.Description)
	aaaConsoleAuthMap, err := aaaConsoleAuth.ToMap()
	if err != nil {
		return nil, err
	}
	d.Set("annotation", aaaConsoleAuthMap["annotation"])
	d.Set("provider_group", aaaConsoleAuthMap["providerGroup"])
	d.Set("realm", aaaConsoleAuthMap["realm"])
	d.Set("realm_sub_type", aaaConsoleAuthMap["realmSubType"])
	d.Set("name_alias", aaaConsoleAuthMap["nameAlias"])
	return d, nil
}

func resourceAciConsoleAuthenticationMethodImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	aaaConsoleAuth, err := GetRemoteConsoleAuthenticationMethod(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setConsoleAuthenticationMethodAttributes(aaaConsoleAuth, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciConsoleAuthenticationMethodCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] ConsoleAuthenticationMethod: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	aaaConsoleAuthAttr := models.ConsoleAuthenticationMethodAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		aaaConsoleAuthAttr.Annotation = Annotation.(string)
	} else {
		aaaConsoleAuthAttr.Annotation = "{}"
	}

	if ProviderGroup, ok := d.GetOk("provider_group"); ok {
		aaaConsoleAuthAttr.ProviderGroup = ProviderGroup.(string)
	}

	if Realm, ok := d.GetOk("realm"); ok {
		aaaConsoleAuthAttr.Realm = Realm.(string)
	}

	if RealmSubType, ok := d.GetOk("realm_sub_type"); ok {
		aaaConsoleAuthAttr.RealmSubType = RealmSubType.(string)
	}
	aaaConsoleAuth := models.NewConsoleAuthenticationMethod(fmt.Sprintf("userext/authrealm/consoleauth"), "uni", desc, nameAlias, aaaConsoleAuthAttr)
	aaaConsoleAuth.Status = "modified"
	err := aciClient.Save(aaaConsoleAuth)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(aaaConsoleAuth.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciConsoleAuthenticationMethodRead(ctx, d, m)
}

func resourceAciConsoleAuthenticationMethodUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] ConsoleAuthenticationMethod: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	aaaConsoleAuthAttr := models.ConsoleAuthenticationMethodAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		aaaConsoleAuthAttr.Annotation = Annotation.(string)
	} else {
		aaaConsoleAuthAttr.Annotation = "{}"
	}

	if ProviderGroup, ok := d.GetOk("provider_group"); ok {
		aaaConsoleAuthAttr.ProviderGroup = ProviderGroup.(string)
	}

	if Realm, ok := d.GetOk("realm"); ok {
		aaaConsoleAuthAttr.Realm = Realm.(string)
	}

	if RealmSubType, ok := d.GetOk("realm_sub_type"); ok {
		aaaConsoleAuthAttr.RealmSubType = RealmSubType.(string)
	}
	aaaConsoleAuth := models.NewConsoleAuthenticationMethod(fmt.Sprintf("userext/authrealm/consoleauth"), "uni", desc, nameAlias, aaaConsoleAuthAttr)
	aaaConsoleAuth.Status = "modified"
	err := aciClient.Save(aaaConsoleAuth)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(aaaConsoleAuth.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciConsoleAuthenticationMethodRead(ctx, d, m)
}

func resourceAciConsoleAuthenticationMethodRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	aaaConsoleAuth, err := GetRemoteConsoleAuthenticationMethod(aciClient, dn)
	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setConsoleAuthenticationMethodAttributes(aaaConsoleAuth, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciConsoleAuthenticationMethodDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	d.SetId("")
	var diags diag.Diagnostics
	diags = append(diags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "Resource with class name aaaConsoleAuth cannot be deleted",
	})
	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	return diags
}
