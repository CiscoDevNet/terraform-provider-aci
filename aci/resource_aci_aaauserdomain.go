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

func resourceAciUserDomain() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciUserDomainCreate,
		UpdateContext: resourceAciUserDomainUpdate,
		ReadContext:   resourceAciUserDomainRead,
		DeleteContext: resourceAciUserDomainDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciUserDomainImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"local_user_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		})),
	}
}

func getRemoteUserDomain(client *client.Client, dn string) (*models.UserDomain, error) {
	aaaUserDomainCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	aaaUserDomain := models.UserDomainFromContainer(aaaUserDomainCont)
	if aaaUserDomain.DistinguishedName == "" {
		return nil, fmt.Errorf("UserDomain %s not found", aaaUserDomain.DistinguishedName)
	}
	return aaaUserDomain, nil
}

func setUserDomainAttributes(aaaUserDomain *models.UserDomain, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(aaaUserDomain.DistinguishedName)
	d.Set("description", aaaUserDomain.Description)
	aaaUserDomainMap, err := aaaUserDomain.ToMap()
	if err != nil {
		return nil, err
	}
	d.Set("local_user_dn", GetParentDn(d.Id(), fmt.Sprintf("/userdomain-%s", aaaUserDomainMap["name"])))
	d.Set("annotation", aaaUserDomainMap["annotation"])
	d.Set("name", aaaUserDomainMap["name"])
	d.Set("name_alias", aaaUserDomainMap["nameAlias"])
	return d, nil
}

func resourceAciUserDomainImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	aaaUserDomain, err := getRemoteUserDomain(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setUserDomainAttributes(aaaUserDomain, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciUserDomainCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] UserDomain: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	LocalUserDn := d.Get("local_user_dn").(string)

	aaaUserDomainAttr := models.UserDomainAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		aaaUserDomainAttr.Annotation = Annotation.(string)
	} else {
		aaaUserDomainAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		aaaUserDomainAttr.Name = Name.(string)
	}
	aaaUserDomain := models.NewUserDomain(fmt.Sprintf("userdomain-%s", name), LocalUserDn, desc, nameAlias, aaaUserDomainAttr)

	err := aciClient.Save(aaaUserDomain)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(aaaUserDomain.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciUserDomainRead(ctx, d, m)
}

func resourceAciUserDomainUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] UserDomain: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	LocalUserDn := d.Get("local_user_dn").(string)
	aaaUserDomainAttr := models.UserDomainAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		aaaUserDomainAttr.Annotation = Annotation.(string)
	} else {
		aaaUserDomainAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		aaaUserDomainAttr.Name = Name.(string)
	}
	aaaUserDomain := models.NewUserDomain(fmt.Sprintf("userdomain-%s", name), LocalUserDn, desc, nameAlias, aaaUserDomainAttr)

	aaaUserDomain.Status = "modified"
	err := aciClient.Save(aaaUserDomain)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(aaaUserDomain.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciUserDomainRead(ctx, d, m)
}

func resourceAciUserDomainRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	aaaUserDomain, err := getRemoteUserDomain(aciClient, dn)
	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setUserDomainAttributes(aaaUserDomain, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciUserDomainDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "aaaUserDomain")
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}
