package aci

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAciLDAPGroupMapRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciLDAPGroupMapRuleCreate,
		UpdateContext: resourceAciLDAPGroupMapRuleUpdate,
		ReadContext:   resourceAciLDAPGroupMapRuleRead,
		DeleteContext: resourceAciLDAPGroupMapRuleDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciLDAPGroupMapRuleImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{

			"groupdn": &schema.Schema{
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
		})),
	}
}

func getRemoteLDAPGroupMapRule(client *client.Client, dn string) (*models.LDAPGroupMapRule, error) {
	aaaLdapGroupMapRuleCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	aaaLdapGroupMapRule := models.LDAPGroupMapRuleFromContainer(aaaLdapGroupMapRuleCont)
	if aaaLdapGroupMapRule.DistinguishedName == "" {
		return nil, fmt.Errorf("LDAP Group Map Rule %s not found", dn)
	}
	return aaaLdapGroupMapRule, nil
}

func setLDAPGroupMapRuleAttributes(ldap_group_map_rule_type string, aaaLdapGroupMapRule *models.LDAPGroupMapRule, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(aaaLdapGroupMapRule.DistinguishedName)
	d.Set("description", aaaLdapGroupMapRule.Description)
	aaaLdapGroupMapRuleMap, err := aaaLdapGroupMapRule.ToMap()
	if err != nil {
		return nil, err
	}
	d.Set("type", ldap_group_map_rule_type)
	d.Set("annotation", aaaLdapGroupMapRuleMap["annotation"])
	d.Set("groupdn", aaaLdapGroupMapRuleMap["groupdn"])
	d.Set("name", aaaLdapGroupMapRuleMap["name"])
	d.Set("name_alias", aaaLdapGroupMapRuleMap["nameAlias"])
	return d, nil
}

func resourceAciLDAPGroupMapRuleImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	aaaLdapGroupMapRule, err := getRemoteLDAPGroupMapRule(aciClient, dn)
	if err != nil {
		return nil, err
	}
	ldap_group_map_rule := strings.Split(dn, "/")

	var ldap_group_type string
	if ldap_group_map_rule[2] == "ldapext" {
		ldap_group_type = "ldap"
	} else {
		ldap_group_type = "duo"
	}
	schemaFilled, err := setLDAPGroupMapRuleAttributes(ldap_group_type, aaaLdapGroupMapRule, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciLDAPGroupMapRuleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] LDAPGroupMapRule: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	ldap_type := d.Get("type").(string)
	aaaLdapGroupMapRuleAttr := models.LDAPGroupMapRuleAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		aaaLdapGroupMapRuleAttr.Annotation = Annotation.(string)
	} else {
		aaaLdapGroupMapRuleAttr.Annotation = "{}"
	}

	if Groupdn, ok := d.GetOk("groupdn"); ok {
		aaaLdapGroupMapRuleAttr.Groupdn = Groupdn.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		aaaLdapGroupMapRuleAttr.Name = Name.(string)
	}
	var aaaLdapGroupMapRule *models.LDAPGroupMapRule
	if ldap_type == "duo" {
		aaaLdapGroupMapRule = models.NewLDAPGroupMapRule(fmt.Sprintf("userext/duoext/ldapgroupmaprule-%s", name), "uni", desc, nameAlias, aaaLdapGroupMapRuleAttr)
	} else {
		aaaLdapGroupMapRule = models.NewLDAPGroupMapRule(fmt.Sprintf("userext/ldapext/ldapgroupmaprule-%s", name), "uni", desc, nameAlias, aaaLdapGroupMapRuleAttr)
	}
	err := aciClient.Save(aaaLdapGroupMapRule)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(aaaLdapGroupMapRule.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciLDAPGroupMapRuleRead(ctx, d, m)
}

func resourceAciLDAPGroupMapRuleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] LDAPGroupMapRule: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	ldap_type := d.Get("type").(string)
	aaaLdapGroupMapRuleAttr := models.LDAPGroupMapRuleAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		aaaLdapGroupMapRuleAttr.Annotation = Annotation.(string)
	} else {
		aaaLdapGroupMapRuleAttr.Annotation = "{}"
	}

	if Groupdn, ok := d.GetOk("groupdn"); ok {
		aaaLdapGroupMapRuleAttr.Groupdn = Groupdn.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		aaaLdapGroupMapRuleAttr.Name = Name.(string)
	}
	var aaaLdapGroupMapRule *models.LDAPGroupMapRule
	if ldap_type == "duo" {
		aaaLdapGroupMapRule = models.NewLDAPGroupMapRule(fmt.Sprintf("userext/duoext/ldapgroupmaprule-%s", name), "uni", desc, nameAlias, aaaLdapGroupMapRuleAttr)
	} else {
		aaaLdapGroupMapRule = models.NewLDAPGroupMapRule(fmt.Sprintf("userext/ldapext/ldapgroupmaprule-%s", name), "uni", desc, nameAlias, aaaLdapGroupMapRuleAttr)
	}
	aaaLdapGroupMapRule.Status = "modified"
	err := aciClient.Save(aaaLdapGroupMapRule)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(aaaLdapGroupMapRule.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciLDAPGroupMapRuleRead(ctx, d, m)
}

func resourceAciLDAPGroupMapRuleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	aaaLdapGroupMapRule, err := getRemoteLDAPGroupMapRule(aciClient, dn)
	if err != nil {
		return errorForObjectNotFound(err, dn, d)
	}
	ldap_group_map_rule := strings.Split(dn, "/")

	var ldap_group_type string
	if ldap_group_map_rule[2] == "ldapext" {
		ldap_group_type = "ldap"
	} else {
		ldap_group_type = "duo"
	}
	_, err = setLDAPGroupMapRuleAttributes(ldap_group_type, aaaLdapGroupMapRule, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciLDAPGroupMapRuleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "aaaLdapGroupMapRule")
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}
