package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"



	


)









func (sm *ServiceManager) CreateX509Certificate(name string ,local_user string , description string, aaaUserCertattr models.X509CertificateAttributes) (*models.X509Certificate, error) {	
	rn := fmt.Sprintf("usercert-%s",name)
	parentDn := fmt.Sprintf("uni/userext/user-%s", local_user )
	aaaUserCert := models.NewX509Certificate(rn, parentDn, description, aaaUserCertattr)
	err := sm.Save(aaaUserCert)
	return aaaUserCert, err
}

func (sm *ServiceManager) ReadX509Certificate(name string ,local_user string ) (*models.X509Certificate, error) {
	dn := fmt.Sprintf("uni/userext/user-%s/usercert-%s", local_user ,name )    
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	aaaUserCert := models.X509CertificateFromContainer(cont)
	return aaaUserCert, nil
}

func (sm *ServiceManager) DeleteX509Certificate(name string ,local_user string ) error {
	dn := fmt.Sprintf("uni/userext/user-%s/usercert-%s", local_user ,name )
	return sm.DeleteByDn(dn, models.AaausercertClassName)
}

func (sm *ServiceManager) UpdateX509Certificate(name string ,local_user string  ,description string, aaaUserCertattr models.X509CertificateAttributes) (*models.X509Certificate, error) {
	rn := fmt.Sprintf("usercert-%s",name)
	parentDn := fmt.Sprintf("uni/userext/user-%s", local_user )
	aaaUserCert := models.NewX509Certificate(rn, parentDn, description, aaaUserCertattr)

    aaaUserCert.Status = "modified"
	err := sm.Save(aaaUserCert)
	return aaaUserCert, err

}

func (sm *ServiceManager) ListX509Certificate(local_user string ) ([]*models.X509Certificate, error) {

	baseurlStr := "/api/node/class"	
	dnUrl := fmt.Sprintf("%s/uni/userext/user-%s/aaaUserCert.json", baseurlStr , local_user )
    
    cont, err := sm.GetViaURL(dnUrl)
	list := models.X509CertificateListFromContainer(cont)

	return list, err
}


