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

func resourceAciUserRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciUserRoleCreate,
		UpdateContext: resourceAciUserRoleUpdate,
		ReadContext:   resourceAciUserRoleRead,
		DeleteContext: resourceAciUserRoleDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciUserRoleImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"user_domain_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"priv_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"readPriv",
					"writePriv",
				}, false),
			},
		})),
	}
}

func getRemoteUserRole(client *client.Client, dn string) (*models.UserRole, error) {
	aaaUserRoleCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	aaaUserRole := models.UserRoleFromContainer(aaaUserRoleCont)
	if aaaUserRole.DistinguishedName == "" {
		return nil, fmt.Errorf("User Role %s not found", dn)
	}
	return aaaUserRole, nil
}

func setUserRoleAttributes(aaaUserRole *models.UserRole, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(aaaUserRole.DistinguishedName)
	d.Set("description", aaaUserRole.Description)

	aaaUserRoleMap, err := aaaUserRole.ToMap()
	if err != nil {
		return nil, err
	}
	d.Set("name", aaaUserRoleMap["name"])
	d.Set("priv_type", aaaUserRoleMap["privType"])
	d.Set("name_alias", aaaUserRoleMap["nameAlias"])
	d.Set("user_domain_dn", GetParentDn(d.Id(), fmt.Sprintf("/role-%s", aaaUserRoleMap["name"])))
	return d, nil
}

func resourceAciUserRoleImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	aaaUserRole, err := getRemoteUserRole(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setUserRoleAttributes(aaaUserRole, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciUserRoleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] UserRole: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	UserDomainDn := d.Get("user_domain_dn").(string)

	aaaUserRoleAttr := models.UserRoleAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		aaaUserRoleAttr.Annotation = Annotation.(string)
	} else {
		aaaUserRoleAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		aaaUserRoleAttr.Name = Name.(string)
	}

	if PrivType, ok := d.GetOk("priv_type"); ok {
		aaaUserRoleAttr.PrivType = PrivType.(string)
	}
	aaaUserRole := models.NewUserRole(fmt.Sprintf("role-%s", name), UserDomainDn, desc, nameAlias, aaaUserRoleAttr)

	err := aciClient.Save(aaaUserRole)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(aaaUserRole.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciUserRoleRead(ctx, d, m)
}

func resourceAciUserRoleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] UserRole: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	UserDomainDn := d.Get("user_domain_dn").(string)
	aaaUserRoleAttr := models.UserRoleAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		aaaUserRoleAttr.Annotation = Annotation.(string)
	} else {
		aaaUserRoleAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		aaaUserRoleAttr.Name = Name.(string)
	}

	if PrivType, ok := d.GetOk("priv_type"); ok {
		aaaUserRoleAttr.PrivType = PrivType.(string)
	}
	aaaUserRole := models.NewUserRole(fmt.Sprintf("role-%s", name), UserDomainDn, desc, nameAlias, aaaUserRoleAttr)

	aaaUserRole.Status = "modified"
	err := aciClient.Save(aaaUserRole)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(aaaUserRole.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciUserRoleRead(ctx, d, m)
}

func resourceAciUserRoleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	aaaUserRole, err := getRemoteUserRole(aciClient, dn)
	if err != nil {
		return errorForObjectNotFound(err, dn, d)
	}
	_, err = setUserRoleAttributes(aaaUserRole, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciUserRoleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "aaaUserRole")
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}
