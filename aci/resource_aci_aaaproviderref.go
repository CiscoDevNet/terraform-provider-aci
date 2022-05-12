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

func resourceAciProviderGroupMember() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciProviderGroupMemberCreate,
		UpdateContext: resourceAciProviderGroupMemberUpdate,
		ReadContext:   resourceAciProviderGroupMemberRead,
		DeleteContext: resourceAciProviderGroupMemberDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciProviderGroupMemberImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"parent_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"order": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.Any(validateIntRange(0, 16), validation.StringInSlice([]string{"lowest-available"}, false)),
			},
		})),
	}
}

func getRemoteProviderGroupMember(client *client.Client, dn string) (*models.ProviderGroupMember, error) {
	aaaProviderRefCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	aaaProviderRef := models.ProviderGroupMemberFromContainer(aaaProviderRefCont)
	if aaaProviderRef.DistinguishedName == "" {
		return nil, fmt.Errorf("ProviderGroupMember %s not found", aaaProviderRef.DistinguishedName)
	}
	return aaaProviderRef, nil
}

func setProviderGroupMemberAttributes(aaaProviderRef *models.ProviderGroupMember, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(aaaProviderRef.DistinguishedName)
	d.Set("description", aaaProviderRef.Description)
	aaaProviderRefMap, err := aaaProviderRef.ToMap()
	if err != nil {
		return nil, err
	}
	d.Set("annotation", aaaProviderRefMap["annotation"])
	d.Set("name", aaaProviderRefMap["name"])
	if Order, ok := d.GetOk("order"); ok {
		if Order.(string) == "lowest-available" {
			d.Set("order", "lowest-available")
		} else {
			d.Set("order", aaaProviderRefMap["order"])
		}
	} else {
		d.Set("order", aaaProviderRefMap["order"])
	}
	d.Set("parent_dn", GetParentDn(d.Id(), fmt.Sprintf("/providerref-%s", aaaProviderRefMap["name"])))
	d.Set("name_alias", aaaProviderRefMap["nameAlias"])
	return d, nil
}

func resourceAciProviderGroupMemberImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	aaaProviderRef, err := getRemoteProviderGroupMember(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setProviderGroupMemberAttributes(aaaProviderRef, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciProviderGroupMemberCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] ProviderGroupMember: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	ParentDn := d.Get("parent_dn").(string)

	aaaProviderRefAttr := models.ProviderGroupMemberAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		aaaProviderRefAttr.Annotation = Annotation.(string)
	} else {
		aaaProviderRefAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		aaaProviderRefAttr.Name = Name.(string)
	}

	if Order, ok := d.GetOk("order"); ok {
		aaaProviderRefAttr.Order = Order.(string)
	}
	aaaProviderRef := models.NewProviderGroupMember(fmt.Sprintf("providerref-%s", name), ParentDn, desc, nameAlias, aaaProviderRefAttr)

	err := aciClient.Save(aaaProviderRef)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(aaaProviderRef.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciProviderGroupMemberRead(ctx, d, m)
}

func resourceAciProviderGroupMemberUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] ProviderGroupMember: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	ParentDn := d.Get("parent_dn").(string)
	aaaProviderRefAttr := models.ProviderGroupMemberAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		aaaProviderRefAttr.Annotation = Annotation.(string)
	} else {
		aaaProviderRefAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		aaaProviderRefAttr.Name = Name.(string)
	}

	if Order, ok := d.GetOk("order"); ok {
		aaaProviderRefAttr.Order = Order.(string)
	}
	aaaProviderRef := models.NewProviderGroupMember(fmt.Sprintf("providerref-%s", name), ParentDn, desc, nameAlias, aaaProviderRefAttr)

	aaaProviderRef.Status = "modified"
	err := aciClient.Save(aaaProviderRef)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(aaaProviderRef.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciProviderGroupMemberRead(ctx, d, m)
}

func resourceAciProviderGroupMemberRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	aaaProviderRef, err := getRemoteProviderGroupMember(aciClient, dn)
	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setProviderGroupMemberAttributes(aaaProviderRef, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciProviderGroupMemberDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "aaaProviderRef")
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}
