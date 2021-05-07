package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAciL4L7ServiceGraphTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciL4L7ServiceGraphTemplateCreate,
		Update: resourceAciL4L7ServiceGraphTemplateUpdate,
		Read:   resourceAciL4L7ServiceGraphTemplateRead,
		Delete: resourceAciL4L7ServiceGraphTemplateDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciL4L7ServiceGraphTemplateImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"tenant_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

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

			"l4_l7_service_graph_template_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"cloud",
					"legacy",
				}, false),
			},

			"ui_template_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ONE_NODE_ADC_ONE_ARM",
					"ONE_NODE_ADC_ONE_ARM_L3EXT",
					"ONE_NODE_ADC_TWO_ARM",
					"ONE_NODE_FW_ROUTED",
					"ONE_NODE_FW_TRANS",
					"TWO_NODE_FW_ROUTED_ADC_ONE_ARM",
					"TWO_NODE_FW_ROUTED_ADC_ONE_ARM_L3EXT",
					"TWO_NODE_FW_ROUTED_ADC_TWO_ARM",
					"TWO_NODE_FW_TRANS_ADC_ONE_ARM",
					"TWO_NODE_FW_TRANS_ADC_ONE_ARM_L3EXT",
					"TWO_NODE_FW_TRANS_ADC_TWO_ARM",
					"UNSPECIFIED",
				}, false),
			},

			"term_cons_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "T1",
			},

			"term_prov_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "T2",
			},

			"term_node_cons_dn": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"term_node_prov_dn": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"term_cons_dn": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"term_prov_dn": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		}),
	}
}

func getRemoteL4L7ServiceGraphTemplate(client *client.Client, dn string) (*models.L4L7ServiceGraphTemplate, error) {
	vnsAbsGraphCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	vnsAbsGraph := models.L4L7ServiceGraphTemplateFromContainer(vnsAbsGraphCont)

	if vnsAbsGraph.DistinguishedName == "" {
		return nil, fmt.Errorf("L4L7ServiceGraphTemplate %s not found", vnsAbsGraph.DistinguishedName)
	}

	return vnsAbsGraph, nil
}

func setL4L7ServiceGraphTemplateAttributes(vnsAbsGraph *models.L4L7ServiceGraphTemplate, d *schema.ResourceData) *schema.ResourceData {
	dn := d.Id()
	d.SetId(vnsAbsGraph.DistinguishedName)
	d.Set("description", vnsAbsGraph.Description)
	if dn != vnsAbsGraph.DistinguishedName {
		d.Set("tenant_dn", "")
	}
	vnsAbsGraphMap, _ := vnsAbsGraph.ToMap()

	d.Set("name", vnsAbsGraphMap["name"])

	d.Set("annotation", vnsAbsGraphMap["annotation"])
	d.Set("name_alias", vnsAbsGraphMap["nameAlias"])
	d.Set("l4_l7_service_graph_template_type", vnsAbsGraphMap["type"])
	d.Set("ui_template_type", vnsAbsGraphMap["uiTemplateType"])

	return d
}

func getRemoteConsumerTerminalNode(client *client.Client, dn string) (*models.ConsumerTerminalNode, error) {
	vnsAbsTermNodeConCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	vnsAbsTermNodeCon := models.ConsumerTerminalNodeFromContainer(vnsAbsTermNodeConCont)

	if vnsAbsTermNodeCon.DistinguishedName == "" {
		return nil, fmt.Errorf("Consumer Terminal Node %s not found", vnsAbsTermNodeCon.DistinguishedName)
	}

	return vnsAbsTermNodeCon, nil
}

func setConsumerTerminalNodeAttributes(vnsAbsTermNodeCon *models.ConsumerTerminalNode, d *schema.ResourceData) *schema.ResourceData {
	vnsAbsTermNodeConMap, _ := vnsAbsTermNodeCon.ToMap()
	d.Set("term_cons_name", vnsAbsTermNodeConMap["name"])
	d.Set("tern_node_cons_dn", vnsAbsTermNodeCon.DistinguishedName)
	return d
}

func getRemoteProviderTerminalNode(client *client.Client, dn string) (*models.ProviderTerminalNode, error) {
	vnsAbsTermNodeProvCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	vnsAbsTermNodeProv := models.ProviderTerminalNodeFromContainer(vnsAbsTermNodeProvCont)

	if vnsAbsTermNodeProv.DistinguishedName == "" {
		return nil, fmt.Errorf("Provider Terminal Node %s not found", vnsAbsTermNodeProv.DistinguishedName)
	}

	return vnsAbsTermNodeProv, nil
}

func setProviderTerminalNodeAttributes(vnsAbsTermNodeProv *models.ProviderTerminalNode, d *schema.ResourceData) *schema.ResourceData {
	vnsAbsTermNodeProvMap, _ := vnsAbsTermNodeProv.ToMap()
	d.Set("term_prov_name", vnsAbsTermNodeProvMap["name"])
	d.Set("tern_node_prov_dn", vnsAbsTermNodeProv.DistinguishedName)
	return d
}

func getRemoteTerminalConnector(client *client.Client, dn string) (*models.TerminalConnector, error) {
	vnsAbsTermConnCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	vnsAbsTermConn := models.TerminalConnectorFromContainer(vnsAbsTermConnCont)

	if vnsAbsTermConn.DistinguishedName == "" {
		return nil, fmt.Errorf("Terminal Connector %s not found", vnsAbsTermConn.DistinguishedName)
	}

	return vnsAbsTermConn, nil
}

func resourceAciL4L7ServiceGraphTemplateImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	vnsAbsGraph, err := getRemoteL4L7ServiceGraphTemplate(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setL4L7ServiceGraphTemplateAttributes(vnsAbsGraph, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciL4L7ServiceGraphTemplateCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] L4L7ServiceGraphTemplate: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	tenantDn := d.Get("tenant_dn").(string)

	vnsAbsGraphAttr := models.L4L7ServiceGraphTemplateAttributes{}
	if annotation, ok := d.GetOk("annotation"); ok {
		vnsAbsGraphAttr.Annotation = annotation.(string)
	} else {
		vnsAbsGraphAttr.Annotation = "{}"
	}
	if nameAlias, ok := d.GetOk("name_alias"); ok {
		vnsAbsGraphAttr.NameAlias = nameAlias.(string)
	}
	if l4l7ServiceGraphTemplateType, ok := d.GetOk("l4_l7_service_graph_template_type"); ok {
		vnsAbsGraphAttr.L4L7ServiceGraphTemplate_type = l4l7ServiceGraphTemplateType.(string)
	}
	if UiTemplateType, ok := d.GetOk("ui_template_type"); ok {
		vnsAbsGraphAttr.UiTemplateType = UiTemplateType.(string)
	}
	vnsAbsGraph := models.NewL4L7ServiceGraphTemplate(fmt.Sprintf("AbsGraph-%s", name), tenantDn, desc, vnsAbsGraphAttr)

	err := aciClient.Save(vnsAbsGraph)
	if err != nil {
		return err
	}
	d.Partial(true)

	d.Partial(false)

	if consName, ok := d.GetOk("term_cons_name"); ok {
		vnsAbsTermNodeConAttr := models.ConsumerTerminalNodeAttributes{}

		vnsAbsTermNodeConAttr.Annotation = "{}"

		vnsAbsTermNodeCon := models.NewConsumerTerminalNode(fmt.Sprintf("AbsTermNodeCon-%s", consName), vnsAbsGraph.DistinguishedName, "", vnsAbsTermNodeConAttr)

		err := aciClient.Save(vnsAbsTermNodeCon)
		if err != nil {
			return err
		}
		d.Set("term_node_cons_dn", vnsAbsTermNodeCon.DistinguishedName)

		vnsAbsTermConnAttr := models.TerminalConnectorAttributes{}
		vnsAbsTermConnAttr.Annotation = "{}"
		vnsAbsTermConn := models.NewTerminalConnector(fmt.Sprintf("AbsTConn"), vnsAbsTermNodeCon.DistinguishedName, "", vnsAbsTermConnAttr)

		err = aciClient.Save(vnsAbsTermConn)
		if err != nil {
			return err
		}

		d.Set("term_cons_dn", vnsAbsTermConn.DistinguishedName)
	}

	if provName, ok := d.GetOk("term_prov_name"); ok {
		vnsAbsTermNodeProvAttr := models.ProviderTerminalNodeAttributes{}

		vnsAbsTermNodeProvAttr.Annotation = "{}"

		vnsAbsTermNodeProv := models.NewProviderTerminalNode(fmt.Sprintf("AbsTermNodeProv-%s", provName), vnsAbsGraph.DistinguishedName, "", vnsAbsTermNodeProvAttr)

		err := aciClient.Save(vnsAbsTermNodeProv)
		if err != nil {
			return err
		}
		d.Set("term_node_prov_dn", vnsAbsTermNodeProv.DistinguishedName)

		vnsAbsTermConnAttr := models.TerminalConnectorAttributes{}
		vnsAbsTermConnAttr.Annotation = "{}"
		vnsAbsTermConn := models.NewTerminalConnector(fmt.Sprintf("AbsTConn"), vnsAbsTermNodeProv.DistinguishedName, "", vnsAbsTermConnAttr)

		err = aciClient.Save(vnsAbsTermConn)
		if err != nil {
			return err
		}

		d.Set("term_prov_dn", vnsAbsTermConn.DistinguishedName)
	}

	d.SetId(vnsAbsGraph.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciL4L7ServiceGraphTemplateRead(d, m)
}

func resourceAciL4L7ServiceGraphTemplateUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] L4L7ServiceGraphTemplate: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	tenantDn := d.Get("tenant_dn").(string)

	vnsAbsGraphAttr := models.L4L7ServiceGraphTemplateAttributes{}
	if annotation, ok := d.GetOk("annotation"); ok {
		vnsAbsGraphAttr.Annotation = annotation.(string)
	} else {
		vnsAbsGraphAttr.Annotation = "{}"
	}
	if nameAlias, ok := d.GetOk("name_alias"); ok {
		vnsAbsGraphAttr.NameAlias = nameAlias.(string)
	}
	if l4l7ServiceGraphTemplateType, ok := d.GetOk("l4_l7_service_graph_template_type"); ok {
		vnsAbsGraphAttr.L4L7ServiceGraphTemplate_type = l4l7ServiceGraphTemplateType.(string)
	}
	if uiTemplateType, ok := d.GetOk("ui_template_type"); ok {
		vnsAbsGraphAttr.UiTemplateType = uiTemplateType.(string)
	}
	vnsAbsGraph := models.NewL4L7ServiceGraphTemplate(fmt.Sprintf("AbsGraph-%s", name), tenantDn, desc, vnsAbsGraphAttr)

	vnsAbsGraph.Status = "modified"

	err := aciClient.Save(vnsAbsGraph)

	if err != nil {
		return err
	}
	d.Partial(true)

	d.Partial(false)

	if d.HasChange("term_cons_name") {
		consNameDn := d.Get("term_node_cons_dn").(string)

		err := aciClient.DeleteByDn(consNameDn, "vnsAbsTermNodeCon")
		if err != nil {
			return err
		}

		consName := d.Get("term_cons_name").(string)
		vnsAbsTermNodeConAttr := models.ConsumerTerminalNodeAttributes{}

		vnsAbsTermNodeConAttr.Annotation = "{}"

		vnsAbsTermNodeCon := models.NewConsumerTerminalNode(fmt.Sprintf("AbsTermNodeCon-%s", consName), vnsAbsGraph.DistinguishedName, "", vnsAbsTermNodeConAttr)

		err = aciClient.Save(vnsAbsTermNodeCon)
		if err != nil {
			return err
		}
		d.Set("term_node_cons_dn", vnsAbsTermNodeCon.DistinguishedName)

		vnsAbsTermConnAttr := models.TerminalConnectorAttributes{}
		vnsAbsTermConnAttr.Annotation = "{}"
		vnsAbsTermConn := models.NewTerminalConnector(fmt.Sprintf("AbsTConn"), vnsAbsTermNodeCon.DistinguishedName, "", vnsAbsTermConnAttr)

		err = aciClient.Save(vnsAbsTermConn)
		if err != nil {
			return err
		}

		d.Set("term_cons_dn", vnsAbsTermConn.DistinguishedName)

	}

	if d.HasChange("term_prov_name") {
		provNameDn := d.Get("term_node_prov_dn").(string)

		err := aciClient.DeleteByDn(provNameDn, "vnsAbsTermNodeProv")
		if err != nil {
			return err
		}

		provName := d.Get("term_prov_name").(string)
		vnsAbsTermNodeProvAttr := models.ProviderTerminalNodeAttributes{}

		vnsAbsTermNodeProvAttr.Annotation = "{}"

		vnsAbsTermNodeProv := models.NewProviderTerminalNode(fmt.Sprintf("AbsTermNodeProv-%s", provName), vnsAbsGraph.DistinguishedName, "", vnsAbsTermNodeProvAttr)

		err = aciClient.Save(vnsAbsTermNodeProv)
		if err != nil {
			return err
		}
		d.Set("term_node_prov_dn", vnsAbsTermNodeProv.DistinguishedName)

		vnsAbsTermConnAttr := models.TerminalConnectorAttributes{}
		vnsAbsTermConnAttr.Annotation = "{}"
		vnsAbsTermConn := models.NewTerminalConnector(fmt.Sprintf("AbsTConn"), vnsAbsTermNodeProv.DistinguishedName, "", vnsAbsTermConnAttr)

		err = aciClient.Save(vnsAbsTermConn)
		if err != nil {
			return err
		}

		d.Set("term_prov_dn", vnsAbsTermConn.DistinguishedName)
	}

	d.SetId(vnsAbsGraph.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciL4L7ServiceGraphTemplateRead(d, m)

}

func resourceAciL4L7ServiceGraphTemplateRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	vnsAbsGraph, err := getRemoteL4L7ServiceGraphTemplate(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setL4L7ServiceGraphTemplateAttributes(vnsAbsGraph, d)

	consDn := d.Get("term_node_cons_dn").(string)
	vnsAbsTermNodeCon, err := getRemoteConsumerTerminalNode(aciClient, consDn)
	if err != nil {
		d.Set("term_node_cons_dn", "")
		d.Set("term_cons_name", "")
		return nil
	}
	setConsumerTerminalNodeAttributes(vnsAbsTermNodeCon, d)

	vnsAbsTermConnDn := d.Get("term_cons_dn").(string)
	vnsAbsTermConn, err := getRemoteTerminalConnector(aciClient, vnsAbsTermConnDn)
	if err != nil {
		d.Set("term_cons_dn", "")
	} else {
		d.Set("term_cons_dn", vnsAbsTermConn.DistinguishedName)
	}

	provDn := d.Get("term_node_prov_dn").(string)
	vnsAbsTermNodeProv, err := getRemoteProviderTerminalNode(aciClient, provDn)
	if err != nil {
		d.Set("term_node_prov_dn", "")
		d.Set("term_prov_name", "")
		return nil
	}
	setProviderTerminalNodeAttributes(vnsAbsTermNodeProv, d)

	vnsAbsTermConnDn = d.Get("term_prov_dn").(string)
	vnsAbsTermConn, err = getRemoteTerminalConnector(aciClient, vnsAbsTermConnDn)
	if err != nil {
		d.Set("term_prov_dn", "")
	} else {
		d.Set("term_prov_dn", vnsAbsTermConn.DistinguishedName)
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciL4L7ServiceGraphTemplateDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "vnsAbsGraph")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}
