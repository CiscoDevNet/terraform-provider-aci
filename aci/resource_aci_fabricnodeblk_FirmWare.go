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

func resourceAciNodeBlockFW() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciNodeBlockCreateFW,
		UpdateContext: resourceAciNodeBlockUpdateFW,
		ReadContext:   resourceAciNodeBlockReadFW,
		DeleteContext: resourceAciNodeBlockDeleteFW,

		Importer: &schema.ResourceImporter{
			State: resourceAciNodeBlockImportFW,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"firmware_group_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"from_": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"to_": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}
func getRemoteNodeBlockFW(client *client.Client, dn string) (*models.NodeBlockFW, error) {
	fabricNodeBlkCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	fabricNodeBlk := models.NodeBlockFromContainer(fabricNodeBlkCont)

	if fabricNodeBlk.DistinguishedName == "" {
		return nil, fmt.Errorf("NodeBlock %s not found", fabricNodeBlk.DistinguishedName)
	}

	return fabricNodeBlk, nil
}

func setNodeBlockAttributesFW(fabricNodeBlk *models.NodeBlockFW, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()
	d.SetId(fabricNodeBlk.DistinguishedName)
	d.Set("description", fabricNodeBlk.Description)
	if dn != fabricNodeBlk.DistinguishedName {
		d.Set("firmware_group_dn", "")
	}
	fabricNodeBlkMap, err := fabricNodeBlk.ToMap()

	if err != nil {
		return d, err
	}

	d.Set("name", fabricNodeBlkMap["name"])
	d.Set("firmware_group_dn", GetParentDn(fabricNodeBlk.DistinguishedName, fmt.Sprintf("/nodeblk-%s", fabricNodeBlkMap["name"])))

	d.Set("annotation", fabricNodeBlkMap["annotation"])
	d.Set("from_", fabricNodeBlkMap["from_"])
	d.Set("name_alias", fabricNodeBlkMap["nameAlias"])
	d.Set("to_", fabricNodeBlkMap["to_"])
	return d, nil
}

func resourceAciNodeBlockImportFW(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	fabricNodeBlk, err := getRemoteNodeBlockFW(aciClient, dn)

	if err != nil {
		return nil, err
	}
	fabricNodeBlkMap, err := fabricNodeBlk.ToMap()

	if err != nil {
		return nil, err
	}

	name := fabricNodeBlkMap["name"]
	pDN := GetParentDn(dn, fmt.Sprintf("/nodeblk-%s", name))
	d.Set("firmware_group_dn", pDN)
	schemaFilled, err := setNodeBlockAttributesFW(fabricNodeBlk, d)

	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciNodeBlockCreateFW(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] NodeBlock: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	FirmwareGroupDn := d.Get("firmware_group_dn").(string)

	fabricNodeBlkAttr := models.NodeBlockAttributesFW{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		fabricNodeBlkAttr.Annotation = Annotation.(string)
	} else {
		fabricNodeBlkAttr.Annotation = "{}"
	}
	if From_, ok := d.GetOk("from_"); ok {
		fabricNodeBlkAttr.From_ = From_.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		fabricNodeBlkAttr.NameAlias = NameAlias.(string)
	}
	if To_, ok := d.GetOk("to_"); ok {
		fabricNodeBlkAttr.To_ = To_.(string)
	}
	fabricNodeBlk := models.NewNodeBlockFW(fmt.Sprintf("nodeblk-%s", name), FirmwareGroupDn, desc, fabricNodeBlkAttr)

	err := aciClient.Save(fabricNodeBlk)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fabricNodeBlk.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciNodeBlockReadFW(ctx, d, m)
}

func resourceAciNodeBlockUpdateFW(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] NodeBlock: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	FirmwareGroupDn := d.Get("firmware_group_dn").(string)

	fabricNodeBlkAttr := models.NodeBlockAttributesFW{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		fabricNodeBlkAttr.Annotation = Annotation.(string)
	} else {
		fabricNodeBlkAttr.Annotation = "{}"
	}
	if From_, ok := d.GetOk("from_"); ok {
		fabricNodeBlkAttr.From_ = From_.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		fabricNodeBlkAttr.NameAlias = NameAlias.(string)
	}
	if To_, ok := d.GetOk("to_"); ok {
		fabricNodeBlkAttr.To_ = To_.(string)
	}
	fabricNodeBlk := models.NewNodeBlockFW(fmt.Sprintf("nodeblk-%s", name), FirmwareGroupDn, desc, fabricNodeBlkAttr)

	fabricNodeBlk.Status = "modified"

	err := aciClient.Save(fabricNodeBlk)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fabricNodeBlk.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciNodeBlockReadFW(ctx, d, m)

}

func resourceAciNodeBlockReadFW(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	fabricNodeBlk, err := getRemoteNodeBlockFW(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setNodeBlockAttributesFW(fabricNodeBlk, d)

	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciNodeBlockDeleteFW(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "fabricNodeBlk")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}
