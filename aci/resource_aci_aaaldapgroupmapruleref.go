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

func resourceAciLDAPGroupMapruleref() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciLDAPGroupMaprulerefCreate,
		UpdateContext: resourceAciLDAPGroupMaprulerefUpdate,
		ReadContext:   resourceAciLDAPGroupMaprulerefRead,
		DeleteContext: resourceAciLDAPGroupMaprulerefDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciLDAPGroupMaprulerefImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"ldap_group_map_dn": &schema.Schema{
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

func getRemoteLDAPGroupMapruleref(client *client.Client, dn string) (*models.LDAPGroupMapruleref, error) {
	aaaLdapGroupMapRuleRefCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	aaaLdapGroupMapRuleRef := models.LDAPGroupMaprulerefFromContainer(aaaLdapGroupMapRuleRefCont)
	if aaaLdapGroupMapRuleRef.DistinguishedName == "" {
		return nil, fmt.Errorf("LDAPGroupMapruleref %s not found", aaaLdapGroupMapRuleRef.DistinguishedName)
	}
	return aaaLdapGroupMapRuleRef, nil
}

func setLDAPGroupMaprulerefAttributes(aaaLdapGroupMapRuleRef *models.LDAPGroupMapruleref, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(aaaLdapGroupMapRuleRef.DistinguishedName)
	d.Set("description", aaaLdapGroupMapRuleRef.Description)
	aaaLdapGroupMapRuleRefMap, err := aaaLdapGroupMapRuleRef.ToMap()
	if err != nil {
		return nil, err
	}
	d.Set("ldap_group_map_dn", GetParentDn(d.Id(), fmt.Sprintf("/ldapgroupmapruleref-%s", aaaLdapGroupMapRuleRefMap["name"])))
	d.Set("name", aaaLdapGroupMapRuleRefMap["name"])
	d.Set("name_alias", aaaLdapGroupMapRuleRefMap["nameAlias"])
	d.Set("annotation", aaaLdapGroupMapRuleRefMap["annotation"])
	return d, nil
}

func resourceAciLDAPGroupMaprulerefImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	aaaLdapGroupMapRuleRef, err := getRemoteLDAPGroupMapruleref(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setLDAPGroupMaprulerefAttributes(aaaLdapGroupMapRuleRef, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciLDAPGroupMaprulerefCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] LDAPGroupMapruleref: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	LDAPGroupMapDn := d.Get("ldap_group_map_dn").(string)

	aaaLdapGroupMapRuleRefAttr := models.LDAPGroupMaprulerefAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		aaaLdapGroupMapRuleRefAttr.Annotation = Annotation.(string)
	} else {
		aaaLdapGroupMapRuleRefAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		aaaLdapGroupMapRuleRefAttr.Name = Name.(string)
	}
	aaaLdapGroupMapRuleRef := models.NewLDAPGroupMapruleref(fmt.Sprintf("ldapgroupmapruleref-%s", name), LDAPGroupMapDn, desc, nameAlias, aaaLdapGroupMapRuleRefAttr)

	err := aciClient.Save(aaaLdapGroupMapRuleRef)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(aaaLdapGroupMapRuleRef.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciLDAPGroupMaprulerefRead(ctx, d, m)
}

func resourceAciLDAPGroupMaprulerefUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] LDAPGroupMapruleref: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	LDAPGroupMapDn := d.Get("ldap_group_map_dn").(string)
	aaaLdapGroupMapRuleRefAttr := models.LDAPGroupMaprulerefAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		aaaLdapGroupMapRuleRefAttr.Annotation = Annotation.(string)
	} else {
		aaaLdapGroupMapRuleRefAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		aaaLdapGroupMapRuleRefAttr.Name = Name.(string)
	}
	aaaLdapGroupMapRuleRef := models.NewLDAPGroupMapruleref(fmt.Sprintf("ldapgroupmapruleref-%s", name), LDAPGroupMapDn, desc, nameAlias, aaaLdapGroupMapRuleRefAttr)

	aaaLdapGroupMapRuleRef.Status = "modified"
	err := aciClient.Save(aaaLdapGroupMapRuleRef)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(aaaLdapGroupMapRuleRef.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciLDAPGroupMaprulerefRead(ctx, d, m)
}

func resourceAciLDAPGroupMaprulerefRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	aaaLdapGroupMapRuleRef, err := getRemoteLDAPGroupMapruleref(aciClient, dn)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	_, err = setLDAPGroupMaprulerefAttributes(aaaLdapGroupMapRuleRef, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciLDAPGroupMaprulerefDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "aaaLdapGroupMapRuleRef")
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}
