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

func resourceAciCloudVpnGateway() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciCloudVpnGatewayCreate,
		UpdateContext: resourceAciCloudVpnGatewayUpdate,
		ReadContext:   resourceAciCloudVpnGatewayRead,
		DeleteContext: resourceAciCloudVpnGatewayDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciCloudVpnGatewayImport,
		},

		SchemaVersion: 1,

		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"cloud_context_profile_dn": &schema.Schema{
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

			"name_alias": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"num_instances": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"cloud_router_profile_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"host-router",
					"vpn-gw",
				}, false),
			},

			"relation_cloud_rs_to_vpn_gw_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_cloud_rs_to_direct_conn_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
			"relation_cloud_rs_to_host_router_pol": &schema.Schema{
				Type: schema.TypeString,

				Optional: true,
			},
		}),
	}
}
func getRemoteCloudVpnGateway(client *client.Client, dn string) (*models.CloudVpnGateway, error) {
	cloudRouterPCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}

	cloudRouterP := models.CloudVpnGatewayFromContainer(cloudRouterPCont)

	if cloudRouterP.DistinguishedName == "" {
		return nil, fmt.Errorf("CloudVpnGateway %s not found", cloudRouterP.DistinguishedName)
	}

	return cloudRouterP, nil
}

func setCloudVpnGatewayAttributes(cloudRouterP *models.CloudVpnGateway, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()

	d.SetId(cloudRouterP.DistinguishedName)
	d.Set("description", cloudRouterP.Description)

	if dn != cloudRouterP.DistinguishedName {
		d.Set("cloud_context_profile_dn", "")
	}

	cloudRouterPMap, err := cloudRouterP.ToMap()
	if err != nil {
		return d, err
	}

	d.Set("name", cloudRouterPMap["name"])
	d.Set("cloud_context_profile_dn", GetParentDn(cloudRouterP.DistinguishedName, fmt.Sprintf("/routerp-%s", cloudRouterPMap["name"])))
	d.Set("annotation", cloudRouterPMap["annotation"])
	d.Set("name_alias", cloudRouterPMap["nameAlias"])
	d.Set("num_instances", cloudRouterPMap["numInstances"])
	d.Set("cloud_router_profile_type", cloudRouterPMap["type"])
	return d, nil
}

func resourceAciCloudVpnGatewayImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)

	dn := d.Id()

	cloudRouterP, err := getRemoteCloudVpnGateway(aciClient, dn)

	if err != nil {
		return nil, err
	}
	schemaFilled, err := setCloudVpnGatewayAttributes(cloudRouterP, d)

	if err != nil {
		return nil, err
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())

	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciCloudVpnGatewayCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] CloudVpnGateway: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	CloudContextProfileDn := d.Get("cloud_context_profile_dn").(string)

	cloudRouterPAttr := models.CloudVpnGatewayAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		cloudRouterPAttr.Annotation = Annotation.(string)
	} else {
		cloudRouterPAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		cloudRouterPAttr.NameAlias = NameAlias.(string)
	}
	if NumInstances, ok := d.GetOk("num_instances"); ok {
		cloudRouterPAttr.NumInstances = NumInstances.(string)
	}
	if CloudVpnGateway_type, ok := d.GetOk("cloud_router_profile_type"); ok {
		cloudRouterPAttr.CloudVpnGateway_type = CloudVpnGateway_type.(string)
	}
	cloudRouterP := models.NewCloudVpnGateway(fmt.Sprintf("routerp-%s", name), CloudContextProfileDn, desc, cloudRouterPAttr)

	err := aciClient.Save(cloudRouterP)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if relationTocloudRsToVpnGwPol, ok := d.GetOk("relation_cloud_rs_to_vpn_gw_pol"); ok {
		relationParam := relationTocloudRsToVpnGwPol.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationTocloudRsToDirectConnPol, ok := d.GetOk("relation_cloud_rs_to_direct_conn_pol"); ok {
		relationParam := relationTocloudRsToDirectConnPol.(string)
		checkDns = append(checkDns, relationParam)
	}

	if relationTocloudRsToHostRouterPol, ok := d.GetOk("relation_cloud_rs_to_host_router_pol"); ok {
		relationParam := relationTocloudRsToHostRouterPol.(string)
		checkDns = append(checkDns, relationParam)
	}

	d.Partial(true)
	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	if relationTocloudRsToVpnGwPol, ok := d.GetOk("relation_cloud_rs_to_vpn_gw_pol"); ok {
		relationParam := relationTocloudRsToVpnGwPol.(string)
		err = aciClient.CreateRelationcloudRsToVpnGwPolFromCloudVpnGateway(cloudRouterP.DistinguishedName, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}

	}

	if relationTocloudRsToDirectConnPol, ok := d.GetOk("relation_cloud_rs_to_direct_conn_pol"); ok {
		relationParam := relationTocloudRsToDirectConnPol.(string)
		err = aciClient.CreateRelationcloudRsToDirectConnPolFromCloudVpnGateway(cloudRouterP.DistinguishedName, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if relationTocloudRsToHostRouterPol, ok := d.GetOk("relation_cloud_rs_to_host_router_pol"); ok {
		relationParam := relationTocloudRsToHostRouterPol.(string)
		err = aciClient.CreateRelationcloudRsToHostRouterPolFromCloudVpnGateway(cloudRouterP.DistinguishedName, relationParam)
		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(cloudRouterP.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())

	return resourceAciCloudVpnGatewayRead(ctx, d, m)
}

func resourceAciCloudVpnGatewayUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] CloudVpnGateway: Beginning Update")

	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)

	name := d.Get("name").(string)

	CloudContextProfileDn := d.Get("cloud_context_profile_dn").(string)

	cloudRouterPAttr := models.CloudVpnGatewayAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		cloudRouterPAttr.Annotation = Annotation.(string)
	} else {
		cloudRouterPAttr.Annotation = "{}"
	}
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		cloudRouterPAttr.NameAlias = NameAlias.(string)
	}
	if NumInstances, ok := d.GetOk("num_instances"); ok {
		cloudRouterPAttr.NumInstances = NumInstances.(string)
	}
	if CloudVpnGateway_type, ok := d.GetOk("cloud_router_profile_type"); ok {
		cloudRouterPAttr.CloudVpnGateway_type = CloudVpnGateway_type.(string)
	}
	cloudRouterP := models.NewCloudVpnGateway(fmt.Sprintf("routerp-%s", name), CloudContextProfileDn, desc, cloudRouterPAttr)

	cloudRouterP.Status = "modified"

	err := aciClient.Save(cloudRouterP)

	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_cloud_rs_to_vpn_gw_pol") {
		_, newRelParam := d.GetChange("relation_cloud_rs_to_vpn_gw_pol")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_cloud_rs_to_direct_conn_pol") {
		_, newRelParam := d.GetChange("relation_cloud_rs_to_direct_conn_pol")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_cloud_rs_to_host_router_pol") {
		_, newRelParam := d.GetChange("relation_cloud_rs_to_host_router_pol")
		checkDns = append(checkDns, newRelParam.(string))
	}

	if d.HasChange("relation_cloud_rs_to_vpn_gw_pol") {
		_, newRelParam := d.GetChange("relation_cloud_rs_to_vpn_gw_pol")
		err = aciClient.CreateRelationcloudRsToVpnGwPolFromCloudVpnGateway(cloudRouterP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_cloud_rs_to_direct_conn_pol") {
		_, newRelParam := d.GetChange("relation_cloud_rs_to_direct_conn_pol")
		err = aciClient.CreateRelationcloudRsToDirectConnPolFromCloudVpnGateway(cloudRouterP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("relation_cloud_rs_to_host_router_pol") {
		_, newRelParam := d.GetChange("relation_cloud_rs_to_host_router_pol")
		err = aciClient.CreateRelationcloudRsToHostRouterPolFromCloudVpnGateway(cloudRouterP.DistinguishedName, newRelParam.(string))
		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(cloudRouterP.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())

	return resourceAciCloudVpnGatewayRead(ctx, d, m)

}

func resourceAciCloudVpnGatewayRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())

	aciClient := m.(*client.Client)

	dn := d.Id()
	cloudRouterP, err := getRemoteCloudVpnGateway(aciClient, dn)

	if err != nil {
		d.SetId("")
		return nil
	}
	_, err = setCloudVpnGatewayAttributes(cloudRouterP, d)

	if err != nil {
		d.SetId("")
		return nil
	}

	cloudRsToVpnGwPolData, err := aciClient.ReadRelationcloudRsToVpnGwPolFromCloudVpnGateway(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation cloudRsToVpnGwPol %v", err)
		d.Set("relation_cloud_rs_to_vpn_gw_pol", "")

	} else {
		if _, ok := d.GetOk("relation_cloud_rs_to_vpn_gw_pol"); ok {
			tfName := d.Get("relation_cloud_rs_to_vpn_gw_pol").(string)
			if tfName != cloudRsToVpnGwPolData {
				d.Set("relation_cloud_rs_to_vpn_gw_pol", "")
			}
		}
		//d.Set("relation_cloud_rs_to_vpn_gw_pol", cloudRsToVpnGwPolData)
	}

	cloudRsToDirectConnPolData, err := aciClient.ReadRelationcloudRsToDirectConnPolFromCloudVpnGateway(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation cloudRsToDirectConnPol %v", err)
		d.Set("relation_cloud_rs_to_direct_conn_pol", "")

	} else {
		if _, ok := d.GetOk("relation_cloud_rs_to_direct_conn_pol"); ok {
			tfName := d.Get("relation_cloud_rs_to_direct_conn_pol").(string)
			if tfName != cloudRsToDirectConnPolData {
				d.Set("relation_cloud_rs_to_direct_conn_pol", "")
			}
		}
		//d.Set("relation_cloud_rs_to_direct_conn_pol", cloudRsToDirectConnPolData)
	}

	cloudRsToHostRouterPolData, err := aciClient.ReadRelationcloudRsToHostRouterPolFromCloudVpnGateway(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation cloudRsToHostRouterPol %v", err)
		d.Set("relation_cloud_rs_to_host_router_pol", "")

	} else {
		if _, ok := d.GetOk("relation_cloud_rs_to_host_router_pol"); ok {
			tfName := d.Get("relation_cloud_rs_to_host_router_pol").(string)
			if tfName != cloudRsToHostRouterPolData {
				d.Set("relation_cloud_rs_to_host_router_pol", "")
			}
		}

		//d.Set("relation_cloud_rs_to_host_router_pol", cloudRsToHostRouterPolData)
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())

	return nil
}

func resourceAciCloudVpnGatewayDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())

	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "cloudRouterP")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())

	d.SetId("")
	return diag.FromErr(err)
}
