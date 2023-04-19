package aci

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/v2/client"
	"github.com/ciscoecosystem/aci-go-client/v2/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAciOutofServiceFabricPath() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciOutofServiceFabricPathCreate,
		UpdateContext: resourceAciOutofServiceFabricPathUpdate,
		ReadContext:   resourceAciOutofServiceFabricPathRead,
		DeleteContext: resourceAciOutofServiceFabricPathDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciOutofServiceFabricPathImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(map[string]*schema.Schema{
			"annotation": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"pod_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"node_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"fex_id": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"interface": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		}),
	}
}

func getRemoteOutofServiceFabricPath(client *client.Client, dn string) (*models.OutofServiceFabricPath, error) {
	fabricRsOosPathCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	fabricRsOosPath := models.OutofServiceFabricPathFromContainer(fabricRsOosPathCont)
	if fabricRsOosPath.DistinguishedName == "" {
		return nil, fmt.Errorf("Interface blacklist %s not found", dn)
	}
	return fabricRsOosPath, nil
}

func setOutofServiceFabricPathAttributes(fabricRsOosPath *models.OutofServiceFabricPath, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(fabricRsOosPath.DistinguishedName)
	fabricRsOosPathMap, err := fabricRsOosPath.ToMap()
	if err != nil {
		return d, err
	}
	d.Set("annotation", fabricRsOosPathMap["annotation"])
	return d, nil
}

func resourceAciOutofServiceFabricPathImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	fabricRsOosPath, err := getRemoteOutofServiceFabricPath(aciClient, dn)
	if err != nil {
		return nil, err
	}

	if fabricRsOosPath.TDn != "" {
		interfaceRegEx := regexp.MustCompile(`topology/pod-([0-9]+?)/paths-([0-9]+?)/pathep-\[(.*?)\]`)
		fexRegEx := regexp.MustCompile(`topology/pod-([0-9]+?)/paths-([0-9]+?)/extpaths-([0-9]+?)/pathep-\[(.*?)\]`)
		matchInterface := interfaceRegEx.FindStringSubmatch(fabricRsOosPath.TDn)
		if len(matchInterface) > 0 {
			podId, err := strconv.Atoi(matchInterface[1])
			if err != nil {
				return nil, err
			}
			d.Set("pod_id", podId)

			nodeId, err := strconv.Atoi(matchInterface[2])
			if err != nil {
				return nil, err
			}
			d.Set("node_id", nodeId)

			d.Set("interface", matchInterface[3])
		} else {
			matchFex := fexRegEx.FindStringSubmatch(fabricRsOosPath.TDn)
			if len(matchFex) > 0 {
				podId, err := strconv.Atoi(matchFex[1])
				if err != nil {
					return nil, err
				}
				d.Set("pod_id", podId)

				nodeId, err := strconv.Atoi(matchFex[2])
				if err != nil {
					return nil, err
				}
				d.Set("node_id", nodeId)

				fexId, err := strconv.Atoi(matchFex[3])
				if err != nil {
					return nil, err
				}
				d.Set("fex_id", fexId)

				d.Set("interface", matchFex[4])
			}
		}
	}

	schemaFilled, err := setOutofServiceFabricPathAttributes(fabricRsOosPath, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciOutofServiceFabricPathCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] OutofServiceFabricPath: Beginning Creation")
	aciClient := m.(*client.Client)

	podId := fmt.Sprintf("%v", d.Get("pod_id"))
	nodeId := fmt.Sprintf("%v", d.Get("node_id"))
	interfaceName := d.Get("interface").(string)

	var interfaceDn string
	if FexId, ok := d.GetOk("fex_id"); ok {
		fexId := fmt.Sprintf("%v", FexId)
		interfaceDn = fmt.Sprintf("topology/pod-%s/paths-%s/extpaths-%s/pathep-[%s]", podId, nodeId, fexId, interfaceName)
	} else {
		interfaceDn = fmt.Sprintf("topology/pod-%s/paths-%s/pathep-[%s]", podId, nodeId, interfaceName)
	}

	fabricRsOosPathAttr := models.OutofServiceFabricPathAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		fabricRsOosPathAttr.Annotation = Annotation.(string)
	} else {
		fabricRsOosPathAttr.Annotation = "{}"
	}

	fabricRsOosPathAttr.Lc = fmt.Sprintf("blacklist")
	fabricRsOosPathAttr.TDn = interfaceDn

	fabricRsOosPath := models.NewOutofServiceFabricPath(fmt.Sprintf(models.RnfabricRsOosPath, interfaceDn), models.ParentDnfabricRsOosPath, fabricRsOosPathAttr)
	err := aciClient.Save(fabricRsOosPath)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fabricRsOosPath.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciOutofServiceFabricPathRead(ctx, d, m)
}

func resourceAciOutofServiceFabricPathUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] OutofServiceFabricPath: Beginning Update")
	aciClient := m.(*client.Client)

	podId := fmt.Sprintf("%v", d.Get("pod_id"))
	nodeId := fmt.Sprintf("%v", d.Get("node_id"))
	interfaceName := d.Get("interface").(string)

	var interfaceDn string
	if FexId, ok := d.GetOk("fex_id"); ok {
		fexId := fmt.Sprintf("%v", FexId)
		interfaceDn = fmt.Sprintf("topology/pod-%s/paths-%s/extpaths-%s/pathep-[%s]", podId, nodeId, fexId, interfaceName)
	} else {
		interfaceDn = fmt.Sprintf("topology/pod-%s/paths-%s/pathep-[%s]", podId, nodeId, interfaceName)
	}

	fabricRsOosPathAttr := models.OutofServiceFabricPathAttributes{}
	if Annotation, ok := d.GetOk("annotation"); ok {
		fabricRsOosPathAttr.Annotation = Annotation.(string)
	} else {
		fabricRsOosPathAttr.Annotation = "{}"
	}

	fabricRsOosPathAttr.Lc = fmt.Sprintf("blacklist")
	fabricRsOosPathAttr.TDn = interfaceDn

	fabricRsOosPath := models.NewOutofServiceFabricPath(fmt.Sprintf(models.RnfabricRsOosPath, interfaceDn), models.ParentDnfabricRsOosPath, fabricRsOosPathAttr)
	fabricRsOosPath.Status = "modified"
	err := aciClient.Save(fabricRsOosPath)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fabricRsOosPath.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciOutofServiceFabricPathRead(ctx, d, m)
}

func resourceAciOutofServiceFabricPathRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	fabricRsOosPath, err := getRemoteOutofServiceFabricPath(aciClient, dn)
	if err != nil {
		return errorForObjectNotFound(err, dn, d)
	}
	_, err = setOutofServiceFabricPathAttributes(fabricRsOosPath, d)
	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciOutofServiceFabricPathDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "fabricRsOosPath")
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}
