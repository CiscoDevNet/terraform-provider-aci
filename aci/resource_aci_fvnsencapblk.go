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

func resourceAciRanges() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciRangesCreate,
		UpdateContext: resourceAciRangesUpdate,
		ReadContext:   resourceAciRangesRead,
		DeleteContext: resourceAciRangesDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciRangesImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"vlan_pool_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"from": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"to": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"alloc_mode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"dynamic",
					"static",
					"inherit",
				}, false),
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"role": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"external",
					"internal",
				}, false),
			},
		}),
	}
}
func getRemoteRanges(client *client.Client, dn string) (*models.Ranges, error) {
	fvnsEncapBlkCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	fvnsEncapBlk := models.RangesFromContainer(fvnsEncapBlkCont)

	if fvnsEncapBlk.DistinguishedName == "" {
		return nil, fmt.Errorf("Ranges %s not found", fvnsEncapBlk.DistinguishedName)
	}

	return fvnsEncapBlk, nil
}

func setRangesAttributes(fvnsEncapBlk *models.Ranges, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()
	d.SetId(fvnsEncapBlk.DistinguishedName)
	d.Set("description", fvnsEncapBlk.Description)
	// d.Set("vlan_pool_dn", GetParentDn(fvnsEncapBlk.DistinguishedName))
	if dn != fvnsEncapBlk.DistinguishedName {
		d.Set("vlan_pool_dn", "")
	}

	fvnsEncapBlkMap, err := fvnsEncapBlk.ToMap()
	if err != nil {
		return d, err
	}

	d.Set("alloc_mode", fvnsEncapBlkMap["allocMode"])
	d.Set("annotation", fvnsEncapBlkMap["annotation"])
	d.Set("from", fvnsEncapBlkMap["from"])
	d.Set("name_alias", fvnsEncapBlkMap["nameAlias"])
	d.Set("role", fvnsEncapBlkMap["role"])
	d.Set("to", fvnsEncapBlkMap["to"])
	d.Set("vlan_pool_dn", GetParentDn(dn, fmt.Sprintf("from-[%s]-to-[%s]", fvnsEncapBlkMap["from"], fvnsEncapBlkMap["to"])))
	return d, nil
}

func resourceAciRangesImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	fvnsEncapBlk, err := getRemoteRanges(aciClient, dn)
	if err != nil {
		return nil, err
	}

	fvnsEncapBlkMap, err := fvnsEncapBlk.ToMap()
	if err != nil {
		return nil, err
	}

	from := fvnsEncapBlkMap["from"]
	to := fvnsEncapBlkMap["to"]
	pDN := GetParentDn(dn, fmt.Sprintf("/from-[%s]-to-[%s]", from, to))
	d.Set("vlan_pool_dn", pDN)

	schemaFilled, err := setRangesAttributes(fvnsEncapBlk, d)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciRangesCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Ranges: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	from := d.Get("from").(string)

	to := d.Get("to").(string)

	VLANPoolDn := d.Get("vlan_pool_dn").(string)

	fvnsEncapBlkAttr := models.RangesAttributes{}
	if AllocMode, ok := d.GetOk("alloc_mode"); ok {
		fvnsEncapBlkAttr.AllocMode = AllocMode.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		fvnsEncapBlkAttr.Annotation = Annotation.(string)
	} else {
		fvnsEncapBlkAttr.Annotation = "{}"
	}
	if From, ok := d.GetOk("from"); ok {
		fvnsEncapBlkAttr.From = From.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		fvnsEncapBlkAttr.NameAlias = NameAlias.(string)
	}
	if Role, ok := d.GetOk("role"); ok {
		fvnsEncapBlkAttr.Role = Role.(string)
	}
	if To, ok := d.GetOk("to"); ok {
		fvnsEncapBlkAttr.To = To.(string)
	}
	fvnsEncapBlk := models.NewRanges(fmt.Sprintf("from-[%s]-to-[%s]", from, to), VLANPoolDn, desc, fvnsEncapBlkAttr)

	err := aciClient.Save(fvnsEncapBlk)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fvnsEncapBlk.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciRangesRead(ctx, d, m)
}

func resourceAciRangesUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Ranges: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	from := d.Get("from").(string)

	to := d.Get("to").(string)

	VLANPoolDn := d.Get("vlan_pool_dn").(string)

	fvnsEncapBlkAttr := models.RangesAttributes{}
	if AllocMode, ok := d.GetOk("alloc_mode"); ok {
		fvnsEncapBlkAttr.AllocMode = AllocMode.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		fvnsEncapBlkAttr.Annotation = Annotation.(string)
	} else {
		fvnsEncapBlkAttr.Annotation = "{}"
	}
	if From, ok := d.GetOk("from"); ok {
		fvnsEncapBlkAttr.From = From.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		fvnsEncapBlkAttr.NameAlias = NameAlias.(string)
	}
	if Role, ok := d.GetOk("role"); ok {
		fvnsEncapBlkAttr.Role = Role.(string)
	}
	if To, ok := d.GetOk("to"); ok {
		fvnsEncapBlkAttr.To = To.(string)
	}
	fvnsEncapBlk := models.NewRanges(fmt.Sprintf("from-[%s]-to-[%s]", from, to), VLANPoolDn, desc, fvnsEncapBlkAttr)

	fvnsEncapBlk.Status = "modified"

	err := aciClient.Save(fvnsEncapBlk)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fvnsEncapBlk.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciRangesRead(ctx, d, m)
}

func resourceAciRangesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	fvnsEncapBlk, err := getRemoteRanges(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setRangesAttributes(fvnsEncapBlk, d)
	if err != nil {
		d.SetId("")
		return nil
	}
	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciRangesDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "fvnsEncapBlk")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}
