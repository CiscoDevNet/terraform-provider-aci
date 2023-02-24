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

func resourceAciFabricWideSettingsPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciFabricWideSettingsPolicyCreate,
		UpdateContext: resourceAciFabricWideSettingsPolicyUpdate,
		ReadContext:   resourceAciFabricWideSettingsPolicyRead,
		DeleteContext: resourceAciFabricWideSettingsPolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciFabricWideSettingsPolicyImport,
		},

		SchemaVersion: 1,
		Schema: AppendAttrSchemas(map[string]*schema.Schema{
			"disable_ep_dampening": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
			},
			"domain_validation": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
			},
			"enable_mo_streaming": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
			},
			"enable_remote_leaf_direct": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
			},
			"enforce_subnet_check": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
			},
			"leaf_opflexp_authenticate_clients": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
			},
			"leaf_opflexp_use_ssl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"opflexp_authenticate_clients": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
			},
			"opflexp_ssl_protocols": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"TLSv1",
						"TLSv1.1",
						"TLSv1.2",
					}, false),
				},
				Optional: true,
				Computed: true,
			},
			"opflexp_use_ssl": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
			},
			"policy_sync_node_bringup": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
			},
			"reallocate_gipo": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
			},
			"restrict_infra_vlan_traffic": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
			},
			"unicast_xr_ep_learn_disable": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
			},
			"validate_overlapping_vlans": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"no",
					"yes",
				}, false),
			},
		}, GetBaseAttrSchema(), GetNameAliasAttrSchema()),
	}
}

func getRemoteFabricWideSettingsPolicy(client *client.Client, dn string) (*models.FabricWideSettingsPolicy, error) {
	infraSetPolCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	infraSetPol := models.FabricWideSettingsPolicyFromContainer(infraSetPolCont)
	if infraSetPol.DistinguishedName == "" {
		return nil, fmt.Errorf("Fabric Wide Settings Policy %s not found", dn)
	}
	return infraSetPol, nil
}

func setFabricWideSettingsPolicyAttributes(infraSetPol *models.FabricWideSettingsPolicy, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(infraSetPol.DistinguishedName)
	d.Set("description", infraSetPol.Description)
	infraSetPolMap, err := infraSetPol.ToMap()
	if err != nil {
		return nil, err
	}
	d.Set("annotation", infraSetPolMap["annotation"])
	d.Set("disable_ep_dampening", infraSetPolMap["disableEpDampening"])
	d.Set("domain_validation", infraSetPolMap["domainValidation"])
	d.Set("enable_mo_streaming", infraSetPolMap["enableMoStreaming"])
	d.Set("enable_remote_leaf_direct", infraSetPolMap["enableRemoteLeafDirect"])
	d.Set("enforce_subnet_check", infraSetPolMap["enforceSubnetCheck"])
	d.Set("leaf_opflexp_authenticate_clients", infraSetPolMap["leafOpflexpAuthenticateClients"])
	d.Set("leaf_opflexp_use_ssl", infraSetPolMap["leafOpflexpUseSsl"])
	d.Set("name", infraSetPolMap["name"])
	d.Set("name_alias", infraSetPolMap["nameAlias"])
	d.Set("opflexp_authenticate_clients", infraSetPolMap["opflexpAuthenticateClients"])
	d.Set("opflexp_ssl_protocols", infraSetPolMap["opflexpSslProtocols"])
	d.Set("opflexp_use_ssl", infraSetPolMap["opflexpUseSsl"])
	d.Set("policy_sync_node_bringup", infraSetPolMap["policySyncNodeBringup"])
	d.Set("reallocate_gipo", infraSetPolMap["reallocateGipo"])
	d.Set("restrict_infra_vlan_traffic", infraSetPolMap["restrictInfraVLANTraffic"])
	d.Set("unicast_xr_ep_learn_disable", infraSetPolMap["unicastXrEpLearnDisable"])
	d.Set("validate_overlapping_vlans", infraSetPolMap["validateOverlappingVlans"])
	return d, nil
}

func resourceAciFabricWideSettingsPolicyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	infraSetPol, err := getRemoteFabricWideSettingsPolicy(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setFabricWideSettingsPolicyAttributes(infraSetPol, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciFabricWideSettingsPolicyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] FabricWideSettingsPolicy: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	infraSetPolAttr := models.FabricWideSettingsPolicyAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		infraSetPolAttr.Annotation = Annotation.(string)
	} else {
		infraSetPolAttr.Annotation = "{}"
	}

	if DisableEpDampening, ok := d.GetOk("disable_ep_dampening"); ok {
		infraSetPolAttr.DisableEpDampening = DisableEpDampening.(string)
	}

	if EnableMoStreaming, ok := d.GetOk("enable_mo_streaming"); ok {
		infraSetPolAttr.EnableMoStreaming = EnableMoStreaming.(string)
	}

	if EnableRemoteLeafDirect, ok := d.GetOk("enable_remote_leaf_direct"); ok {
		infraSetPolAttr.EnableRemoteLeafDirect = EnableRemoteLeafDirect.(string)
	}

	if EnforceSubnetCheck, ok := d.GetOk("enforce_subnet_check"); ok {
		infraSetPolAttr.EnforceSubnetCheck = EnforceSubnetCheck.(string)
	}

	if LeafOpflexpAuthenticateClients, ok := d.GetOk("leaf_opflexp_authenticate_clients"); ok {
		infraSetPolAttr.LeafOpflexpAuthenticateClients = LeafOpflexpAuthenticateClients.(string)
	}

	if LeafOpflexpUseSsl, ok := d.GetOk("leaf_opflexp_use_ssl"); ok {
		infraSetPolAttr.LeafOpflexpUseSsl = LeafOpflexpUseSsl.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		infraSetPolAttr.Name = Name.(string)
	}

	if OpflexpAuthenticateClients, ok := d.GetOk("opflexp_authenticate_clients"); ok {
		infraSetPolAttr.OpflexpAuthenticateClients = OpflexpAuthenticateClients.(string)
	}

	if OpflexpSslProtocols, ok := d.GetOk("opflexp_ssl_protocols"); ok {
		OpflexpSslProtocolsList := make([]string, 0, 1)
		for _, value := range OpflexpSslProtocols.([]interface{}) {
			OpflexpSslProtocolsList = append(OpflexpSslProtocolsList, value.(string))
		}
		infraSetPolAttr.OpflexpSslProtocols = strings.Join(OpflexpSslProtocolsList, ",")
	}

	if OpflexpUseSsl, ok := d.GetOk("opflexp_use_ssl"); ok {
		infraSetPolAttr.OpflexpUseSsl = OpflexpUseSsl.(string)
	}

	if RestrictInfraVLANTraffic, ok := d.GetOk("restrict_infra_vlan_traffic"); ok {
		infraSetPolAttr.RestrictInfraVLANTraffic = RestrictInfraVLANTraffic.(string)
	}

	if UnicastXrEpLearnDisable, ok := d.GetOk("unicast_xr_ep_learn_disable"); ok {
		infraSetPolAttr.UnicastXrEpLearnDisable = UnicastXrEpLearnDisable.(string)
	}

	if ValidateOverlappingVlans, ok := d.GetOk("validate_overlapping_vlans"); ok {
		infraSetPolAttr.ValidateOverlappingVlans = ValidateOverlappingVlans.(string)
	}

	if DomainValidation, ok := d.GetOk("domain_validation"); ok {
		infraSetPolAttr.DomainValidation = DomainValidation.(string)
	}

	if ReallocateGipo, ok := d.GetOk("reallocate_gipo"); ok {
		infraSetPolAttr.ReallocateGipo = ReallocateGipo.(string)
	}

	if PolicySyncNodeBringup, ok := d.GetOk("policy_sync_node_bringup"); ok {
		infraSetPolAttr.PolicySyncNodeBringup = PolicySyncNodeBringup.(string)
	}

	infraSetPol := models.NewFabricWideSettingsPolicy(fmt.Sprintf("infra/settings"), "uni", desc, nameAlias, infraSetPolAttr)
	infraSetPol.Status = "modified"
	err := aciClient.Save(infraSetPol)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(infraSetPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciFabricWideSettingsPolicyRead(ctx, d, m)
}

func resourceAciFabricWideSettingsPolicyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] FabricWideSettingsPolicy: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	infraSetPolAttr := models.FabricWideSettingsPolicyAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		infraSetPolAttr.Annotation = Annotation.(string)
	} else {
		infraSetPolAttr.Annotation = "{}"
	}

	if DisableEpDampening, ok := d.GetOk("disable_ep_dampening"); ok {
		infraSetPolAttr.DisableEpDampening = DisableEpDampening.(string)
	}

	if EnableMoStreaming, ok := d.GetOk("enable_mo_streaming"); ok {
		infraSetPolAttr.EnableMoStreaming = EnableMoStreaming.(string)
	}

	if EnableRemoteLeafDirect, ok := d.GetOk("enable_remote_leaf_direct"); ok {
		infraSetPolAttr.EnableRemoteLeafDirect = EnableRemoteLeafDirect.(string)
	}

	if EnforceSubnetCheck, ok := d.GetOk("enforce_subnet_check"); ok {
		infraSetPolAttr.EnforceSubnetCheck = EnforceSubnetCheck.(string)
	}

	if LeafOpflexpAuthenticateClients, ok := d.GetOk("leaf_opflexp_authenticate_clients"); ok {
		infraSetPolAttr.LeafOpflexpAuthenticateClients = LeafOpflexpAuthenticateClients.(string)
	}

	if LeafOpflexpUseSsl, ok := d.GetOk("leaf_opflexp_use_ssl"); ok {
		infraSetPolAttr.LeafOpflexpUseSsl = LeafOpflexpUseSsl.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		infraSetPolAttr.Name = Name.(string)
	}

	if OpflexpAuthenticateClients, ok := d.GetOk("opflexp_authenticate_clients"); ok {
		infraSetPolAttr.OpflexpAuthenticateClients = OpflexpAuthenticateClients.(string)
	}

	if OpflexpSslProtocols, ok := d.GetOk("opflexp_ssl_protocols"); ok {
		OpflexpSslProtocolsList := make([]string, 0, 1)
		for _, value := range OpflexpSslProtocols.([]interface{}) {
			OpflexpSslProtocolsList = append(OpflexpSslProtocolsList, value.(string))
		}
		infraSetPolAttr.OpflexpSslProtocols = strings.Join(OpflexpSslProtocolsList, ",")
	}

	if OpflexpUseSsl, ok := d.GetOk("opflexp_use_ssl"); ok {
		infraSetPolAttr.OpflexpUseSsl = OpflexpUseSsl.(string)
	}

	if RestrictInfraVLANTraffic, ok := d.GetOk("restrict_infra_vlan_traffic"); ok {
		infraSetPolAttr.RestrictInfraVLANTraffic = RestrictInfraVLANTraffic.(string)
	}

	if UnicastXrEpLearnDisable, ok := d.GetOk("unicast_xr_ep_learn_disable"); ok {
		infraSetPolAttr.UnicastXrEpLearnDisable = UnicastXrEpLearnDisable.(string)
	}

	if ValidateOverlappingVlans, ok := d.GetOk("validate_overlapping_vlans"); ok {
		infraSetPolAttr.ValidateOverlappingVlans = ValidateOverlappingVlans.(string)
	}

	if DomainValidation, ok := d.GetOk("domain_validation"); ok {
		infraSetPolAttr.DomainValidation = DomainValidation.(string)
	}

	if ReallocateGipo, ok := d.GetOk("reallocate_gipo"); ok {
		infraSetPolAttr.ReallocateGipo = ReallocateGipo.(string)
	}

	if PolicySyncNodeBringup, ok := d.GetOk("policy_sync_node_bringup"); ok {
		infraSetPolAttr.PolicySyncNodeBringup = PolicySyncNodeBringup.(string)
	}

	infraSetPol := models.NewFabricWideSettingsPolicy(fmt.Sprintf("infra/settings"), "uni", desc, nameAlias, infraSetPolAttr)
	infraSetPol.Status = "modified"
	err := aciClient.Save(infraSetPol)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(infraSetPol.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciFabricWideSettingsPolicyRead(ctx, d, m)
}

func resourceAciFabricWideSettingsPolicyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	infraSetPol, err := getRemoteFabricWideSettingsPolicy(aciClient, dn)
	if err != nil {
		return errorForObjectNotFound(err, dn, d)
	}
	_, err = setFabricWideSettingsPolicyAttributes(infraSetPol, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciFabricWideSettingsPolicyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	d.SetId("")
	var diags diag.Diagnostics
	diags = append(diags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "Resource with class name infraSetPol cannot be deleted",
	})
	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	return diags
}
