package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAciAccessSubPortBlock() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciAccessSubPortBlockCreate,
		Update: resourceAciAccessSubPortBlockUpdate,
		Read:   resourceAciAccessSubPortBlockRead,
		Delete: resourceAciAccessSubPortBlockDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciAccessSubPortBlockImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"access_port_selector_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"from_card": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"from_port": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"from_sub_port": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"to_card": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"to_port": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"to_sub_port": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}
func getRemoteAccessSubPortBlock(client *client.Client, dn string) (*models.AccessSubPortBlock, error) {
	infraSubPortBlkCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	infraSubPortBlk := models.AccessSubPortBlockFromContainer(infraSubPortBlkCont)

	if infraSubPortBlk.DistinguishedName == "" {
		return nil, fmt.Errorf("AccessSubPortBlock %s not found", infraSubPortBlk.DistinguishedName)
	}

	return infraSubPortBlk, nil
}

func setAccessSubPortBlockAttributes(infraSubPortBlk *models.AccessSubPortBlock, d *schema.ResourceData) *schema.ResourceData {
	dn := d.Id()
	d.SetId(infraSubPortBlk.DistinguishedName)
	d.Set("description", infraSubPortBlk.Description)
	// d.Set("access_port_selector_dn", GetParentDn(infraSubPortBlk.DistinguishedName))
	if dn != infraSubPortBlk.DistinguishedName {
		d.Set("access_port_selector_dn", "")
	}
	infraSubPortBlkMap, _ := infraSubPortBlk.ToMap()

	d.Set("name", infraSubPortBlkMap["name"])

	d.Set("annotation", infraSubPortBlkMap["annotation"])
	d.Set("from_card", infraSubPortBlkMap["fromCard"])
	d.Set("from_port", infraSubPortBlkMap["fromPort"])
	d.Set("from_sub_port", infraSubPortBlkMap["fromSubPort"])
	d.Set("name_alias", infraSubPortBlkMap["nameAlias"])
	d.Set("to_card", infraSubPortBlkMap["toCard"])
	d.Set("to_port", infraSubPortBlkMap["toPort"])
	d.Set("to_sub_port", infraSubPortBlkMap["toSubPort"])
	return d
}

func resourceAciAccessSubPortBlockImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	infraSubPortBlk, err := getRemoteAccessSubPortBlock(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setAccessSubPortBlockAttributes(infraSubPortBlk, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciAccessSubPortBlockCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] AccessSubPortBlock: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	AccessPortSelectorDn := d.Get("access_port_selector_dn").(string)

	infraSubPortBlkAttr := models.AccessSubPortBlockAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		infraSubPortBlkAttr.Annotation = Annotation.(string)
	} else {
		infraSubPortBlkAttr.Annotation = "{}"
	}
	if FromCard, ok := d.GetOk("from_card"); ok {
		infraSubPortBlkAttr.FromCard = FromCard.(string)
	}
	if FromPort, ok := d.GetOk("from_port"); ok {
		infraSubPortBlkAttr.FromPort = FromPort.(string)
	}
	if FromSubPort, ok := d.GetOk("from_sub_port"); ok {
		infraSubPortBlkAttr.FromSubPort = FromSubPort.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		infraSubPortBlkAttr.NameAlias = NameAlias.(string)
	}
	if ToCard, ok := d.GetOk("to_card"); ok {
		infraSubPortBlkAttr.ToCard = ToCard.(string)
	}
	if ToPort, ok := d.GetOk("to_port"); ok {
		infraSubPortBlkAttr.ToPort = ToPort.(string)
	}
	if ToSubPort, ok := d.GetOk("to_sub_port"); ok {
		infraSubPortBlkAttr.ToSubPort = ToSubPort.(string)
	}
	infraSubPortBlk := models.NewAccessSubPortBlock(fmt.Sprintf("subportblk-%s", name), AccessPortSelectorDn, desc, infraSubPortBlkAttr)

	err := aciClient.Save(infraSubPortBlk)
	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	d.SetId(infraSubPortBlk.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciAccessSubPortBlockRead(d, m)
}

func resourceAciAccessSubPortBlockUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] AccessSubPortBlock: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	AccessPortSelectorDn := d.Get("access_port_selector_dn").(string)

	infraSubPortBlkAttr := models.AccessSubPortBlockAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		infraSubPortBlkAttr.Annotation = Annotation.(string)
	} else {
		infraSubPortBlkAttr.Annotation = "{}"
	}
	if FromCard, ok := d.GetOk("from_card"); ok {
		infraSubPortBlkAttr.FromCard = FromCard.(string)
	}
	if FromPort, ok := d.GetOk("from_port"); ok {
		infraSubPortBlkAttr.FromPort = FromPort.(string)
	}
	if FromSubPort, ok := d.GetOk("from_sub_port"); ok {
		infraSubPortBlkAttr.FromSubPort = FromSubPort.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		infraSubPortBlkAttr.NameAlias = NameAlias.(string)
	}
	if ToCard, ok := d.GetOk("to_card"); ok {
		infraSubPortBlkAttr.ToCard = ToCard.(string)
	}
	if ToPort, ok := d.GetOk("to_port"); ok {
		infraSubPortBlkAttr.ToPort = ToPort.(string)
	}
	if ToSubPort, ok := d.GetOk("to_sub_port"); ok {
		infraSubPortBlkAttr.ToSubPort = ToSubPort.(string)
	}
	infraSubPortBlk := models.NewAccessSubPortBlock(fmt.Sprintf("subportblk-%s", name), AccessPortSelectorDn, desc, infraSubPortBlkAttr)

	infraSubPortBlk.Status = "modified"

	err := aciClient.Save(infraSubPortBlk)

	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	d.SetId(infraSubPortBlk.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciAccessSubPortBlockRead(d, m)

}

func resourceAciAccessSubPortBlockRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	infraSubPortBlk, err := getRemoteAccessSubPortBlock(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setAccessSubPortBlockAttributes(infraSubPortBlk, d)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciAccessSubPortBlockDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "infraSubPortBlk")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}
