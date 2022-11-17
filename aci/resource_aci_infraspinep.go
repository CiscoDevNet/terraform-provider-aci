package aci

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAciSpineProfile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciSpineProfileCreate,
		UpdateContext: resourceAciSpineProfileUpdate,
		ReadContext:   resourceAciSpineProfileRead,
		DeleteContext: resourceAciSpineProfileDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciSpineProfileImport,
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

			"spine_selector": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},

						"description": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"id": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"switch_association_type": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"ALL",
								"range",
								"ALL_IN_POD",
							}, false),
						},

						"node_block": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": &schema.Schema{
										Type:     schema.TypeString,
										Required: true,
									},

									"description": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},

									"id": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},

									"from_": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},

									"to_": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},

			"spine_selector_ids": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"node_block_ids": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"relation_infra_rs_sp_acc_port_p": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},
		}),
	}
}
func getRemoteSpineProfile(client *client.Client, dn string) (*models.SpineProfile, error) {
	infraSpinePCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	infraSpineP := models.SpineProfileFromContainer(infraSpinePCont)

	if infraSpineP.DistinguishedName == "" {
		return nil, fmt.Errorf("SpineProfile %s not found", infraSpineP.DistinguishedName)
	}

	return infraSpineP, nil
}

func getRemoteSwitchSpineAssociationFromSpineP(client *client.Client, dn string) (*models.SwitchSpineAssociation, error) {
	infraSpineSCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	infraSpineS := models.SwitchSpineAssociationFromContainer(infraSpineSCont)

	if infraSpineS.DistinguishedName == "" {
		return nil, fmt.Errorf("SwitchSpineAssociation %s not found", infraSpineS.DistinguishedName)
	}

	return infraSpineS, nil
}

func getRemoteNodeBlockFromSpineP(client *client.Client, dn string) (*models.NodeBlock, error) {
	infraNodeBlkCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	infraNodeBlk := models.NodeBlockFromContainerBLK(infraNodeBlkCont)

	if infraNodeBlk.DistinguishedName == "" {
		return nil, fmt.Errorf("NodeBlock %s not found", infraNodeBlk.DistinguishedName)
	}

	return infraNodeBlk, nil
}

func setSpineProfileAttributes(infraSpineP *models.SpineProfile, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(infraSpineP.DistinguishedName)
	d.Set("description", infraSpineP.Description)
	infraSpinePMap, err := infraSpineP.ToMap()
	if err != nil {
		return d, err
	}

	d.Set("name", infraSpinePMap["name"])
	d.Set("annotation", infraSpinePMap["annotation"])
	d.Set("name_alias", infraSpinePMap["nameAlias"])
	return d, nil
}

func setSpineSelectorAttributesFromSpineP(selectors []*models.SwitchSpineAssociation, nodeBlocks []*models.NodeBlock, d *schema.ResourceData) (*schema.ResourceData, error) {
	selectorSet := make([]interface{}, 0, 1)

	for _, selector := range selectors {
		selMap := make(map[string]interface{})
		selMap["description"] = selector.Description
		selMap["id"] = selector.DistinguishedName

		infraSpineSMap, err := selector.ToMap()
		if err != nil {
			return d, err
		}
		selMap["name"] = infraSpineSMap["name"]
		selMap["switch_association_type"] = infraSpineSMap["type"]

		nodeSet := make([]interface{}, 0, 1)
		for _, nodeBlock := range nodeBlocks {
			if strings.Contains(nodeBlock.DistinguishedName, selector.DistinguishedName) {
				nodeBlockMap, err := setNodeBlockAttributesFromSpineP(nodeBlock)
				if err != nil {
					return d, err
				}
				nodeSet = append(nodeSet, nodeBlockMap)
			}
		}
		selMap["node_block"] = nodeSet
		selectorSet = append(selectorSet, selMap)
	}
	d.Set("spine_selector", selectorSet)
	return d, nil
}

func setNodeBlockAttributesFromSpineP(nodeBlock *models.NodeBlock) (map[string]interface{}, error) {
	nodeMap := make(map[string]interface{})
	nodeMap["description"] = nodeBlock.Description
	nodeMap["id"] = nodeBlock.DistinguishedName

	infraNodeBlkMap, err := nodeBlock.ToMap()
	if err != nil {
		return nodeMap, err
	}
	nodeMap["name"] = infraNodeBlkMap["name"]
	nodeMap["from_"] = infraNodeBlkMap["from_"]
	nodeMap["to_"] = infraNodeBlkMap["to_"]
	return nodeMap, nil
}

func resourceAciSpineProfileImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	infraSpineP, err := getRemoteSpineProfile(aciClient, dn)
	if err != nil {
		return nil, err
	}

	schemaFilled, err := setSpineProfileAttributes(infraSpineP, d)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciSpineProfileCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] SpineProfile: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	infraSpinePAttr := models.SpineProfileAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		infraSpinePAttr.Annotation = Annotation.(string)
	} else {
		infraSpinePAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		infraSpinePAttr.NameAlias = NameAlias.(string)
	}
	infraSpineP := models.NewSpineProfile(fmt.Sprintf("infra/spprof-%s", name), "uni", desc, infraSpinePAttr)

	err := aciClient.Save(infraSpineP)
	if err != nil {
		return diag.FromErr(err)
	}

	spineSelectorIDs := make([]string, 0, 1)
	nodeBlockIDs := make([]string, 0, 1)
	if spineSelectors, ok := d.GetOk("spine_selector"); ok {
		selectors := spineSelectors.([]interface{})
		spinePDN := infraSpineP.DistinguishedName

		for _, val := range selectors {
			selector := val.(map[string]interface{})

			name := selector["name"].(string)
			desc := selector["description"].(string)
			switchAssType := selector["switch_association_type"].(string)

			infraSpineSAttr := models.SwitchSpineAssociationAttributes{}
			infraSpineSAttr.Annotation = "{}"
			infraSpineSAttr.SwitchAssociationType = switchAssType

			infraSpineS := models.NewSwitchSpineAssociation(fmt.Sprintf("spines-%s-typ-%s", name, switchAssType), spinePDN, desc, infraSpineSAttr)
			err := aciClient.Save(infraSpineS)
			if err != nil {
				return diag.FromErr(err)
			}
			spineSelectorIDs = append(spineSelectorIDs, infraSpineS.DistinguishedName)

			if selector["node_block"] != nil {
				nodeBlocks := selector["node_block"].([]interface{})
				selectorDN := infraSpineS.DistinguishedName

				for _, block := range nodeBlocks {
					nodeBlock := block.(map[string]interface{})

					name := nodeBlock["name"].(string)
					desc := nodeBlock["description"].(string)

					infraNodeBlkAttr := models.NodeBlockAttributes{}
					infraNodeBlkAttr.Annotation = "{}"
					if nodeBlock["from_"] != nil {
						infraNodeBlkAttr.From_ = nodeBlock["from_"].(string)
					}
					if nodeBlock["to_"] != nil {
						infraNodeBlkAttr.To_ = nodeBlock["to_"].(string)
					}
					infraNodeBlk := models.NewNodeBlock(fmt.Sprintf("nodeblk-%s", name), selectorDN, desc, infraNodeBlkAttr)

					err := aciClient.Save(infraNodeBlk)
					if err != nil {
						return diag.FromErr(err)
					}

					nodeBlockIDs = append(nodeBlockIDs, infraNodeBlk.DistinguishedName)
				}
			}
		}
	}
	d.Set("spine_selector_ids", spineSelectorIDs)
	d.Set("node_block_ids", nodeBlockIDs)
	checkDns := make([]string, 0, 1)

	if relationToinfraRsSpAccPortP, ok := d.GetOk("relation_infra_rs_sp_acc_port_p"); ok {
		relationParamList := toStringList(relationToinfraRsSpAccPortP.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			checkDns = append(checkDns, relationParam)
		}
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if relationToinfraRsSpAccPortP, ok := d.GetOk("relation_infra_rs_sp_acc_port_p"); ok {
		relationParamList := toStringList(relationToinfraRsSpAccPortP.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			err = aciClient.CreateRelationinfraRsSpAccPortPFromSpineProfile(infraSpineP.DistinguishedName, relationParam)

			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	d.SetId(infraSpineP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciSpineProfileRead(ctx, d, m)
}

func resourceAciSpineProfileUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] SpineProfile: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	infraSpinePAttr := models.SpineProfileAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		infraSpinePAttr.Annotation = Annotation.(string)
	} else {
		infraSpinePAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		infraSpinePAttr.NameAlias = NameAlias.(string)
	}
	infraSpineP := models.NewSpineProfile(fmt.Sprintf("infra/spprof-%s", name), "uni", desc, infraSpinePAttr)

	infraSpineP.Status = "modified"

	err := aciClient.Save(infraSpineP)

	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("spine_selector") {
		for _, selectorDn := range d.Get("spine_selector_ids").([]interface{}) {
			err := aciClient.DeleteByDn(selectorDn.(string), "infraSpineS")
			if err != nil {
				return diag.FromErr(err)
			}
		}

		spineSelectorIDs := make([]string, 0, 1)
		nodeBlockIDs := make([]string, 0, 1)
		if spineSelectors, ok := d.GetOk("spine_selector"); ok {
			selectors := spineSelectors.([]interface{})
			spinePDN := infraSpineP.DistinguishedName

			for _, val := range selectors {
				selector := val.(map[string]interface{})

				name := selector["name"].(string)
				desc := selector["description"].(string)
				switchAssType := selector["switch_association_type"].(string)

				infraSpineSAttr := models.SwitchSpineAssociationAttributes{}
				infraSpineSAttr.Annotation = "{}"
				infraSpineSAttr.SwitchAssociationType = switchAssType

				infraSpineS := models.NewSwitchSpineAssociation(fmt.Sprintf("spines-%s-typ-%s", name, switchAssType), spinePDN, desc, infraSpineSAttr)
				err := aciClient.Save(infraSpineS)
				if err != nil {
					return diag.FromErr(err)
				}
				spineSelectorIDs = append(spineSelectorIDs, infraSpineS.DistinguishedName)

				if selector["node_block"] != nil {
					nodeBlocks := selector["node_block"].([]interface{})
					selectorDN := infraSpineS.DistinguishedName

					for _, block := range nodeBlocks {
						nodeBlock := block.(map[string]interface{})

						name := nodeBlock["name"].(string)
						desc := nodeBlock["description"].(string)

						infraNodeBlkAttr := models.NodeBlockAttributes{}
						infraNodeBlkAttr.Annotation = "{}"
						if nodeBlock["from_"] != nil {
							infraNodeBlkAttr.From_ = nodeBlock["from_"].(string)
						}
						if nodeBlock["to_"] != nil {
							infraNodeBlkAttr.To_ = nodeBlock["to_"].(string)
						}
						infraNodeBlk := models.NewNodeBlock(fmt.Sprintf("nodeblk-%s", name), selectorDN, desc, infraNodeBlkAttr)

						err := aciClient.Save(infraNodeBlk)
						if err != nil {
							return diag.FromErr(err)
						}

						nodeBlockIDs = append(nodeBlockIDs, infraNodeBlk.DistinguishedName)
					}
				}
			}
		}
		d.Set("spine_selector_ids", spineSelectorIDs)
		d.Set("node_block_ids", nodeBlockIDs)
	}
	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_infra_rs_sp_acc_port_p") {
		oldRel, newRel := d.GetChange("relation_infra_rs_sp_acc_port_p")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToCreate {
			checkDns = append(checkDns, relDn)
		}
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if d.HasChange("relation_infra_rs_sp_acc_port_p") {
		oldRel, newRel := d.GetChange("relation_infra_rs_sp_acc_port_p")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			err = aciClient.DeleteRelationinfraRsSpAccPortPFromSpineProfile(infraSpineP.DistinguishedName, relDn)
			if err != nil {
				return diag.FromErr(err)
			}

		}

		for _, relDn := range relToCreate {
			err = aciClient.CreateRelationinfraRsSpAccPortPFromSpineProfile(infraSpineP.DistinguishedName, relDn)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	d.SetId(infraSpineP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciSpineProfileRead(ctx, d, m)

}

func resourceAciSpineProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	infraSpineP, err := getRemoteSpineProfile(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setSpineProfileAttributes(infraSpineP, d)
	if err != nil {
		d.SetId("")
		return nil
	}
	spineSelectors := make([]*models.SwitchSpineAssociation, 0, 1)
	nodeBlocks := make([]*models.NodeBlock, 0, 1)
	selectors := d.Get("spine_selector_ids").([]interface{})
	if _, ok := d.GetOk("spine_selector_ids"); !ok {
		d.Set("spine_selector_ids", make([]string, 0, 1))
	}
	if _, ok := d.GetOk("node_block_ids"); !ok {
		d.Set("node_block_ids", make([]string, 0, 1))
	}
	for _, val := range selectors {
		selectorDn := val.(string)
		selector, err := getRemoteSwitchSpineAssociationFromSpineP(aciClient, selectorDn)
		if err == nil {
			for _, node := range d.Get("node_block_ids").([]interface{}) {
				if strings.Contains(node.(string), selectorDn) {
					nodeBlock, err := getRemoteNodeBlockFromSpineP(aciClient, node.(string))
					if err == nil {
						nodeBlocks = append(nodeBlocks, nodeBlock)
					}
				}
			}
			spineSelectors = append(spineSelectors, selector)
		}
	}
	_, err = setSpineSelectorAttributesFromSpineP(spineSelectors, nodeBlocks, d)

	if err != nil {
		d.SetId("")
		return nil
	}

	infraRsSpAccPortPData, err := aciClient.ReadRelationinfraRsSpAccPortPFromSpineProfile(dn)
	log.Printf("[TRACE] infraRsSpAccPortP %v", infraRsSpAccPortPData)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsSpAccPortP %v", err)
		setRelationAttribute(d, "relation_infra_rs_sp_acc_port_p", make([]interface{}, 0, 1))
	} else {
		setRelationAttribute(d, "relation_infra_rs_sp_acc_port_p", toStringList(infraRsSpAccPortPData.(*schema.Set).List()))
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciSpineProfileDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "infraSpineP")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}
