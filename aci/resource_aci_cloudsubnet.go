package aci

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAciCloudSubnet() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciCloudSubnetCreate,
		UpdateContext: resourceAciCloudSubnetUpdate,
		ReadContext:   resourceAciCloudSubnetRead,
		DeleteContext: resourceAciCloudSubnetDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciCloudSubnetImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"cloud_cidr_pool_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"ip": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"scope": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringInSlice([]string{
						"private",
						"public",
						"shared",
					}, false),
				},
			},

			"usage": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"infra-router",
					"user",
					"gateway",
					"transit",
				}, false),
			},

			"zone": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"relation_cloud_rs_subnet_to_flow_log": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
		}),
	}
}
func getRemoteCloudSubnet(client *client.Client, dn string) (*models.CloudSubnet, error) {
	cloudSubnetCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	cloudSubnet := models.CloudSubnetFromContainer(cloudSubnetCont)

	if cloudSubnet.DistinguishedName == "" {
		return nil, fmt.Errorf("CloudSubnet %s not found", cloudSubnet.DistinguishedName)
	}

	return cloudSubnet, nil
}

func setCloudSubnetAttributes(cloudSubnet *models.CloudSubnet, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()
	d.SetId(cloudSubnet.DistinguishedName)
	d.Set("description", cloudSubnet.Description)
	if dn != cloudSubnet.DistinguishedName {
		d.Set("cloud_cidr_pool_dn", "")
	}
	cloudSubnetMap, err := cloudSubnet.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("ip", cloudSubnetMap["ip"])
	d.Set("name", cloudSubnetMap["name"])

	d.Set("annotation", cloudSubnetMap["annotation"])
	d.Set("name_alias", cloudSubnetMap["nameAlias"])
	scopeGet := make([]string, 0, 1)
	for _, val := range strings.Split(cloudSubnetMap["scope"], ",") {
		scopeGet = append(scopeGet, strings.Trim(val, " "))
	}
	sort.Strings(scopeGet)
	if len(scopeGet) == 1 && scopeGet[0] == "" {
		d.Set("scope", make([]string, 0, 1))
	} else {
		d.Set("scope", scopeGet)
	}
	//	d.Set("scope", cloudSubnetMap["scope"])
	d.Set("usage", cloudSubnetMap["usage"])
	return d, nil
}

func resourceAciCloudSubnetImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	cloudSubnet, err := getRemoteCloudSubnet(aciClient, dn)

	if err != nil {
		return nil, err
	}
	cloudSubnetMap, err := cloudSubnet.ToMap()
	if err != nil {
		return nil, err
	}
	ip := cloudSubnetMap["ip"]
	pDN := GetParentDn(dn, fmt.Sprintf("/subnet-[%s]", ip))
	d.Set("cloud_cidr_pool_dn", pDN)
	schemaFilled, err := setCloudSubnetAttributes(cloudSubnet, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciCloudSubnetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] CloudSubnet: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	ip := d.Get("ip").(string)

	CloudCIDRPoolDn := d.Get("cloud_cidr_pool_dn").(string)

	cloudSubnetAttr := models.CloudSubnetAttributes{}
	if name, ok := d.GetOk("name"); ok {
		cloudSubnetAttr.Name = name.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		cloudSubnetAttr.Annotation = Annotation.(string)
	} else {
		cloudSubnetAttr.Annotation = "{}"
	}
	if Ip, ok := d.GetOk("ip"); ok {
		cloudSubnetAttr.Ip = Ip.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		cloudSubnetAttr.NameAlias = NameAlias.(string)
	}
	if Scope, ok := d.GetOk("scope"); ok {
		scopeList := make([]string, 0, 1)
		for _, val := range Scope.([]interface{}) {
			scopeList = append(scopeList, val.(string))
		}
		Scope := strings.Join(scopeList, ",")
		cloudSubnetAttr.Scope = Scope
	}
	if Usage, ok := d.GetOk("usage"); ok {
		cloudSubnetAttr.Usage = Usage.(string)
	}

	checkDns := make([]string, 0, 1)

	if relationTocloudRsZoneAttach, ok := d.GetOk("zone"); ok {
		relationParam := relationTocloudRsZoneAttach.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationTocloudRsSubnetToFlowLog, ok := d.GetOk("relation_cloud_rs_subnet_to_flow_log"); ok {
		relationParam := relationTocloudRsSubnetToFlowLog.(string)
		checkDns = append(checkDns, relationParam)
	}

	d.Partial(true)
	err := checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	var zoneDn string
	if zone, ok := d.GetOk("zone"); ok {
		zoneDn = zone.(string)
	} else {
		zoneDn = ""
	}

	cloudSubnet, err := aciClient.CreateCloudSubnet(ip, CloudCIDRPoolDn, desc, cloudSubnetAttr, zoneDn)
	if err != nil {
		return diag.FromErr(err)
	}

	if relationTocloudRsSubnetToFlowLog, ok := d.GetOk("relation_cloud_rs_subnet_to_flow_log"); ok {
		relationParam := relationTocloudRsSubnetToFlowLog.(string)
		relationParamName := GetMOName(relationParam)
		err = aciClient.CreateRelationcloudRsSubnetToFlowLogFromCloudSubnet(cloudSubnet.DistinguishedName, relationParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(cloudSubnet.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciCloudSubnetRead(ctx, d, m)
}

func resourceAciCloudSubnetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] CloudSubnet: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	ip := d.Get("ip").(string)

	CloudCIDRPoolDn := d.Get("cloud_cidr_pool_dn").(string)

	cloudSubnetAttr := models.CloudSubnetAttributes{}
	if name, ok := d.GetOk("name"); ok {
		cloudSubnetAttr.Name = name.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		cloudSubnetAttr.Annotation = Annotation.(string)
	} else {
		cloudSubnetAttr.Annotation = "{}"
	}
	if Ip, ok := d.GetOk("ip"); ok {
		cloudSubnetAttr.Ip = Ip.(string)
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		cloudSubnetAttr.NameAlias = NameAlias.(string)
	}
	if Scope, ok := d.GetOk("scope"); ok {
		scopeList := make([]string, 0, 1)
		for _, val := range Scope.([]interface{}) {
			scopeList = append(scopeList, val.(string))
		}
		Scope := strings.Join(scopeList, ",")
		cloudSubnetAttr.Scope = Scope
	}
	if Usage, ok := d.GetOk("usage"); ok {
		cloudSubnetAttr.Usage = Usage.(string)
	}

	checkDns := make([]string, 0, 1)

	if d.HasChange("zone") {
		_, newRelParam := d.GetChange("zone")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_cloud_rs_subnet_to_flow_log") {
		_, newRelParam := d.GetChange("relation_cloud_rs_subnet_to_flow_log")
		checkDns = append(checkDns, newRelParam.(string))
	}

	d.Partial(true)
	err := checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	var zoneDn string
	if zone, ok := d.GetOk("zone"); ok {
		zoneDn = zone.(string)
	} else {
		zoneDn = ""
	}

	cloudSubnet, err := aciClient.UpdateCloudSubnet(ip, CloudCIDRPoolDn, desc, cloudSubnetAttr, zoneDn)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("relation_cloud_rs_subnet_to_flow_log") {
		_, newRelParam := d.GetChange("relation_cloud_rs_subnet_to_flow_log")
		newRelParamName := GetMOName(newRelParam.(string))
		err = aciClient.DeleteRelationcloudRsSubnetToFlowLogFromCloudSubnet(cloudSubnet.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationcloudRsSubnetToFlowLogFromCloudSubnet(cloudSubnet.DistinguishedName, newRelParamName)
		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(cloudSubnet.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciCloudSubnetRead(ctx, d, m)

}

func resourceAciCloudSubnetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	cloudSubnet, err := getRemoteCloudSubnet(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setCloudSubnetAttributes(cloudSubnet, d)
	if err != nil {
		d.SetId("")
		return nil
	}
	cloudRsZoneAttachData, err := aciClient.ReadRelationcloudRsZoneAttachFromCloudSubnet(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation cloudRsZoneAttach %v", err)
		d.Set("zone", "")

	} else {
		d.Set("zone", cloudRsZoneAttachData.(string))
	}

	cloudRsSubnetToFlowLogData, err := aciClient.ReadRelationcloudRsSubnetToFlowLogFromCloudSubnet(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation cloudRsSubnetToFlowLog %v", err)
		d.Set("relation_cloud_rs_subnet_to_flow_log", "")

	} else {
		d.Set("relation_cloud_rs_subnet_to_flow_log", cloudRsSubnetToFlowLogData.(string))
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciCloudSubnetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "cloudSubnet")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}
