package aci

import (
	"context"
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAciErrorDisabledRecoveryPolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext:   dataSourceAciErrorDisabledRecoveryPolicyReadContext,
		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"err_dis_recov_intvl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"edr_event": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
						"event": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"recover": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					})),
				},
			},
		})),
	}
}

func dataSourceAciErrorDisabledRecoveryPolicyReadContext(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	aciClient := m.(*client.Client)
	name := "default"

	rn := fmt.Sprintf("infra/edrErrDisRecoverPol-%s", name)
	dn := fmt.Sprintf("uni/%s", rn)
	edrErrDisRecoverPol, err := getRemoteErrorDisabledRecoveryPolicy(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dn)
	_, err = setErrorDisabledRecoveryPolicyAttributes(edrErrDisRecoverPol, d)
	if err != nil {
		return diag.FromErr(err)
	}
	edrEventIDS, err := getErroDisabledRecoveryChildrenIDS(aciClient, dn)
	if err != nil {
		return diag.FromErr(err)
	}
	events := edrEventIDS
	edrEvents := make([]*models.ErrorDisabledRecoveryEvent, 0, 1)

	for _, val := range events {
		edrEventDn := val
		edrEvent, err := getRemoteErrorDisabledRecoveryEvent(aciClient, edrEventDn)
		if err != nil {
			return diag.FromErr(err)
		}
		edrEvents = append(edrEvents, edrEvent)
	}

	_, err = setErrorDisabledRecoveryEventAttributes(edrEvents, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	return nil
}

func getErroDisabledRecoveryChildrenIDS(aciClient *client.Client, edrErrorDisabledDn string) ([]string, error) {
	edrEventURL := fmt.Sprintf("/api/node/mo/%s.json?query-target=children", edrErrorDisabledDn)
	edrEventCont, err := aciClient.GetViaURL(edrEventURL)
	if err != nil {
		return nil, err
	}
	childCont := edrEventCont.S("imdata")
	edrEventIDS := make([]string, 0, 1)
	for i := 0; i < len(childCont.Data().([]interface{})); i++ {
		edrEventID := G(childCont.Index(i).S("edrEventP", "attributes"), "dn")
		if edrEventID != "{}" {
			edrEventIDS = append(edrEventIDS, edrEventID)
		}
	}
	return edrEventIDS, nil
}
