package aci

import (
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAciFunctionNode() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciFunctionNodeCreate,
		Update: resourceAciFunctionNodeUpdate,
		Read:   resourceAciFunctionNodeRead,
		Delete: resourceAciFunctionNodeDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciFunctionNodeImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"l4_l7_service_graph_template_dn": &schema.Schema{
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

			"func_template_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"OTHER",
					"FW_TRANS",
					"FW_ROUTED",
					"CLOUD_VENDOR_LB",
					"CLOUD_VENDOR_FW",
					"CLOUD_NATIVE_LB",
					"CLOUD_NATIVE_FW",
					"ADC_TWO_ARM",
					"ADC_ONE_ARM",
				}, false),
			},

			"func_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"GoThrough",
					"GoTo",
					"L1",
					"L2",
					"None",
				}, false),
			},

			"is_copy": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"managed": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"routing_mode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Redirect",
				}, false),
			},

			"sequence_number": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"share_encap": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"relation_vns_rs_node_to_abs_func_prof": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_vns_rs_node_to_l_dev": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_vns_rs_node_to_m_func": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_vns_rs_default_scope_to_term": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_vns_rs_node_to_cloud_l_dev": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
		}),
	}
}
func getRemoteFunctionNode(client *client.Client, dn string) (*models.FunctionNode, error) {
	vnsAbsNodeCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	vnsAbsNode := models.FunctionNodeFromContainer(vnsAbsNodeCont)

	if vnsAbsNode.DistinguishedName == "" {
		return nil, fmt.Errorf("FunctionNode %s not found", vnsAbsNode.DistinguishedName)
	}

	return vnsAbsNode, nil
}

func setFunctionNodeAttributes(vnsAbsNode *models.FunctionNode, d *schema.ResourceData) *schema.ResourceData {

	dn := d.Id()
	d.SetId(vnsAbsNode.DistinguishedName)
	d.Set("description", vnsAbsNode.Description)

	if dn != vnsAbsNode.DistinguishedName {
		d.Set("l4_l7_service_graph_template_dn", "")
	}

	vnsAbsNodeMap, _ := vnsAbsNode.ToMap()
	d.Set("name", vnsAbsNodeMap["name"])

	d.Set("annotation", vnsAbsNodeMap["annotation"])
	d.Set("func_template_type", vnsAbsNodeMap["funcTemplateType"])
	d.Set("func_type", vnsAbsNodeMap["funcType"])
	d.Set("is_copy", vnsAbsNodeMap["isCopy"])
	d.Set("managed", vnsAbsNodeMap["managed"])
	d.Set("name_alias", vnsAbsNodeMap["nameAlias"])
	d.Set("routing_mode", vnsAbsNodeMap["routingMode"])
	d.Set("sequence_number", vnsAbsNodeMap["sequenceNumber"])
	d.Set("share_encap", vnsAbsNodeMap["shareEncap"])
	return d
}

func resourceAciFunctionNodeImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	vnsAbsNode, err := getRemoteFunctionNode(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setFunctionNodeAttributes(vnsAbsNode, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciFunctionNodeCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] FunctionNode: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	l4l7ServiceGraphTemplateDn := d.Get("l4_l7_service_graph_template_dn").(string)

	vnsAbsNodeAttr := models.FunctionNodeAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		vnsAbsNodeAttr.Annotation = Annotation.(string)
	} else {
		vnsAbsNodeAttr.Annotation = "{}"
	}
	if FuncTemplateType, ok := d.GetOk("func_template_type"); ok {
		vnsAbsNodeAttr.FuncTemplateType = FuncTemplateType.(string)
	}
	if FuncType, ok := d.GetOk("func_type"); ok {
		vnsAbsNodeAttr.FuncType = FuncType.(string)
	}
	if IsCopy, ok := d.GetOk("is_copy"); ok {
		vnsAbsNodeAttr.IsCopy = IsCopy.(string)
	}
	if Managed, ok := d.GetOk("managed"); ok {
		vnsAbsNodeAttr.Managed = Managed.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		vnsAbsNodeAttr.NameAlias = NameAlias.(string)
	}
	if RoutingMode, ok := d.GetOk("routing_mode"); ok {
		vnsAbsNodeAttr.RoutingMode = RoutingMode.(string)
	}
	if SequenceNumber, ok := d.GetOk("sequence_number"); ok {
		vnsAbsNodeAttr.SequenceNumber = SequenceNumber.(string)
	}
	if ShareEncap, ok := d.GetOk("share_encap"); ok {
		vnsAbsNodeAttr.ShareEncap = ShareEncap.(string)
	}
	vnsAbsNode := models.NewFunctionNode(fmt.Sprintf("AbsNode-%s", name), l4l7ServiceGraphTemplateDn, desc, vnsAbsNodeAttr)

	err := aciClient.Save(vnsAbsNode)
	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	checkDns := make([]string, 0, 1)

	if relationTovnsRsNodeToAbsFuncProf, ok := d.GetOk("relation_vns_rs_node_to_abs_func_prof"); ok {
		relationParam := relationTovnsRsNodeToAbsFuncProf.(string)
		checkDns = append(checkDns, relationParam)
	}
	if relationTovnsRsNodeToLDev, ok := d.GetOk("relation_vns_rs_node_to_l_dev"); ok {
		relationParam := relationTovnsRsNodeToLDev.(string)
		checkDns = append(checkDns, relationParam)
	}
	if relationTovnsRsNodeToMFunc, ok := d.GetOk("relation_vns_rs_node_to_m_func"); ok {
		relationParam := relationTovnsRsNodeToMFunc.(string)
		checkDns = append(checkDns, relationParam)
	}
	if relationTovnsRsDefaultScopeToTerm, ok := d.GetOk("relation_vns_rs_default_scope_to_term"); ok {
		relationParam := relationTovnsRsDefaultScopeToTerm.(string)
		checkDns = append(checkDns, relationParam)

	}
	if relationTovnsRsNodeToCloudLDev, ok := d.GetOk("relation_vns_rs_node_to_cloud_l_dev"); ok {
		relationParam := relationTovnsRsNodeToCloudLDev.(string)
		checkDns = append(checkDns, relationParam)

	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return err
	}
	d.Partial(false)

	if relationTovnsRsNodeToAbsFuncProf, ok := d.GetOk("relation_vns_rs_node_to_abs_func_prof"); ok {
		relationParam := relationTovnsRsNodeToAbsFuncProf.(string)
		err = aciClient.CreateRelationvnsRsNodeToAbsFuncProfFromFunctionNode(vnsAbsNode.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_vns_rs_node_to_abs_func_prof")
		d.Partial(false)

	}
	if relationTovnsRsNodeToLDev, ok := d.GetOk("relation_vns_rs_node_to_l_dev"); ok {
		relationParam := relationTovnsRsNodeToLDev.(string)
		err = aciClient.CreateRelationvnsRsNodeToLDevFromFunctionNode(vnsAbsNode.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_vns_rs_node_to_l_dev")
		d.Partial(false)

	}
	if relationTovnsRsNodeToMFunc, ok := d.GetOk("relation_vns_rs_node_to_m_func"); ok {
		relationParam := relationTovnsRsNodeToMFunc.(string)
		err = aciClient.CreateRelationvnsRsNodeToMFuncFromFunctionNode(vnsAbsNode.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_vns_rs_node_to_m_func")
		d.Partial(false)

	}
	if relationTovnsRsDefaultScopeToTerm, ok := d.GetOk("relation_vns_rs_default_scope_to_term"); ok {
		relationParam := relationTovnsRsDefaultScopeToTerm.(string)
		err = aciClient.CreateRelationvnsRsDefaultScopeToTermFromFunctionNode(vnsAbsNode.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_vns_rs_default_scope_to_term")
		d.Partial(false)

	}
	if relationTovnsRsNodeToCloudLDev, ok := d.GetOk("relation_vns_rs_node_to_cloud_l_dev"); ok {
		relationParam := relationTovnsRsNodeToCloudLDev.(string)
		err = aciClient.CreateRelationvnsRsNodeToCloudLDevFromFunctionNode(vnsAbsNode.DistinguishedName, relationParam)
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_vns_rs_node_to_cloud_l_dev")
		d.Partial(false)

	}

	d.SetId(vnsAbsNode.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciFunctionNodeRead(d, m)
}

func resourceAciFunctionNodeUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] FunctionNode: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	l4l7ServiceGraphTemplateDn := d.Get("l4_l7_service_graph_template_dn").(string)

	vnsAbsNodeAttr := models.FunctionNodeAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		vnsAbsNodeAttr.Annotation = Annotation.(string)
	} else {
		vnsAbsNodeAttr.Annotation = "{}"
	}
	if FuncTemplateType, ok := d.GetOk("func_template_type"); ok {
		vnsAbsNodeAttr.FuncTemplateType = FuncTemplateType.(string)
	}
	if FuncType, ok := d.GetOk("func_type"); ok {
		vnsAbsNodeAttr.FuncType = FuncType.(string)
	}
	if IsCopy, ok := d.GetOk("is_copy"); ok {
		vnsAbsNodeAttr.IsCopy = IsCopy.(string)
	}
	if Managed, ok := d.GetOk("managed"); ok {
		vnsAbsNodeAttr.Managed = Managed.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		vnsAbsNodeAttr.NameAlias = NameAlias.(string)
	}
	if RoutingMode, ok := d.GetOk("routing_mode"); ok {
		vnsAbsNodeAttr.RoutingMode = RoutingMode.(string)
	}
	if SequenceNumber, ok := d.GetOk("sequence_number"); ok {
		vnsAbsNodeAttr.SequenceNumber = SequenceNumber.(string)
	}
	if ShareEncap, ok := d.GetOk("share_encap"); ok {
		vnsAbsNodeAttr.ShareEncap = ShareEncap.(string)
	}
	vnsAbsNode := models.NewFunctionNode(fmt.Sprintf("AbsNode-%s", name), l4l7ServiceGraphTemplateDn, desc, vnsAbsNodeAttr)

	vnsAbsNode.Status = "modified"

	err := aciClient.Save(vnsAbsNode)

	if err != nil {
		return err
	}
	d.Partial(true)

	d.SetPartial("name")

	d.Partial(false)

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_vns_rs_node_to_abs_func_prof") {
		_, newRelParam := d.GetChange("relation_vns_rs_node_to_abs_func_prof")
		checkDns = append(checkDns, newRelParam.(string))
	}
	if d.HasChange("relation_vns_rs_node_to_l_dev") {
		_, newRelParam := d.GetChange("relation_vns_rs_node_to_l_dev")
		checkDns = append(checkDns, newRelParam.(string))
	}
	if d.HasChange("relation_vns_rs_node_to_m_func") {
		_, newRelParam := d.GetChange("relation_vns_rs_node_to_m_func")
		checkDns = append(checkDns, newRelParam.(string))
	}
	if d.HasChange("relation_vns_rs_default_scope_to_term") {
		_, newRelParam := d.GetChange("relation_vns_rs_default_scope_to_term")
		checkDns = append(checkDns, newRelParam.(string))
	}
	if d.HasChange("relation_vns_rs_node_to_cloud_l_dev") {
		_, newRelParam := d.GetChange("relation_vns_rs_node_to_cloud_l_dev")
		checkDns = append(checkDns, newRelParam.(string))
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return err
	}
	d.Partial(false)

	if d.HasChange("relation_vns_rs_node_to_abs_func_prof") {
		_, newRelParam := d.GetChange("relation_vns_rs_node_to_abs_func_prof")
		err = aciClient.DeleteRelationvnsRsNodeToAbsFuncProfFromFunctionNode(vnsAbsNode.DistinguishedName)
		if err != nil {
			return err
		}
		err = aciClient.CreateRelationvnsRsNodeToAbsFuncProfFromFunctionNode(vnsAbsNode.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_vns_rs_node_to_abs_func_prof")
		d.Partial(false)

	}
	if d.HasChange("relation_vns_rs_node_to_l_dev") {
		_, newRelParam := d.GetChange("relation_vns_rs_node_to_l_dev")
		err = aciClient.DeleteRelationvnsRsNodeToLDevFromFunctionNode(vnsAbsNode.DistinguishedName)
		if err != nil {
			return err
		}
		err = aciClient.CreateRelationvnsRsNodeToLDevFromFunctionNode(vnsAbsNode.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_vns_rs_node_to_l_dev")
		d.Partial(false)

	}
	if d.HasChange("relation_vns_rs_node_to_m_func") {
		_, newRelParam := d.GetChange("relation_vns_rs_node_to_m_func")
		err = aciClient.DeleteRelationvnsRsNodeToMFuncFromFunctionNode(vnsAbsNode.DistinguishedName)
		if err != nil {
			return err
		}
		err = aciClient.CreateRelationvnsRsNodeToMFuncFromFunctionNode(vnsAbsNode.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_vns_rs_node_to_m_func")
		d.Partial(false)

	}
	if d.HasChange("relation_vns_rs_default_scope_to_term") {
		_, newRelParam := d.GetChange("relation_vns_rs_default_scope_to_term")
		err = aciClient.DeleteRelationvnsRsDefaultScopeToTermFromFunctionNode(vnsAbsNode.DistinguishedName)
		if err != nil {
			return err
		}
		err = aciClient.CreateRelationvnsRsDefaultScopeToTermFromFunctionNode(vnsAbsNode.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_vns_rs_default_scope_to_term")
		d.Partial(false)

	}
	if d.HasChange("relation_vns_rs_node_to_cloud_l_dev") {
		_, newRelParam := d.GetChange("relation_vns_rs_node_to_cloud_l_dev")
		err = aciClient.DeleteRelationvnsRsNodeToCloudLDevFromFunctionNode(vnsAbsNode.DistinguishedName)
		if err != nil {
			return err
		}
		err = aciClient.CreateRelationvnsRsNodeToCloudLDevFromFunctionNode(vnsAbsNode.DistinguishedName, newRelParam.(string))
		if err != nil {
			return err
		}
		d.Partial(true)
		d.SetPartial("relation_vns_rs_node_to_cloud_l_dev")
		d.Partial(false)

	}

	d.SetId(vnsAbsNode.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciFunctionNodeRead(d, m)

}

func resourceAciFunctionNodeRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	vnsAbsNode, err := getRemoteFunctionNode(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setFunctionNodeAttributes(vnsAbsNode, d)

	vnsRsNodeToAbsFuncProfData, err := aciClient.ReadRelationvnsRsNodeToAbsFuncProfFromFunctionNode(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vnsRsNodeToAbsFuncProf %v", err)
		d.Set("relation_vns_rs_node_to_abs_func_prof", "")
	} else {
		if _, ok := d.GetOk("relation_vns_rs_node_to_abs_func_prof"); ok {
			tfName := d.Get("relation_vns_rs_node_to_abs_func_prof").(string)
			if tfName != vnsRsNodeToAbsFuncProfData {
				d.Set("relation_vns_rs_node_to_abs_func_prof", "")
			}
		}
	}

	vnsRsNodeToLDevData, err := aciClient.ReadRelationvnsRsNodeToLDevFromFunctionNode(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vnsRsNodeToLDev %v", err)
		d.Set("relation_vns_rs_node_to_l_dev", "")
	} else {
		if _, ok := d.GetOk("relation_vns_rs_node_to_l_dev"); ok {
			tfName := d.Get("relation_vns_rs_node_to_l_dev").(string)
			if tfName != vnsRsNodeToLDevData {
				d.Set("relation_vns_rs_node_to_l_dev", "")
			}
		}
	}

	vnsRsNodeToMFuncData, err := aciClient.ReadRelationvnsRsNodeToMFuncFromFunctionNode(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vnsRsNodeToMFunc %v", err)
		d.Set("relation_vns_rs_node_to_m_func", "")
	} else {
		if _, ok := d.GetOk("relation_vns_rs_node_to_m_func"); ok {
			tfName := d.Get("relation_vns_rs_node_to_m_func").(string)
			if tfName != vnsRsNodeToMFuncData {
				d.Set("relation_vns_rs_node_to_m_func", "")
			}
		}
	}

	vnsRsDefaultScopeToTermData, err := aciClient.ReadRelationvnsRsDefaultScopeToTermFromFunctionNode(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vnsRsDefaultScopeToTerm %v", err)
		d.Set("relation_vns_rs_default_scope_to_term", "")
	} else {
		if _, ok := d.GetOk("relation_vns_rs_default_scope_to_term"); ok {
			tfName := d.Get("relation_vns_rs_default_scope_to_term").(string)
			if tfName != vnsRsDefaultScopeToTermData {
				d.Set("relation_vns_rs_default_scope_to_term", "")
			}
		}
	}

	vnsRsNodeToCloudLDevData, err := aciClient.ReadRelationvnsRsNodeToCloudLDevFromFunctionNode(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vnsRsNodeToCloudLDev %v", err)
		d.Set("relation_vns_rs_node_to_cloud_l_dev", "")
	} else {
		if _, ok := d.GetOk("relation_vns_rs_node_to_cloud_l_dev"); ok {
			tfName := d.Get("relation_vns_rs_node_to_cloud_l_dev").(string)
			if tfName != vnsRsNodeToCloudLDevData {
				d.Set("relation_vns_rs_node_to_cloud_l_dev", "")
			}
		}
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciFunctionNodeDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "vnsAbsNode")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}
