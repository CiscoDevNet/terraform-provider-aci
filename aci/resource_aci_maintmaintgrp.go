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

func resourceAciPODMaintenanceGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciPODMaintenanceGroupCreate,
		UpdateContext: resourceAciPODMaintenanceGroupUpdate,
		ReadContext:   resourceAciPODMaintenanceGroupRead,
		DeleteContext: resourceAciPODMaintenanceGroupDelete,

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

			"fwtype": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"controller",
					"switch",
					"catalog",
					"plugin",
					"pluginPackage",
					"config",
					"vpod",
				}, false),
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

func setPODMaintenanceGroupAttributes(maintMaintGrp *models.PODMaintenanceGroup, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(maintMaintGrp.DistinguishedName)
	d.Set("description", maintMaintGrp.Description)
	maintMaintGrpMap, err := maintMaintGrp.ToMap()
	if err != nil {
		return d, err
	}

	d.Set("name", maintMaintGrpMap["name"])

	d.Set("annotation", maintMaintGrpMap["annotation"])
	d.Set("fwtype", maintMaintGrpMap["fwtype"])
	d.Set("name_alias", maintMaintGrpMap["nameAlias"])
	d.Set("pod_maintenance_group_type", maintMaintGrpMap["type"])
	return d, nil
}

func resourceAciPODMaintenanceGroupImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	maintMaintGrp, err := getRemotePODMaintenanceGroup(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled, err := setPODMaintenanceGroupAttributes(maintMaintGrp, d)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciPODMaintenanceGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] PODMaintenanceGroup: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	maintMaintGrpAttr := models.PODMaintenanceGroupAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		maintMaintGrpAttr.Annotation = Annotation.(string)
	} else {
		maintMaintGrpAttr.Annotation = "{}"
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
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if relationTomaintRsMgrpp, ok := d.GetOk("relation_maint_rs_mgrpp"); ok {
		relationParam := relationTomaintRsMgrpp.(string)
		checkDns = append(checkDns, relationParam)
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if relationTomaintRsMgrpp, ok := d.GetOk("relation_maint_rs_mgrpp"); ok {
		relationParam := relationTomaintRsMgrpp.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationmaintRsMgrppFromPODMaintenanceGroup(maintMaintGrp.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(maintMaintGrp.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciPODMaintenanceGroupRead(ctx, d, m)
}

func resourceAciPODMaintenanceGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] PODMaintenanceGroup: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	maintMaintGrpAttr := models.PODMaintenanceGroupAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		maintMaintGrpAttr.Annotation = Annotation.(string)
	} else {
		maintMaintGrpAttr.Annotation = "{}"
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
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_maint_rs_mgrpp") {
		_, newRelParam := d.GetChange("relation_maint_rs_mgrpp")
		checkDns = append(checkDns, newRelParam.(string))
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if d.HasChange("relation_maint_rs_mgrpp") {
		_, newRelParam := d.GetChange("relation_maint_rs_mgrpp")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.CreateRelationmaintRsMgrppFromPODMaintenanceGroup(maintMaintGrp.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(maintMaintGrp.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciPODMaintenanceGroupRead(ctx, d, m)

}

func resourceAciPODMaintenanceGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	maintMaintGrp, err := getRemotePODMaintenanceGroup(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setPODMaintenanceGroupAttributes(maintMaintGrp, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	maintRsMgrppData, err := aciClient.ReadRelationmaintRsMgrppFromPODMaintenanceGroup(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation maintRsMgrpp %v", err)
		d.Set("relation_maint_rs_mgrpp", "")

	} else {
		d.Set("relation_maint_rs_mgrpp", maintRsMgrppData.(string))
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciPODMaintenanceGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "maintMaintGrp")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}
