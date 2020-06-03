package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAciAccessPortBlock() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciAccessPortBlockCreate,
		Update: resourceAciAccessPortBlockUpdate,
		Read:   resourceAciAccessPortBlockRead,
		Delete: resourceAciAccessPortBlockDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciAccessPortBlockImport,
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

			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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

			"relation_infra_rs_acc_bndl_subgrp": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
		}),
	}
}
func getRemoteAccessPortBlock(client *client.Client, dn string) (*models.AccessPortBlock, error) {
	infraPortBlkCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	infraPortBlk := models.AccessPortBlockFromContainer(infraPortBlkCont)

	if infraPortBlk.DistinguishedName == "" {
		return nil, fmt.Errorf("AccessPortBlock %s not found", infraPortBlk.DistinguishedName)
	}

	return infraPortBlk, nil
}

func setAccessPortBlockAttributes(infraPortBlk *models.AccessPortBlock, d *schema.ResourceData) *schema.ResourceData {
	dn := d.Id()
	d.SetId(infraPortBlk.DistinguishedName)
	d.Set("description", infraPortBlk.Description)
	// d.Set("access_port_selector_dn", GetParentDn(infraPortBlk.DistinguishedName))
	if dn != infraPortBlk.DistinguishedName {
		d.Set("access_port_selector_dn", "")
	}
	infraPortBlkMap, _ := infraPortBlk.ToMap()

	d.Set("name", infraPortBlkMap["name"])

	d.Set("annotation", infraPortBlkMap["annotation"])
	d.Set("from_card", infraPortBlkMap["fromCard"])
	d.Set("from_port", infraPortBlkMap["fromPort"])
	d.Set("name_alias", infraPortBlkMap["nameAlias"])
	d.Set("to_card", infraPortBlkMap["toCard"])
	d.Set("to_port", infraPortBlkMap["toPort"])
	return d
}

func resourceAciAccessPortBlockImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	infraPortBlk, err := getRemoteAccessPortBlock(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setAccessPortBlockAttributes(infraPortBlk, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciAccessPortBlockCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] AccessPortBlock: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	AccessPortSelectorDn := d.Get("access_port_selector_dn").(string)

	infraPortBlkAttr := models.AccessPortBlockAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		infraPortBlkAttr.Annotation = Annotation.(string)
	}
	if FromCard, ok := d.GetOk("from_card"); ok {
		infraPortBlkAttr.FromCard = FromCard.(string)
	}
	if FromPort, ok := d.GetOk("from_port"); ok {
		infraPortBlkAttr.FromPort = FromPort.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		infraPortBlkAttr.NameAlias = NameAlias.(string)
	}
	if ToCard, ok := d.GetOk("to_card"); ok {
		infraPortBlkAttr.ToCard = ToCard.(string)
	}
	if ToPort, ok := d.GetOk("to_port"); ok {
		infraPortBlkAttr.ToPort = ToPort.(string)
	}
	infraPortBlk := models.NewAccessPortBlock(fmt.Sprintf("portblk-%s", name), AccessPortSelectorDn, desc, infraPortBlkAttr)

	err := aciClient.Save(infraPortBlk)
	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	if relationToinfraRsAccBndlSubgrp, ok := d.GetOk("relation_infra_rs_acc_bndl_subgrp"); ok {
		relationParam := relationToinfraRsAccBndlSubgrp.(string)
		err = aciClient.CreateRelationinfraRsAccBndlSubgrpFromAccessPortBlock(infraPortBlk.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_acc_bndl_subgrp")
		d.Partial(false)

	}

	d.SetId(infraPortBlk.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciAccessPortBlockRead(d, m)
}

func resourceAciAccessPortBlockUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] AccessPortBlock: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	AccessPortSelectorDn := d.Get("access_port_selector_dn").(string)

	infraPortBlkAttr := models.AccessPortBlockAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		infraPortBlkAttr.Annotation = Annotation.(string)
	}
	if FromCard, ok := d.GetOk("from_card"); ok {
		infraPortBlkAttr.FromCard = FromCard.(string)
	}
	if FromPort, ok := d.GetOk("from_port"); ok {
		infraPortBlkAttr.FromPort = FromPort.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		infraPortBlkAttr.NameAlias = NameAlias.(string)
	}
	if ToCard, ok := d.GetOk("to_card"); ok {
		infraPortBlkAttr.ToCard = ToCard.(string)
	}
	if ToPort, ok := d.GetOk("to_port"); ok {
		infraPortBlkAttr.ToPort = ToPort.(string)
	}
	infraPortBlk := models.NewAccessPortBlock(fmt.Sprintf("portblk-%s", name), AccessPortSelectorDn, desc, infraPortBlkAttr)

	infraPortBlk.Status = "modified"

	err := aciClient.Save(infraPortBlk)

	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	if d.HasChange("relation_infra_rs_acc_bndl_subgrp") {
		_, newRelParam := d.GetChange("relation_infra_rs_acc_bndl_subgrp")
		err = aciClient.DeleteRelationinfraRsAccBndlSubgrpFromAccessPortBlock(infraPortBlk.DistinguishedName)
		if err != nil {
			return err
		}
		err = aciClient.CreateRelationinfraRsAccBndlSubgrpFromAccessPortBlock(infraPortBlk.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_infra_rs_acc_bndl_subgrp")
		d.Partial(false)

	}

	d.SetId(infraPortBlk.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciAccessPortBlockRead(d, m)

}

func resourceAciAccessPortBlockRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	infraPortBlk, err := getRemoteAccessPortBlock(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setAccessPortBlockAttributes(infraPortBlk, d)

	infraRsAccBndlSubgrpData, err := aciClient.ReadRelationinfraRsAccBndlSubgrpFromAccessPortBlock(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsAccBndlSubgrp %v", err)

	} else {
		d.Set("relation_infra_rs_acc_bndl_subgrp", infraRsAccBndlSubgrpData)
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciAccessPortBlockDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "infraPortBlk")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}
