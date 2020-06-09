package models


import (
    "strconv"

    "github.com/ciscoecosystem/aci-go-client/container"
)

const MaintUpgJobClassName = "maintUpgJob"

type UpgJob struct {
    BaseAttributes
    UpgJobAttributes
}

type UpgJobAttributes struct {
    CreationDate string `json:",omitempty"`
    DesiredVersion string `json:",omitempty"`
    DnldPercent string `json:",omitempty"`
    DnldStatus string `json:",omitempty"`
    DnldStatusStr string `json:",omitempty"`
    EndDate string `json:",omitempty"`
    FwGrp string `json:",omitempty"`
    FwPolName string `json:",omitempty"`
    GrpPriKey string `json:",omitempty"`
    IUrl string `json:",omitempty"`
    InstallId string `json:",omitempty"`
    InstallStage string `json:",omitempty"`
    InstlProgPct string `json:",omitempty"`
    InternalLabel string `json:",omitempty"`
    LastInstallDate string `json:",omitempty"`
    MaintGrp string `json:",omitempty"`
    NumAttempts string `json:",omitempty"`
    NumAttemptsToBootupReady string `json:",omitempty"`
    PolName string `json:",omitempty"`
    SrDesiredVer string `json:",omitempty"`
    SrUpg string `json:",omitempty"`
    StartDate string `json:",omitempty"`
    UpgradeStatus string `json:",omitempty"`
    UpgradeStatusStr string `json:",omitempty"`
}

/*
 * No NewUpgJob as this is a non-configurable MO
func NewUpgJob(maintUpgJobRn, parentDn, description string, maintUpgJobattr UpgJobAttributes) *UpgJob {
    dn := fmt.Sprintf("%s/%s", parentDn, maintUpgJobRn)
    return &UpgJob{
        BaseAttributes: BaseAttributes{
            DistinguishedName: dn,
            Description:       description,
            Status:            "created, modified",
            ClassName:         MaintUpgJobClassName,
            Rn:                maintUpgJobRn,
        },

        UpgJobAttributes: maintUpgJobattr,

    }
}
 */

func (maintUpgJob *UpgJob) ToMap() (map[string]string, error) {
    maintUpgJobMap, err := maintUpgJob.BaseAttributes.ToMap()
    if err != nil {
        return nil, err
    }

    A(maintUpgJobMap, "creationDate",maintUpgJob.CreationDate)
    A(maintUpgJobMap, "desiredVersion",maintUpgJob.DesiredVersion)
    A(maintUpgJobMap, "dnldPercent",maintUpgJob.DnldPercent)
    A(maintUpgJobMap, "dnldStatus",maintUpgJob.DnldStatus)
    A(maintUpgJobMap, "dnldStatusStr",maintUpgJob.DnldStatusStr)
    A(maintUpgJobMap, "endDate",maintUpgJob.EndDate)
    A(maintUpgJobMap, "fwGrp",maintUpgJob.FwGrp)
    A(maintUpgJobMap, "fwPolName",maintUpgJob.FwPolName)
    A(maintUpgJobMap, "grpPriKey",maintUpgJob.GrpPriKey)
    A(maintUpgJobMap, "iUrl",maintUpgJob.IUrl)
    A(maintUpgJobMap, "installId",maintUpgJob.InstallId)
    A(maintUpgJobMap, "installStage",maintUpgJob.InstallStage)
    A(maintUpgJobMap, "instlProgPct",maintUpgJob.InstlProgPct)
    A(maintUpgJobMap, "internalLabel",maintUpgJob.InternalLabel)
    A(maintUpgJobMap, "lastInstallDate",maintUpgJob.LastInstallDate)
    A(maintUpgJobMap, "maintGrp",maintUpgJob.MaintGrp)
    A(maintUpgJobMap, "numAttempts",maintUpgJob.NumAttempts)
    A(maintUpgJobMap, "numAttemptsToBootupReady",maintUpgJob.NumAttemptsToBootupReady)
    A(maintUpgJobMap, "polName",maintUpgJob.PolName)
    A(maintUpgJobMap, "srDesiredVer",maintUpgJob.SrDesiredVer)
    A(maintUpgJobMap, "srUpg",maintUpgJob.SrUpg)
    A(maintUpgJobMap, "startDate",maintUpgJob.StartDate)
    A(maintUpgJobMap, "upgradeStatus",maintUpgJob.UpgradeStatus)
    A(maintUpgJobMap, "upgradeStatusStr",maintUpgJob.UpgradeStatusStr)

    return maintUpgJobMap, err
}

func UpgJobFromContainerList(cont *container.Container, index int) *UpgJob {

    UpgJobCont := cont.S("imdata").Index(index).S(MaintUpgJobClassName, "attributes")
    return &UpgJob{
        BaseAttributes{
            DistinguishedName: G(UpgJobCont, "dn"),
            Description:       G(UpgJobCont, "descr"),
            Status:            G(UpgJobCont, "status"),
            ClassName:         MaintUpgJobClassName,
            Rn:                G(UpgJobCont, "rn"),
        },

        UpgJobAttributes{

        CreationDate : G(UpgJobCont,"creationDate"),
        DesiredVersion : G(UpgJobCont,"desiredVersion"),
        DnldPercent : G(UpgJobCont,"dnldPercent"),
        DnldStatus : G(UpgJobCont,"dnldStatus"),
        DnldStatusStr : G(UpgJobCont,"dnldStatusStr"),
        EndDate : G(UpgJobCont,"endDate"),
        FwGrp : G(UpgJobCont,"fwGrp"),
        FwPolName : G(UpgJobCont,"fwPolName"),
        GrpPriKey : G(UpgJobCont,"grpPriKey"),
        IUrl : G(UpgJobCont,"iUrl"),
        InstallId : G(UpgJobCont,"installId"),
        InstallStage : G(UpgJobCont,"installStage"),
        InstlProgPct : G(UpgJobCont,"instlProgPct"),
        InternalLabel : G(UpgJobCont,"internalLabel"),
        LastInstallDate : G(UpgJobCont,"lastInstallDate"),
        MaintGrp : G(UpgJobCont,"maintGrp"),
        NumAttempts : G(UpgJobCont,"numAttempts"),
        NumAttemptsToBootupReady : G(UpgJobCont,"numAttemptsToBootupReady"),
        PolName : G(UpgJobCont,"polName"),
        SrDesiredVer : G(UpgJobCont,"srDesiredVer"),
        SrUpg : G(UpgJobCont,"srUpg"),
        StartDate : G(UpgJobCont,"startDate"),
        UpgradeStatus : G(UpgJobCont,"upgradeStatus"),
        UpgradeStatusStr : G(UpgJobCont,"upgradeStatusStr"),
        },

    }
}

func UpgJobFromContainer(cont *container.Container) *UpgJob {

    return UpgJobFromContainerList(cont, 0)
}

func UpgJobListFromContainer(cont *container.Container) []*UpgJob {
    length, _ := strconv.Atoi(G(cont, "totalCount"))

    arr := make([]*UpgJob, length)

    for i := 0; i < length; i++ {

        arr[i] = UpgJobFromContainerList(cont, i)
    }

    return arr
}
