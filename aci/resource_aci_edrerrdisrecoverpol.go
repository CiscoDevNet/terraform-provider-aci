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

func resourceAciErrorDisabledRecoveryPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciErrorDisabledRecoveryPolicyCreate,
		UpdateContext: resourceAciErrorDisabledRecoveryPolicyUpdate,
		ReadContext:   resourceAciErrorDisabledRecoveryPolicyRead,
		DeleteContext: resourceAciErrorDisabledRecoveryPolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciErrorDisabledRecoveryPolicyImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{

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
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"event-arp-inspection",
								"event-bpduguard",
								"event-debug-1",
								"event-debug-2",
								"event-debug-3",
								"event-debug-4",
								"event-debug-5",
								"event-dhcp-rate-lim",
								"event-ep-move",
								"event-ethpm",
								"event-ip-addr-conflict",
								"event-ipqos-dcbxp-compat-failure",
								"event-ipqos-mgr-error",
								"event-link-flap",
								"event-loopback",
								"event-mcp-loop",
								"event-psec-violation",
								"event-sec-violation",
								"event-set-port-state-failed",
								"event-storm-ctrl",
								"event-stp-inconsist-vpc-peerlink",
								"event-syserr-based",
								"event-udld",
								"unknown",
							}, false),
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
							ValidateFunc: validation.StringInSlice([]string{
								"no",
								"yes",
							}, false),
						},
						"id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					})),
				},
			},
			"edr_event_ids": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
		})),
	}
}

func getRemoteErrorDisabledRecoveryPolicy(client *client.Client, dn string) (*models.ErrorDisabledRecoveryPolicy, error) {
	edrErrDisRecoverPolCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	edrErrDisRecoverPol := models.ErrorDisabledRecoveryPolicyFromContainer(edrErrDisRecoverPolCont)
	if edrErrDisRecoverPol.DistinguishedName == "" {
		return nil, fmt.Errorf("ErrorDisabledRecoveryPolicy %s not found", edrErrDisRecoverPol.DistinguishedName)
	}
	return edrErrDisRecoverPol, nil
}

func getRemoteErrorDisabledRecoveryEvent(client *client.Client, dn string) (*models.ErrorDisabledRecoveryEvent, error) {
	edrEventPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	edrEventP := models.ErrorDisabledRecoveryEventFromContainer(edrEventPCont)
	if edrEventP.DistinguishedName == "" {
		return nil, fmt.Errorf("ErrorDisabledRecoveryEvent %s not found", edrEventP.DistinguishedName)
	}
	return edrEventP, nil
}

func setErrorDisabledRecoveryPolicyAttributes(edrErrDisRecoverPol *models.ErrorDisabledRecoveryPolicy, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(edrErrDisRecoverPol.DistinguishedName)
	d.Set("description", edrErrDisRecoverPol.Description)
	edrErrDisRecoverPolMap, err := edrErrDisRecoverPol.ToMap()
	if err != nil {
		return nil, err
	}
	d.Set("annotation", edrErrDisRecoverPolMap["annotation"])
	d.Set("err_dis_recov_intvl", edrErrDisRecoverPolMap["errDisRecovIntvl"])
	d.Set("name_alias", edrErrDisRecoverPolMap["nameAlias"])
	return d, nil
}

func setErrorDisabledRecoveryEventAttributes(edrEventPs []*models.ErrorDisabledRecoveryEvent, d *schema.ResourceData) (*schema.ResourceData, error) {
	errorDisabledRecoveryEventSet := make([]interface{}, 0, 1)
	for _, edrEvent := range edrEventPs {

		edrMap := make(map[string]interface{})
		edrMap["id"] = edrEvent.DistinguishedName
		edrEventMap, err := edrEvent.ToMap()
		if err != nil {
			return d, err
		}
		edrMap["name"] = edrEventMap["name"]
		edrMap["annotation"] = edrEventMap["annotation"]
		edrMap["name_alias"] = edrEventMap["nameAlias"]
		edrMap["event"] = edrEventMap["event"]
		edrMap["recover"] = edrEventMap["recover"]
		edrMap["description"] = edrEventMap["descr"]

		errorDisabledRecoveryEventSet = append(errorDisabledRecoveryEventSet, edrMap)
	}

	d.Set("edr_event", errorDisabledRecoveryEventSet)
	return d, nil
}
func resourceAciErrorDisabledRecoveryPolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	edrErrDisRecoverPol, err := getRemoteErrorDisabledRecoveryPolicy(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setErrorDisabledRecoveryPolicyAttributes(edrErrDisRecoverPol, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciErrorDisabledRecoveryPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] ErrorDisabledRecoveryPolicy: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := "default"
	edrErrDisRecoverPolAttr := models.ErrorDisabledRecoveryPolicyAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		edrErrDisRecoverPolAttr.Annotation = Annotation.(string)
	} else {
		edrErrDisRecoverPolAttr.Annotation = "{}"
	}

	if ErrDisRecovIntvl, ok := d.GetOk("err_dis_recov_intvl"); ok {
		edrErrDisRecoverPolAttr.ErrDisRecovIntvl = ErrDisRecovIntvl.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		edrErrDisRecoverPolAttr.Name = Name.(string)
	}
	edrErrDisRecoverPol := models.NewErrorDisabledRecoveryPolicy(fmt.Sprintf("infra/edrErrDisRecoverPol-%s", name), "uni", desc, nameAlias, edrErrDisRecoverPolAttr)
	edrErrDisRecoverPol.Status = "modified"
	err := aciClient.Save(edrErrDisRecoverPol)
	if err != nil {
		return diag.FromErr(err)
	}

	edrEventIDS := make([]string, 0, 1)
	if events, ok := d.GetOk("edr_event"); ok {

		edrEvents := events.([]interface{})
		for _, val := range edrEvents {
			edrEventAttr := models.ErrorDisabledRecoveryEventAttributes{}
			edrEvent := val.(map[string]interface{})
			edrEventAttr.Event = edrEvent["event"].(string)
			edrEventID := edrEvent["event"].(string)
			EDRPolicyDn := edrErrDisRecoverPol.DistinguishedName
			nameAlias := ""
			description := ""

			if edrEvent["annotation"] != nil {
				edrEventAttr.Annotation = edrEvent["annotation"].(string)
			} else {
				edrEventAttr.Annotation = "{}"
			}
			if edrEvent["name"] != nil {
				edrEventAttr.Name = edrEvent["name"].(string)
			} else {
				edrEventAttr.Name = ""
			}
			if edrEvent["recover"] != nil {
				edrEventAttr.Recover = edrEvent["recover"].(string)
			} else {
				edrEventAttr.Recover = ""
			}
			if edrEvent["name_alias"] != nil {
				nameAlias = edrEvent["name_alias"].(string)
			}
			if edrEvent["description"] != nil {
				description = edrEvent["description"].(string)
			}

			edrEventModel := models.NewErrorDisabledRecoveryEvent(fmt.Sprintf("edrEventP-%s", edrEventID), EDRPolicyDn, description, nameAlias, edrEventAttr)
			edrEventModel.Status = "modified"
			err := aciClient.Save(edrEventModel)
			if err != nil {
				return diag.FromErr(err)
			}
			edrEventIDS = append(edrEventIDS, edrEventModel.DistinguishedName)
		}

		d.Set("edr_event_ids", edrEventIDS)
	} else {
		d.Set("edr_event_ids", edrEventIDS)
	}

	d.SetId(edrErrDisRecoverPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciErrorDisabledRecoveryPolicyRead(ctx, d, m)
}

func resourceAciErrorDisabledRecoveryPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] ErrorDisabledRecoveryPolicy: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := "default"
	edrErrDisRecoverPolAttr := models.ErrorDisabledRecoveryPolicyAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		edrErrDisRecoverPolAttr.Annotation = Annotation.(string)
	} else {
		edrErrDisRecoverPolAttr.Annotation = "{}"
	}

	if ErrDisRecovIntvl, ok := d.GetOk("err_dis_recov_intvl"); ok {
		edrErrDisRecoverPolAttr.ErrDisRecovIntvl = ErrDisRecovIntvl.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		edrErrDisRecoverPolAttr.Name = Name.(string)
	}
	edrErrDisRecoverPol := models.NewErrorDisabledRecoveryPolicy(fmt.Sprintf("infra/edrErrDisRecoverPol-%s", name), "uni", desc, nameAlias, edrErrDisRecoverPolAttr)
	edrErrDisRecoverPol.Status = "modified"
	err := aciClient.Save(edrErrDisRecoverPol)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("edr_event") {
		edrEvent := d.Get("edr_event_ids").([]interface{})
		for _, val := range edrEvent {
			edrEventDn := val.(string)
			err := aciClient.DeleteByDn(edrEventDn, "edrEventP")
			if err != nil {
				return diag.FromErr(err)
			}
		}

		events := d.Get("edr_event")
		edrEventIDS := make([]string, 0, 1)

		edrEvents := events.([]interface{})
		for _, val := range edrEvents {
			edrEventAttr := models.ErrorDisabledRecoveryEventAttributes{}
			edrEvent := val.(map[string]interface{})
			edrEventAttr.Event = edrEvent["event"].(string)
			edrEventID := edrEvent["event"].(string)
			EDRPolicyDn := edrErrDisRecoverPol.DistinguishedName
			nameAlias := ""
			description := ""

			if edrEvent["annotation"] != nil {
				edrEventAttr.Annotation = edrEvent["annotation"].(string)
			} else {
				edrEventAttr.Annotation = "{}"
			}
			if edrEvent["name"] != nil {
				edrEventAttr.Name = edrEvent["name"].(string)
			} else {
				edrEventAttr.Name = ""
			}
			if edrEvent["recover"] != nil {
				edrEventAttr.Recover = edrEvent["recover"].(string)
			} else {
				edrEventAttr.Recover = ""
			}
			if edrEvent["name_alias"] != nil {
				nameAlias = edrEvent["name_alias"].(string)
			}
			if edrEvent["description"] != nil {
				description = edrEvent["description"].(string)
			}

			edrEventModel := models.NewErrorDisabledRecoveryEvent(fmt.Sprintf("edrEventP-%s", edrEventID), EDRPolicyDn, description, nameAlias, edrEventAttr)
			edrEventModel.Status = "modified"
			err := aciClient.Save(edrEventModel)
			if err != nil {
				return diag.FromErr(err)
			}
			edrEventIDS = append(edrEventIDS, edrEventModel.DistinguishedName)
		}
		d.Set("edr_event_ids", edrEventIDS)
	}
	d.SetId(edrErrDisRecoverPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciErrorDisabledRecoveryPolicyRead(ctx, d, m)
}

func resourceAciErrorDisabledRecoveryPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	edrErrDisRecoverPol, err := getRemoteErrorDisabledRecoveryPolicy(aciClient, dn)
	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setErrorDisabledRecoveryPolicyAttributes(edrErrDisRecoverPol, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	events := d.Get("edr_event_ids").([]interface{})
	edrEvents := make([]*models.ErrorDisabledRecoveryEvent, 0, 1)

	for _, val := range events {
		edrEventDn := val.(string)
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

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciErrorDisabledRecoveryPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	d.SetId("")
	var diags diag.Diagnostics
	diags = append(diags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "Resource with class name edrErrDisRecoverPol cannot be deleted",
	})
	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	return diags
}
