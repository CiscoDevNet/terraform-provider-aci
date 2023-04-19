package aci

import (
	"context"
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAciTACACSMonitoringDestinationGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciTACACSMonitoringDestinationGroupCreate,
		UpdateContext: resourceAciTACACSMonitoringDestinationGroupUpdate,
		ReadContext:   resourceAciTACACSMonitoringDestinationGroupRead,
		DeleteContext: resourceAciTACACSMonitoringDestinationGroupDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciTACACSMonitoringDestinationGroupImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		})),
	}
}

func getRemoteTACACSMonitoringDestinationGroup(client *client.Client, dn string) (*models.TACACSMonitoringDestinationGroup, error) {
	tacacsGroupCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	tacacsGroup := models.TACACSMonitoringDestinationGroupFromContainer(tacacsGroupCont)
	if tacacsGroup.DistinguishedName == "" {
		return nil, fmt.Errorf("TACACS Monitoring Destination Group %s not found", dn)
	}
	return tacacsGroup, nil
}

func setTACACSMonitoringDestinationGroupAttributes(tacacsGroup *models.TACACSMonitoringDestinationGroup, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(tacacsGroup.DistinguishedName)
	d.Set("description", tacacsGroup.Description)
	tacacsGroupMap, err := tacacsGroup.ToMap()
	if err != nil {
		return nil, err
	}
	d.Set("annotation", tacacsGroupMap["annotation"])
	d.Set("name", tacacsGroupMap["name"])
	d.Set("name_alias", tacacsGroupMap["nameAlias"])
	return d, nil
}

func resourceAciTACACSMonitoringDestinationGroupImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	tacacsGroup, err := getRemoteTACACSMonitoringDestinationGroup(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setTACACSMonitoringDestinationGroupAttributes(tacacsGroup, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciTACACSMonitoringDestinationGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] TACACSMonitoringDestinationGroup: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	tacacsGroupAttr := models.TACACSMonitoringDestinationGroupAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		tacacsGroupAttr.Annotation = Annotation.(string)
	} else {
		tacacsGroupAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		tacacsGroupAttr.Name = Name.(string)
	}
	tacacsGroup := models.NewTACACSMonitoringDestinationGroup(fmt.Sprintf("fabric/tacacsgroup-%s", name), "uni", desc, nameAlias, tacacsGroupAttr)
	err := aciClient.Save(tacacsGroup)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(tacacsGroup.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciTACACSMonitoringDestinationGroupRead(ctx, d, m)
}

func resourceAciTACACSMonitoringDestinationGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] TACACSMonitoringDestinationGroup: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	tacacsGroupAttr := models.TACACSMonitoringDestinationGroupAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		tacacsGroupAttr.Annotation = Annotation.(string)
	} else {
		tacacsGroupAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		tacacsGroupAttr.Name = Name.(string)
	}
	tacacsGroup := models.NewTACACSMonitoringDestinationGroup(fmt.Sprintf("fabric/tacacsgroup-%s", name), "uni", desc, nameAlias, tacacsGroupAttr)
	tacacsGroup.Status = "modified"
	err := aciClient.Save(tacacsGroup)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(tacacsGroup.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciTACACSMonitoringDestinationGroupRead(ctx, d, m)
}

func resourceAciTACACSMonitoringDestinationGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	tacacsGroup, err := getRemoteTACACSMonitoringDestinationGroup(aciClient, dn)
	if err != nil {
		return errorForObjectNotFound(err, dn, d)
	}
	_, err = setTACACSMonitoringDestinationGroupAttributes(tacacsGroup, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciTACACSMonitoringDestinationGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "tacacsGroup")
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}
