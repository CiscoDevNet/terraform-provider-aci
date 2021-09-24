package models

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/aci-go-client/container"
)

const (
	DnfileRemotePath        = "uni/fabric/path-%s"
	RnfileRemotePath        = "path-%s"
	ParentDnfileRemotePath  = "uni/fabric"
	FileremotepathClassName = "fileRemotePath"
)

type RemotePathofaFile struct {
	BaseAttributes
	NameAliasAttribute
	RemotePathofaFileAttributes
}

type RemotePathofaFileAttributes struct {
	Annotation                   string `json:",omitempty"`
	AuthType                     string `json:",omitempty"`
	Host                         string `json:",omitempty"`
	IdentityPrivateKeyContents   string `json:",omitempty"`
	IdentityPrivateKeyPassphrase string `json:",omitempty"`
	Name                         string `json:",omitempty"`
	Protocol                     string `json:",omitempty"`
	RemotePath                   string `json:",omitempty"`
	RemotePort                   string `json:",omitempty"`
	UserName                     string `json:",omitempty"`
	UserPasswd                   string `json:",omitempty"`
}

func NewRemotePathofaFile(fileRemotePathRn, parentDn, description, nameAlias string, fileRemotePathAttr RemotePathofaFileAttributes) *RemotePathofaFile {
	dn := fmt.Sprintf("%s/%s", parentDn, fileRemotePathRn)
	return &RemotePathofaFile{
		BaseAttributes: BaseAttributes{
			DistinguishedName: dn,
			Description:       description,
			Status:            "created, modified",
			ClassName:         FileremotepathClassName,
			Rn:                fileRemotePathRn,
		},
		NameAliasAttribute: NameAliasAttribute{
			NameAlias: nameAlias,
		},
		RemotePathofaFileAttributes: fileRemotePathAttr,
	}
}

func (fileRemotePath *RemotePathofaFile) ToMap() (map[string]string, error) {
	fileRemotePathMap, err := fileRemotePath.BaseAttributes.ToMap()
	if err != nil {
		return nil, err
	}
	alias, err := fileRemotePath.NameAliasAttribute.ToMap()
	if err != nil {
		return nil, err
	}
	for key, value := range alias {
		A(fileRemotePathMap, key, value)
	}
	A(fileRemotePathMap, "annotation", fileRemotePath.Annotation)
	A(fileRemotePathMap, "authType", fileRemotePath.AuthType)
	A(fileRemotePathMap, "host", fileRemotePath.Host)
	A(fileRemotePathMap, "identityPrivateKeyContents", fileRemotePath.IdentityPrivateKeyContents)
	A(fileRemotePathMap, "identityPrivateKeyPassphrase", fileRemotePath.IdentityPrivateKeyPassphrase)
	A(fileRemotePathMap, "name", fileRemotePath.Name)
	A(fileRemotePathMap, "protocol", fileRemotePath.Protocol)
	A(fileRemotePathMap, "remotePath", fileRemotePath.RemotePath)
	A(fileRemotePathMap, "remotePort", fileRemotePath.RemotePort)
	A(fileRemotePathMap, "userName", fileRemotePath.UserName)
	A(fileRemotePathMap, "userPasswd", fileRemotePath.UserPasswd)
	return fileRemotePathMap, err
}

func RemotePathofaFileFromContainerList(cont *container.Container, index int) *RemotePathofaFile {
	RemotePathofaFileCont := cont.S("imdata").Index(index).S(FileremotepathClassName, "attributes")
	return &RemotePathofaFile{
		BaseAttributes{
			DistinguishedName: G(RemotePathofaFileCont, "dn"),
			Description:       G(RemotePathofaFileCont, "descr"),
			Status:            G(RemotePathofaFileCont, "status"),
			ClassName:         FileremotepathClassName,
			Rn:                G(RemotePathofaFileCont, "rn"),
		},
		NameAliasAttribute{
			NameAlias: G(RemotePathofaFileCont, "nameAlias"),
		},
		RemotePathofaFileAttributes{
			Annotation:                   G(RemotePathofaFileCont, "annotation"),
			AuthType:                     G(RemotePathofaFileCont, "authType"),
			Host:                         G(RemotePathofaFileCont, "host"),
			IdentityPrivateKeyContents:   G(RemotePathofaFileCont, "identityPrivateKeyContents"),
			IdentityPrivateKeyPassphrase: G(RemotePathofaFileCont, "identityPrivateKeyPassphrase"),
			Name:                         G(RemotePathofaFileCont, "name"),
			Protocol:                     G(RemotePathofaFileCont, "protocol"),
			RemotePath:                   G(RemotePathofaFileCont, "remotePath"),
			RemotePort:                   G(RemotePathofaFileCont, "remotePort"),
			UserName:                     G(RemotePathofaFileCont, "userName"),
			UserPasswd:                   G(RemotePathofaFileCont, "userPasswd"),
		},
	}
}

func RemotePathofaFileFromContainer(cont *container.Container) *RemotePathofaFile {
	return RemotePathofaFileFromContainerList(cont, 0)
}

func RemotePathofaFileListFromContainer(cont *container.Container) []*RemotePathofaFile {
	length, _ := strconv.Atoi(G(cont, "totalCount"))
	arr := make([]*RemotePathofaFile, length)
	for i := 0; i < length; i++ {
		arr[i] = RemotePathofaFileFromContainerList(cont, i)
	}
	return arr
}
