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

func resourceAciLDAPGroupMap() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciLDAPGroupMapCreate,
		UpdateContext: resourceAciLDAPGroupMapUpdate,
		ReadContext:   resourceAciLDAPGroupMapRead,
		DeleteContext: resourceAciLDAPGroupMapDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciLDAPGroupMapImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ldap",
					"duo",
				}, false),
			},
		})),
	}
}

func getRemoteLDAPGroupMap(client *client.Client, dn string) (*models.LDAPGroupMap, error) {
	aaaLdapGroupMapCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	aaaLdapGroupMap := models.LDAPGroupMapFromContainer(aaaLdapGroupMapCont)
	if aaaLdapGroupMap.DistinguishedName == "" {
		return nil, fmt.Errorf("LDAP Group Map %s not found", dn)
	}
	return aaaLdapGroupMap, nil
}

func setLDAPGroupMapAttributes(aaaLdapGroupMap *models.LDAPGroupMap, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(aaaLdapGroupMap.DistinguishedName)
	d.Set("description", aaaLdapGroupMap.Description)
	aaaLdapGroupMapMap, err := aaaLdapGroupMap.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("name", aaaLdapGroupMapMap["name"])
	d.Set("name_alias", aaaLdapGroupMapMap["nameAlias"])
	return d, nil
}

func resourceAciLDAPGroupMapImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	aaaLdapGroupMap, err := getRemoteLDAPGroupMap(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setLDAPGroupMapAttributes(aaaLdapGroupMap, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciLDAPGroupMapCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] LDAPGroupMap: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	groupType := d.Get("type").(string)
	aaaLdapGroupMapAttr := models.LDAPGroupMapAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		aaaLdapGroupMapAttr.Annotation = Annotation.(string)
	} else {
		aaaLdapGroupMapAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		aaaLdapGroupMapAttr.Name = Name.(string)
	}
	aaaLdapGroupMap := models.NewLDAPGroupMap(fmt.Sprintf("userext/%sext/ldapgroupmap-%s", groupType, name), "uni", desc, nameAlias, aaaLdapGroupMapAttr)
	err := aciClient.Save(aaaLdapGroupMap)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(aaaLdapGroupMap.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciLDAPGroupMapRead(ctx, d, m)
}

func resourceAciLDAPGroupMapUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] LDAPGroupMap: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	groupType := d.Get("type").(string)
	aaaLdapGroupMapAttr := models.LDAPGroupMapAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		aaaLdapGroupMapAttr.Annotation = Annotation.(string)
	} else {
		aaaLdapGroupMapAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		aaaLdapGroupMapAttr.Name = Name.(string)
	}
	aaaLdapGroupMap := models.NewLDAPGroupMap(fmt.Sprintf("userext/%sext/ldapgroupmap-%s", groupType, name), "uni", desc, nameAlias, aaaLdapGroupMapAttr)
	aaaLdapGroupMap.Status = "modified"
	err := aciClient.Save(aaaLdapGroupMap)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(aaaLdapGroupMap.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciLDAPGroupMapRead(ctx, d, m)
}

func resourceAciLDAPGroupMapRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	aaaLdapGroupMap, err := getRemoteLDAPGroupMap(aciClient, dn)
	if err != nil {
		return errorForObjectNotFound(err, dn, d)
	}
	_, err = setLDAPGroupMapAttributes(aaaLdapGroupMap, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciLDAPGroupMapDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "aaaLdapGroupMap")
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}
