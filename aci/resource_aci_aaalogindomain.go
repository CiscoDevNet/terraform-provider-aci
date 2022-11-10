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

func resourceAciLoginDomain() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciLoginDomainCreate,
		UpdateContext: resourceAciLoginDomainUpdate,
		ReadContext:   resourceAciLoginDomainRead,
		DeleteContext: resourceAciLoginDomainDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciLoginDomainImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
					"none",
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

func getRemoteAuthenticationMethodfortheDomain(client *client.Client, dn string) (*models.AuthenticationMethodfortheDomain, error) {
	aaaDomainAuthCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	aaaDomainAuth := models.AuthenticationMethodfortheDomainFromContainer(aaaDomainAuthCont)
	if aaaDomainAuth.DistinguishedName == "" {
		return nil, fmt.Errorf("AuthenticationMethodfortheDomain %s not found", aaaDomainAuth.DistinguishedName)
	}
	return aaaDomainAuth, nil
}

func getRemoteLoginDomain(client *client.Client, dn string) (*models.LoginDomain, error) {
	aaaLoginDomainCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	aaaLoginDomain := models.LoginDomainFromContainer(aaaLoginDomainCont)
	if aaaLoginDomain.DistinguishedName == "" {
		return nil, fmt.Errorf("LoginDomain %s not found", aaaLoginDomain.DistinguishedName)
	}
	return aaaLoginDomain, nil
}

func setAuthenticationMethodfortheDomainAttributes(aaaDomainAuth *models.AuthenticationMethodfortheDomain, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.Set("description", aaaDomainAuth.Description)
	aaaDomainAuthMap, err := aaaDomainAuth.ToMap()
	if err != nil {
		return nil, err
	}
	d.Set("annotation", aaaDomainAuthMap["annotation"])
	d.Set("provider_group", aaaDomainAuthMap["providerGroup"])
	d.Set("realm", aaaDomainAuthMap["realm"])
	d.Set("realm_sub_type", aaaDomainAuthMap["realmSubType"])
	d.Set("name_alias", aaaDomainAuthMap["nameAlias"])
	return d, nil
}

func setLoginDomainAttributes(aaaLoginDomain *models.LoginDomain, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(aaaLoginDomain.DistinguishedName)
	d.Set("description", aaaLoginDomain.Description)
	aaaLoginDomainMap, err := aaaLoginDomain.ToMap()
	if err != nil {
		return nil, err
	}
	d.Set("annotation", aaaLoginDomainMap["annotation"])
	d.Set("name", aaaLoginDomainMap["name"])
	d.Set("name_alias", aaaLoginDomainMap["nameAlias"])
	return d, nil
}

func resourceAciLoginDomainImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	aaaLoginDomain, err := getRemoteLoginDomain(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setLoginDomainAttributes(aaaLoginDomain, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciLoginDomainCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] LoginDomain: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	aaaLoginDomainAttr := models.LoginDomainAttributes{}
	aaaDomainAuthAttr := models.AuthenticationMethodfortheDomainAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		aaaLoginDomainAttr.Annotation = Annotation.(string)
		aaaDomainAuthAttr.Annotation = Annotation.(string)
	} else {
		aaaLoginDomainAttr.Annotation = "{}"
		aaaDomainAuthAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		aaaLoginDomainAttr.Name = Name.(string)
	}

	if ProviderGroup, ok := d.GetOk("provider_group"); ok {
		aaaDomainAuthAttr.ProviderGroup = ProviderGroup.(string)
	}

	if Realm, ok := d.GetOk("realm"); ok {
		aaaDomainAuthAttr.Realm = Realm.(string)
	}

	if RealmSubType, ok := d.GetOk("realm_sub_type"); ok {
		aaaDomainAuthAttr.RealmSubType = RealmSubType.(string)
	}
	aaaLoginDomain := models.NewLoginDomain(fmt.Sprintf("userext/logindomain-%s", name), "uni", desc, nameAlias, aaaLoginDomainAttr)
	err := aciClient.Save(aaaLoginDomain)
	if err != nil {
		return diag.FromErr(err)
	}

	aaaDomainAuth := models.NewAuthenticationMethodfortheDomain(fmt.Sprintf("domainauth"), aaaLoginDomain.DistinguishedName, desc, nameAlias, aaaDomainAuthAttr)
	err = aciClient.Save(aaaDomainAuth)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(aaaLoginDomain.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciLoginDomainRead(ctx, d, m)
}

func resourceAciLoginDomainUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] LoginDomain: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	aaaLoginDomainAttr := models.LoginDomainAttributes{}
	aaaDomainAuthAttr := models.AuthenticationMethodfortheDomainAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		aaaLoginDomainAttr.Annotation = Annotation.(string)
		aaaDomainAuthAttr.Annotation = Annotation.(string)
	} else {
		aaaLoginDomainAttr.Annotation = "{}"
		aaaDomainAuthAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		aaaLoginDomainAttr.Name = Name.(string)
	}

	if ProviderGroup, ok := d.GetOk("provider_group"); ok {
		aaaDomainAuthAttr.ProviderGroup = ProviderGroup.(string)
	}

	if Realm, ok := d.GetOk("realm"); ok {
		aaaDomainAuthAttr.Realm = Realm.(string)
	}

	if RealmSubType, ok := d.GetOk("realm_sub_type"); ok {
		aaaDomainAuthAttr.RealmSubType = RealmSubType.(string)
	}
	aaaLoginDomain := models.NewLoginDomain(fmt.Sprintf("userext/logindomain-%s", name), "uni", desc, nameAlias, aaaLoginDomainAttr)
	aaaLoginDomain.Status = "modified"
	err := aciClient.Save(aaaLoginDomain)
	if err != nil {
		return diag.FromErr(err)
	}

	aaaDomainAuth := models.NewAuthenticationMethodfortheDomain(fmt.Sprintf("domainauth"), aaaLoginDomain.DistinguishedName, desc, nameAlias, aaaDomainAuthAttr)
	aaaDomainAuth.Status = "modified"
	err = aciClient.Save(aaaDomainAuth)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(aaaLoginDomain.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciLoginDomainRead(ctx, d, m)
}

func resourceAciLoginDomainRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	log.Printf("%s", dn)
	childDn := fmt.Sprintf("%s/domainauth", dn)
	log.Printf("%s", childDn)
	aaaLoginDomain, err := getRemoteLoginDomain(aciClient, dn)
	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setLoginDomainAttributes(aaaLoginDomain, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	aaaDomainAuth, err := getRemoteAuthenticationMethodfortheDomain(aciClient, childDn)
	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setAuthenticationMethodfortheDomainAttributes(aaaDomainAuth, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciLoginDomainDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	childDn := fmt.Sprintf("%s/domainauth", dn)
	err := aciClient.DeleteByDn(dn, "aaaLoginDomain")
	if err != nil {
		return diag.FromErr(err)
	}
	err = aciClient.DeleteByDn(childDn, "aaaDomainAuth")
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}
