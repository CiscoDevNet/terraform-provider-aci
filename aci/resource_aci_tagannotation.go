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

func resourceAciAnnotation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciAnnotationCreate,
		UpdateContext: resourceAciAnnotationUpdate,
		ReadContext:   resourceAciAnnotationRead,
		DeleteContext: resourceAciAnnotationDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciAnnotationImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"parent_dn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"value": {
				Type:     schema.TypeString,
				Required: true,
			},
		}),
	}
}

func getRemoteAnnotation(client *client.Client, dn string) (*models.Annotation, error) {
	tagAnnotationCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	tagAnnotation := models.AnnotationFromContainer(tagAnnotationCont)
	if tagAnnotation.DistinguishedName == "" {
		return nil, fmt.Errorf("Annotation %s not found", tagAnnotation.DistinguishedName)
	}
	return tagAnnotation, nil
}

func setAnnotationAttributes(tagAnnotation *models.Annotation, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(tagAnnotation.DistinguishedName)
	tagAnnotationMap, err := tagAnnotation.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("key", tagAnnotationMap["key"])
	d.Set("value", tagAnnotationMap["value"])
	return d, nil
}

func resourceAciAnnotationImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	tagAnnotation, err := getRemoteAnnotation(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setAnnotationAttributes(tagAnnotation, d)
	if err != nil {
		return nil, err
	}
	d.Set("parent_dn", GetParentDn(dn, fmt.Sprintf("/"+models.RnTagAnnotation, d.Get("key"))))
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciAnnotationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Annotation: Beginning Creation")
	aciClient := m.(*client.Client)
	key := d.Get("key").(string)
	parentDn := d.Get("parent_dn").(string)
	tagAnnotationAttr := models.AnnotationAttributes{}

	if Key, ok := d.GetOk("key"); ok {
		tagAnnotationAttr.Key = Key.(string)
	}

	if Value, ok := d.GetOk("value"); ok {
		tagAnnotationAttr.Value = Value.(string)
	}

	tagAnnotation := models.NewAnnotation(fmt.Sprintf(models.RnTagAnnotation, key), parentDn, tagAnnotationAttr)

	err := aciClient.Save(tagAnnotation)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(tagAnnotation.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciAnnotationRead(ctx, d, m)
}

func resourceAciAnnotationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Annotation: Beginning Update")
	aciClient := m.(*client.Client)
	key := d.Get("key").(string)
	parentDn := d.Get("parent_dn").(string)
	tagAnnotationAttr := models.AnnotationAttributes{}

	if Key, ok := d.GetOk("key"); ok {
		tagAnnotationAttr.Key = Key.(string)
	}

	if Value, ok := d.GetOk("value"); ok {
		tagAnnotationAttr.Value = Value.(string)
	}

	tagAnnotation := models.NewAnnotation(fmt.Sprintf(models.RnTagAnnotation, key), parentDn, tagAnnotationAttr)
	tagAnnotation.Status = "modified"

	err := aciClient.Save(tagAnnotation)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(tagAnnotation.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciAnnotationRead(ctx, d, m)
}

func resourceAciAnnotationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	tagAnnotation, err := getRemoteAnnotation(aciClient, dn)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	_, err = setAnnotationAttributes(tagAnnotation, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciAnnotationDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "tagAnnotation")
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}
