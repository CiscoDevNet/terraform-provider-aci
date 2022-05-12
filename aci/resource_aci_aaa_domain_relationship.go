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

func resourceAciDomainRelationship() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciDomainRelationshipCreate,
		ReadContext:   resourceAciDomainRelationshipRead,
		DeleteContext: resourceAciDomainRelationshipDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciDomainRelationshipImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"parent_dn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"aaa_domain_dn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		})),
	}
}

func getRemoteAaaDomainRelationship(client *client.Client, dn string) (*models.AaaDomainRef, error) {
	aaaDomainRefCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	aaaDomainRef := models.AaaDomainRefFromContainer(aaaDomainRefCont)
	if aaaDomainRef.DistinguishedName == "" {
		return nil, fmt.Errorf("AaaDomainRef %s not found", aaaDomainRef.DistinguishedName)
	}
	return aaaDomainRef, nil
}

func setAaaDomainRelationshipAttributes(aaaDomainRef *models.AaaDomainRef, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()
	d.SetId(aaaDomainRef.DistinguishedName)
	if dn != aaaDomainRef.DistinguishedName {
		d.Set("parent_dn", "")
	}
	aaaDomainRefMap, err := aaaDomainRef.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("name", aaaDomainRef.Name)
	d.Set("name_alias", aaaDomainRefMap["nameAlias"])
	return d, nil
}

func resourceAciDomainRelationshipImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	aaaDomainRef, err := getRemoteAaaDomainRelationship(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setAaaDomainRelationshipAttributes(aaaDomainRef, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciDomainRelationshipCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] AaaDomainRef: Beginning Creation")
	aciClient := m.(*client.Client)
	name := GetMOName(d.Get("aaa_domain_dn").(string))
	parent_dn := d.Get("parent_dn").(string)

	aaaDomainRefAttr := models.AaaDomainRefAttributes{}
	aaaDomainRefAttr.Name = name

	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		aaaDomainRefAttr.Annotation = Annotation.(string)
	} else {
		aaaDomainRefAttr.Annotation = "{}"
	}

	aaaDomainRef := models.NewAaaDomainRef(fmt.Sprintf(models.RnaaaDomainRef, name), parent_dn, nameAlias, aaaDomainRefAttr)

	err := aciClient.Save(aaaDomainRef)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(aaaDomainRef.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciDomainRelationshipRead(ctx, d, m)
}

func resourceAciDomainRelationshipRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	aaaDomainRef, err := getRemoteAaaDomainRelationship(aciClient, dn)
	if err != nil {
		d.SetId("")
		return nil
	}

	_, err = setAaaDomainRelationshipAttributes(aaaDomainRef, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciDomainRelationshipDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	err := aciClient.DeleteByDn(dn, "aaaDomainRef")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}
