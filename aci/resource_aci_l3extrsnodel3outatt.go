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

func resourceAciFabricNode() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciFabricNodeCreate,
		UpdateContext: resourceAciFabricNodeUpdate,
		ReadContext:   resourceAciFabricNodeRead,
		DeleteContext: resourceAciFabricNodeDelete,

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
					"anchor-node-mismatch",
					"bd-profile-missmatch",
					"missing-mpls-infra-l3out",
					"missing-rs-export-route-profile",
					"node-vlif-misconfig",
					"subnet-mismatch",
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
					"true",
					"false",
				}, false),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if old == "yes" && new == "true" {
						return true
					} else if old == "no" && new == "false" {
						return true
					}
					return false
				},
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

func setFabricNodeAttributes(l3extRsNodeL3OutAtt *models.FabricNode, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()
	d.SetId(l3extRsNodeL3OutAtt.DistinguishedName)
	if dn != l3extRsNodeL3OutAtt.DistinguishedName {
		d.Set("logical_node_profile_dn", "")
	}
	l3extRsNodeL3OutAttMap, err := l3extRsNodeL3OutAtt.ToMap()

	if err != nil {
		return d, err
	}
	d.Set("logical_node_profile_dn", GetParentDn(l3extRsNodeL3OutAtt.DistinguishedName, fmt.Sprintf("/rsnodeL3OutAtt-[%s]", l3extRsNodeL3OutAttMap["tDn"])))

	d.Set("tdn", l3extRsNodeL3OutAttMap["tDn"])

	d.Set("annotation", l3extRsNodeL3OutAttMap["annotation"])
	d.Set("config_issues", l3extRsNodeL3OutAttMap["configIssues"])

	if l3extRsNodeL3OutAttMap["configIssues"] == "" {
		d.Set("config_issues", "none")
	}
	d.Set("rtr_id", l3extRsNodeL3OutAttMap["rtrId"])
	d.Set("rtr_id_loop_back", l3extRsNodeL3OutAttMap["rtrIdLoopBack"])
	return d, nil
}

func resourceAciFabricNodeImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	l3extRsNodeL3OutAtt, err := getRemoteFabricNode(aciClient, dn)

	if err != nil {
		return nil, err
	}
	l3extRsNodeL3OutAttMap, err := l3extRsNodeL3OutAtt.ToMap()

	if err != nil {
		return nil, err
	}

	tdn := l3extRsNodeL3OutAttMap["tDn"]
	pDN := GetParentDn(dn, fmt.Sprintf("/rsnodeL3OutAtt-[%s]", tdn))
	d.Set("logical_node_profile_dn", pDN)
	schemaFilled, err := setFabricNodeAttributes(l3extRsNodeL3OutAtt, d)

	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciFabricNodeCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] FabricNode: Beginning Creation")
	aciClient := m.(*client.Client)

	tdn := d.Get("tdn").(string)

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

	l3extRsNodeL3OutAtt := models.NewFabricNode(fmt.Sprintf("rsnodeL3OutAtt-[%s]", tdn), LogicalNodeProfileDn, l3extRsNodeL3OutAttAttr)

	err := aciClient.Save(l3extRsNodeL3OutAtt)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(l3extRsNodeL3OutAtt.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciFabricNodeRead(ctx, d, m)
}

func resourceAciFabricNodeUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] FabricNode: Beginning Update")

	aciClient := m.(*client.Client)

	tdn := d.Get("tdn").(string)

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

	l3extRsNodeL3OutAtt := models.NewFabricNode(fmt.Sprintf("rsnodeL3OutAtt-[%s]", tdn), LogicalNodeProfileDn, l3extRsNodeL3OutAttAttr)

	l3extRsNodeL3OutAtt.Status = "modified"

	err := aciClient.Save(l3extRsNodeL3OutAtt)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(l3extRsNodeL3OutAtt.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciFabricNodeRead(ctx, d, m)

}

func resourceAciFabricNodeRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	l3extRsNodeL3OutAtt, err := getRemoteFabricNode(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setFabricNodeAttributes(l3extRsNodeL3OutAtt, d)

	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciFabricNodeDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "l3extRsNodeL3OutAtt")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}
