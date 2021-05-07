package aci

import (
	"fmt"
	"log"
	"reflect"
	"sort"
	"strings"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAciLeafProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceAciLeafProfileCreate,
		Update: resourceAciLeafProfileUpdate,
		Read:   resourceAciLeafProfileRead,
		Delete: resourceAciLeafProfileDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciLeafProfileImport,
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

			"leaf_selector": &schema.Schema{
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

			"leaf_selector_ids": &schema.Schema{
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

			"relation_infra_rs_acc_card_p": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},
			"relation_infra_rs_acc_port_p": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				Set:      schema.HashString,
			},
		}),
	}
}
func getRemoteLeafProfile(client *client.Client, dn string) (*models.LeafProfile, error) {
	infraNodePCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	infraNodeP := models.LeafProfileFromContainer(infraNodePCont)

	if infraNodeP.DistinguishedName == "" {
		return nil, fmt.Errorf("LeafProfile %s not found", infraNodeP.DistinguishedName)
	}

	return infraNodeP, nil
}

func getRemoteSwitchAssociationFromLeafP(client *client.Client, dn string) (*models.SwitchAssociation, error) {
	infraLeafSCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	infraLeafS := models.SwitchAssociationFromContainer(infraLeafSCont)

	if infraLeafS.DistinguishedName == "" {
		return nil, fmt.Errorf("SwitchAssociation %s not found", infraLeafS.DistinguishedName)
	}

	return infraLeafS, nil
}

func getRemoteNodeBlockFromLeafP(client *client.Client, dn string) (*models.NodeBlock, error) {
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

func setLeafProfileAttributes(infraNodeP *models.LeafProfile, d *schema.ResourceData) *schema.ResourceData {
	d.SetId(infraNodeP.DistinguishedName)
	d.Set("description", infraNodeP.Description)
	infraNodePMap, _ := infraNodeP.ToMap()

	d.Set("name", infraNodePMap["name"])

	d.Set("annotation", infraNodePMap["annotation"])
	d.Set("name_alias", infraNodePMap["nameAlias"])
	return d
}

func setLeafSelectorAttributesFromLeafP(selectors []*models.SwitchAssociation, nodeBlocks []*models.NodeBlock, d *schema.ResourceData) *schema.ResourceData {
	selectorSet := make([]interface{}, 0, 1)

	for _, selector := range selectors {
		selMap := make(map[string]interface{})
		selMap["description"] = selector.Description
		selMap["id"] = selector.DistinguishedName

		infraLeafSMap, _ := selector.ToMap()
		selMap["name"] = infraLeafSMap["name"]
		selMap["switch_association_type"] = infraLeafSMap["type"]

		nodeSet := make([]interface{}, 0, 1)
		for _, nodeBlock := range nodeBlocks {
			if strings.Contains(nodeBlock.DistinguishedName, selector.DistinguishedName) {
				nodeBlockMap := setNodeBlockAttributesFromLeafP(nodeBlock)
				nodeSet = append(nodeSet, nodeBlockMap)
			}
		}
		selMap["node_block"] = nodeSet
		selectorSet = append(selectorSet, selMap)
	}
	d.Set("leaf_selector", selectorSet)
	return d
}

func setNodeBlockAttributesFromLeafP(nodeBlock *models.NodeBlock) map[string]interface{} {
	nodeMap := make(map[string]interface{})
	nodeMap["description"] = nodeBlock.Description
	nodeMap["id"] = nodeBlock.DistinguishedName

	infraNodeBlkMap, _ := nodeBlock.ToMap()
	nodeMap["name"] = infraNodeBlkMap["name"]
	nodeMap["from_"] = infraNodeBlkMap["from_"]
	nodeMap["to_"] = infraNodeBlkMap["to_"]
	return nodeMap
}

func resourceAciLeafProfileImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	infraNodeP, err := getRemoteLeafProfile(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled := setLeafProfileAttributes(infraNodeP, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciLeafProfileCreate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] LeafProfile: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	infraNodePAttr := models.LeafProfileAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		infraNodePAttr.Annotation = Annotation.(string)
	} else {
		infraNodePAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		infraNodePAttr.NameAlias = NameAlias.(string)
	}
	infraNodeP := models.NewLeafProfile(fmt.Sprintf("infra/nprof-%s", name), "uni", desc, infraNodePAttr)

	err := aciClient.Save(infraNodeP)
	if err != nil {
		return err
	}

	leafSelectorIDs := make([]string, 0, 1)
	nodeBlockIDs := make([]string, 0, 1)
	if leafSelectors, ok := d.GetOk("leaf_selector"); ok {
		selectors := leafSelectors.([]interface{})
		leafPDN := infraNodeP.DistinguishedName

		for _, val := range selectors {
			selector := val.(map[string]interface{})

			name := selector["name"].(string)
			desc := selector["description"].(string)
			switchAssType := selector["switch_association_type"].(string)

			infraLeafSAttr := models.SwitchAssociationAttributes{}
			infraLeafSAttr.Annotation = "{}"
			infraLeafSAttr.Switch_association_type = switchAssType

			infraLeafS := models.NewSwitchAssociation(fmt.Sprintf("leaves-%s-typ-%s", name, switchAssType), leafPDN, desc, infraLeafSAttr)
			err := aciClient.Save(infraLeafS)
			if err != nil {
				return err
			}
			leafSelectorIDs = append(leafSelectorIDs, infraLeafS.DistinguishedName)

			if selector["node_block"] != nil {
				nodeBlocks := selector["node_block"].([]interface{})
				selectorDN := infraLeafS.DistinguishedName

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
						return err
					}

					nodeBlockIDs = append(nodeBlockIDs, infraNodeBlk.DistinguishedName)
				}
			}
		}
	}
	d.Set("leaf_selector_ids", leafSelectorIDs)
	d.Set("node_block_ids", nodeBlockIDs)

	d.Partial(true)

	d.Partial(false)

	checkDns := make([]string, 0, 1)

	if relationToinfraRsAccCardP, ok := d.GetOk("relation_infra_rs_acc_card_p"); ok {
		relationParamList := toStringList(relationToinfraRsAccCardP.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			checkDns = append(checkDns, relationParam)
		}
	}

	if relationToinfraRsAccPortP, ok := d.GetOk("relation_infra_rs_acc_port_p"); ok {
		relationParamList := toStringList(relationToinfraRsAccPortP.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			checkDns = append(checkDns, relationParam)
		}
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return err
	}
	d.Partial(false)

	if relationToinfraRsAccCardP, ok := d.GetOk("relation_infra_rs_acc_card_p"); ok {
		relationParamList := toStringList(relationToinfraRsAccCardP.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			err = aciClient.CreateRelationinfraRsAccCardPFromLeafProfile(infraNodeP.DistinguishedName, relationParam)

			if err != nil {
				return err
			}
			d.Partial(true)
			d.Partial(false)
		}
	}
	if relationToinfraRsAccPortP, ok := d.GetOk("relation_infra_rs_acc_port_p"); ok {
		relationParamList := toStringList(relationToinfraRsAccPortP.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			err = aciClient.CreateRelationinfraRsAccPortPFromLeafProfile(infraNodeP.DistinguishedName, relationParam)

			if err != nil {
				return err
			}
			d.Partial(true)
			d.Partial(false)
		}
	}

	d.SetId(infraNodeP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciLeafProfileRead(d, m)
}

func resourceAciLeafProfileUpdate(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] LeafProfile: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	infraNodePAttr := models.LeafProfileAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		infraNodePAttr.Annotation = Annotation.(string)
	} else {
		infraNodePAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		infraNodePAttr.NameAlias = NameAlias.(string)
	}
	infraNodeP := models.NewLeafProfile(fmt.Sprintf("infra/nprof-%s", name), "uni", desc, infraNodePAttr)

	infraNodeP.Status = "modified"

	err := aciClient.Save(infraNodeP)

	if err != nil {
		return err
	}

	if d.HasChange("leaf_selector") {
		for _, selectorDn := range d.Get("leaf_selector_ids").([]interface{}) {
			err := aciClient.DeleteByDn(selectorDn.(string), "infraLeafS")
			if err != nil {
				return err
			}
		}

		leafSelectorIDs := make([]string, 0, 1)
		nodeBlockIDs := make([]string, 0, 1)
		if leafSelectors, ok := d.GetOk("leaf_selector"); ok {
			selectors := leafSelectors.([]interface{})
			leafPDN := infraNodeP.DistinguishedName

			for _, val := range selectors {
				selector := val.(map[string]interface{})

				name := selector["name"].(string)
				desc := selector["description"].(string)
				switchAssType := selector["switch_association_type"].(string)

				infraLeafSAttr := models.SwitchAssociationAttributes{}
				infraLeafSAttr.Annotation = "{}"
				infraLeafSAttr.Switch_association_type = switchAssType

				infraLeafS := models.NewSwitchAssociation(fmt.Sprintf("leaves-%s-typ-%s", name, switchAssType), leafPDN, desc, infraLeafSAttr)
				err := aciClient.Save(infraLeafS)
				if err != nil {
					return err
				}
				leafSelectorIDs = append(leafSelectorIDs, infraLeafS.DistinguishedName)

				if selector["node_block"] != nil {
					nodeBlocks := selector["node_block"].([]interface{})
					selectorDN := infraLeafS.DistinguishedName

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
							return err
						}

						nodeBlockIDs = append(nodeBlockIDs, infraNodeBlk.DistinguishedName)
					}
				}
			}
		}
		d.Set("leaf_selector_ids", leafSelectorIDs)
		d.Set("node_block_ids", nodeBlockIDs)
	}

	d.Partial(true)

	d.Partial(false)

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_infra_rs_acc_card_p") {
		oldRel, newRel := d.GetChange("relation_infra_rs_acc_card_p")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToCreate {
			checkDns = append(checkDns, relDn)
		}
	}

	if d.HasChange("relation_infra_rs_acc_port_p") {
		oldRel, newRel := d.GetChange("relation_infra_rs_acc_port_p")
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
		return err
	}
	d.Partial(false)

	if d.HasChange("relation_infra_rs_acc_card_p") {
		oldRel, newRel := d.GetChange("relation_infra_rs_acc_card_p")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			err = aciClient.DeleteRelationinfraRsAccCardPFromLeafProfile(infraNodeP.DistinguishedName, relDn)
			if err != nil {
				return err
			}

		}

		for _, relDn := range relToCreate {
			err = aciClient.CreateRelationinfraRsAccCardPFromLeafProfile(infraNodeP.DistinguishedName, relDn)
			if err != nil {
				return err
			}
			d.Partial(true)
			d.Partial(false)

		}

	}
	if d.HasChange("relation_infra_rs_acc_port_p") {
		oldRel, newRel := d.GetChange("relation_infra_rs_acc_port_p")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			err = aciClient.DeleteRelationinfraRsAccPortPFromLeafProfile(infraNodeP.DistinguishedName, relDn)
			if err != nil {
				return err
			}

		}

		for _, relDn := range relToCreate {
			err = aciClient.CreateRelationinfraRsAccPortPFromLeafProfile(infraNodeP.DistinguishedName, relDn)
			if err != nil {
				return err
			}
			d.Partial(true)
			d.Partial(false)

		}

	}

	d.SetId(infraNodeP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciLeafProfileRead(d, m)

}

func resourceAciLeafProfileRead(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	infraNodeP, err := getRemoteLeafProfile(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	setLeafProfileAttributes(infraNodeP, d)

	leafSelectors := make([]*models.SwitchAssociation, 0, 1)
	nodeBlocks := make([]*models.NodeBlock, 0, 1)
	selectors := d.Get("leaf_selector_ids").([]interface{})
	for _, val := range selectors {
		selectorDn := val.(string)
		selector, err := getRemoteSwitchAssociationFromLeafP(aciClient, selectorDn)
		if err == nil {
			for _, node := range d.Get("node_block_ids").([]interface{}) {
				if strings.Contains(node.(string), selectorDn) {
					nodeBlock, err := getRemoteNodeBlockFromLeafP(aciClient, node.(string))
					if err == nil {
						nodeBlocks = append(nodeBlocks, nodeBlock)
					}
				}
			}
			leafSelectors = append(leafSelectors, selector)
		}
	}
	setLeafSelectorAttributesFromLeafP(leafSelectors, nodeBlocks, d)

	infraRsAccCardPData, err := aciClient.ReadRelationinfraRsAccCardPFromLeafProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsAccCardP %v", err)
		d.Set("relation_infra_rs_acc_card_p", make([]string, 0, 1))

	} else {
		if _, ok := d.GetOk("relation_infra_rs_acc_card_p"); ok {
			relationParamList := toStringList(d.Get("relation_infra_rs_acc_card_p").(*schema.Set).List())
			infraRsAccCardPDataList := toStringList(infraRsAccCardPData.(*schema.Set).List())
			sort.Strings(relationParamList)
			sort.Strings(infraRsAccCardPDataList)
			if !reflect.DeepEqual(relationParamList, infraRsAccCardPDataList) {
				d.Set("relation_infra_rs_acc_card_p", make([]string, 0, 1))
			}
		}
	}

	infraRsAccPortPData, err := aciClient.ReadRelationinfraRsAccPortPFromLeafProfile(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation infraRsAccPortP %v", err)
		d.Set("relation_infra_rs_acc_port_p", make([]string, 0, 1))

	} else {
		if _, ok := d.GetOk("relation_infra_rs_acc_port_p"); ok {
			relationParamList := toStringList(d.Get("relation_infra_rs_acc_port_p").(*schema.Set).List())
			infraRsAccPortPDataList := toStringList(infraRsAccPortPData.(*schema.Set).List())
			sort.Strings(relationParamList)
			sort.Strings(infraRsAccPortPDataList)
			if !reflect.DeepEqual(relationParamList, infraRsAccPortPDataList) {
				d.Set("relation_infra_rs_acc_port_p", make([]string, 0, 1))
			}
		}
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciLeafProfileDelete(d *schema.ResourceData, m interface{}) error {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "infraNodeP")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return err
}
