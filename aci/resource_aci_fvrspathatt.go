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

func resourceAciStaticPath() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciStaticPathCreate,
		UpdateContext: resourceAciStaticPathUpdate,
		ReadContext:   resourceAciStaticPathRead,
		DeleteContext: resourceAciStaticPathDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciStaticPathImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"application_epg_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"tdn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"encap": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"instr_imedcy": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"immediate",
					"lazy",
				}, false),
			},

			"mode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"regular",
					"native",
					"untagged",
				}, false),
			},

			"primary_encap": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}
func getRemoteStaticPath(client *client.Client, dn string) (*models.StaticPath, error) {
	fvRsPathAttCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	fvRsPathAtt := models.StaticPathFromContainer(fvRsPathAttCont)

	if fvRsPathAtt.DistinguishedName == "" {
		return nil, fmt.Errorf("StaticPath %s not found", fvRsPathAtt.DistinguishedName)
	}

	return fvRsPathAtt, nil
}

func setStaticPathAttributes(fvRsPathAtt *models.StaticPath, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()
	d.SetId(fvRsPathAtt.DistinguishedName)
	fvRsPathAttMap, err := fvRsPathAtt.ToMap()

	if err != nil {
		return d, err
	}
	d.Set("application_epg_dn", GetParentDn(fvRsPathAtt.DistinguishedName, fmt.Sprintf("/rspathAtt-[%s]", fvRsPathAttMap["tDn"])))
	if dn != fvRsPathAtt.DistinguishedName {
		d.Set("application_epg_dn", "")
	}

	d.Set("tdn", fvRsPathAttMap["tDn"])

	d.Set("annotation", fvRsPathAttMap["annotation"])
	d.Set("encap", fvRsPathAttMap["encap"])
	d.Set("instr_imedcy", fvRsPathAttMap["instrImedcy"])
	d.Set("mode", fvRsPathAttMap["mode"])
	d.Set("primary_encap", fvRsPathAttMap["primaryEncap"])
	return d, nil
}

func resourceAciStaticPathImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	fvRsPathAtt, err := getRemoteStaticPath(aciClient, dn)

	if err != nil {
		return nil, err
	}
	fvRsPathAttMap, err := fvRsPathAtt.ToMap()

	if err != nil {
		return nil, err
	}

	tDn := fvRsPathAttMap["tDn"]
	pDN := GetParentDn(dn, fmt.Sprintf("/rspathAtt-[%s]", tDn))
	d.Set("application_epg_dn", pDN)
	schemaFilled, err := setStaticPathAttributes(fvRsPathAtt, d)

	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciStaticPathCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] StaticPath: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	tDn := d.Get("tdn").(string)

	ApplicationEPGDn := d.Get("application_epg_dn").(string)

	fvRsPathAttAttr := models.StaticPathAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		fvRsPathAttAttr.Annotation = Annotation.(string)
	} else {
		fvRsPathAttAttr.Annotation = "{}"
	}
	if Encap, ok := d.GetOk("encap"); ok {
		fvRsPathAttAttr.Encap = Encap.(string)
	}
	if InstrImedcy, ok := d.GetOk("instr_imedcy"); ok {
		fvRsPathAttAttr.InstrImedcy = InstrImedcy.(string)
	}
	if Mode, ok := d.GetOk("mode"); ok {
		fvRsPathAttAttr.Mode = Mode.(string)
	}
	if PrimaryEncap, ok := d.GetOk("primary_encap"); ok {
		fvRsPathAttAttr.PrimaryEncap = PrimaryEncap.(string)
	}

	fvRsPathAtt := models.NewStaticPath(fmt.Sprintf("rspathAtt-[%s]", tDn), ApplicationEPGDn, desc, fvRsPathAttAttr)

	err := aciClient.Save(fvRsPathAtt)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fvRsPathAtt.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciStaticPathRead(ctx, d, m)
}

func resourceAciStaticPathUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] StaticPath: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	tDn := d.Get("tdn").(string)

	ApplicationEPGDn := d.Get("application_epg_dn").(string)

	fvRsPathAttAttr := models.StaticPathAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		fvRsPathAttAttr.Annotation = Annotation.(string)
	} else {
		fvRsPathAttAttr.Annotation = "{}"
	}
	if Encap, ok := d.GetOk("encap"); ok {
		fvRsPathAttAttr.Encap = Encap.(string)
	}
	if InstrImedcy, ok := d.GetOk("instr_imedcy"); ok {
		fvRsPathAttAttr.InstrImedcy = InstrImedcy.(string)
	}
	if Mode, ok := d.GetOk("mode"); ok {
		fvRsPathAttAttr.Mode = Mode.(string)
	}
	if PrimaryEncap, ok := d.GetOk("primary_encap"); ok {
		fvRsPathAttAttr.PrimaryEncap = PrimaryEncap.(string)
	}

	fvRsPathAtt := models.NewStaticPath(fmt.Sprintf("rspathAtt-[%s]", tDn), ApplicationEPGDn, desc, fvRsPathAttAttr)

	fvRsPathAtt.Status = "modified"

	err := aciClient.Save(fvRsPathAtt)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fvRsPathAtt.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciStaticPathRead(ctx, d, m)

}

func resourceAciStaticPathRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	fvRsPathAtt, err := getRemoteStaticPath(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setStaticPathAttributes(fvRsPathAtt, d)

	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciStaticPathDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "fvRsPathAtt")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}
