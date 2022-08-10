package aci

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciLeafProfile() *schema.Resource {
	return &schema.Resource{

		ReadContext: dataSourceAciLeafProfileRead,

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
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"description": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"id": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
						"switch_association_type": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"node_block": &schema.Schema{
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": &schema.Schema{
										Type:     schema.TypeString,
										Required: true,
									},
									"description": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"id": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"from_": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"to_": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
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
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"node_block_ids": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
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

func dataSourceAciLeafProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)

	name := d.Get("name").(string)

	rn := fmt.Sprintf("infra/nprof-%s", name)

	dn := fmt.Sprintf("uni/%s", rn)
	log.Printf("[DEBUG] %s: Data Source - Beginning Read", dn)

	infraNodeP, err := getRemoteLeafProfile(aciClient, dn)

	if err != nil {
		return diag.FromErr(err)
	}
	_, err = setLeafProfileAttributes(infraNodeP, d)
	if err != nil {
		return diag.FromErr(err)
	}

	leafSelectors := make([]*models.SwitchAssociation, 0, 1)
	nodeBlocks := make([]*models.NodeBlock, 0, 1)
	selectors := d.Get("leaf_selector_ids").([]interface{})
	if _, ok := d.GetOk("leaf_selector_ids"); !ok {
		d.Set("leaf_selector_ids", make([]string, 0, 1))
	}
	if _, ok := d.GetOk("node_block_ids"); !ok {
		d.Set("node_block_ids", make([]string, 0, 1))
	}
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
	_, err = setLeafSelectorAttributesFromLeafP(leafSelectors, nodeBlocks, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	// infraRsAccCardP - Beginning Read
	log.Printf("[DEBUG] %s: infraRsAccCardP - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsAccCardPFromLeafProfile(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsAccCardP - Read finished successfully", d.Get("relation_infra_rs_acc_card_p"))
	}
	// infraRsAccCardP - Read finished successfully

	// infraRsAccPortP - Beginning Read
	log.Printf("[DEBUG] %s: infraRsAccPortP - Beginning Read with parent DN", dn)
	_, err = getAndSetReadRelationinfraRsAccPortPFromLeafProfile(aciClient, dn, d)
	if err == nil {
		log.Printf("[DEBUG] %s: infraRsAccPortP - Read finished successfully", d.Get("relation_infra_rs_acc_port_p"))
	}
	// infraRsAccPortP - Read finished successfully

	log.Printf("[DEBUG] %s: Data Source - Read finished successfully", dn)
	return nil
}
