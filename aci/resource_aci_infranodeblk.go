package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAciNodeBlock() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciNodeBlockCreate,
		Update: resourceAciNodeBlockUpdate,
		Read:   resourceAciNodeBlockRead,
		Delete: resourceAciNodeBlockDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciNodeBlockImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"switch_association_dn": &schema.Schema{
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
func getRemoteNodeBlock(client *client.Client, dn string) (*models.NodeBlock, error) {
	infraNodeBlkCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	infraNodeBlk := models.NodeBlockFromContainerBLK(infraNodeBlkCont)

	if infraNodeBlk.DistinguishedName == "" {
		return nil, fmt.Errorf("NodeBlock %s not found", infraNodeBlk.DistinguishedName)
	}

	return infraNodeBlk, nil
}

func setNodeBlockAttributes(infraNodeBlk *models.NodeBlock, d *schema.ResourceData) *schema.ResourceData {
	dn := d.Id()
	d.SetId(infraNodeBlk.DistinguishedName)
	d.Set("description", infraNodeBlk.Description)
	// d.Set("switch_association_dn", GetParentDn(infraNodeBlk.DistinguishedName))
	if dn != infraNodeBlk.DistinguishedName {
		d.Set("switch_association_dn", "")
	}
	infraNodeBlkMap, _ := infraNodeBlk.ToMap()

	d.Set("name", infraNodeBlkMap["name"])

	d.Set("annotation", infraNodeBlkMap["annotation"])
	d.Set("from_", infraNodeBlkMap["from_"])
	d.Set("name_alias", infraNodeBlkMap["nameAlias"])
	d.Set("to_", infraNodeBlkMap["to_"])
	return d
}

func resourceAciNodeBlockImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	infraNodeBlk, err := getRemoteNodeBlock(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setNodeBlockAttributes(infraNodeBlk, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciNodeBlockCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] NodeBlock: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	SwitchAssociationDn := d.Get("switch_association_dn").(string)

	infraNodeBlkAttr := models.NodeBlockAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		infraNodeBlkAttr.Annotation = Annotation.(string)
	} else {
		infraNodeBlkAttr.Annotation = "{}"
	}
	if From_, ok := d.GetOk("from_"); ok {
		infraNodeBlkAttr.From_ = From_.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		infraNodeBlkAttr.NameAlias = NameAlias.(string)
	}
	if To_, ok := d.GetOk("to_"); ok {
		infraNodeBlkAttr.To_ = To_.(string)
	}
	infraNodeBlk := models.NewNodeBlock(fmt.Sprintf("nodeblk-%s", name), SwitchAssociationDn, desc, infraNodeBlkAttr)

	err := aciClient.Save(infraNodeBlk)
	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	d.SetId(infraNodeBlk.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciNodeBlockRead(d, m)
}

func resourceAciNodeBlockUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] NodeBlock: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	SwitchAssociationDn := d.Get("switch_association_dn").(string)

	infraNodeBlkAttr := models.NodeBlockAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		infraNodeBlkAttr.Annotation = Annotation.(string)
	} else {
		infraNodeBlkAttr.Annotation = "{}"
	}
	if From_, ok := d.GetOk("from_"); ok {
		infraNodeBlkAttr.From_ = From_.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		infraNodeBlkAttr.NameAlias = NameAlias.(string)
	}
	if To_, ok := d.GetOk("to_"); ok {
		infraNodeBlkAttr.To_ = To_.(string)
	}
	infraNodeBlk := models.NewNodeBlock(fmt.Sprintf("nodeblk-%s", name), SwitchAssociationDn, desc, infraNodeBlkAttr)

	infraNodeBlk.Status = "modified"

	err := aciClient.Save(infraNodeBlk)

	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	d.SetId(infraNodeBlk.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciNodeBlockRead(d, m)

}

func resourceAciNodeBlockRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	infraNodeBlk, err := getRemoteNodeBlock(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setNodeBlockAttributes(infraNodeBlk, d)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciNodeBlockDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "infraNodeBlk")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}
