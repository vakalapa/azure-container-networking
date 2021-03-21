package cns

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"
)

// Container Network Service DNC Contract
const (
	SetOrchestratorType                      = "/network/setorchestratortype"
	CreateOrUpdateNetworkContainer           = "/network/createorupdatenetworkcontainer"
	DeleteNetworkContainer                   = "/network/deletenetworkcontainer"
	GetNetworkContainerStatus                = "/network/getnetworkcontainerstatus"
	PublishNetworkContainer                  = "/network/publishnetworkcontainer"
	UnpublishNetworkContainer                = "/network/unpublishnetworkcontainer"
	GetInterfaceForContainer                 = "/network/getinterfaceforcontainer"
	GetNetworkContainerByOrchestratorContext = "/network/getnetworkcontainerbyorchestratorcontext"
	AttachContainerToNetwork                 = "/network/attachcontainertonetwork"
	DetachContainerFromNetwork               = "/network/detachcontainerfromnetwork"
	RequestIPConfig                          = "/network/requestipconfig"
	ReleaseIPConfig                          = "/network/releaseipconfig"
	GetIPAddresses                           = "/debug/getipaddresses"
)

// NetworkContainer Prefixes
const (
	SwiftPrefix = "Swift_"
)

// NetworkContainer Types
const (
	AzureContainerInstance = "AzureContainerInstance"
	WebApps                = "WebApps"
	Docker                 = "Docker"
	Basic                  = "Basic"
	JobObject              = "JobObject"
	COW                    = "COW" // Container on Windows
)

// Orchestrator Types
const (
	Kubernetes      = "Kubernetes"
	ServiceFabric   = "ServiceFabric"
	Batch           = "Batch"
	DBforPostgreSQL = "DBforPostgreSQL"
	AzureFirstParty = "AzureFirstParty"
	KubernetesCRD   = "KubernetesCRD"
	// TODO: Add OrchastratorType as CRD: https://msazure.visualstudio.com/One/_workitems/edit/7711872
)

// Encap Types
const (
	Vlan  = "Vlan"
	Vxlan = "Vxlan"
)

// IPConfig States for CNS IPAM
const (
	Available          = "Available"
	Allocated          = "Allocated"
	PendingRelease     = "PendingRelease"
	PendingProgramming = "PendingProgramming"
)

// ChannelMode :- CNS channel modes
const (
	Direct  = "Direct"
	Managed = "Managed"
	CRD     = "CRD"
)

// CreateNetworkContainerRequest specifies request to create a network container or network isolation boundary.
type CreateNetworkContainerRequest struct {
	Version                    string
	NetworkContainerType       string
	NetworkContainerid         string // Mandatory input.
	PrimaryInterfaceIdentifier string // Primary CA.
	AuthorizationToken         string
	LocalIPConfiguration       IPConfiguration
	OrchestratorContext        json.RawMessage
	IPConfiguration            IPConfiguration
	SecondaryIPConfigs         map[string]SecondaryIPConfig //uuid is key
	MultiTenancyInfo           MultiTenancyInfo
	CnetAddressSpace           []IPSubnet // To setup SNAT (should include service endpoint vips).
	Routes                     []Route
	AllowHostToNCCommunication bool
	AllowNCToHostCommunication bool
	EndpointPolicies           []NetworkContainerRequestPolicies
}

// NetworkContainerRequestPolicies - specifies policies associated with create network request
type NetworkContainerRequestPolicies struct {
	Type         string
	EndpointType string
	Settings     json.RawMessage
}

// ConfigureContainerNetworkingRequest - specifies request to attach/detach container to network.
type ConfigureContainerNetworkingRequest struct {
	Containerid        string
	NetworkContainerid string
}

// KubernetesPodInfo is an OrchestratorContext that holds PodName and PodNamespace.
type KubernetesPodInfo struct {
	PodName      string
	PodNamespace string
}

// GetOrchestratorContext will return the orchestratorcontext as a string
// TODO - should use a hashed name or can this be PODUid?
func (podinfo *KubernetesPodInfo) GetOrchestratorContextKey() string {
	return podinfo.PodName + ":" + podinfo.PodNamespace
}

// MultiTenancyInfo contains encap type and id.
type MultiTenancyInfo struct {
	EncapType string
	ID        int // This can be vlanid, vxlanid, gre-key etc. (depends on EnacapType).
}

// IPConfiguration contains details about ip config to provision in the VM.
type IPConfiguration struct {
	IPSubnet         IPSubnet
	DNSServers       []string
	GatewayIPAddress string
}

// SecondaryIPConfig contains IP info of SecondaryIP
type SecondaryIPConfig struct {
	IPAddress string
	// NCVesion will help in determining whether IP is in pending programming or available when reconciling.
	NCVersion int
}

// IPSubnet contains ip subnet.
type IPSubnet struct {
	IPAddress    string
	PrefixLength uint8
}

//GetIPNet converts the IPSubnet to the standard net type
func (ips *IPSubnet) GetIPNet() (net.IP, *net.IPNet, error) {
	prefix := strconv.Itoa(int(ips.PrefixLength))
	return net.ParseCIDR(ips.IPAddress + "/" + prefix)
}

// Route describes an entry in routing table.
type Route struct {
	IPAddress        string
	GatewayIPAddress string
	InterfaceToUse   string
}

// SetOrchestratorTypeRequest specifies the orchestrator type for the node.
type SetOrchestratorTypeRequest struct {
	OrchestratorType string
	DncPartitionKey  string
	NodeID           string
}

// CreateNetworkContainerResponse specifies response of creating a network container.
type CreateNetworkContainerResponse struct {
	Response Response
}

// GetNetworkContainerStatusRequest specifies the details about the request to retrieve status of a specifc network container.
type GetNetworkContainerStatusRequest struct {
	NetworkContainerid string
}

// GetNetworkContainerStatusResponse specifies response of retriving a network container status.
type GetNetworkContainerStatusResponse struct {
	NetworkContainerid string
	Version            string
	AzureHostVersion   string
	Response           Response
}

// GetNetworkContainerRequest specifies the details about the request to retrieve a specifc network container.
type GetNetworkContainerRequest struct {
	NetworkContainerid  string
	OrchestratorContext json.RawMessage
}

// GetNetworkContainerResponse describes the response to retrieve a specifc network container.
type GetNetworkContainerResponse struct {
	NetworkContainerID         string
	IPConfiguration            IPConfiguration
	Routes                     []Route
	CnetAddressSpace           []IPSubnet
	MultiTenancyInfo           MultiTenancyInfo
	PrimaryInterfaceIdentifier string
	LocalIPConfiguration       IPConfiguration
	Response                   Response
	AllowHostToNCCommunication bool
	AllowNCToHostCommunication bool
}

// DeleteNetworkContainerRequest specifies the details about the request to delete a specifc network container.
type PodIpInfo struct {
	PodIPConfig                     IPSubnet
	NetworkContainerPrimaryIPConfig IPConfiguration
	HostPrimaryIPInfo               HostIPInfo
}

// DeleteNetworkContainerRequest specifies the details about the request to delete a specifc network container.
type HostIPInfo struct {
	Gateway   string
	PrimaryIP string
	Subnet    string
}

type IPConfigRequest struct {
	DesiredIPAddress    string
	OrchestratorContext json.RawMessage
}

func (i IPConfigRequest) String() string {
	return fmt.Sprintf("[IPConfigRequest: DesiredIPAddress %s, OrchestratorContext %s]",
		i.DesiredIPAddress, string(i.OrchestratorContext))
}

// IPConfigResponse is used in CNS IPAM mode as a response to CNI ADD
type IPConfigResponse struct {
	PodIpInfo PodIpInfo
	Response  Response
}

// GetIPAddressesRequest is used in CNS IPAM mode to get the states of IPConfigs
// The IPConfigStateFilter is a slice of IP's to fetch from CNS that match those states
type GetIPAddressesRequest struct {
	IPConfigStateFilter []string
}

// GetIPAddressStateResponse is used in CNS IPAM mode as a response to get IP address state
type GetIPAddressStateResponse struct {
	IPAddresses []IPAddressState
	Response    Response
}

// GetIPAddressStatusResponse is used in CNS IPAM mode as a response to get IP address, state and Pod info
type GetIPAddressStatusResponse struct {
	IPConfigurationStatus[] IPConfigurationStatus
	Response Response
}

// IPAddressState Only used in the GetIPConfig API to return IP's that match a filter
type IPAddressState struct {
	IPAddress string
	State     string
}

// DeleteNetworkContainerRequest specifies the details about the request to delete a specifc network container.
type DeleteNetworkContainerRequest struct {
	NetworkContainerid string
}

// DeleteNetworkContainerResponse describes the response to delete a specifc network container.
type DeleteNetworkContainerResponse struct {
	Response Response
}

// GetInterfaceForContainerRequest specifies the container ID for which interface needs to be identified.
type GetInterfaceForContainerRequest struct {
	NetworkContainerID string
}

// GetInterfaceForContainerResponse specifies the interface for a given container ID.
type GetInterfaceForContainerResponse struct {
	NetworkContainerVersion string
	NetworkInterface        NetworkInterface
	CnetAddressSpace        []IPSubnet
	DNSServers              []string
	Response                Response
}

// AttachContainerToNetworkResponse specifies response of attaching network container to network.
type AttachContainerToNetworkResponse struct {
	Response Response
}

// DetachContainerFromNetworkResponse specifies response of detaching network container from network.
type DetachContainerFromNetworkResponse struct {
	Response Response
}

// NetworkInterface specifies the information that can be used to unquely identify an interface.
type NetworkInterface struct {
	Name      string
	IPAddress string
}

// PublishNetworkContainerRequest specifies request to publish network container via NMAgent.
type PublishNetworkContainerRequest struct {
	NetworkID                         string
	NetworkContainerID                string
	JoinNetworkURL                    string
	CreateNetworkContainerURL         string
	CreateNetworkContainerRequestBody []byte
}

// PublishNetworkContainerResponse specifies the response to publish network container request.
type PublishNetworkContainerResponse struct {
	Response            Response
	PublishErrorStr     string
	PublishStatusCode   int
	PublishResponseBody []byte
}

// UnpublishNetworkContainerRequest specifies request to unpublish network container via NMAgent.
type UnpublishNetworkContainerRequest struct {
	NetworkID                 string
	NetworkContainerID        string
	JoinNetworkURL            string
	DeleteNetworkContainerURL string
}

// UnpublishNetworkContainerResponse specifies the response to unpublish network container request.
type UnpublishNetworkContainerResponse struct {
	Response              Response
	UnpublishErrorStr     string
	UnpublishStatusCode   int
	UnpublishResponseBody []byte
}

// ValidAclPolicySetting - Used to validate ACL policy
type ValidAclPolicySetting struct {
	Protocols       string `json:","`
	Action          string `json:","`
	Direction       string `json:","`
	LocalAddresses  string `json:","`
	RemoteAddresses string `json:","`
	LocalPorts      string `json:","`
	RemotePorts     string `json:","`
	RuleType        string `json:","`
	Priority        uint16 `json:","`
}

const (
	ActionTypeAllow  string = "Allow"
	ActionTypeBlock  string = "Block"
	DirectionTypeIn  string = "In"
	DirectionTypeOut string = "Out"
)

// Validate - Validates network container request policies
func (networkContainerRequestPolicy *NetworkContainerRequestPolicies) Validate() error {
	// validate ACL policy
	if networkContainerRequestPolicy != nil {
		if strings.EqualFold(networkContainerRequestPolicy.Type, "ACLPolicy") && strings.EqualFold(networkContainerRequestPolicy.EndpointType, "APIPA") {
			var requestedAclPolicy ValidAclPolicySetting
			if err := json.Unmarshal(networkContainerRequestPolicy.Settings, &requestedAclPolicy); err != nil {
				return fmt.Errorf("ACL policy failed to pass validation with error: %+v ", err)
			}
			//Deny request if ACL Action is empty
			if len(strings.TrimSpace(string(requestedAclPolicy.Action))) == 0 {
				return fmt.Errorf("Action field cannot be empty in ACL Policy")
			}
			//Deny request if ACL Action is not Allow or Deny
			if !strings.EqualFold(requestedAclPolicy.Action, ActionTypeAllow) && !strings.EqualFold(requestedAclPolicy.Action, ActionTypeBlock) {
				return fmt.Errorf("Only Allow or Block is supported in Action field")
			}
			//Deny request if ACL Direction is empty
			if len(strings.TrimSpace(string(requestedAclPolicy.Direction))) == 0 {
				return fmt.Errorf("Direction field cannot be empty in ACL Policy")
			}
			//Deny request if ACL direction is not In or Out
			if !strings.EqualFold(requestedAclPolicy.Direction, DirectionTypeIn) && !strings.EqualFold(requestedAclPolicy.Direction, DirectionTypeOut) {
				return fmt.Errorf("Only In or Out is supported in Direction field")
			}
			if requestedAclPolicy.Priority == 0 {
				return fmt.Errorf("Priority field cannot be empty in ACL Policy")
			}
		} else {
			return fmt.Errorf("Only ACL Policies on APIPA endpoint supported")
		}
	}
	return nil
}

// NodeInfoResponse - Struct to hold the node info response.
type NodeInfoResponse struct {
	NetworkContainers []CreateNetworkContainerRequest
}

// NodeRegisterRequest - Struct to hold the node register request.
type NodeRegisterRequest struct {
	NumCPU               int
	NmAgentSupportedApis []string
}
