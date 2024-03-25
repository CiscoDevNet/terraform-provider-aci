package aci

import (
	"context"
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAciFunctionNode() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciFunctionNodeCreate,
		UpdateContext: resourceAciFunctionNodeUpdate,
		ReadContext:   resourceAciFunctionNodeRead,
		DeleteContext: resourceAciFunctionNodeDelete,

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
				ValidateFunc: validation.StringInSlice([]string{
					"yes",
					"no",
				}, false),
			},
			"conn_copy_dn": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"managed": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"yes",
					"no",
				}, false),
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
					"unspecified",
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
				ValidateFunc: validation.StringInSlice([]string{
					"yes",
					"no",
				}, false),
			},
			"relation_vns_rs_node_to_abs_func_prof": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"relation_vns_rs_node_to_l_dev": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"relation_vns_rs_node_to_m_func": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"relation_vns_rs_default_scope_to_term": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"relation_vns_rs_node_to_cloud_l_dev": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"l4_l7_device_interface_consumer_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"conn_consumer_dn": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"l4_l7_device_interface_consumer_connector_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"none",
					"redir",
				}, false),
			},
			"l4_l7_device_interface_consumer_attachment_notification": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
			},
			"l4_l7_device_interface_provider_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"conn_provider_dn": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"l4_l7_device_interface_provider_connector_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"none",
					"redir",
					"dnat",
					"snat",
					"snat_dnat",
				}, false),
			},
			"l4_l7_device_interface_provider_attachment_notification": &schema.Schema{
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
func getRemoteFunctionNode(client *client.Client, dn string) (*models.FunctionNode, error) {
	vnsAbsNodeCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	vnsAbsNode := models.FunctionNodeFromContainer(vnsAbsNodeCont)

	if vnsAbsNode.DistinguishedName == "" {
		return nil, fmt.Errorf("Function Node %s not found", dn)
	}

	return vnsAbsNode, nil
}

func setFunctionNodeAttributes(vnsAbsNode *models.FunctionNode, d *schema.ResourceData) (*schema.ResourceData, error) {

	dn := d.Id()
	d.SetId(vnsAbsNode.DistinguishedName)
	d.Set("description", vnsAbsNode.Description)

	if dn != vnsAbsNode.DistinguishedName {
		d.Set("l4_l7_service_graph_template_dn", "")
	}

	vnsAbsNodeMap, err := vnsAbsNode.ToMap()
	if err != nil {
		return d, err
	}
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
	return d, nil
}

func getAndSetFunctionNodeRelationalAttributes(client *client.Client, dn string, d *schema.ResourceData) error {
	if d.Get("is_copy").(string) == "yes" { // Copy Part
		d.Set("conn_copy_dn", fmt.Sprintf("%s/AbsFConn-copy", dn))
	} else { // Consumer Part
		consDn := fmt.Sprintf("%s/AbsFConn-consumer", dn)
		vnsAbsFuncConnCont, err := client.Get(consDn)
		if err != nil {
			return err
		}
		vnsAbsFuncConn := models.FunctionConnectorFromContainer(vnsAbsFuncConnCont)
		if vnsAbsFuncConn.DistinguishedName == "" {
			return fmt.Errorf("Function Connector %s not found", consDn)
		}
		d.Set("conn_consumer_dn", vnsAbsFuncConn.DistinguishedName)
		d.Set("l4_l7_device_interface_consumer_name", vnsAbsFuncConn.DeviceLIfName)
		d.Set("l4_l7_device_interface_consumer_connector_type", vnsAbsFuncConn.ConnType)
		d.Set("l4_l7_device_interface_consumer_attachment_notification", vnsAbsFuncConn.AttNotify)

		// Provider Part
		provDn := fmt.Sprintf("%s/AbsFConn-provider", dn)
		vnsAbsFuncConnCont, err = client.Get(provDn)
		if err != nil {
			return err
		}
		vnsAbsFuncConn = models.FunctionConnectorFromContainer(vnsAbsFuncConnCont)
		if vnsAbsFuncConn.DistinguishedName == "" {
			return fmt.Errorf("Function Connector %s not found", provDn)
		}
		d.Set("conn_provider_dn", vnsAbsFuncConn.DistinguishedName)
		d.Set("l4_l7_device_interface_provider_name", vnsAbsFuncConn.DeviceLIfName)
		d.Set("l4_l7_device_interface_provider_connector_type", vnsAbsFuncConn.ConnType)
		d.Set("l4_l7_device_interface_provider_attachment_notification", vnsAbsFuncConn.AttNotify)
	}

	vnsRsNodeToAbsFuncProfData, err := client.ReadRelationvnsRsNodeToAbsFuncProfFromFunctionNode(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vnsRsNodeToAbsFuncProf %v", err)
		d.Set("relation_vns_rs_node_to_abs_func_prof", "")
	} else {
		d.Set("relation_vns_rs_node_to_abs_func_prof", vnsRsNodeToAbsFuncProfData.(string))
	}

	vnsRsNodeToLDevData, err := client.ReadRelationvnsRsNodeToLDevFromFunctionNode(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vnsRsNodeToLDev %v", err)
		d.Set("relation_vns_rs_node_to_l_dev", "")
	} else {
		d.Set("relation_vns_rs_node_to_l_dev", vnsRsNodeToLDevData.(string))
	}

	vnsRsNodeToMFuncData, err := client.ReadRelationvnsRsNodeToMFuncFromFunctionNode(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vnsRsNodeToMFunc %v", err)
		d.Set("relation_vns_rs_node_to_m_func", "")
	} else {
		d.Set("relation_vns_rs_node_to_m_func", vnsRsNodeToMFuncData.(string))
	}

	vnsRsDefaultScopeToTermData, err := client.ReadRelationvnsRsDefaultScopeToTermFromFunctionNode(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vnsRsDefaultScopeToTerm %v", err)
		d.Set("relation_vns_rs_default_scope_to_term", "")
	} else {
		d.Set("relation_vns_rs_default_scope_to_term", vnsRsDefaultScopeToTermData.(string))
	}

	vnsRsNodeToCloudLDevData, err := client.ReadRelationvnsRsNodeToCloudLDevFromFunctionNode(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vnsRsNodeToCloudLDev %v", err)
		d.Set("relation_vns_rs_node_to_cloud_l_dev", "")
	} else {
		d.Set("relation_vns_rs_node_to_cloud_l_dev", vnsRsNodeToCloudLDevData.(string))
	}
	return nil
}

func resourceAciFunctionNodeImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	vnsAbsNode, err := getRemoteFunctionNode(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled, err := setFunctionNodeAttributes(vnsAbsNode, d)
	if err != nil {
		return nil, err
	}

	err = getAndSetFunctionNodeRelationalAttributes(aciClient, dn, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciFunctionNodeCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
		return diag.FromErr(err)
	}

	vnsAbsFuncConnAttr := models.FunctionConnectorAttributes{}
	vnsAbsFuncConnAttr.Annotation = "{}"
	if vnsAbsNodeAttr.IsCopy == "yes" {
		vnsAbsFuncConnAttr.DeviceLIfName = ""
		vnsAbsFuncConn := models.NewFunctionConnector(fmt.Sprintf("AbsFConn-%s", "copy"), vnsAbsNode.DistinguishedName, "", vnsAbsFuncConnAttr)
		err = aciClient.Save(vnsAbsFuncConn)
		if err != nil {
			return diag.FromErr(err)
		}
		d.Set("conn_copy_dn", vnsAbsFuncConn.DistinguishedName)
	} else {
		vnsAbsFuncConnAttr.DeviceLIfName = d.Get("l4_l7_device_interface_consumer_name").(string)
		vnsAbsFuncConnAttr.AttNotify = d.Get("l4_l7_device_interface_consumer_attachment_notification").(string)
		vnsAbsFuncConnAttr.ConnType = d.Get("l4_l7_device_interface_consumer_connector_type").(string)
		vnsAbsFuncConn := models.NewFunctionConnector(fmt.Sprintf("AbsFConn-%s", "consumer"), vnsAbsNode.DistinguishedName, "", vnsAbsFuncConnAttr)
		err = aciClient.Save(vnsAbsFuncConn)
		if err != nil {
			return diag.FromErr(err)
		}
		d.Set("conn_consumer_dn", vnsAbsFuncConn.DistinguishedName)

		vnsAbsFuncConnAttr.DeviceLIfName = d.Get("l4_l7_device_interface_provider_name").(string)
		vnsAbsFuncConnAttr.AttNotify = d.Get("l4_l7_device_interface_provider_attachment_notification").(string)
		vnsAbsFuncConnAttr.ConnType = d.Get("l4_l7_device_interface_provider_connector_type").(string)
		vnsAbsFuncConn = models.NewFunctionConnector(fmt.Sprintf("AbsFConn-%s", "provider"), vnsAbsNode.DistinguishedName, "", vnsAbsFuncConnAttr)
		err = aciClient.Save(vnsAbsFuncConn)
		if err != nil {
			return diag.FromErr(err)
		}
		d.Set("conn_provider_dn", vnsAbsFuncConn.DistinguishedName)
	}

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
		return diag.FromErr(err)
	}
	d.Partial(false)

	if relationTovnsRsNodeToAbsFuncProf, ok := d.GetOk("relation_vns_rs_node_to_abs_func_prof"); ok {
		relationParam := relationTovnsRsNodeToAbsFuncProf.(string)
		err = aciClient.CreateRelationvnsRsNodeToAbsFuncProfFromFunctionNode(vnsAbsNode.DistinguishedName, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if relationTovnsRsNodeToLDev, ok := d.GetOk("relation_vns_rs_node_to_l_dev"); ok {
		relationParam := relationTovnsRsNodeToLDev.(string)
		err = aciClient.CreateRelationvnsRsNodeToLDevFromFunctionNode(vnsAbsNode.DistinguishedName, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if relationTovnsRsNodeToMFunc, ok := d.GetOk("relation_vns_rs_node_to_m_func"); ok {
		relationParam := relationTovnsRsNodeToMFunc.(string)
		err = aciClient.CreateRelationvnsRsNodeToMFuncFromFunctionNode(vnsAbsNode.DistinguishedName, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if relationTovnsRsDefaultScopeToTerm, ok := d.GetOk("relation_vns_rs_default_scope_to_term"); ok {
		relationParam := relationTovnsRsDefaultScopeToTerm.(string)
		err = aciClient.CreateRelationvnsRsDefaultScopeToTermFromFunctionNode(vnsAbsNode.DistinguishedName, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if relationTovnsRsNodeToCloudLDev, ok := d.GetOk("relation_vns_rs_node_to_cloud_l_dev"); ok {
		relationParam := relationTovnsRsNodeToCloudLDev.(string)
		err = aciClient.CreateRelationvnsRsNodeToCloudLDevFromFunctionNode(vnsAbsNode.DistinguishedName, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(vnsAbsNode.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciFunctionNodeRead(ctx, d, m)
}

func resourceAciFunctionNodeUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
		return diag.FromErr(err)
	}

	vnsAbsFuncConnAttr := models.FunctionConnectorAttributes{}
	vnsAbsFuncConnAttr.Annotation = "{}"
	if vnsAbsNodeAttr.IsCopy == "yes" {
		vnsAbsFuncConnAttr.DeviceLIfName = ""
		vnsAbsFuncConn := models.NewFunctionConnector(fmt.Sprintf("AbsFConn-%s", "copy"), vnsAbsNode.DistinguishedName, "", vnsAbsFuncConnAttr)
		err = aciClient.Save(vnsAbsFuncConn)
		if err != nil {
			return diag.FromErr(err)
		}
		d.Set("conn_copy_dn", vnsAbsFuncConn.DistinguishedName)
	}

	if d.HasChange("l4_l7_device_interface_consumer_name") {
		vnsAbsFuncConnAttr.DeviceLIfName = d.Get("l4_l7_device_interface_consumer_name").(string)
		vnsAbsFuncConnAttr.AttNotify = d.Get("l4_l7_device_interface_consumer_attachment_notification").(string)
		vnsAbsFuncConnAttr.ConnType = d.Get("l4_l7_device_interface_consumer_connector_type").(string)
		vnsAbsFuncConn := models.NewFunctionConnector(fmt.Sprintf("AbsFConn-%s", "consumer"), vnsAbsNode.DistinguishedName, "", vnsAbsFuncConnAttr)
		err = aciClient.Save(vnsAbsFuncConn)
		if err != nil {
			return diag.FromErr(err)
		}
		d.Set("conn_consumer_dn", vnsAbsFuncConn.DistinguishedName)
	}

	if d.HasChange("l4_l7_device_interface_provider_name") {
		vnsAbsFuncConnAttr.DeviceLIfName = d.Get("l4_l7_device_interface_provider_name").(string)
		vnsAbsFuncConnAttr.AttNotify = d.Get("l4_l7_device_interface_provider_attachment_notification").(string)
		vnsAbsFuncConnAttr.ConnType = d.Get("l4_l7_device_interface_provider_connector_type").(string)
		vnsAbsFuncConn := models.NewFunctionConnector(fmt.Sprintf("AbsFConn-%s", "provider"), vnsAbsNode.DistinguishedName, "", vnsAbsFuncConnAttr)
		err = aciClient.Save(vnsAbsFuncConn)
		if err != nil {
			return diag.FromErr(err)
		}
		d.Set("conn_provider_dn", vnsAbsFuncConn.DistinguishedName)
	}

	if d.HasChange("conn_consumer_dn") || d.HasChange("conn_provider_dn") || d.HasChange("conn_copy_dn") {
		consOld, _ := d.GetChange("conn_consumer_dn")
		d.Set("conn_consumer_dn", consOld.(string))

		provOld, _ := d.GetChange("conn_provider_dn")
		d.Set("conn_provider_dn", provOld.(string))

		copyOld, _ := d.GetChange("conn_copy_dn")
		d.Set("conn_copy_dn", copyOld.(string))
		return diag.FromErr(fmt.Errorf("conn_consumer_dn, conn_provider_dn and conn_copy_dn are not user configurable"))
	}

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
		return diag.FromErr(err)
	}
	d.Partial(false)

	if d.HasChange("relation_vns_rs_node_to_abs_func_prof") {
		_, newRelParam := d.GetChange("relation_vns_rs_node_to_abs_func_prof")
		err = aciClient.DeleteRelationvnsRsNodeToAbsFuncProfFromFunctionNode(vnsAbsNode.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationvnsRsNodeToAbsFuncProfFromFunctionNode(vnsAbsNode.DistinguishedName, newRelParam.(string))
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_vns_rs_node_to_l_dev") {
		_, newRelParam := d.GetChange("relation_vns_rs_node_to_l_dev")
		err = aciClient.DeleteRelationvnsRsNodeToLDevFromFunctionNode(vnsAbsNode.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationvnsRsNodeToLDevFromFunctionNode(vnsAbsNode.DistinguishedName, newRelParam.(string))
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_vns_rs_node_to_m_func") {
		_, newRelParam := d.GetChange("relation_vns_rs_node_to_m_func")
		err = aciClient.DeleteRelationvnsRsNodeToMFuncFromFunctionNode(vnsAbsNode.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationvnsRsNodeToMFuncFromFunctionNode(vnsAbsNode.DistinguishedName, newRelParam.(string))
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_vns_rs_default_scope_to_term") {
		_, newRelParam := d.GetChange("relation_vns_rs_default_scope_to_term")
		err = aciClient.DeleteRelationvnsRsDefaultScopeToTermFromFunctionNode(vnsAbsNode.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationvnsRsDefaultScopeToTermFromFunctionNode(vnsAbsNode.DistinguishedName, newRelParam.(string))
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_vns_rs_node_to_cloud_l_dev") {
		_, newRelParam := d.GetChange("relation_vns_rs_node_to_cloud_l_dev")
		err = aciClient.DeleteRelationvnsRsNodeToCloudLDevFromFunctionNode(vnsAbsNode.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationvnsRsNodeToCloudLDevFromFunctionNode(vnsAbsNode.DistinguishedName, newRelParam.(string))
		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(vnsAbsNode.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciFunctionNodeRead(ctx, d, m)

}

func resourceAciFunctionNodeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	vnsAbsNode, err := getRemoteFunctionNode(aciClient, dn)

	if err != nil {
		return errorForObjectNotFound(err, dn, d)
	}
	_, err = setFunctionNodeAttributes(vnsAbsNode, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	err = getAndSetFunctionNodeRelationalAttributes(aciClient, dn, d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciFunctionNodeDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "vnsAbsNode")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}
