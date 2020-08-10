package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAciFirmwareGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciFirmwareGroupCreate,
		Update: resourceAciFirmwareGroupUpdate,
		Read:   resourceAciFirmwareGroupRead,
		Delete: resourceAciFirmwareGroupDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciFirmwareGroupImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{

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

			"firmware_group_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"relation_firmware_rs_fwgrpp": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
		}),
	}
}
func getRemoteFirmwareGroup(client *client.Client, dn string) (*models.FirmwareGroup, error) {
	firmwareFwGrpCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	firmwareFwGrp := models.FirmwareGroupFromContainer(firmwareFwGrpCont)

	if firmwareFwGrp.DistinguishedName == "" {
		return nil, fmt.Errorf("FirmwareGroup %s not found", firmwareFwGrp.DistinguishedName)
	}

	return firmwareFwGrp, nil
}

func setFirmwareGroupAttributes(firmwareFwGrp *models.FirmwareGroup, d *schema.ResourceData) *schema.ResourceData {
	d.SetId(firmwareFwGrp.DistinguishedName)
	d.Set("description", firmwareFwGrp.Description)
	firmwareFwGrpMap, _ := firmwareFwGrp.ToMap()

	d.Set("name", firmwareFwGrpMap["name"])

	d.Set("annotation", firmwareFwGrpMap["annotation"])
	d.Set("name_alias", firmwareFwGrpMap["nameAlias"])
	d.Set("firmware_group_type", firmwareFwGrpMap["type"])
	return d
}

func resourceAciFirmwareGroupImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	firmwareFwGrp, err := getRemoteFirmwareGroup(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setFirmwareGroupAttributes(firmwareFwGrp, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciFirmwareGroupCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] FirmwareGroup: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	firmwareFwGrpAttr := models.FirmwareGroupAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		firmwareFwGrpAttr.Annotation = Annotation.(string)
	} else {
		firmwareFwGrpAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		firmwareFwGrpAttr.NameAlias = NameAlias.(string)
	}
	if FirmwareGroup_type, ok := d.GetOk("firmware_group_type"); ok {
		firmwareFwGrpAttr.FirmwareGroup_type = FirmwareGroup_type.(string)
	}
	firmwareFwGrp := models.NewFirmwareGroup(fmt.Sprintf("fabric/fwgrp-%s", name), "uni", desc, firmwareFwGrpAttr)

	err := aciClient.Save(firmwareFwGrp)
	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	if relationTofirmwareRsFwgrpp, ok := d.GetOk("relation_firmware_rs_fwgrpp"); ok {
		relationParam := relationTofirmwareRsFwgrpp.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationfirmwareRsFwgrppFromFirmwareGroup(firmwareFwGrp.DistinguishedName, relationParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_firmware_rs_fwgrpp")
		d.Partial(false)

	}

	d.SetId(firmwareFwGrp.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciFirmwareGroupRead(d, m)
}

func resourceAciFirmwareGroupUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] FirmwareGroup: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	firmwareFwGrpAttr := models.FirmwareGroupAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		firmwareFwGrpAttr.Annotation = Annotation.(string)
	} else {
		firmwareFwGrpAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		firmwareFwGrpAttr.NameAlias = NameAlias.(string)
	}
	if FirmwareGroup_type, ok := d.GetOk("firmware_group_type"); ok {
		firmwareFwGrpAttr.FirmwareGroup_type = FirmwareGroup_type.(string)
	}
	firmwareFwGrp := models.NewFirmwareGroup(fmt.Sprintf("fabric/fwgrp-%s", name), "uni", desc, firmwareFwGrpAttr)

	firmwareFwGrp.Status = "modified"

	err := aciClient.Save(firmwareFwGrp)

	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	if d.HasChange("relation_firmware_rs_fwgrpp") {
		_, newRelParam := d.GetChange("relation_firmware_rs_fwgrpp")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationfirmwareRsFwgrppFromFirmwareGroup(firmwareFwGrp.DistinguishedName, newRelParamName)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_firmware_rs_fwgrpp")
		d.Partial(false)

	}

	d.SetId(firmwareFwGrp.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciFirmwareGroupRead(d, m)

}

func resourceAciFirmwareGroupRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	firmwareFwGrp, err := getRemoteFirmwareGroup(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setFirmwareGroupAttributes(firmwareFwGrp, d)

	firmwareRsFwgrppData, err := aciClient.ReadRelationfirmwareRsFwgrppFromFirmwareGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation firmwareRsFwgrpp %v", err)

	} else {
		if _, ok := d.GetOk("relation_firmware_rs_fwgrpp"); ok {
			tfName := GetMOName(d.Get("relation_firmware_rs_fwgrpp").(string))
			if tfName != firmwareRsFwgrppData {
				d.Set("relation_firmware_rs_fwgrpp", "")
			}
		}
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciFirmwareGroupDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "firmwareFwGrp")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}
