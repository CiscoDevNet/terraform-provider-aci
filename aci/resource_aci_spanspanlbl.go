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

func resourceAciSPANSourcedestinationGroupMatchLabel() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciSPANSourcedestinationGroupMatchLabelCreate,
		UpdateContext: resourceAciSPANSourcedestinationGroupMatchLabelUpdate,
		ReadContext:   resourceAciSPANSourcedestinationGroupMatchLabelRead,
		DeleteContext: resourceAciSPANSourcedestinationGroupMatchLabelDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciSPANSourcedestinationGroupMatchLabelImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"span_source_group_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"tag": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}
func getRemoteSPANSourcedestinationGroupMatchLabel(client *client.Client, dn string) (*models.SPANSourcedestinationGroupMatchLabel, error) {
	spanSpanLblCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	spanSpanLbl := models.SPANSourcedestinationGroupMatchLabelFromContainer(spanSpanLblCont)

	if spanSpanLbl.DistinguishedName == "" {
		return nil, fmt.Errorf("SPANSourcedestinationGroupMatchLabel %s not found", spanSpanLbl.DistinguishedName)
	}

	return spanSpanLbl, nil
}

func setSPANSourcedestinationGroupMatchLabelAttributes(spanSpanLbl *models.SPANSourcedestinationGroupMatchLabel, d *schema.ResourceData) *schema.ResourceData {
	dn := d.Id()
	d.SetId(spanSpanLbl.DistinguishedName)
	d.Set("description", spanSpanLbl.Description)
	if dn != spanSpanLbl.DistinguishedName {
		d.Set("span_source_group_dn", "")
	}
	spanSpanLblMap, _ := spanSpanLbl.ToMap()

	d.Set("name", spanSpanLblMap["name"])
	d.Set("annotation", spanSpanLblMap["annotation"])
	d.Set("name_alias", spanSpanLblMap["nameAlias"])
	d.Set("tag", spanSpanLblMap["tag"])
	return d
}

func resourceAciSPANSourcedestinationGroupMatchLabelImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	spanSpanLbl, err := getRemoteSPANSourcedestinationGroupMatchLabel(aciClient, dn)

	if err != nil {
		return nil, err
	}
	spanSpanLblMap, _ := spanSpanLbl.ToMap()
	name := spanSpanLblMap["name"]
	pDN := GetParentDn(dn, fmt.Sprintf("/spanlbl-%s", name))
	d.Set("span_source_group_dn", pDN)
	schemaFilled := setSPANSourcedestinationGroupMatchLabelAttributes(spanSpanLbl, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciSPANSourcedestinationGroupMatchLabelCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] SPANSourcedestinationGroupMatchLabel: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	SPANSourceGroupDn := d.Get("span_source_group_dn").(string)

	spanSpanLblAttr := models.SPANSourcedestinationGroupMatchLabelAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		spanSpanLblAttr.Annotation = Annotation.(string)
	} else {
		spanSpanLblAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		spanSpanLblAttr.NameAlias = NameAlias.(string)
	}
	if Tag, ok := d.GetOk("tag"); ok {
		spanSpanLblAttr.Tag = Tag.(string)
	}
	spanSpanLbl := models.NewSPANSourcedestinationGroupMatchLabel(fmt.Sprintf("spanlbl-%s", name), SPANSourceGroupDn, desc, spanSpanLblAttr)

	err := aciClient.Save(spanSpanLbl)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(spanSpanLbl.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciSPANSourcedestinationGroupMatchLabelRead(ctx, d, m)
}

func resourceAciSPANSourcedestinationGroupMatchLabelUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] SPANSourcedestinationGroupMatchLabel: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	SPANSourceGroupDn := d.Get("span_source_group_dn").(string)

	spanSpanLblAttr := models.SPANSourcedestinationGroupMatchLabelAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		spanSpanLblAttr.Annotation = Annotation.(string)
	} else {
		spanSpanLblAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		spanSpanLblAttr.NameAlias = NameAlias.(string)
	}
	if Tag, ok := d.GetOk("tag"); ok {
		spanSpanLblAttr.Tag = Tag.(string)
	}
	spanSpanLbl := models.NewSPANSourcedestinationGroupMatchLabel(fmt.Sprintf("spanlbl-%s", name), SPANSourceGroupDn, desc, spanSpanLblAttr)

	spanSpanLbl.Status = "modified"

	err := aciClient.Save(spanSpanLbl)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(spanSpanLbl.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciSPANSourcedestinationGroupMatchLabelRead(ctx, d, m)

}

func resourceAciSPANSourcedestinationGroupMatchLabelRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	spanSpanLbl, err := getRemoteSPANSourcedestinationGroupMatchLabel(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setSPANSourcedestinationGroupMatchLabelAttributes(spanSpanLbl, d)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciSPANSourcedestinationGroupMatchLabelDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "spanSpanLbl")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}
