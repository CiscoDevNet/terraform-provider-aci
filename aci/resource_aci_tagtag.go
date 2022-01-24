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

func resourceAciTag() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciTagCreate,
		UpdateContext: resourceAciTagUpdate,
		ReadContext:   resourceAciTagRead,
		DeleteContext: resourceAciTagDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciTagImport,
		},

		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
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
		},
	}
}

func getRemoteTag(client *client.Client, dn string) (*models.Tag, error) {
	tagTagCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	tagTag := models.TagFromContainer(tagTagCont)
	if tagTag.DistinguishedName == "" {
		return nil, fmt.Errorf("Tag %s not found", tagTag.DistinguishedName)
	}
	return tagTag, nil
}

func setTagAttributes(tagTag *models.Tag, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(tagTag.DistinguishedName)
	tagTagMap, err := tagTag.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("key", tagTagMap["key"])
	d.Set("value", tagTagMap["value"])
	return d, nil
}

func resourceAciTagImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	tagTag, err := getRemoteTag(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setTagAttributes(tagTag, d)
	if err != nil {
		return nil, err
	}
	d.Set("parent_dn", GetParentDn(dn, fmt.Sprintf("/"+models.RnTagTag, d.Get("key"))))
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciTagCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Tag: Beginning Creation")
	aciClient := m.(*client.Client)
	key := d.Get("key").(string)
	parentDn := d.Get("parent_dn").(string)
	tagTagAttr := models.TagAttributes{}

	if Key, ok := d.GetOk("key"); ok {
		tagTagAttr.Key = Key.(string)
	}

	if Value, ok := d.GetOk("value"); ok {
		tagTagAttr.Value = Value.(string)
	}

	tagTag := models.NewTag(fmt.Sprintf(models.RnTagTag, key), parentDn, tagTagAttr)

	err := aciClient.Save(tagTag)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(tagTag.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciTagRead(ctx, d, m)
}

func resourceAciTagUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Tag: Beginning Update")
	aciClient := m.(*client.Client)
	key := d.Get("key").(string)
	parentDn := d.Get("parent_dn").(string)
	tagTagAttr := models.TagAttributes{}

	if Key, ok := d.GetOk("key"); ok {
		tagTagAttr.Key = Key.(string)
	}

	if Value, ok := d.GetOk("value"); ok {
		tagTagAttr.Value = Value.(string)
	}

	tagTag := models.NewTag(fmt.Sprintf(models.RnTagTag, key), parentDn, tagTagAttr)

	tagTag.Status = "modified"
	err := aciClient.Save(tagTag)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(tagTag.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciTagRead(ctx, d, m)
}

func resourceAciTagRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	tagTag, err := getRemoteTag(aciClient, dn)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	_, err = setTagAttributes(tagTag, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciTagDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "tagTag")
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}
