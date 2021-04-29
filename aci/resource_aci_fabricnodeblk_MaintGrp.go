package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAciNodeBlockMG() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciNodeBlockCreateMG,
		Update: resourceAciNodeBlockUpdateMG,
		Read:   resourceAciNodeBlockReadMG,
		Delete: resourceAciNodeBlockDeleteMG,

		Importer: &schema.ResourceImporter{
			State: resourceAciNodeBlockImportMG,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"pod_maintenance_group_dn": &schema.Schema{
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
func getRemoteNodeBlockMG(client *client.Client, dn string) (*models.NodeBlockMG, error) {
	fabricNodeBlkCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	fabricNodeBlk := models.NodeBlockFromContainerMG(fabricNodeBlkCont)

	if fabricNodeBlk.DistinguishedName == "" {
		return nil, fmt.Errorf("NodeBlock %s not found", fabricNodeBlk.DistinguishedName)
	}

	return fabricNodeBlk, nil
}

func setNodeBlockAttributesMG(fabricNodeBlk *models.NodeBlockMG, d *schema.ResourceData) *schema.ResourceData {
	dn := d.Id()
	d.SetId(fabricNodeBlk.DistinguishedName)
	d.Set("description", fabricNodeBlk.Description)
	// d.Set("pod_maintenance_group_dn", GetParentDn(fabricNodeBlk.DistinguishedName))
	if dn != fabricNodeBlk.DistinguishedName {
		d.Set("pod_maintenance_group_dn", "")
	}
	fabricNodeBlkMap, _ := fabricNodeBlk.ToMap()

	d.Set("name", fabricNodeBlkMap["name"])

	d.Set("annotation", fabricNodeBlkMap["annotation"])
	d.Set("from_", fabricNodeBlkMap["from_"])
	d.Set("name_alias", fabricNodeBlkMap["nameAlias"])
	d.Set("to_", fabricNodeBlkMap["to_"])
	return d
}

func resourceAciNodeBlockImportMG(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	fabricNodeBlk, err := getRemoteNodeBlockMG(aciClient, dn)

	if err != nil {
		return nil, err
	}
	fabricNodeBlkMap, _ := fabricNodeBlk.ToMap()

	name := fabricNodeBlkMap["name"]
	pDN := GetParentDn(dn, fmt.Sprintf("/nodeblk-%s", name))
	d.Set("pod_maintenance_group_dn", pDN)
	schemaFilled := setNodeBlockAttributesMG(fabricNodeBlk, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciNodeBlockCreateMG(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] NodeBlock: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	PODMaintenanceGroupDn := d.Get("pod_maintenance_group_dn").(string)

	fabricNodeBlkAttr := models.NodeBlockAttributesMG{}
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
	fabricNodeBlk := models.NewNodeBlockMG(fmt.Sprintf("nodeblk-%s", name), PODMaintenanceGroupDn, desc, fabricNodeBlkAttr)

	err := aciClient.Save(fabricNodeBlk)
	if err != nil {
		return err
	}
	d.Partial(true)

	d.Partial(false)

	d.SetId(fabricNodeBlk.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciNodeBlockReadMG(d, m)
}

func resourceAciNodeBlockUpdateMG(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] NodeBlock: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	PODMaintenanceGroupDn := d.Get("pod_maintenance_group_dn").(string)

	fabricNodeBlkAttr := models.NodeBlockAttributesMG{}
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
	fabricNodeBlk := models.NewNodeBlockMG(fmt.Sprintf("nodeblk-%s", name), PODMaintenanceGroupDn, desc, fabricNodeBlkAttr)

	fabricNodeBlk.Status = "modified"

	err := aciClient.Save(fabricNodeBlk)

	if err != nil {
		return err
	}
	d.Partial(true)

	d.Partial(false)

	d.SetId(fabricNodeBlk.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciNodeBlockReadMG(d, m)

}

func resourceAciNodeBlockReadMG(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	fabricNodeBlk, err := getRemoteNodeBlockMG(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setNodeBlockAttributesMG(fabricNodeBlk, d)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciNodeBlockDeleteMG(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "fabricNodeBlk")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}
