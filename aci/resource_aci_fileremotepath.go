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

func resourceAciRemotePathofaFile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAciRemotePathofaFileCreate,
		UpdateContext: resourceAciRemotePathofaFileUpdate,
		ReadContext:   resourceAciRemotePathofaFileRead,
		DeleteContext: resourceAciRemotePathofaFileDelete,

		Importer: &schema.ResourceImporter{
			State: resourceAciRemotePathofaFileImport,
		},

		SchemaVersion: 1,
		Schema: AppendBaseAttrSchema(AppendNameAliasAttrSchema(map[string]*schema.Schema{

			"auth_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "usePassword",
				ValidateFunc: validation.StringInSlice([]string{
					"usePassword",
					"useSshKeyContents",
				}, false),
			},
			"host": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"identity_private_key_contents": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"identity_private_key_passphrase": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"protocol": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "sftp",
				ValidateFunc: validation.StringInSlice([]string{
					"ftp",
					"scp",
					"sftp",
				}, false),
			},
			"remote_path": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateRemoteFilePath(),
			},
			"remote_port": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"user_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"user_passwd": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"relation_file_rs_a_remote_host_to_epg": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Create relation to fv:ATg",
			},
			"relation_file_rs_a_remote_host_to_epp": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Create relation to fv:AREpP",
			},
		})),
	}
}

func getRemoteRemotePathofaFile(client *client.Client, dn string) (*models.RemotePathofaFile, error) {
	fileRemotePathCont, err := client.Get(dn)
	if err != nil {
		return nil, err
	}
	fileRemotePath := models.RemotePathofaFileFromContainer(fileRemotePathCont)
	if fileRemotePath.DistinguishedName == "" {
		return nil, fmt.Errorf("RemotePathofaFile %s not found", fileRemotePath.DistinguishedName)
	}
	return fileRemotePath, nil
}

func setRemotePathofaFileAttributes(fileRemotePath *models.RemotePathofaFile, d *schema.ResourceData) (*schema.ResourceData, error) {
	d.SetId(fileRemotePath.DistinguishedName)
	d.Set("description", fileRemotePath.Description)
	fileRemotePathMap, err := fileRemotePath.ToMap()
	if err != nil {
		return nil, err
	}
	d.Set("annotation", fileRemotePathMap["annotation"])
	d.Set("auth_type", fileRemotePathMap["authType"])
	d.Set("host", fileRemotePathMap["host"])
	d.Set("name", fileRemotePathMap["name"])
	d.Set("protocol", fileRemotePathMap["protocol"])
	d.Set("remote_path", fileRemotePathMap["remotePath"])
	d.Set("remote_port", fileRemotePathMap["remotePort"])
	d.Set("user_name", fileRemotePathMap["userName"])
	d.Set("name_alias", fileRemotePathMap["nameAlias"])
	return d, nil
}

func resourceAciRemotePathofaFileImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[DEBUG] %s: Beginning Import", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	fileRemotePath, err := getRemoteRemotePathofaFile(aciClient, dn)
	if err != nil {
		return nil, err
	}
	schemaFilled, err := setRemotePathofaFileAttributes(fileRemotePath, d)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] %s: Import finished successfully", d.Id())
	return []*schema.ResourceData{schemaFilled}, nil
}

func resourceAciRemotePathofaFileCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] RemotePathofaFile: Beginning Creation")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	var PrivateKeyOk, PasswordOk, PassPhraseOk bool
	fileRemotePathAttr := models.RemotePathofaFileAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}
	if Annotation, ok := d.GetOk("annotation"); ok {
		fileRemotePathAttr.Annotation = Annotation.(string)
	} else {
		fileRemotePathAttr.Annotation = "{}"
	}

	if AuthType, ok := d.GetOk("auth_type"); ok {
		fileRemotePathAttr.AuthType = AuthType.(string)
	}

	if Host, ok := d.GetOk("host"); ok {
		fileRemotePathAttr.Host = Host.(string)
	}

	if IdentityPrivateKeyContents, ok := d.GetOk("identity_private_key_contents"); ok {
		PrivateKeyOk = ok
		fileRemotePathAttr.IdentityPrivateKeyContents = IdentityPrivateKeyContents.(string)
	}

	if IdentityPrivateKeyPassphrase, ok := d.GetOk("identity_private_key_passphrase"); ok {
		PassPhraseOk = ok
		fileRemotePathAttr.IdentityPrivateKeyPassphrase = IdentityPrivateKeyPassphrase.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		fileRemotePathAttr.Name = Name.(string)
	}

	if Protocol, ok := d.GetOk("protocol"); ok {
		fileRemotePathAttr.Protocol = Protocol.(string)
	}

	if RemotePath, ok := d.GetOk("remote_path"); ok {
		fileRemotePathAttr.RemotePath = RemotePath.(string)
	}

	if RemotePort, ok := d.GetOk("remote_port"); ok {
		fileRemotePathAttr.RemotePort = RemotePort.(string)
	}

	if UserName, ok := d.GetOk("user_name"); ok {
		fileRemotePathAttr.UserName = UserName.(string)
	}

	if UserPasswd, ok := d.GetOk("user_passwd"); ok {
		PasswordOk = ok
		fileRemotePathAttr.UserPasswd = UserPasswd.(string)
	}

	if fileRemotePathAttr.Protocol == "ftp" && fileRemotePathAttr.AuthType == "useSshKeyContents" {
		return diag.FromErr(fmt.Errorf("auth_type should be usePassword when protocol is ftp"))
	}

	if fileRemotePathAttr.AuthType == "useSshKeyContents" && !PrivateKeyOk {
		return diag.FromErr(fmt.Errorf("identity_private_key_contents must be set when auth_type is useSshKeyContents"))
	} else if fileRemotePathAttr.AuthType == "usePassword" && !PasswordOk {
		return diag.FromErr(fmt.Errorf("user_passwd must be set when auth_type is usePassword"))
	}

	if PassPhraseOk && !PrivateKeyOk {
		return diag.FromErr(fmt.Errorf("identity_private_key_passphrase must be set if and only if identity_private_key_contents is set"))
	}
	fileRemotePath := models.NewRemotePathofaFile(fmt.Sprintf("fabric/path-%s", name), "uni", desc, nameAlias, fileRemotePathAttr)
	err := aciClient.Save(fileRemotePath)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if relationTofileRsARemoteHostToEpg, ok := d.GetOk("relation_file_rs_a_remote_host_to_epg"); ok {
		relationParam := relationTofileRsARemoteHostToEpg.(string)
		checkDns = append(checkDns, relationParam)

	}

	if relationTofileRsARemoteHostToEpp, ok := d.GetOk("relation_file_rs_a_remote_host_to_epp"); ok {
		relationParam := relationTofileRsARemoteHostToEpp.(string)
		checkDns = append(checkDns, relationParam)

	}

	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}

	if relationTofileRsARemoteHostToEpg, ok := d.GetOk("relation_file_rs_a_remote_host_to_epg"); ok {
		relationParam := relationTofileRsARemoteHostToEpg.(string)
		err = aciClient.CreateRelationfileRsARemoteHostToEpg(fileRemotePath.DistinguishedName, fileRemotePathAttr.Annotation, relationParam)

		if err != nil {
			return diag.FromErr(err)
		}

	}

	if relationTofileRsARemoteHostToEpp, ok := d.GetOk("relation_file_rs_a_remote_host_to_epp"); ok {
		relationParam := relationTofileRsARemoteHostToEpp.(string)
		err = aciClient.CreateRelationfileRsARemoteHostToEpp(fileRemotePath.DistinguishedName, fileRemotePathAttr.Annotation, relationParam)

		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(fileRemotePath.DistinguishedName)
	log.Printf("[DEBUG] %s: Creation finished successfully", d.Id())
	return resourceAciRemotePathofaFileRead(ctx, d, m)
}

func resourceAciRemotePathofaFileUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] RemotePathofaFile: Beginning Update")
	aciClient := m.(*client.Client)
	desc := d.Get("description").(string)
	name := d.Get("name").(string)
	var PrivateKeyOk, PasswordOk, PassPhraseOk bool
	fileRemotePathAttr := models.RemotePathofaFileAttributes{}
	nameAlias := ""
	if NameAlias, ok := d.GetOk("name_alias"); ok {
		nameAlias = NameAlias.(string)
	}

	if Annotation, ok := d.GetOk("annotation"); ok {
		fileRemotePathAttr.Annotation = Annotation.(string)
	} else {
		fileRemotePathAttr.Annotation = "{}"
	}

	if AuthType, ok := d.GetOk("auth_type"); ok {
		fileRemotePathAttr.AuthType = AuthType.(string)
	}

	if Host, ok := d.GetOk("host"); ok {
		fileRemotePathAttr.Host = Host.(string)
	}

	if IdentityPrivateKeyContents, ok := d.GetOk("identity_private_key_contents"); ok {
		PrivateKeyOk = ok
		fileRemotePathAttr.IdentityPrivateKeyContents = IdentityPrivateKeyContents.(string)
	}

	if IdentityPrivateKeyPassphrase, ok := d.GetOk("identity_private_key_passphrase"); ok {
		PassPhraseOk = ok
		fileRemotePathAttr.IdentityPrivateKeyPassphrase = IdentityPrivateKeyPassphrase.(string)
	}

	if Name, ok := d.GetOk("name"); ok {
		fileRemotePathAttr.Name = Name.(string)
	}

	if Protocol, ok := d.GetOk("protocol"); ok {
		fileRemotePathAttr.Protocol = Protocol.(string)
	}

	if RemotePath, ok := d.GetOk("remote_path"); ok {
		fileRemotePathAttr.RemotePath = RemotePath.(string)
	}

	if RemotePort, ok := d.GetOk("remote_port"); ok {
		fileRemotePathAttr.RemotePort = RemotePort.(string)
	}

	if UserName, ok := d.GetOk("user_name"); ok {
		fileRemotePathAttr.UserName = UserName.(string)
	}

	if UserPasswd, ok := d.GetOk("user_passwd"); ok {
		PasswordOk = ok
		fileRemotePathAttr.UserPasswd = UserPasswd.(string)
	}

	if fileRemotePathAttr.Protocol == "ftp" && fileRemotePathAttr.AuthType == "useSshKeyContents" {
		return diag.FromErr(fmt.Errorf("auth_type should be usePassword when protocol is ftp"))
	}

	if fileRemotePathAttr.AuthType == "useSshKeyContents" && !PrivateKeyOk {
		return diag.FromErr(fmt.Errorf("identity_private_key_contents must be set when auth_type is useSshKeyContents"))
	} else if fileRemotePathAttr.AuthType == "usePassword" && !PasswordOk {
		return diag.FromErr(fmt.Errorf("user_passwd must be set when auth_type is usePassword"))
	}

	if PassPhraseOk && !PrivateKeyOk {
		return diag.FromErr(fmt.Errorf("identity_private_key_passphrase must be set if and only if identity_private_key_contents is set"))
	}

	fileRemotePath := models.NewRemotePathofaFile(fmt.Sprintf("fabric/path-%s", name), "uni", desc, nameAlias, fileRemotePathAttr)
	fileRemotePath.Status = "modified"
	err := aciClient.Save(fileRemotePath)
	if err != nil {
		return diag.FromErr(err)
	}

	checkDns := make([]string, 0, 1)

	if d.HasChange("relation_file_rs_a_remote_host_to_epg") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_file_rs_a_remote_host_to_epg")
		checkDns = append(checkDns, newRelParam.(string))

	}

	if d.HasChange("relation_file_rs_a_remote_host_to_epp") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_file_rs_a_remote_host_to_epp")
		checkDns = append(checkDns, newRelParam.(string))

	}

	err = checkTDn(aciClient, checkDns)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("relation_file_rs_a_remote_host_to_epg") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_file_rs_a_remote_host_to_epg")
		err = aciClient.DeleteRelationfileRsARemoteHostToEpg(fileRemotePath.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationfileRsARemoteHostToEpg(fileRemotePath.DistinguishedName, fileRemotePathAttr.Annotation, newRelParam.(string))

		if err != nil {
			return diag.FromErr(err)
		}

	}
	if d.HasChange("relation_file_rs_a_remote_host_to_epp") || d.HasChange("annotation") {
		_, newRelParam := d.GetChange("relation_file_rs_a_remote_host_to_epp")
		err = aciClient.DeleteRelationfileRsARemoteHostToEpp(fileRemotePath.DistinguishedName)
		if err != nil {
			return diag.FromErr(err)
		}
		err = aciClient.CreateRelationfileRsARemoteHostToEpp(fileRemotePath.DistinguishedName, fileRemotePathAttr.Annotation, newRelParam.(string))

		if err != nil {
			return diag.FromErr(err)
		}

	}

	d.SetId(fileRemotePath.DistinguishedName)
	log.Printf("[DEBUG] %s: Update finished successfully", d.Id())
	return resourceAciRemotePathofaFileRead(ctx, d, m)
}

func resourceAciRemotePathofaFileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Read", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	fileRemotePath, err := getRemoteRemotePathofaFile(aciClient, dn)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	setRemotePathofaFileAttributes(fileRemotePath, d)

	fileRsARemoteHostToEpgData, err := aciClient.ReadRelationfileRsARemoteHostToEpg(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fileRsARemoteHostToEpg %v", err)
		d.Set("relation_file_rs_a_remote_host_to_epg", "")
	} else {
		setRelationAttribute(d, "relation_file_rs_a_remote_host_to_epg", fileRsARemoteHostToEpgData)
	}

	fileRsARemoteHostToEppData, err := aciClient.ReadRelationfileRsARemoteHostToEpp(dn)
	if err != nil {
		log.Printf("[DEBUG] Error while reading relation fileRsARemoteHostToEpp %v", err)
		d.Set("relation_file_rs_a_remote_host_to_epp", "")
	} else {
		setRelationAttribute(d, "relation_file_rs_a_remote_host_to_epp", fileRsARemoteHostToEppData)
	}
	log.Printf("[DEBUG] %s: Read finished successfully", d.Id())
	return nil
}

func resourceAciRemotePathofaFileDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] %s: Beginning Destroy", d.Id())
	aciClient := m.(*client.Client)
	dn := d.Id()
	err := aciClient.DeleteByDn(dn, "fileRemotePath")
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] %s: Destroy finished successfully", d.Id())
	d.SetId("")
	return diag.FromErr(err)
}
