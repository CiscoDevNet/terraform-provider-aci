package aci

import (
	"context"
	"fmt"
	"log"
	"sort"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAciCloudServiceEPg() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciCloudServiceEPgCreate,
		UpdateContext: resourceAciCloudServiceEPgUpdate,
		ReadContext:   resourceAciCloudServiceEPgRead,
		DeleteContext: resourceAciCloudServiceEPgDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciCloudServiceEPgImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"cloud_applicationcontainer_dn": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"access_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Private",
					"Public",
					"PublicAndPrivate",
					"Unknown",
				}, false),
			},
			"azure_private_endpoint": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"custom_service_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"deployment_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"CloudNative",
					"CloudNativeManaged",
					"Third-party",
					"Third-partyManaged",
					"Unknown",
				}, false),
			},
			"exception_tag": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"flood_on_encap": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"disabled",
					"enabled",
				}, false),
			},
			"label_match_criteria": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"All",
					"AtleastOne",
					"AtmostOne",
					"None",
				}, false),
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"preferred_group_member": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"exclude",
					"include",
				}, false),
			},
			"prio": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"level1",
					"level2",
					"level3",
					"level4",
					"level5",
					"level6",
					"unspecified",
				}, false),
			},
			"cloud_service_epg_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Azure-ADDS",
					"Azure-AksCluster",
					"Azure-ApiManagement",
					"Azure-ContainerRegistry",
					"Azure-CosmosDB",
					"Azure-Databricks",
					"Azure-KeyVault",
					"Azure-Redis",
					"Azure-SqlServer",
					"Azure-Storage",
					"Azure-StorageBlob",
					"Azure-StorageFile",
					"Azure-StorageQueue",
					"Azure-StorageTable",
					"Custom",
					"Unknown",
				}, false),
			},

			"relation_cloudrs_cloud_epg_ctx": {
				Type: schema.TypeString,

				Optional:    true,
				Description: "Create relation to fv:Ctx",
			},
			"relation_fvrs_cons": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Create relation to vzBrCP",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"prio": {
							Optional: true,
							Type:     schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								"level1",
								"level2",
								"level3",
								"level4",
								"level5",
								"level6",
								"unspecified",
							}, false),
						},
						"target_dn": {
							Required: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
			"relation_fvrs_cons_if": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Create relation to vzCPIf",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"prio": {
							Optional: true,
							Type:     schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								"level1",
								"level2",
								"level3",
								"level4",
								"level5",
								"level6",
								"unspecified",
							}, false),
						},
						"target_dn": {
							Required: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
			"relation_fvrs_cust_qos_pol": {
				Type: schema.TypeString,

				Optional:    true,
				Description: "Create relation to qos:CustomPol",
			},
			"relation_fvrs_graph_def": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Create relation to vz:GraphCont",
				Set:         schema.HashString,
			},
			"relation_fvrs_intra_epg": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Create relation to vz:BrCP",
				Set:         schema.HashString,
			},
			"relation_fvrs_prot_by": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Create relation to vz:Taboo",
				Set:         schema.HashString,
			},
			"relation_fvrs_prov": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Create relation to vzBrCP",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"label_match_criteria": {
							Optional: true,
							Type:     schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								"All",
								"AtleastOne",
								"AtmostOne",
								"None",
							}, false),
						},
						"prio": {
							Optional: true,
							Type:     schema.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								"level1",
								"level2",
								"level3",
								"level4",
								"level5",
								"level6",
								"unspecified",
							}, false),
						},
						"target_dn": {
							Required: true,
							Type:     schema.TypeString,
						},
					},
				},
			},
			"relation_fvrs_prov_def": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Create relation to vz:CtrctEPgCont",
				Set:         schema.HashString,
			},
			"relation_fvrs_sec_inherited": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Create relation to fv:EPg",
				Set:         schema.HashString,
			}})),
	}
}

func getRemoteCloudServiceEPg(client *client.Client, dn string) (*models.CloudServiceEPg, error) {
	cloudSvcEPgCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	cloudSvcEPg := models.CloudServiceEPgFromContainer(cloudSvcEPgCont)
	if cloudSvcEPg.DistinguishedName == "" {
		return nil, fmt.Errorf("CloudServiceEPg %s not found", dn)
	}
	return cloudSvcEPg, nil
}

func setCloudServiceEPgAttributes(cloudSvcEPg *models.CloudServiceEPg, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(cloudSvcEPg.DistinguishedName)
	d.Set("description", cloudSvcEPg.Description)
	cloudSvcEPgMap, err := cloudSvcEPg.ToMap()
	if err != nil {
		return d, err
	}
	dn := d.Id()
	if dn != cloudSvcEPg.DistinguishedName {
		d.Set("cloud_applicationcontainer_dn", "")
	} else {
		d.Set("cloud_applicationcontainer_dn", GetParentDn(cloudSvcEPg.DistinguishedName, fmt.Sprintf("/"+models.RnCloudSvcEPg, cloudSvcEPgMap["name"])))
	}
	d.Set("access_type", cloudSvcEPgMap["accessType"])
	d.Set("annotation", cloudSvcEPgMap["annotation"])
	d.Set("azure_private_endpoint", cloudSvcEPgMap["azPrivateEndpoint"])
	d.Set("custom_service_type", cloudSvcEPgMap["customSvcType"])
	d.Set("deployment_type", cloudSvcEPgMap["deploymentType"])
	d.Set("exception_tag", cloudSvcEPgMap["exceptionTag"])
	d.Set("flood_on_encap", cloudSvcEPgMap["floodOnEncap"])
	d.Set("label_match_criteria", cloudSvcEPgMap["matchT"])
	d.Set("name", cloudSvcEPgMap["name"])
	d.Set("name_alias", cloudSvcEPgMap["nameAlias"])
	d.Set("preferred_group_member", cloudSvcEPgMap["prefGrMemb"])
	d.Set("prio", cloudSvcEPgMap["prio"])
	d.Set("cloud_service_epg_type", cloudSvcEPgMap["type"])
	return d, nil
}

func getAndSetCloudServiceEPgRelationalAttributes(client *client.Client, dn string, d *schema.ResourceData) {

	log.Printf("[DEBUG] cloudRsCloudEPgCtx: Beginning Read")

	cloudRsCloudEPgCtxData, err := client.ReadRelationcloudRsCloudEPgCtxFromCloudServiceEpg(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation cloudRsCloudEPgCtx %v", err)
		d.Set("relation_cloudrs_cloud_epg_ctx", make([]string, 0, 1))
	} else {
		cloudRsCloudEPgCtxResultData := make([]string, 0, 1)
		for _, obj := range cloudRsCloudEPgCtxData.([]map[string]string) {
			cloudRsCloudEPgCtxResultData = append(cloudRsCloudEPgCtxResultData, obj["tDn"])
		}
		sort.Strings(cloudRsCloudEPgCtxResultData)
		d.Set("relation_cloudrs_cloud_epg_ctx", cloudRsCloudEPgCtxResultData)
		log.Printf("[DEBUG]: cloudRsCloudEPgCtx: Reading finished successfully")
	}

	log.Printf("[DEBUG] fvRsCons: Beginning Read")

	fvRsConsData, err := client.ReadRelationfvRsConsFromCloudServiceEpg(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsCons %v", err)
		d.Set("relation_fvrs_cons", make([]interface{}, 0, 1))
	} else {
		fvRsConsResultData := make([]map[string]string, 0, 1)
		for _, obj := range fvRsConsData.([]map[string]string) {
			fvRsConsResultData = append(fvRsConsResultData, map[string]string{
				"prio":      obj["prio"],
				"target_dn": obj["tDn"],
			})
		}
		d.Set("relation_fvrs_cons", fvRsConsResultData)
		log.Printf("[DEBUG]: fvRsCons: Reading finished successfully")
	}

	log.Printf("[DEBUG] fvRsConsIf: Beginning Read")

	fvRsConsIfData, err := client.ReadRelationfvRsConsIfFromCloudServiceEpg(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsConsIf %v", err)
		d.Set("relation_fvrs_cons_if", make([]interface{}, 0, 1))
	} else {
		fvRsConsIfResultData := make([]map[string]string, 0, 1)
		for _, obj := range fvRsConsIfData.([]map[string]string) {
			fvRsConsIfResultData = append(fvRsConsIfResultData, map[string]string{
				"prio":      obj["prio"],
				"target_dn": obj["tDn"],
			})
		}
		d.Set("relation_fvrs_cons_if", fvRsConsIfResultData)
		log.Printf("[DEBUG]: fvRsConsIf: Reading finished successfully")
	}

	log.Printf("[DEBUG] fvRsCustQosPol: Beginning Read")

	fvRsCustQosPolData, err := client.ReadRelationfvRsCustQosPolFromCloudServiceEpg(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsCustQosPol %v", err)
		d.Set("relation_fvrs_cust_qos_pol", make([]string, 0, 1))
	} else {
		fvRsCustQosPolResultData := make([]string, 0, 1)
		for _, obj := range fvRsCustQosPolData.([]map[string]string) {
			fvRsCustQosPolResultData = append(fvRsCustQosPolResultData, obj["tDn"])
		}
		sort.Strings(fvRsCustQosPolResultData)
		d.Set("relation_fvrs_cust_qos_pol", fvRsCustQosPolResultData)
		log.Printf("[DEBUG]: fvRsCustQosPol: Reading finished successfully")
	}

	log.Printf("[DEBUG] fvRsGraphDef: Beginning Read")
	fvRsGraphDefData, err := client.ReadRelationfvRsGraphDefFromCloudServiceEpg(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsGraphDef %v", err)
		d.Set("relation_fvrs_graph_def", make([]string, 0, 1))
	} else {
		fvRsGraphDefResultData := make([]string, 0, 1)
		for _, obj := range fvRsGraphDefData.([]map[string]string) {
			fvRsGraphDefResultData = append(fvRsGraphDefResultData, obj["tDn"])
		}
		sort.Strings(fvRsGraphDefResultData)
		d.Set("relation_fvrs_graph_def", fvRsGraphDefResultData)
		log.Printf("[DEBUG]: fvRsGraphDef: Reading finished successfully")
	}

	log.Printf("[DEBUG] fvRsIntraEpg: Beginning Read")
	fvRsIntraEpgData, err := client.ReadRelationfvRsIntraEpgFromCloudServiceEpg(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsIntraEpg %v", err)
		d.Set("relation_fvrs_intra_epg", make([]string, 0, 1))
	} else {
		fvRsIntraEpgResultData := make([]string, 0, 1)
		for _, obj := range fvRsIntraEpgData.([]map[string]string) {
			fvRsIntraEpgResultData = append(fvRsIntraEpgResultData, obj["tDn"])
		}
		sort.Strings(fvRsIntraEpgResultData)
		d.Set("relation_fvrs_intra_epg", fvRsIntraEpgResultData)
		log.Printf("[DEBUG]: fvRsIntraEpg: Reading finished successfully")
	}

	log.Printf("[DEBUG] fvRsProtBy: Beginning Read")
	fvRsProtByData, err := client.ReadRelationfvRsProtByFromCloudServiceEpg(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsProtBy %v", err)
		d.Set("relation_fvrs_prot_by", make([]string, 0, 1))
	} else {
		fvRsProtByResultData := make([]string, 0, 1)
		for _, obj := range fvRsProtByData.([]map[string]string) {
			fvRsProtByResultData = append(fvRsProtByResultData, obj["tDn"])
		}
		sort.Strings(fvRsProtByResultData)
		d.Set("relation_fvrs_prot_by", fvRsProtByResultData)
		log.Printf("[DEBUG]: fvRsProtBy: Reading finished successfully")
	}

	log.Printf("[DEBUG] fvRsProv: Beginning Read")

	fvRsProvData, err := client.ReadRelationfvRsProvFromCloudServiceEpg(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsProv %v", err)
		d.Set("relation_fvrs_prov", make([]interface{}, 0, 1))
	} else {
		fvRsProvResultData := make([]map[string]string, 0, 1)
		for _, obj := range fvRsProvData.([]map[string]string) {
			fvRsProvResultData = append(fvRsProvResultData, map[string]string{
				"label_match_criteria": obj["matchT"],
				"prio":                 obj["prio"],
				"target_dn":            obj["tDn"],
			})
		}
		d.Set("relation_fvrs_prov", fvRsProvResultData)
		log.Printf("[DEBUG]: fvRsProv: Reading finished successfully")
	}

	log.Printf("[DEBUG] fvRsProvDef: Beginning Read")
	fvRsProvDefData, err := client.ReadRelationfvRsProvDefFromCloudServiceEpg(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsProvDef %v", err)
		d.Set("relation_fvrs_prov_def", make([]string, 0, 1))
	} else {
		fvRsProvDefResultData := make([]string, 0, 1)
		for _, obj := range fvRsProvDefData.([]map[string]string) {
			fvRsProvDefResultData = append(fvRsProvDefResultData, obj["tDn"])
		}
		sort.Strings(fvRsProvDefResultData)
		d.Set("relation_fvrs_prov_def", fvRsProvDefResultData)
		log.Printf("[DEBUG]: fvRsProvDef: Reading finished successfully")
	}

	log.Printf("[DEBUG] fvRsSecInherited: Beginning Read")
	fvRsSecInheritedData, err := client.ReadRelationfvRsSecInheritedFromCloudServiceEpg(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fvRsSecInherited %v", err)
		d.Set("relation_fvrs_sec_inherited", make([]string, 0, 1))
	} else {
		fvRsSecInheritedResultData := make([]string, 0, 1)
		for _, obj := range fvRsSecInheritedData.([]map[string]string) {
			fvRsSecInheritedResultData = append(fvRsSecInheritedResultData, obj["tDn"])
		}
		sort.Strings(fvRsSecInheritedResultData)
		d.Set("relation_fvrs_sec_inherited", fvRsSecInheritedResultData)
		log.Printf("[DEBUG]: fvRsSecInherited: Reading finished successfully")
	}
}

func resourceAciCloudServiceEPgImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	cloudSvcEPg, err := getRemoteCloudServiceEPg(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setCloudServiceEPgAttributes(cloudSvcEPg, d)
	if err != nil {
		return nil, err
	}

	// Get and Set Relational Attributes
	getAndSetCloudServiceEPgRelationalAttributes(aciClient, dn, d)

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciCloudServiceEPgCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] CloudServiceEPg: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	CloudApplicationcontainerDn := d.Get("cloud_applicationcontainer_dn").(string)

	cloudSvcEPgAttr := models.CloudServiceEPgAttributes{}

	if Annotation, ok := d.GetOk("annotation"); ok {
		cloudSvcEPgAttr.Annotation = Annotation.(string)
	} else {
		cloudSvcEPgAttr.Annotation = "{}"
	}

	if AccessType, ok := d.GetOk("access_type"); ok {
		cloudSvcEPgAttr.AccessType = AccessType.(string)
	}

	if AzPrivateEndpoint, ok := d.GetOk("azure_private_endpoint"); ok {
		cloudSvcEPgAttr.AzPrivateEndpoint = AzPrivateEndpoint.(string)
	}

	if CustomSvcType, ok := d.GetOk("custom_service_type"); ok {
		cloudSvcEPgAttr.CustomSvcType = CustomSvcType.(string)
	}

	if DeploymentType, ok := d.GetOk("deployment_type"); ok {
		cloudSvcEPgAttr.DeploymentType = DeploymentType.(string)
	}

	if ExceptionTag, ok := d.GetOk("exception_tag"); ok {
		cloudSvcEPgAttr.ExceptionTag = ExceptionTag.(string)
	}

	if FloodOnEncap, ok := d.GetOk("flood_on_encap"); ok {
		cloudSvcEPgAttr.FloodOnEncap = FloodOnEncap.(string)
	}

	if MatchT, ok := d.GetOk("label_match_criteria"); ok {
		cloudSvcEPgAttr.MatchT = MatchT.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		cloudSvcEPgAttr.Name = Name.(string)
	}

	if NameAlias, ok := d.GetOk("name_alias"); ok {
		cloudSvcEPgAttr.NameAlias = NameAlias.(string)
	}

	if PrefGrMemb, ok := d.GetOk("preferred_group_member"); ok {
		cloudSvcEPgAttr.PrefGrMemb = PrefGrMemb.(string)
	}

	if Prio, ok := d.GetOk("prio"); ok {
		cloudSvcEPgAttr.Prio = Prio.(string)
	}

	if CloudServiceEPg_type, ok := d.GetOk("cloud_service_epg_type"); ok {
		cloudSvcEPgAttr.CloudServiceEPg_type = CloudServiceEPg_type.(string)
	}
	cloudSvcEPg := models.NewCloudServiceEPg(fmt.Sprintf(models.RnCloudSvcEPg, name), CloudApplicationcontainerDn, desc, cloudSvcEPgAttr)

	err := aciClient.Save(cloudSvcEPg)
	if err != nil {
		return diag.FromErr(err)
	}
	checkDns := make([]string, 0, 1)

	if relationTocloudRsCloudEPgCtx, ok := d.GetOk("relation_cloudrs_cloud_epg_ctx"); ok {
		relationParam := relationTocloudRsCloudEPgCtx.(string)
		checkDns = append(checkDns, relationParam)

	}

	if relationTofvRsCons, ok := d.GetOk("relation_fvrs_cons"); ok {
		relationParamList := toStringList(relationTofvRsCons.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			checkDns = append(checkDns, relationParam)
		}
	}

	if relationTofvRsConsIf, ok := d.GetOk("relation_fvrs_cons_if"); ok {
		relationParamList := toStringList(relationTofvRsConsIf.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			checkDns = append(checkDns, relationParam)
		}
	}

	if relationTofvRsCustQosPol, ok := d.GetOk("relation_fvrs_cust_qos_pol"); ok {
		relationParam := relationTofvRsCustQosPol.(string)
		checkDns = append(checkDns, relationParam)

	}

	if relationTofvRsGraphDef, ok := d.GetOk("relation_fvrs_graph_def"); ok {
		relationParamList := toStringList(relationTofvRsGraphDef.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			checkDns = append(checkDns, relationParam)
		}
	}

	if relationTofvRsIntraEpg, ok := d.GetOk("relation_fvrs_intra_epg"); ok {
		relationParamList := toStringList(relationTofvRsIntraEpg.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			checkDns = append(checkDns, relationParam)
		}
	}

	if relationTofvRsProtBy, ok := d.GetOk("relation_fvrs_prot_by"); ok {
		relationParamList := toStringList(relationTofvRsProtBy.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			checkDns = append(checkDns, relationParam)
		}
	}

	if relationTofvRsProv, ok := d.GetOk("relation_fvrs_prov"); ok {
		relationParamList := toStringList(relationTofvRsProv.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			checkDns = append(checkDns, relationParam)
		}
	}

	if relationTofvRsProvDef, ok := d.GetOk("relation_fvrs_prov_def"); ok {
		relationParamList := toStringList(relationTofvRsProvDef.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			checkDns = append(checkDns, relationParam)
		}
	}

	if relationTofvRsSecInherited, ok := d.GetOk("relation_fvrs_sec_inherited"); ok {
		relationParamList := toStringList(relationTofvRsSecInherited.(*schema.Set).List())
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

	if relationTocloudRsCloudEPgCtx, ok := d.GetOk("relation_cloudrs_cloud_epg_ctx"); ok {
		relationParam := relationTocloudRsCloudEPgCtx.(string)
		err = aciClient.CreateRelationcloudRsCloudEPgCtxFromCloudServiceEpg(cloudSvcEPg.DistinguishedName, cloudSvcEPgAttr.Annotation, GetMOName(relationParam))
		if err != nil {
			return diag.FromErr(err)
		}

	}

	if relationTofvRsCons, ok := d.GetOk("relation_fvrs_cons"); ok {
		relationParamList := relationTofvRsCons.(*schema.Set).List()
		for _, relationParam := range relationParamList {
			paramMap := relationParam.(map[string]interface{})

			err = aciClient.CreateRelationfvRsConsFromCloudServiceEpg(cloudSvcEPg.DistinguishedName, cloudSvcEPgAttr.Annotation, paramMap["prio"].(string), GetMOName(paramMap["target_dn"].(string)))
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if relationTofvRsConsIf, ok := d.GetOk("relation_fvrs_cons_if"); ok {
		relationParamList := relationTofvRsConsIf.(*schema.Set).List()
		for _, relationParam := range relationParamList {
			paramMap := relationParam.(map[string]interface{})

			err = aciClient.CreateRelationfvRsConsIfFromCloudServiceEpg(cloudSvcEPg.DistinguishedName, cloudSvcEPgAttr.Annotation, paramMap["prio"].(string), GetMOName(paramMap["target_dn"].(string)))
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if relationTofvRsCustQosPol, ok := d.GetOk("relation_fvrs_cust_qos_pol"); ok {
		relationParam := relationTofvRsCustQosPol.(string)
		err = aciClient.CreateRelationfvRsCustQosPolFromCloudServiceEpg(cloudSvcEPg.DistinguishedName, cloudSvcEPgAttr.Annotation, GetMOName(relationParam))
		if err != nil {
			return diag.FromErr(err)
		}

	}

	if relationTofvRsGraphDef, ok := d.GetOk("relation_fvrs_graph_def"); ok {
		relationParamList := toStringList(relationTofvRsGraphDef.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			err = aciClient.CreateRelationfvRsGraphDefFromCloudServiceEpg(cloudSvcEPg.DistinguishedName, cloudSvcEPgAttr.Annotation, relationParam)

			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if relationTofvRsIntraEpg, ok := d.GetOk("relation_fvrs_intra_epg"); ok {
		relationParamList := toStringList(relationTofvRsIntraEpg.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			err = aciClient.CreateRelationfvRsIntraEpgFromCloudServiceEpg(cloudSvcEPg.DistinguishedName, cloudSvcEPgAttr.Annotation, GetMOName(relationParam))
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if relationTofvRsProtBy, ok := d.GetOk("relation_fvrs_prot_by"); ok {
		relationParamList := toStringList(relationTofvRsProtBy.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			err = aciClient.CreateRelationfvRsProtByFromCloudServiceEpg(cloudSvcEPg.DistinguishedName, cloudSvcEPgAttr.Annotation, GetMOName(relationParam))
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if relationTofvRsProv, ok := d.GetOk("relation_fvrs_prov"); ok {
		relationParamList := relationTofvRsProv.(*schema.Set).List()
		for _, relationParam := range relationParamList {
			paramMap := relationParam.(map[string]interface{})

			err = aciClient.CreateRelationfvRsProvFromCloudServiceEpg(cloudSvcEPg.DistinguishedName, cloudSvcEPgAttr.Annotation, paramMap["label_match_criteria"].(string), paramMap["prio"].(string), GetMOName(paramMap["target_dn"].(string)))
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if relationTofvRsProvDef, ok := d.GetOk("relation_fvrs_prov_def"); ok {
		relationParamList := toStringList(relationTofvRsProvDef.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			err = aciClient.CreateRelationfvRsProvDefFromCloudServiceEpg(cloudSvcEPg.DistinguishedName, cloudSvcEPgAttr.Annotation, relationParam)

			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if relationTofvRsSecInherited, ok := d.GetOk("relation_fvrs_sec_inherited"); ok {
		relationParamList := toStringList(relationTofvRsSecInherited.(*schema.Set).List())
		for _, relationParam := range relationParamList {
			err = aciClient.CreateRelationfvRsSecInheritedFromCloudServiceEpg(cloudSvcEPg.DistinguishedName, cloudSvcEPgAttr.Annotation, relationParam)

			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	d.SetId(cloudSvcEPg.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciCloudServiceEPgRead(ctx, d, m)
}
func resourceAciCloudServiceEPgUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] CloudServiceEPg: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	CloudApplicationcontainerDn := d.Get("cloud_applicationcontainer_dn").(string)

	cloudSvcEPgAttr := models.CloudServiceEPgAttributes{}

	if Annotation, ok := d.GetOk("annotation"); ok {
		cloudSvcEPgAttr.Annotation = Annotation.(string)
	} else {
		cloudSvcEPgAttr.Annotation = "{}"
	}

	if AccessType, ok := d.GetOk("access_type"); ok {
		cloudSvcEPgAttr.AccessType = AccessType.(string)
	}

	if AzPrivateEndpoint, ok := d.GetOk("azure_private_endpoint"); ok {
		cloudSvcEPgAttr.AzPrivateEndpoint = AzPrivateEndpoint.(string)
	}

	if CustomSvcType, ok := d.GetOk("custom_service_type"); ok {
		cloudSvcEPgAttr.CustomSvcType = CustomSvcType.(string)
	}

	if DeploymentType, ok := d.GetOk("deployment_type"); ok {
		cloudSvcEPgAttr.DeploymentType = DeploymentType.(string)
	}

	if FloodOnEncap, ok := d.GetOk("flood_on_encap"); ok {
		cloudSvcEPgAttr.FloodOnEncap = FloodOnEncap.(string)
	}

	if MatchT, ok := d.GetOk("label_match_criteria"); ok {
		cloudSvcEPgAttr.MatchT = MatchT.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		cloudSvcEPgAttr.Name = Name.(string)
	}

	if NameAlias, ok := d.GetOk("name_alias"); ok {
		cloudSvcEPgAttr.NameAlias = NameAlias.(string)
	}

	if PrefGrMemb, ok := d.GetOk("preferred_group_member"); ok {
		cloudSvcEPgAttr.PrefGrMemb = PrefGrMemb.(string)
	}

	if Prio, ok := d.GetOk("prio"); ok {
		cloudSvcEPgAttr.Prio = Prio.(string)
	}

	if CloudServiceEPg_type, ok := d.GetOk("cloud_service_epg_type"); ok {
		cloudSvcEPgAttr.CloudServiceEPg_type = CloudServiceEPg_type.(string)
	}
	cloudSvcEPg := models.NewCloudServiceEPg(fmt.Sprintf(models.RnCloudSvcEPg, name), CloudApplicationcontainerDn, desc, cloudSvcEPgAttr)

	cloudSvcEPg.Status = "modified"

	err := aciClient.Save(cloudSvcEPg)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_cloudrs_cloud_epg_ctx") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_cloudrs_cloud_epg_ctx")
		checkDns = append(checkDns, newRelParam.(string))

	}

	if d.HasChange("relation_fvrs_cons") || d.HasChange("annotation") {
		oldRel, newRel := d.GetChange("relation_fvrs_cons")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())
		for _, relDn := range relToCreate {
			checkDns = append(checkDns, relDn)
		}
	}

	if d.HasChange("relation_fvrs_cons_if") || d.HasChange("annotation") {
		oldRel, newRel := d.GetChange("relation_fvrs_cons_if")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())
		for _, relDn := range relToCreate {
			checkDns = append(checkDns, relDn)
		}
	}

	if d.HasChange("relation_fvrs_cust_qos_pol") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_fvrs_cust_qos_pol")
		checkDns = append(checkDns, newRelParam.(string))

	}

	if d.HasChange("relation_fvrs_graph_def") || d.HasChange("annotation") {
		oldRel, newRel := d.GetChange("relation_fvrs_graph_def")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())
		for _, relDn := range relToCreate {
			checkDns = append(checkDns, relDn)
		}
	}

	if d.HasChange("relation_fvrs_intra_epg") || d.HasChange("annotation") {
		oldRel, newRel := d.GetChange("relation_fvrs_intra_epg")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())
		for _, relDn := range relToCreate {
			checkDns = append(checkDns, relDn)
		}
	}

	if d.HasChange("relation_fvrs_prot_by") || d.HasChange("annotation") {
		oldRel, newRel := d.GetChange("relation_fvrs_prot_by")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())
		for _, relDn := range relToCreate {
			checkDns = append(checkDns, relDn)
		}
	}

	if d.HasChange("relation_fvrs_prov") || d.HasChange("annotation") {
		oldRel, newRel := d.GetChange("relation_fvrs_prov")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())
		for _, relDn := range relToCreate {
			checkDns = append(checkDns, relDn)
		}
	}

	if d.HasChange("relation_fvrs_prov_def") || d.HasChange("annotation") {
		oldRel, newRel := d.GetChange("relation_fvrs_prov_def")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())
		for _, relDn := range relToCreate {
			checkDns = append(checkDns, relDn)
		}
	}

	if d.HasChange("relation_fvrs_sec_inherited") || d.HasChange("annotation") {
		oldRel, newRel := d.GetChange("relation_fvrs_sec_inherited")
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

	if d.HasChange("relation_cloudrs_cloud_epg_ctx") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_cloudrs_cloud_epg_ctx")
		err = aciClient.DeleteRelationcloudRsCloudEPgCtxFromCloudServiceEpg(cloudSvcEPg.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationcloudRsCloudEPgCtxFromCloudServiceEpg(cloudSvcEPg.DistinguishedName, cloudSvcEPgAttr.Annotation, GetMOName(newRelParam.(string)))
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_fvrs_cons") || d.HasChange("annotation") {
		oldRel, newRel := d.GetChange("relation_fvrs_cons")
		oldRelList := oldRel.(*schema.Set).List()
		newRelList := newRel.(*schema.Set).List()
		for _, relationParam := range oldRelList {
			paramMap := relationParam.(map[string]interface{})

			err = aciClient.DeleteRelationfvRsConsFromCloudServiceEpg(cloudSvcEPg.DistinguishedName, GetMOName(paramMap["target_dn"].(string)))
			if err != nil {
				return diag.FromErr(err)
			}
		}
		for _, relationParam := range newRelList {
			paramMap := relationParam.(map[string]interface{})

			err = aciClient.CreateRelationfvRsConsFromCloudServiceEpg(cloudSvcEPg.DistinguishedName, cloudSvcEPgAttr.Annotation, paramMap["prio"].(string), GetMOName(paramMap["target_dn"].(string)))
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if d.HasChange("relation_fvrs_cons_if") || d.HasChange("annotation") {
		oldRel, newRel := d.GetChange("relation_fvrs_cons_if")
		oldRelList := oldRel.(*schema.Set).List()
		newRelList := newRel.(*schema.Set).List()
		for _, relationParam := range oldRelList {
			paramMap := relationParam.(map[string]interface{})

			err = aciClient.DeleteRelationfvRsConsIfFromCloudServiceEpg(cloudSvcEPg.DistinguishedName, GetMOName(paramMap["target_dn"].(string)))
			if err != nil {
				return diag.FromErr(err)
			}
		}
		for _, relationParam := range newRelList {
			paramMap := relationParam.(map[string]interface{})

			err = aciClient.CreateRelationfvRsConsIfFromCloudServiceEpg(cloudSvcEPg.DistinguishedName, cloudSvcEPgAttr.Annotation, paramMap["prio"].(string), GetMOName(paramMap["target_dn"].(string)))
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if d.HasChange("relation_fvrs_cust_qos_pol") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_fvrs_cust_qos_pol")
		err = aciClient.DeleteRelationfvRsCustQosPolFromCloudServiceEpg(cloudSvcEPg.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationfvRsCustQosPolFromCloudServiceEpg(cloudSvcEPg.DistinguishedName, cloudSvcEPgAttr.Annotation, GetMOName(newRelParam.(string)))
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_fvrs_graph_def") || d.HasChange("annotation") {
		oldRel, newRel := d.GetChange("relation_fvrs_graph_def")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			err = aciClient.DeleteRelationfvRsGraphDefFromCloudServiceEpg(cloudSvcEPg.DistinguishedName, relDn)

			if err != nil {
				return diag.FromErr(err)
			}
		}
		for _, relDn := range relToCreate {
			err = aciClient.CreateRelationfvRsGraphDefFromCloudServiceEpg(cloudSvcEPg.DistinguishedName, cloudSvcEPgAttr.Annotation, relDn)

			if err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if d.HasChange("relation_fvrs_intra_epg") || d.HasChange("annotation") {
		oldRel, newRel := d.GetChange("relation_fvrs_intra_epg")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			err = aciClient.DeleteRelationfvRsIntraEpgFromCloudServiceEpg(cloudSvcEPg.DistinguishedName, GetMOName(relDn))
			if err != nil {
				return diag.FromErr(err)
			}
		}
		for _, relDn := range relToCreate {
			err = aciClient.CreateRelationfvRsIntraEpgFromCloudServiceEpg(cloudSvcEPg.DistinguishedName, cloudSvcEPgAttr.Annotation, GetMOName(relDn))
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if d.HasChange("relation_fvrs_prot_by") || d.HasChange("annotation") {
		oldRel, newRel := d.GetChange("relation_fvrs_prot_by")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			err = aciClient.DeleteRelationfvRsProtByFromCloudServiceEpg(cloudSvcEPg.DistinguishedName, GetMOName(relDn))
			if err != nil {
				return diag.FromErr(err)
			}
		}
		for _, relDn := range relToCreate {
			err = aciClient.CreateRelationfvRsProtByFromCloudServiceEpg(cloudSvcEPg.DistinguishedName, cloudSvcEPgAttr.Annotation, GetMOName(relDn))
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if d.HasChange("relation_fvrs_prov") || d.HasChange("annotation") {
		oldRel, newRel := d.GetChange("relation_fvrs_prov")
		oldRelList := oldRel.(*schema.Set).List()
		newRelList := newRel.(*schema.Set).List()
		for _, relationParam := range oldRelList {
			paramMap := relationParam.(map[string]interface{})

			err = aciClient.DeleteRelationfvRsProvFromCloudServiceEpg(cloudSvcEPg.DistinguishedName, GetMOName(paramMap["target_dn"].(string)))
			if err != nil {
				return diag.FromErr(err)
			}
		}
		for _, relationParam := range newRelList {
			paramMap := relationParam.(map[string]interface{})

			err = aciClient.CreateRelationfvRsProvFromCloudServiceEpg(cloudSvcEPg.DistinguishedName, cloudSvcEPgAttr.Annotation, paramMap["label_match_criteria"].(string), paramMap["prio"].(string), GetMOName(paramMap["target_dn"].(string)))
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if d.HasChange("relation_fvrs_prov_def") || d.HasChange("annotation") {
		oldRel, newRel := d.GetChange("relation_fvrs_prov_def")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			err = aciClient.DeleteRelationfvRsProvDefFromCloudServiceEpg(cloudSvcEPg.DistinguishedName, relDn)

			if err != nil {
				return diag.FromErr(err)
			}
		}
		for _, relDn := range relToCreate {
			err = aciClient.CreateRelationfvRsProvDefFromCloudServiceEpg(cloudSvcEPg.DistinguishedName, cloudSvcEPgAttr.Annotation, relDn)

			if err != nil {
				return diag.FromErr(err)
			}
		}
	}
	if d.HasChange("relation_fvrs_sec_inherited") || d.HasChange("annotation") {
		oldRel, newRel := d.GetChange("relation_fvrs_sec_inherited")
		oldRelSet := oldRel.(*schema.Set)
		newRelSet := newRel.(*schema.Set)
		relToDelete := toStringList(oldRelSet.Difference(newRelSet).List())
		relToCreate := toStringList(newRelSet.Difference(oldRelSet).List())

		for _, relDn := range relToDelete {
			err = aciClient.DeleteRelationfvRsSecInheritedFromCloudServiceEpg(cloudSvcEPg.DistinguishedName, relDn)

			if err != nil {
				return diag.FromErr(err)
			}
		}
		for _, relDn := range relToCreate {
			err = aciClient.CreateRelationfvRsSecInheritedFromCloudServiceEpg(cloudSvcEPg.DistinguishedName, cloudSvcEPgAttr.Annotation, relDn)

			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	d.SetId(cloudSvcEPg.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciCloudServiceEPgRead(ctx, d, m)
}

func resourceAciCloudServiceEPgRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	cloudSvcEPg, err := getRemoteCloudServiceEPg(aciClient, dn)
	if err != nil {
		return errorForObjectNotFound(err, dn, d)
	}

	_, err = setCloudServiceEPgAttributes(cloudSvcEPg, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	// Get and Set Relational Attributes
	getAndSetCloudServiceEPgRelationalAttributes(aciClient, dn, d)

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciCloudServiceEPgDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	err := aciClient.DeleteByDn(dn, "cloudSvcEPg")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}
