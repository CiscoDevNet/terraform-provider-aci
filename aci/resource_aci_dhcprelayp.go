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

func resourceAciDHCPRelayPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciDHCPRelayPolicyCreate,
		UpdateContext: resourceAciDHCPRelayPolicyUpdate,
		ReadContext:   resourceAciDHCPRelayPolicyRead,
		DeleteContext: resourceAciDHCPRelayPolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciDHCPRelayPolicyImport,
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

			"annotation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"mode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"visible",
					"not-visible",
				}, false),
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"owner": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"infra",
					"tenant",
				}, false),
			},

			"relation_dhcp_rs_prov": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tdn": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"addr": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsIPv4Address,
						},
					},
				},
			},
		}),
	}
}
func getRemoteDHCPRelayPolicy(client *client.Client, dn string) (*models.DHCPRelayPolicy, error) {
	dhcpRelayPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	dhcpRelayP := models.DHCPRelayPolicyFromContainer(dhcpRelayPCont)

	if dhcpRelayP.DistinguishedName == "" {
		return nil, fmt.Errorf("DHCPRelayPolicy %s not found", dhcpRelayP.DistinguishedName)
	}

	return dhcpRelayP, nil
}

func setDHCPRelayPolicyAttributes(dhcpRelayP *models.DHCPRelayPolicy, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()

	d.SetId(dhcpRelayP.DistinguishedName)
	d.Set("description", dhcpRelayP.Description)

	if dn != dhcpRelayP.DistinguishedName {
		d.Set("tenant_dn", "")
	}
	dhcpRelayPMap, err := dhcpRelayP.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("tenant_dn", GetParentDn(dn, fmt.Sprintf("/relayp-%s", dhcpRelayPMap["name"])))
	d.Set("name", dhcpRelayPMap["name"])

	d.Set("annotation", dhcpRelayPMap["annotation"])
	d.Set("mode", dhcpRelayPMap["mode"])
	d.Set("name_alias", dhcpRelayPMap["nameAlias"])
	d.Set("owner", dhcpRelayPMap["owner"])
	return d, nil
}

func resourceAciDHCPRelayPolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	dhcpRelayP, err := getRemoteDHCPRelayPolicy(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled, err := setDHCPRelayPolicyAttributes(dhcpRelayP, d)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciDHCPRelayPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] DHCPRelayPolicy: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	dhcpRelayPAttr := models.DHCPRelayPolicyAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		dhcpRelayPAttr.Annotation = Annotation.(string)
	} else {
		dhcpRelayPAttr.Annotation = "{}"
	}
	if Mode, ok := d.GetOk("mode"); ok {
		dhcpRelayPAttr.Mode = Mode.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		dhcpRelayPAttr.NameAlias = NameAlias.(string)
	}
	if Owner, ok := d.GetOk("owner"); ok {
		dhcpRelayPAttr.Owner = Owner.(string)
	}
	dhcpRelayP := models.NewDHCPRelayPolicy(fmt.Sprintf("relayp-%s", name), TenantDn, desc, dhcpRelayPAttr)

	err := aciClient.Save(dhcpRelayP)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if relationTodhcpRsProv, ok := d.GetOk("relation_dhcp_rs_prov"); ok {
		relationParamList := relationTodhcpRsProv.(*schema.Set).List()
		for _, relationParam := range relationParamList {
			paramMap := relationParam.(map[string]interface{})
			checkDns = append(checkDns, paramMap["tdn"].(string))
		}
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if relationTodhcpRsProv, ok := d.GetOk("relation_dhcp_rs_prov"); ok {

		relationParamList := relationTodhcpRsProv.(*schema.Set).List()
		for _, relationParam := range relationParamList {
			paramMap := relationParam.(map[string]interface{})
			err = aciClient.CreateRelationdhcpRsProvFromDHCPRelayPolicy(dhcpRelayP.DistinguishedName, paramMap["tdn"].(string), paramMap["addr"].(string))
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	d.SetId(dhcpRelayP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciDHCPRelayPolicyRead(ctx, d, m)
}

func resourceAciDHCPRelayPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] DHCPRelayPolicy: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	dhcpRelayPAttr := models.DHCPRelayPolicyAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		dhcpRelayPAttr.Annotation = Annotation.(string)
	} else {
		dhcpRelayPAttr.Annotation = "{}"
	}
	if Mode, ok := d.GetOk("mode"); ok {
		dhcpRelayPAttr.Mode = Mode.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		dhcpRelayPAttr.NameAlias = NameAlias.(string)
	}
	if Owner, ok := d.GetOk("owner"); ok {
		dhcpRelayPAttr.Owner = Owner.(string)
	}
	dhcpRelayP := models.NewDHCPRelayPolicy(fmt.Sprintf("relayp-%s", name), TenantDn, desc, dhcpRelayPAttr)

	dhcpRelayP.Status = "modified"

	err := aciClient.Save(dhcpRelayP)

	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_dhcp_rs_prov") {
		newRel := d.Get("relation_dhcp_rs_prov")
		newRelList := newRel.(*schema.Set).List()
		for _, relationParam := range newRelList {
			paramMap := relationParam.(map[string]interface{})
			checkDns = append(checkDns, paramMap["tdn"].(string))
		}

	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if d.HasChange("relation_dhcp_rs_prov") {
		oldRel, newRel := d.GetChange("relation_dhcp_rs_prov")
		oldRelList := oldRel.(*schema.Set).List()
		newRelList := newRel.(*schema.Set).List()

		for _, relationParam := range oldRelList {
			paramMap := relationParam.(map[string]interface{})
			err = aciClient.DeleteRelationdhcpRsProvFromDHCPRelayPolicy(dhcpRelayP.DistinguishedName, paramMap["tdn"].(string))
			if err != nil {
				return diag.FromErr(err)
			}
		}

		for _, relationParam := range newRelList {
			paramMap := relationParam.(map[string]interface{})
			err = aciClient.CreateRelationdhcpRsProvFromDHCPRelayPolicy(dhcpRelayP.DistinguishedName, paramMap["tdn"].(string), paramMap["addr"].(string))
			if err != nil {
				return diag.FromErr(err)
			}
		}

	}
	d.SetId(dhcpRelayP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciDHCPRelayPolicyRead(ctx, d, m)

}

func resourceAciDHCPRelayPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	dhcpRelayP, err := getRemoteDHCPRelayPolicy(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setDHCPRelayPolicyAttributes(dhcpRelayP, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	dhcpRsProvData, err := aciClient.ReadRelationdhcpRsProvFromDHCPRelayPolicy(dn)

	if err != nil {
		log.Printf("[DEBUG] Error while reading relation dhcpRsProv %v", err)
		d.Set("relation_dhcp_rs_prov", make([]map[string]string, 0))

	} else {
		dhcpRsProvMap := dhcpRsProvData.([]map[string]string)
		st := make([]map[string]string, 0)
		for _, dhcpRsProvObj := range dhcpRsProvMap {
			obj := make(map[string]string, 0)
			obj["addr"] = dhcpRsProvObj["addr"]
			obj["tdn"] = dhcpRsProvObj["tDn"]
			st = append(st, obj)
		}
		d.Set("relation_dhcp_rs_prov", st)
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciDHCPRelayPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "dhcpRelayP")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}
