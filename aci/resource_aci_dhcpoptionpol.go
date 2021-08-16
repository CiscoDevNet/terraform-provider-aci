package aci

import (
	"context"
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAciDHCPOptionPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciDHCPOptionPolicyCreate,
		UpdateContext: resourceAciDHCPOptionPolicyUpdate,
		ReadContext:   resourceAciDHCPOptionPolicyRead,
		DeleteContext: resourceAciDHCPOptionPolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciDHCPOptionPolicyImport,
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

			"dhcp_option": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{

						"id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},

						"annotation": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Default:  "orchestrator:terraform",
						},

						"name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},

						"data": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"dhcp_option_id": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"name_alias": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"dhcp_option_ids": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
		}),
	}
}
func getRemoteDHCPOptionPolicy(client *client.Client, dn string) (*models.DHCPOptionPolicy, error) {
	dhcpOptionPolCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	dhcpOptionPol := models.DHCPOptionPolicyFromContainer(dhcpOptionPolCont)

	if dhcpOptionPol.DistinguishedName == "" {
		return nil, fmt.Errorf("DHCPOptionPolicy %s not found", dhcpOptionPol.DistinguishedName)
	}

	return dhcpOptionPol, nil
}

func getRemoteDHCPOptionFromDHCPOptionPolicy(client *client.Client, dn string) (*models.DHCPOption, error) {
	dhcpOptionCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	dhcpOption := models.DHCPOptionFromContainer(dhcpOptionCont)

	if dhcpOption.DistinguishedName == "" {
		return nil, fmt.Errorf("DHCPOption %s not found", dhcpOption.DistinguishedName)
	}

	return dhcpOption, nil
}

func setDHCPOptionPolicyAttributes(dhcpOptionPol *models.DHCPOptionPolicy, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()
	d.SetId(dhcpOptionPol.DistinguishedName)
	d.Set("description", dhcpOptionPol.Description)

	if dn != dhcpOptionPol.DistinguishedName {
		d.Set("tenant_dn", "")
	}

	dhcpOptionPolMap, err := dhcpOptionPol.ToMap()
	if err != nil {
		return d, err
	}

	d.Set("tenant_dn", GetParentDn(dn, fmt.Sprintf("/dhcpoptpol-%s", dhcpOptionPolMap["name"])))

	d.Set("name", dhcpOptionPolMap["name"])

	d.Set("annotation", dhcpOptionPolMap["annotation"])
	d.Set("name_alias", dhcpOptionPolMap["nameAlias"])
	return d, nil
}

func setDHCPOptionAttributesFromDHCPOptionPolicy(dhcpOptions []*models.DHCPOption, d *schema.ResourceData) (*schema.ResourceData, error) {

	dhcpOptionSet := make([]interface{}, 0, 1)
	for _, dhcpOption := range dhcpOptions {

		opMap := make(map[string]interface{})
		opMap["id"] = dhcpOption.DistinguishedName
		dhcpOptionMap, err := dhcpOption.ToMap()
		if err != nil {
			return d, err
		}
		opMap["name"] = dhcpOptionMap["name"]
		opMap["annotation"] = dhcpOptionMap["annotation"]
		opMap["name_alias"] = dhcpOptionMap["nameAlias"]
		opMap["dhcp_option_id"] = dhcpOptionMap["id"]
		opMap["data"] = dhcpOptionMap["data"]
		dhcpOptionSet = append(dhcpOptionSet, opMap)
	}

	d.Set("dhcp_option", dhcpOptionSet)
	return d, nil
}

func resourceAciDHCPOptionPolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	dhcpOptionPol, err := getRemoteDHCPOptionPolicy(aciClient, dn)

	if err != nil {
		return nil, err
	}

	dhcpOptionPolMap, _ := dhcpOptionPol.ToMap()

	name := dhcpOptionPolMap["name"]
	pDN := GetParentDn(dn, fmt.Sprintf("/dhcpoptpol-%s", name))
	d.Set("tenant_dn", pDN)
	schemaFilled, err := setDHCPOptionPolicyAttributes(dhcpOptionPol, d)
	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciDHCPOptionPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] DHCPOptionPolicy: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	dhcpOptionPolAttr := models.DHCPOptionPolicyAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		dhcpOptionPolAttr.Annotation = Annotation.(string)
	} else {
		dhcpOptionPolAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		dhcpOptionPolAttr.NameAlias = NameAlias.(string)
	}
	dhcpOptionPol := models.NewDHCPOptionPolicy(fmt.Sprintf("dhcpoptpol-%s", name), TenantDn, desc, dhcpOptionPolAttr)

	err := aciClient.Save(dhcpOptionPol)
	if err != nil {
		return diag.FromErr(err)
	}

	dhcpOptionIDS := make([]string, 0, 1)
	if options, ok := d.GetOk("dhcp_option"); ok {

		dhcpOptions := options.([]interface{})
		for _, val := range dhcpOptions {
			dhcpOptionAttr := models.DHCPOptionAttributes{}
			dhcpOption := val.(map[string]interface{})

			name := dhcpOption["name"].(string)

			DHCPOptionPolicyDn := dhcpOptionPol.DistinguishedName

			if dhcpOption["annotation"] != nil {
				dhcpOptionAttr.Annotation = dhcpOption["annotation"].(string)
			} else {
				dhcpOptionAttr.Annotation = "{}"
			}
			if dhcpOption["data"] != nil {
				dhcpOptionAttr.Data = dhcpOption["data"].(string)
			}
			if dhcpOption["dhcp_option_id"] != nil {
				dhcpOptionAttr.DHCPOption_id = dhcpOption["dhcp_option_id"].(string)
			}
			if dhcpOption["name_alias"] != nil {
				dhcpOptionAttr.NameAlias = dhcpOption["name_alias"].(string)
			}
			dhcpOptionModel := models.NewDHCPOption(fmt.Sprintf("opt-%s", name), DHCPOptionPolicyDn, dhcpOptionAttr)
			err := aciClient.Save(dhcpOptionModel)
			if err != nil {
				return diag.FromErr(err)
			}
			dhcpOptionIDS = append(dhcpOptionIDS, dhcpOptionModel.DistinguishedName)
		}
		d.Set("dhcp_option_ids", dhcpOptionIDS)
	} else {
		d.Set("dhcp_option_ids", dhcpOptionIDS)
	}

	d.SetId(dhcpOptionPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciDHCPOptionPolicyRead(ctx, d, m)
}

func resourceAciDHCPOptionPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] DHCPOptionPolicy: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	TenantDn := d.Get("tenant_dn").(string)

	dhcpOptionPolAttr := models.DHCPOptionPolicyAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		dhcpOptionPolAttr.Annotation = Annotation.(string)
	} else {
		dhcpOptionPolAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		dhcpOptionPolAttr.NameAlias = NameAlias.(string)
	}
	dhcpOptionPol := models.NewDHCPOptionPolicy(fmt.Sprintf("dhcpoptpol-%s", name), TenantDn, desc, dhcpOptionPolAttr)

	dhcpOptionPol.Status = "modified"

	err := aciClient.Save(dhcpOptionPol)

	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("dhcp_option") {
		dhcpOption := d.Get("dhcp_option_ids").([]interface{})
		for _, val := range dhcpOption {
			dhcpOptionDN := val.(string)
			err := aciClient.DeleteByDn(dhcpOptionDN, "dhcpOption")
			if err != nil {
				return diag.FromErr(err)
			}
		}

		options := d.Get("dhcp_option")
		dhcpOptionIDS := make([]string, 0, 1)

		dhcpOptions := options.([]interface{})
		for _, val := range dhcpOptions {
			dhcpOptionAttr := models.DHCPOptionAttributes{}
			dhcpOption := val.(map[string]interface{})

			name := dhcpOption["name"].(string)

			DHCPOptionPolicyDn := dhcpOptionPol.DistinguishedName

			if dhcpOption["annotation"] != nil {
				dhcpOptionAttr.Annotation = dhcpOption["annotation"].(string)
			} else {
				dhcpOptionAttr.Annotation = "{}"
			}
			if dhcpOption["data"] != nil {
				dhcpOptionAttr.Data = dhcpOption["data"].(string)
			}
			if dhcpOption["dhcp_option_id"] != nil {
				dhcpOptionAttr.DHCPOption_id = dhcpOption["dhcp_option_id"].(string)
			}
			if dhcpOption["name_alias"] != nil {
				dhcpOptionAttr.NameAlias = dhcpOption["name_alias"].(string)
			}
			dhcpOptionModel := models.NewDHCPOption(fmt.Sprintf("opt-%s", name), DHCPOptionPolicyDn, dhcpOptionAttr)
			err := aciClient.Save(dhcpOptionModel)
			if err != nil {
				return diag.FromErr(err)
			}
			dhcpOptionIDS = append(dhcpOptionIDS, dhcpOptionModel.DistinguishedName)
		}

		d.Set("dhcp_option_ids", dhcpOptionIDS)
	}

	d.SetId(dhcpOptionPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciDHCPOptionPolicyRead(ctx, d, m)

}

func resourceAciDHCPOptionPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	dhcpOptionPol, err := getRemoteDHCPOptionPolicy(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setDHCPOptionPolicyAttributes(dhcpOptionPol, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	options := d.Get("dhcp_option_ids").([]interface{})
	dhcpOptions := make([]*models.DHCPOption, 0, 1)

	for _, val := range options {
		dhcpOptionDN := val.(string)
		dhcpOption, err := getRemoteDHCPOption(aciClient, dhcpOptionDN)
		if err != nil {
			return diag.FromErr(err)
		}
		dhcpOptions = append(dhcpOptions, dhcpOption)

	}
	_, err = setDHCPOptionAttributesFromDHCPOptionPolicy(dhcpOptions, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciDHCPOptionPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "dhcpOptionPol")
	if err != nil {
		return diag.FromErr(err)
	}

	options := d.Get("dhcp_option_ids").([]interface{})
	for _, val := range options {
		dhcpOptionDN := val.(string)
		err := aciClient.DeleteByDn(dhcpOptionDN, "dhcpOption")
		if err != nil {
			return diag.FromErr(err)
		}
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}
