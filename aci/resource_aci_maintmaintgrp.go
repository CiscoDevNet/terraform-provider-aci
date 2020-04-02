package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAciPODMaintenanceGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciPODMaintenanceGroupCreate,
		Update: resourceAciPODMaintenanceGroupUpdate,
		Read:   resourceAciPODMaintenanceGroupRead,
		Delete: resourceAciPODMaintenanceGroupDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciPODMaintenanceGroupImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{

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

			"fwtype": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"pod_maintenance_group_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"relation_maint_rs_mgrpp": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
		}),
	}
}
func getRemotePODMaintenanceGroup(client *client.Client, dn string) (*models.PODMaintenanceGroup, error) {
	maintMaintGrpCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	maintMaintGrp := models.PODMaintenanceGroupFromContainer(maintMaintGrpCont)

	if maintMaintGrp.DistinguishedName == "" {
		return nil, fmt.Errorf("PODMaintenanceGroup %s not found", maintMaintGrp.DistinguishedName)
	}

	return maintMaintGrp, nil
}

func setPODMaintenanceGroupAttributes(maintMaintGrp *models.PODMaintenanceGroup, d *schema.ResourceData) *schema.ResourceData {
	d.SetId(maintMaintGrp.DistinguishedName)
	d.Set("description", maintMaintGrp.Description)
	maintMaintGrpMap, _ := maintMaintGrp.ToMap()

	d.Set("name", maintMaintGrpMap["name"])

	d.Set("annotation", maintMaintGrpMap["annotation"])
	d.Set("fwtype", maintMaintGrpMap["fwtype"])
	d.Set("name_alias", maintMaintGrpMap["nameAlias"])
	d.Set("pod_maintenance_group_type", maintMaintGrpMap["type"])
	return d
}

func resourceAciPODMaintenanceGroupImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	maintMaintGrp, err := getRemotePODMaintenanceGroup(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setPODMaintenanceGroupAttributes(maintMaintGrp, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciPODMaintenanceGroupCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] PODMaintenanceGroup: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	maintMaintGrpAttr := models.PODMaintenanceGroupAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		maintMaintGrpAttr.Annotation = Annotation.(string)
	}
	if Fwtype, ok := d.GetOk("fwtype"); ok {
		maintMaintGrpAttr.Fwtype = Fwtype.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		maintMaintGrpAttr.NameAlias = NameAlias.(string)
	}
	if PODMaintenanceGroup_type, ok := d.GetOk("pod_maintenance_group_type"); ok {
		maintMaintGrpAttr.PODMaintenanceGroup_type = PODMaintenanceGroup_type.(string)
	}
	maintMaintGrp := models.NewPODMaintenanceGroup(fmt.Sprintf("fabric/maintgrp-%s", name), "uni", desc, maintMaintGrpAttr)

	err := aciClient.Save(maintMaintGrp)
	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	if relationTomaintRsMgrpp, ok := d.GetOk("relation_maint_rs_mgrpp"); ok {
		relationParam := relationTomaintRsMgrpp.(string)
		err = aciClient.CreateRelationmaintRsMgrppFromPODMaintenanceGroup(maintMaintGrp.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_maint_rs_mgrpp")
		d.Partial(false)

	}

	d.SetId(maintMaintGrp.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciPODMaintenanceGroupRead(d, m)
}

func resourceAciPODMaintenanceGroupUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] PODMaintenanceGroup: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	maintMaintGrpAttr := models.PODMaintenanceGroupAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		maintMaintGrpAttr.Annotation = Annotation.(string)
	}
	if Fwtype, ok := d.GetOk("fwtype"); ok {
		maintMaintGrpAttr.Fwtype = Fwtype.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		maintMaintGrpAttr.NameAlias = NameAlias.(string)
	}
	if PODMaintenanceGroup_type, ok := d.GetOk("pod_maintenance_group_type"); ok {
		maintMaintGrpAttr.PODMaintenanceGroup_type = PODMaintenanceGroup_type.(string)
	}
	maintMaintGrp := models.NewPODMaintenanceGroup(fmt.Sprintf("fabric/maintgrp-%s", name), "uni", desc, maintMaintGrpAttr)

	maintMaintGrp.Status = "modified"

	err := aciClient.Save(maintMaintGrp)

	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	if d.HasChange("relation_maint_rs_mgrpp") {
		_, newRelParam := d.GetChange("relation_maint_rs_mgrpp")
		err = aciClient.CreateRelationmaintRsMgrppFromPODMaintenanceGroup(maintMaintGrp.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_maint_rs_mgrpp")
		d.Partial(false)

	}

	d.SetId(maintMaintGrp.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciPODMaintenanceGroupRead(d, m)

}

func resourceAciPODMaintenanceGroupRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	maintMaintGrp, err := getRemotePODMaintenanceGroup(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setPODMaintenanceGroupAttributes(maintMaintGrp, d)

	maintRsMgrppData, err := aciClient.ReadRelationmaintRsMgrppFromPODMaintenanceGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation maintRsMgrpp %v", err)

	} else {
		d.Set("relation_maint_rs_mgrpp", maintRsMgrppData)
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciPODMaintenanceGroupDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "maintMaintGrp")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}
