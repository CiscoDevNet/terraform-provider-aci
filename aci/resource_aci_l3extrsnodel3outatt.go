package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAciFabricNode() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciFabricNodeCreate,
		Update: resourceAciFabricNodeUpdate,
		Read:   resourceAciFabricNodeRead,
		Delete: resourceAciFabricNodeDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciFabricNodeImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"logical_node_profile_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"tdn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"config_issues": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"none",
					"node-path-misconfig",
					"routerid-not-changable-with-mcast",
					"loopback-ip-missing",
				}, false),
			},

			"rtr_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"rtr_id_loop_back": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
			},
		}),
	}
}
func getRemoteFabricNode(client *client.Client, dn string) (*models.FabricNode, error) {
	l3extRsNodeL3OutAttCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	l3extRsNodeL3OutAtt := models.FabricNodeFromContainer(l3extRsNodeL3OutAttCont)

	if l3extRsNodeL3OutAtt.DistinguishedName == "" {
		return nil, fmt.Errorf("FabricNode %s not found", l3extRsNodeL3OutAtt.DistinguishedName)
	}

	return l3extRsNodeL3OutAtt, nil
}

func setFabricNodeAttributes(l3extRsNodeL3OutAtt *models.FabricNode, d *schema.ResourceData) *schema.ResourceData {
	dn := d.Id()
	d.SetId(l3extRsNodeL3OutAtt.DistinguishedName)
	// d.Set("logical_node_profile_dn", GetParentDn(l3extRsNodeL3OutAtt.DistinguishedName))
	if dn != l3extRsNodeL3OutAtt.DistinguishedName {
		d.Set("logical_node_profile_dn", "")
	}
	l3extRsNodeL3OutAttMap, _ := l3extRsNodeL3OutAtt.ToMap()

	d.Set("tdn", l3extRsNodeL3OutAttMap["tDn"])

	d.Set("annotation", l3extRsNodeL3OutAttMap["annotation"])
	d.Set("config_issues", l3extRsNodeL3OutAttMap["configIssues"])
	d.Set("rtr_id", l3extRsNodeL3OutAttMap["rtrId"])
	d.Set("rtr_id_loop_back", l3extRsNodeL3OutAttMap["rtrIdLoopBack"])
	return d
}

func resourceAciFabricNodeImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	l3extRsNodeL3OutAtt, err := getRemoteFabricNode(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setFabricNodeAttributes(l3extRsNodeL3OutAtt, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciFabricNodeCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] FabricNode: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	tDn := d.Get("tdn").(string)

	LogicalNodeProfileDn := d.Get("logical_node_profile_dn").(string)

	l3extRsNodeL3OutAttAttr := models.FabricNodeAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		l3extRsNodeL3OutAttAttr.Annotation = Annotation.(string)
	} else {
		l3extRsNodeL3OutAttAttr.Annotation = "{}"
	}
	if ConfigIssues, ok := d.GetOk("config_issues"); ok {
		l3extRsNodeL3OutAttAttr.ConfigIssues = ConfigIssues.(string)
	}
	if RtrId, ok := d.GetOk("rtr_id"); ok {
		l3extRsNodeL3OutAttAttr.RtrId = RtrId.(string)
	}
	if RtrIdLoopBack, ok := d.GetOk("rtr_id_loop_back"); ok {
		l3extRsNodeL3OutAttAttr.RtrIdLoopBack = RtrIdLoopBack.(string)
	}

	l3extRsNodeL3OutAtt := models.NewFabricNode(fmt.Sprintf("rsnodeL3OutAtt-[%s]", tDn), LogicalNodeProfileDn, desc, l3extRsNodeL3OutAttAttr)

	err := aciClient.Save(l3extRsNodeL3OutAtt)
	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("tdn")

	d.Partial(false)

	d.SetId(l3extRsNodeL3OutAtt.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciFabricNodeRead(d, m)
}

func resourceAciFabricNodeUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] FabricNode: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	tDn := d.Get("tdn").(string)

	LogicalNodeProfileDn := d.Get("logical_node_profile_dn").(string)

	l3extRsNodeL3OutAttAttr := models.FabricNodeAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		l3extRsNodeL3OutAttAttr.Annotation = Annotation.(string)
	} else {
		l3extRsNodeL3OutAttAttr.Annotation = "{}"
	}
	if ConfigIssues, ok := d.GetOk("config_issues"); ok {
		l3extRsNodeL3OutAttAttr.ConfigIssues = ConfigIssues.(string)
	}
	if RtrId, ok := d.GetOk("rtr_id"); ok {
		l3extRsNodeL3OutAttAttr.RtrId = RtrId.(string)
	}
	if RtrIdLoopBack, ok := d.GetOk("rtr_id_loop_back"); ok {
		l3extRsNodeL3OutAttAttr.RtrIdLoopBack = RtrIdLoopBack.(string)
	}

	l3extRsNodeL3OutAtt := models.NewFabricNode(fmt.Sprintf("rsnodeL3OutAtt-[%s]", tDn), LogicalNodeProfileDn, desc, l3extRsNodeL3OutAttAttr)

	l3extRsNodeL3OutAtt.Status = "modified"

	err := aciClient.Save(l3extRsNodeL3OutAtt)

	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("tdn")

	d.Partial(false)

	d.SetId(l3extRsNodeL3OutAtt.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciFabricNodeRead(d, m)

}

func resourceAciFabricNodeRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	l3extRsNodeL3OutAtt, err := getRemoteFabricNode(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setFabricNodeAttributes(l3extRsNodeL3OutAtt, d)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciFabricNodeDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "l3extRsNodeL3OutAtt")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}
