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
)

func resourceAciL4L7RedirectHealthGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciL4L7RedirectHealthGroupCreate,
		ReadContext:   resourceAciL4L7RedirectHealthGroupRead,
		DeleteContext: resourceAciL4L7RedirectHealthGroupDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciL4L7RedirectHealthGroupImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"tenant_dn": &schema.Schema{
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

func getRemoteL4L7RedirectHealthGroup(client *client.Client, dn string) (*models.L4L7RedirectHealthGroup, error) {
	vnsRedirectHealthGroupCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	vnsRedirectHealthGroup := models.L4L7RedirectHealthGroupFromContainer(vnsRedirectHealthGroupCont)
	if vnsRedirectHealthGroup.DistinguishedName == "" {
		return nil, fmt.Errorf("L4 L7Redirect Health Group %s not found", dn)
	}
	return vnsRedirectHealthGroup, nil
}

func setL4L7RedirectHealthGroupAttributes(vnsRedirectHealthGroup *models.L4L7RedirectHealthGroup, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(vnsRedirectHealthGroup.DistinguishedName)
	d.Set("description", vnsRedirectHealthGroup.Description)
	vnsRedirectHealthGroupMap, err := vnsRedirectHealthGroup.ToMap()
	if err != nil {
		return d, err
	}
	dn := d.Id()
	if dn != vnsRedirectHealthGroup.DistinguishedName {
		d.Set("tenant_dn", "")
	} else {
		d.Set("tenant_dn", strings.Split(GetParentDn(vnsRedirectHealthGroup.DistinguishedName, fmt.Sprintf("/"+models.RnvnsRedirectHealthGroup, vnsRedirectHealthGroupMap["name"])), "/svcCont")[0])
	}
	d.Set("annotation", vnsRedirectHealthGroupMap["annotation"])
	d.Set("name", vnsRedirectHealthGroupMap["name"])
	d.Set("name_alias", vnsRedirectHealthGroupMap["nameAlias"])
	return d, nil
}

func resourceAciL4L7RedirectHealthGroupImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	vnsRedirectHealthGroup, err := getRemoteL4L7RedirectHealthGroup(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setL4L7RedirectHealthGroupAttributes(vnsRedirectHealthGroup, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciL4L7RedirectHealthGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] L4L7RedirectHealthGroup: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	TenantDn := d.Get("tenant_dn").(string)

	vnsRedirectHealthGroupAttr := models.L4L7RedirectHealthGroupAttributes{}

	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		vnsRedirectHealthGroupAttr.Annotation = Annotation.(string)
	} else {
		vnsRedirectHealthGroupAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		vnsRedirectHealthGroupAttr.Name = Name.(string)
	}
	vnsRedirectHealthGroup := models.NewL4L7RedirectHealthGroup(fmt.Sprintf("svcCont/"+models.RnvnsRedirectHealthGroup, name), TenantDn, desc, nameAlias, vnsRedirectHealthGroupAttr)

	err := aciClient.Save(vnsRedirectHealthGroup)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(vnsRedirectHealthGroup.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciL4L7RedirectHealthGroupRead(ctx, d, m)
}

func resourceAciL4L7RedirectHealthGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	vnsRedirectHealthGroup, err := getRemoteL4L7RedirectHealthGroup(aciClient, dn)
	if err != nil {
		d.SetId("")
		return nil
	}

	_, err = setL4L7RedirectHealthGroupAttributes(vnsRedirectHealthGroup, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciL4L7RedirectHealthGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	err := aciClient.DeleteByDn(dn, "vnsRedirectHealthGroup")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}
