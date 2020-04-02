package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAciTriggerScheduler() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciTriggerSchedulerCreate,
		Update: resourceAciTriggerSchedulerUpdate,
		Read:   resourceAciTriggerSchedulerRead,
		Delete: resourceAciTriggerSchedulerDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciTriggerSchedulerImport,
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

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}
func getRemoteTriggerScheduler(client *client.Client, dn string) (*models.TriggerScheduler, error) {
	trigSchedPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	trigSchedP := models.TriggerSchedulerFromContainer(trigSchedPCont)

	if trigSchedP.DistinguishedName == "" {
		return nil, fmt.Errorf("TriggerScheduler %s not found", trigSchedP.DistinguishedName)
	}

	return trigSchedP, nil
}

func setTriggerSchedulerAttributes(trigSchedP *models.TriggerScheduler, d *schema.ResourceData) *schema.ResourceData {
	d.SetId(trigSchedP.DistinguishedName)
	d.Set("description", trigSchedP.Description)
	trigSchedPMap, _ := trigSchedP.ToMap()

	d.Set("name", trigSchedPMap["name"])

	d.Set("annotation", trigSchedPMap["annotation"])
	d.Set("name_alias", trigSchedPMap["nameAlias"])
	return d
}

func resourceAciTriggerSchedulerImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	trigSchedP, err := getRemoteTriggerScheduler(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setTriggerSchedulerAttributes(trigSchedP, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciTriggerSchedulerCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] TriggerScheduler: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	trigSchedPAttr := models.TriggerSchedulerAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		trigSchedPAttr.Annotation = Annotation.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		trigSchedPAttr.NameAlias = NameAlias.(string)
	}
	trigSchedP := models.NewTriggerScheduler(fmt.Sprintf("fabric/schedp-%s", name), "uni", desc, trigSchedPAttr)

	err := aciClient.Save(trigSchedP)
	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	d.SetId(trigSchedP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciTriggerSchedulerRead(d, m)
}

func resourceAciTriggerSchedulerUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] TriggerScheduler: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	trigSchedPAttr := models.TriggerSchedulerAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		trigSchedPAttr.Annotation = Annotation.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		trigSchedPAttr.NameAlias = NameAlias.(string)
	}
	trigSchedP := models.NewTriggerScheduler(fmt.Sprintf("fabric/schedp-%s", name), "uni", desc, trigSchedPAttr)

	trigSchedP.Status = "modified"

	err := aciClient.Save(trigSchedP)

	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	d.SetId(trigSchedP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciTriggerSchedulerRead(d, m)

}

func resourceAciTriggerSchedulerRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	trigSchedP, err := getRemoteTriggerScheduler(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setTriggerSchedulerAttributes(trigSchedP, d)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciTriggerSchedulerDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "trigSchedP")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}
