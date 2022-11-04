package aci

import (
	"context"
	"fmt"
	"log"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAciSubnetPoolforIpSecTunnels() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciSubnetPoolforIpSecTunnelsCreate,
		UpdateContext: resourceAciSubnetPoolforIpSecTunnelsUpdate,
		ReadContext:   resourceAciSubnetPoolforIpSecTunnelsRead,
		DeleteContext: resourceAciSubnetPoolforIpSecTunnelsDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciSubnetPoolforIpSecTunnelsImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"subnet_pool_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"subnet_pool": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		}),
	}
}

func getRemoteSubnetPoolforIpSecTunnels(client *client.Client, dn string) (*models.SubnetPoolforIpSecTunnels, error) {
	cloudtemplateIpSecTunnelSubnetPoolCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	cloudtemplateIpSecTunnelSubnetPool := models.SubnetPoolforIpSecTunnelsFromContainer(cloudtemplateIpSecTunnelSubnetPoolCont)
	if cloudtemplateIpSecTunnelSubnetPool.DistinguishedName == "" {
		return nil, fmt.Errorf("SubnetPoolforIpSecTunnels %s not found", cloudtemplateIpSecTunnelSubnetPool.DistinguishedName)
	}
	return cloudtemplateIpSecTunnelSubnetPool, nil
}

func setSubnetPoolforIpSecTunnelsAttributes(cloudtemplateIpSecTunnelSubnetPool *models.SubnetPoolforIpSecTunnels, d *schema.ResourceData) (*schema.ResourceData, error) {
	dn := d.Id()
	d.SetId(cloudtemplateIpSecTunnelSubnetPool.DistinguishedName)
	cloudtemplateIpSecTunnelSubnetPoolMap, err := cloudtemplateIpSecTunnelSubnetPool.ToMap()
	if err != nil {
		return d, err
	}

	if dn != cloudtemplateIpSecTunnelSubnetPool.DistinguishedName {
		d.Set("infra_network_template_dn", "")
	} else {
		d.Set("infra_network_template_dn", GetParentDn(dn, "/"+fmt.Sprintf(models.RncloudtemplateIpSecTunnelSubnetPool, cloudtemplateIpSecTunnelSubnetPoolMap["subnetpool"])))
	}
	d.Set("annotation", cloudtemplateIpSecTunnelSubnetPoolMap["annotation"])
	d.Set("subnet_pool_name", cloudtemplateIpSecTunnelSubnetPoolMap["poolname"])
	d.Set("subnet_pool", cloudtemplateIpSecTunnelSubnetPoolMap["subnetpool"])
	return d, nil
}

func resourceAciSubnetPoolforIpSecTunnelsImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	cloudtemplateIpSecTunnelSubnetPool, err := getRemoteSubnetPoolforIpSecTunnels(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setSubnetPoolforIpSecTunnelsAttributes(cloudtemplateIpSecTunnelSubnetPool, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciSubnetPoolforIpSecTunnelsCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] SubnetPoolforIpSecTunnels: Beginning Creation")
	aciClient := m.(*client.Client)
	CloudInfraNetworkTemplateDn := "uni/tn-infra/infranetwork-default"

	cloudtemplateIpSecTunnelSubnetPoolAttr := models.SubnetPoolforIpSecTunnelsAttributes{}

	if Annotation, ok := d.GetOk("annotation"); ok {
		cloudtemplateIpSecTunnelSubnetPoolAttr.Annotation = Annotation.(string)
	} else {
		cloudtemplateIpSecTunnelSubnetPoolAttr.Annotation = "{}"
	}

	if Poolname, ok := d.GetOk("subnet_pool_name"); ok {
		cloudtemplateIpSecTunnelSubnetPoolAttr.Poolname = Poolname.(string)
	}

	if Subnetpool, ok := d.GetOk("subnet_pool"); ok {
		cloudtemplateIpSecTunnelSubnetPoolAttr.Subnetpool = Subnetpool.(string)
	}
	cloudtemplateIpSecTunnelSubnetPool := models.NewSubnetPoolforIpSecTunnels(fmt.Sprintf(models.RncloudtemplateIpSecTunnelSubnetPool, cloudtemplateIpSecTunnelSubnetPoolAttr.Subnetpool), CloudInfraNetworkTemplateDn, cloudtemplateIpSecTunnelSubnetPoolAttr)

	err := aciClient.Save(cloudtemplateIpSecTunnelSubnetPool)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(cloudtemplateIpSecTunnelSubnetPool.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciSubnetPoolforIpSecTunnelsRead(ctx, d, m)
}

func resourceAciSubnetPoolforIpSecTunnelsUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] SubnetPoolforIpSecTunnels: Beginning Update")
	aciClient := m.(*client.Client)
	CloudInfraNetworkTemplateDn := "uni/tn-infra/infranetwork-default"

	cloudtemplateIpSecTunnelSubnetPoolAttr := models.SubnetPoolforIpSecTunnelsAttributes{}

	if Annotation, ok := d.GetOk("annotation"); ok {
		cloudtemplateIpSecTunnelSubnetPoolAttr.Annotation = Annotation.(string)
	} else {
		cloudtemplateIpSecTunnelSubnetPoolAttr.Annotation = "{}"
	}

	if Poolname, ok := d.GetOk("subnet_pool_name"); ok {
		cloudtemplateIpSecTunnelSubnetPoolAttr.Poolname = Poolname.(string)
	}

	if Subnetpool, ok := d.GetOk("subnet_pool"); ok {
		cloudtemplateIpSecTunnelSubnetPoolAttr.Subnetpool = Subnetpool.(string)
	}
	cloudtemplateIpSecTunnelSubnetPool := models.NewSubnetPoolforIpSecTunnels(fmt.Sprintf("ipsecsubnetpool-[%s]", cloudtemplateIpSecTunnelSubnetPoolAttr.Subnetpool), CloudInfraNetworkTemplateDn, cloudtemplateIpSecTunnelSubnetPoolAttr)

	cloudtemplateIpSecTunnelSubnetPool.Status = "modified"

	err := aciClient.Save(cloudtemplateIpSecTunnelSubnetPool)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(cloudtemplateIpSecTunnelSubnetPool.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciSubnetPoolforIpSecTunnelsRead(ctx, d, m)
}

func resourceAciSubnetPoolforIpSecTunnelsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	cloudtemplateIpSecTunnelSubnetPool, err := getRemoteSubnetPoolforIpSecTunnels(aciClient, dn)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	_, err = setSubnetPoolforIpSecTunnelsAttributes(cloudtemplateIpSecTunnelSubnetPool, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciSubnetPoolforIpSecTunnelsDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	err := aciClient.DeleteByDn(dn, "cloudtemplateIpSecTunnelSubnetPool")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}
