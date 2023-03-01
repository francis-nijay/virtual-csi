package service

import (
	"github.com/container-storage-interface/spec/lib/go/csi"

	powerscale "github.com/dell/csi-isilon/service"
	powermax "github.com/dell/csi-powermax/v2/service"
	powerstoreCtrl "github.com/dell/csi-powerstore/pkg/controller"
	powerstoreNode "github.com/dell/csi-powerstore/pkg/node"
	unity "github.com/dell/csi-unity/service"
	powerflex "github.com/dell/csi-vxflexos/v2/service"
)

const (
	// Name is the name of this CSI SP.
	Name = "virtual-csi"

	// VendorVersion is the version of this CSP SP.
	VendorVersion = "0.0.1"
)

var (
	syncNodeInfoChan chan bool
	// DriverConfig driver config
	DriverConfig string
	// DriverSecret driver secret
	DriverSecret string
	// Name name of the driver
	Name                  string
	syncMutex             sync.Mutex
	unityService          = unity.New()
	powerscaleService     = powerscale.New()
	powerflexService      = powerflex.New()
	powermaxService       = powermax.New()
	powerstoreCtrlService = powerstoreCtrl.Service{}
	powerstoreNodeService = powerstoreNode.Service{}
)

// StorageArrayList To parse the secret json file
type StorageArrayList struct {
	StorageArrayList []StorageArrayConfig `json:"storageArrayList"`
}

// StorageArrayConfig storage array config
type StorageArrayConfig struct {
	ArrayType        string `json:"arrayType"`
	ArrayID          string `json:"arrayId"`
	ProvisioningType string `json:"provisioningType"`
	Username         string `json:"username"`
	Password         string `json:"password"`
	RestGateway      string `json:"restGateway"`
	Insecure         bool   `json:"insecure"`
	IsDefaultArray   bool   `json:"isDefaultArray"`
	IsProbeSuccess   bool   `json:"isProbeSuccess"`
	IsHostAdded      bool   `json:"isHostAdded"`
	unityClient      *gounity.Client
	powerscaleClient *goisilon.Client
	powermaxClient   *pmax.Pmax
	powerflexClient  *goscaleio.Client
	powerstoreClient *gopowerstore.Client
}

// Service is a CSI SP and idempotency.Provider.
type Service interface {
	csi.ControllerServer
	csi.IdentityServer
	csi.NodeServer
}

type service struct{
	opts   Opts
	arrays *sync.Map
	mode   string
	BeforeServe(context.Context, *gocsi.StoragePlugin, net.Listener) error
	RegisterAdditionalServers(*grpc.Server)
}

// Opts defines service configuration options.
type Opts struct {
	NodeName                 string
	LongNodeName             string
	Chroot                   string
	Thick                    bool
	AutoProbe                bool
	PvtMountDir              string
	Debug                    bool
	SyncNodeInfoTimeInterval int
}

// New returns a new Service.
func New() Service {
	return &service{}
}
