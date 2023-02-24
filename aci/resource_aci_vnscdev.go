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

func resourceAciConcreteDevice() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciConcreteDeviceCreate,
		UpdateContext: resourceAciConcreteDeviceUpdate,
		ReadContext:   resourceAciConcreteDeviceRead,
		DeleteContext: resourceAciConcreteDeviceDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciConcreteDeviceImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{
			"l4_l7_device_dn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vmm_controller_dn": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vm_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		})),
	}
}

func getRemoteConcreteDevice(client *client.Client, dn string) (*models.ConcreteDevice, error) {
	vnsCDevCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	vnsCDev := models.ConcreteDeviceFromContainer(vnsCDevCont)
	if vnsCDev.DistinguishedName == "" {
		return nil, fmt.Errorf("Concrete Device %s not found", dn)
	}
	return vnsCDev, nil
}

func setConcreteDeviceAttributes(vnsCDev *models.ConcreteDevice, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(vnsCDev.DistinguishedName)
	vnsCDevMap, err := vnsCDev.ToMap()
	if err != nil {
		return d, err
	}
	dn := d.Id()
	if dn != vnsCDev.DistinguishedName {
		d.Set("l4_l7_device_dn", "")
	} else {
		d.Set("l4_l7_device_dn", GetParentDn(vnsCDev.DistinguishedName, fmt.Sprintf("/"+models.RnvnsCDev, vnsCDevMap["name"])))
	}
	d.Set("name", vnsCDevMap["name"])
	d.Set("vm_name", vnsCDevMap["vmName"])
	d.Set("name_alias", vnsCDevMap["nameAlias"])
	return d, nil
}

func resourceAciConcreteDeviceImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	vnsCDev, err := getRemoteConcreteDevice(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setConcreteDeviceAttributes(vnsCDev, d)
	if err != nil {
		return nil, err
	}
	vnsRsCDevToCtrlrPData, err := aciClient.ReadRelationvnsRsCDevToCtrlrP(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vnsRsCDevToCtrlrP %v", err)
		d.Set("vmm_controller_dn", "")
	} else {
		d.Set("vmm_controller_dn", vnsRsCDevToCtrlrPData.(string))
	}

	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciConcreteDeviceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] ConcreteDevice: Beginning Creation")
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)
	L4_L7DeviceDn := d.Get("l4_l7_device_dn").(string)

	vnsCDevAttr := models.ConcreteDeviceAttributes{}

	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		vnsCDevAttr.Annotation = Annotation.(string)
	} else {
		vnsCDevAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		vnsCDevAttr.Name = Name.(string)
	}

	if VcenterName, ok := d.GetOk("vmm_controller_dn"); ok {
		vnsCDevAttr.VcenterName = GetMOName(VcenterName.(string))
	}

	if VmName, ok := d.GetOk("vm_name"); ok {
		vnsCDevAttr.VmName = VmName.(string)
	}
	vnsCDev := models.NewConcreteDevice(fmt.Sprintf(models.RnvnsCDev, name), L4_L7DeviceDn, nameAlias, vnsCDevAttr)

	err := aciClient.Save(vnsCDev)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(vnsCDev.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciConcreteDeviceRead(ctx, d, m)
}
func resourceAciConcreteDeviceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] ConcreteDevice: Beginning Update")
	aciClient := m.(*client.Client)
	name := d.Get("name").(string)
	L4_L7DeviceDn := d.Get("l4_l7_device_dn").(string)

	vnsCDevAttr := models.ConcreteDeviceAttributes{}

	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		vnsCDevAttr.Annotation = Annotation.(string)
	} else {
		vnsCDevAttr.Annotation = "{}"
	}

	if Name, ok := d.GetOk("name"); ok {
		vnsCDevAttr.Name = Name.(string)
	}

	if VcenterName, ok := d.GetOk("vmm_controller_dn"); ok {
		vnsCDevAttr.VcenterName = GetMOName(VcenterName.(string))
	}

	if VmName, ok := d.GetOk("vm_name"); ok {
		vnsCDevAttr.VmName = VmName.(string)
	}
	vnsCDev := models.NewConcreteDevice(fmt.Sprintf(models.RnvnsCDev, name), L4_L7DeviceDn, nameAlias, vnsCDevAttr)

	vnsCDev.Status = "modified"

	err := aciClient.Save(vnsCDev)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(vnsCDev.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciConcreteDeviceRead(ctx, d, m)
}

func resourceAciConcreteDeviceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	vnsCDev, err := getRemoteConcreteDevice(aciClient, dn)
	if err != nil {
		return errorForObjectNotFound(err, dn, d)
	}

	_, err = setConcreteDeviceAttributes(vnsCDev, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	vnsRsCDevToCtrlrPData, err := aciClient.ReadRelationvnsRsCDevToCtrlrP(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation vnsRsCDevToCtrlrP %v", err)
		d.Set("vmm_controller_dn", "")
	} else {
		d.Set("vmm_controller_dn", vnsRsCDevToCtrlrPData.(string))
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciConcreteDeviceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()

	err := aciClient.DeleteByDn(dn, "vnsCDev")
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}
