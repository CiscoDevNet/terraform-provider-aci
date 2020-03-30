package client

import (
	"fmt"

	"github.com/ciscoecosystem/aci-go-client/models"
)

func (sm *ServiceManager) CreateSPANSourcedestinationGroupMatchLabel(name string, span_source_group string, tenant string, description string, spanSpanLblattr models.SPANSourcedestinationGroupMatchLabelAttributes) (*models.SPANSourcedestinationGroupMatchLabel, error) {
	rn := fmt.Sprintf("spanlbl-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s/srcgrp-%s", tenant, span_source_group)
	spanSpanLbl := models.NewSPANSourcedestinationGroupMatchLabel(rn, parentDn, description, spanSpanLblattr)
	err := sm.Save(spanSpanLbl)
	return spanSpanLbl, err
}

func (sm *ServiceManager) ReadSPANSourcedestinationGroupMatchLabel(name string, span_source_group string, tenant string) (*models.SPANSourcedestinationGroupMatchLabel, error) {
	dn := fmt.Sprintf("uni/tn-%s/srcgrp-%s/spanlbl-%s", tenant, span_source_group, name)
	cont, err := sm.Get(dn)
	if err != nil {
		return nil, err
	}

	spanSpanLbl := models.SPANSourcedestinationGroupMatchLabelFromContainer(cont)
	return spanSpanLbl, nil
}

func (sm *ServiceManager) DeleteSPANSourcedestinationGroupMatchLabel(name string, span_source_group string, tenant string) error {
	dn := fmt.Sprintf("uni/tn-%s/srcgrp-%s/spanlbl-%s", tenant, span_source_group, name)
	return sm.DeleteByDn(dn, models.SpanspanlblClassName)
}

func (sm *ServiceManager) UpdateSPANSourcedestinationGroupMatchLabel(name string, span_source_group string, tenant string, description string, spanSpanLblattr models.SPANSourcedestinationGroupMatchLabelAttributes) (*models.SPANSourcedestinationGroupMatchLabel, error) {
	rn := fmt.Sprintf("spanlbl-%s", name)
	parentDn := fmt.Sprintf("uni/tn-%s/srcgrp-%s", tenant, span_source_group)
	spanSpanLbl := models.NewSPANSourcedestinationGroupMatchLabel(rn, parentDn, description, spanSpanLblattr)

	spanSpanLbl.Status = "modified"
	err := sm.Save(spanSpanLbl)
	return spanSpanLbl, err

}

func (sm *ServiceManager) ListSPANSourcedestinationGroupMatchLabel(span_source_group string, tenant string) ([]*models.SPANSourcedestinationGroupMatchLabel, error) {

	baseurlStr := "/api/node/class"
	dnUrl := fmt.Sprintf("%s/uni/tn-%s/srcgrp-%s/spanSpanLbl.json", baseurlStr, tenant, span_source_group)

	cont, err := sm.GetViaURL(dnUrl)
	list := models.SPANSourcedestinationGroupMatchLabelListFromContainer(cont)

	return list, err
}
