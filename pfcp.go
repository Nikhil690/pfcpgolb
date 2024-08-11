package pfcpgolb

import (
    "net"
	logger "github.com/sirupsen/logrus"
)

type MessageType uint8

type PFCPMessage struct {
	Header Header
	Body   interface{}
}

const (
    PFCP_HEARTBEAT_REQUEST              MessageType = 1
    PFCP_HEARTBEAT_RESPONSE             MessageType = 2
    PFCP_PFD_MANAGEMENT_REQUEST         MessageType = 3
    PFCP_PFD_MANAGEMENT_RESPONSE        MessageType = 4
    PFCP_ASSOCIATION_SETUP_REQUEST      MessageType = 5
    PFCP_ASSOCIATION_SETUP_RESPONSE     MessageType = 6
    PFCP_ASSOCIATION_UPDATE_REQUEST     MessageType = 7
    PFCP_ASSOCIATION_UPDATE_RESPONSE    MessageType = 8
    PFCP_ASSOCIATION_RELEASE_REQUEST    MessageType = 9
    PFCP_ASSOCIATION_RELEASE_RESPONSE   MessageType = 10
    PFCP_VERSION_NOT_SUPPORTED_RESPONSE MessageType = 11
    PFCP_NODE_REPORT_REQUEST            MessageType = 12
    PFCP_NODE_REPORT_RESPONSE           MessageType = 13
    PFCP_SESSION_SET_DELETION_REQUEST   MessageType = 14
    PFCP_SESSION_SET_DELETION_RESPONSE  MessageType = 15

    PFCP_SESSION_ESTABLISHMENT_REQUEST  MessageType = 50
    PFCP_SESSION_ESTABLISHMENT_RESPONSE MessageType = 51
    PFCP_SESSION_MODIFICATION_REQUEST   MessageType = 52
    PFCP_SESSION_MODIFICATION_RESPONSE  MessageType = 53
    PFCP_SESSION_DELETION_REQUEST       MessageType = 54
    PFCP_SESSION_DELETION_RESPONSE      MessageType = 55
    PFCP_SESSION_REPORT_REQUEST         MessageType = 56
    PFCP_SESSION_REPORT_RESPONSE        MessageType = 57
)

type PFCPAssociationSetupRequest struct {
    NodeID                         *NodeID                         `tlv:"60"`
    RecoveryTimeStamp              *RecoveryTimeStamp              `tlv:"96"`
    UPFunctionFeatures             *UPFunctionFeatures             `tlv:"43"`
    CPFunctionFeatures             *CPFunctionFeatures             `tlv:"89"`
    UserPlaneIPResourceInformation *UserPlaneIPResourceInformation `tlv:"116"`
}


type PFCPAssociationSetupResponse struct {
    NodeID                         *NodeID                         `tlv:"60"`
    Cause                          *Cause                          `tlv:"19"`
    RecoveryTimeStamp              *RecoveryTimeStamp              `tlv:"96"`
    UPFunctionFeatures             *UPFunctionFeatures             `tlv:"43"`
    CPFunctionFeatures             *CPFunctionFeatures             `tlv:"89"`
    UserPlaneIPResourceInformation *UserPlaneIPResourceInformation `tlv:"116"`
}

type PFCPAssociationReleaseRequest struct {
    NodeID *NodeID `tlv:"60"`
}

type PFCPAssociationReleaseResponse struct {
    NodeID *NodeID `tlv:"60"`
    Cause  *Cause  `tlv:"19"`
}

type CreatePDR struct {
    PDRID                   *PacketDetectionRuleID   `tlv:"56"`
    Precedence              *Precedence              `tlv:"29"`
    PDI                     *PDI                              `tlv:"2"`
    OuterHeaderRemoval      *OuterHeaderRemoval      `tlv:"95"`
    FARID                   *FARID                   `tlv:"108"`
    URRID                   []*URRID                 `tlv:"81"`
    QERID                   []*QERID                 `tlv:"109"`
    ActivatePredefinedRules *ActivatePredefinedRules `tlv:"106"`
}

type PDI struct {
    SourceInterface               *SourceInterface               `tlv:"20"`
    LocalFTEID                    *FTEID                         `tlv:"21"`
    NetworkInstance               *NetworkInstance               `tlv:"22"`
    UEIPAddress                   *UEIPAddress                   `tlv:"93"`
    TrafficEndpointID             *TrafficEndpointID             `tlv:"131"`
    SDFFilter                     *SDFFilter                     `tlv:"23"`
    ApplicationID                 *ApplicationID                 `tlv:"24"`
    EthernetPDUSessionInformation *EthernetPDUSessionInformation `tlv:"142"`
    EthernetPacketFilter          *EthernetPacketFilter                   `tlv:"132"`
    QFI                           []*QFI                         `tlv:"124"`
    FramedRoute                   *FramedRoute                   `tlv:"153"`
    FramedRouting                 *FramedRouting                 `tlv:"154"`
    FramedIPv6Route               *FramedIPv6Route               `tlv:"155"`
}

type CreateFAR struct {
    FARID                 *FARID                 `tlv:"108"`
    ApplyAction           *ApplyAction           `tlv:"44"`
    ForwardingParameters  *ForwardingParametersIEInFAR    `tlv:"4"`
    DuplicatingParameters *DuplicatingParameters `tlv:"5"`
    BARID                 *BARID                 `tlv:"88"`
}

type ForwardingParametersIEInFAR struct {
    DestinationInterface    *DestinationInterface  `tlv:"42"`
    NetworkInstance         *NetworkInstance       `tlv:"22"`
    RedirectInformation     *RedirectInformation   `tlv:"38"`
    OuterHeaderCreation     *OuterHeaderCreation   `tlv:"84"`
    TransportLevelMarking   *TransportLevelMarking `tlv:"30"`
    ForwardingPolicy        *ForwardingPolicy      `tlv:"41"`
    HeaderEnrichment        *HeaderEnrichment      `tlv:"98"`
    LinkedTrafficEndpointID *TrafficEndpointID     `tlv:"131"`
    Proxying                *Proxying              `tlv:"137"`
}

type CreateQER struct {
    QERID              *QERID              `tlv:"109"`
    QERCorrelationID   *QERCorrelationID   `tlv:"28"`
    GateStatus         *GateStatus         `tlv:"25"`
    MaximumBitrate     *MBR                `tlv:"26"`
    GuaranteedBitrate  *GBR                `tlv:"27"`
    PacketRate         *PacketRate         `tlv:"94"`
    DLFlowLevelMarking *DLFlowLevelMarking `tlv:"97"`
    QoSFlowIdentifier  *QFI                `tlv:"124"`
    ReflectiveQoS      *RQI                `tlv:"123"`
}

type UpdatePDR struct {
    PDRID                     *PacketDetectionRuleID     `tlv:"56"`
    OuterHeaderRemoval        *OuterHeaderRemoval        `tlv:"95"`
    Precedence                *Precedence                `tlv:"29"`
    PDI                       *PDI                                `tlv:"2"`
    FARID                     *FARID                     `tlv:"108"`
    URRID                     []*URRID                   `tlv:"81"`
    QERID                     []*QERID                   `tlv:"109"`
    ActivatePredefinedRules   *ActivatePredefinedRules   `tlv:"106"`
    DeactivatePredefinedRules *DeactivatePredefinedRules `tlv:"107"`
}

type UpdateFAR struct {
    FARID                       *FARID                       `tlv:"108"`
    ApplyAction                 *ApplyAction                 `tlv:"44"`
    UpdateForwardingParameters  *UpdateForwardingParametersIEInFAR    `tlv:"11"`
    UpdateDuplicatingParameters *UpdateDuplicatingParameters `tlv:"105"`
    BARID                       *BARID                       `tlv:"88"`
}

type UpdateForwardingParametersIEInFAR struct {
    DestinationInterface    *DestinationInterface  `tlv:"42"`
    NetworkInstance         *NetworkInstance       `tlv:"22"`
    RedirectInformation     *RedirectInformation   `tlv:"38"`
    OuterHeaderCreation     *OuterHeaderCreation   `tlv:"84"`
    TransportLevelMarking   *TransportLevelMarking `tlv:"30"`
    ForwardingPolicy        *ForwardingPolicy      `tlv:"41"`
    HeaderEnrichment        *HeaderEnrichment      `tlv:"98"`
    PFCPSMReqFlags          *PFCPSMReqFlags        `tlv:"49"`
    LinkedTrafficEndpointID *TrafficEndpointID     `tlv:"131"`
}

type CreateTrafficEndpoint struct {
	TrafficEndpointID             *TrafficEndpointID             `tlv:"131"`
	LocalFTEID                    *FTEID                         `tlv:"21"`
	NetworkInstance               *NetworkInstance               `tlv:"22"`
	UEIPAddress                   *UEIPAddress                   `tlv:"93"`
	EthernetPDUSessionInformation *EthernetPDUSessionInformation `tlv:"142"`
	FramedRoute                   *FramedRoute                   `tlv:"153"`
	FramedRouting                 *FramedRouting                 `tlv:"154"`
	FramedIPv6Route               *FramedIPv6Route               `tlv:"155"`
}

type PFCPSessionEstablishmentRequest struct {
    NodeID                   *NodeID                   `tlv:"60"`
    CPFSEID                  *FSEID                    `tlv:"57"`
    CreatePDR                []*CreatePDR                       `tlv:"1"`
    CreateFAR                []*CreateFAR                       `tlv:"3"`
    CreateQER                []*CreateQER                       `tlv:"7"`
    CreateTrafficEndpoint    *CreateTrafficEndpoint             `tlv:"127"`
    PDNType                  *PDNType                  `tlv:"113"`
    UserPlaneInactivityTimer *UserPlaneInactivityTimer `tlv:"117"`
    UserID                   *UserID                   `tlv:"141"`
    TraceInformation         *TraceInformation         `tlv:"152"`
}

type LoadControlInformation struct {
    LoadControlSequenceNumber *SequenceNumber `tlv:"52"`
}

type CreatedTrafficEndpoint struct {
    TrafficEndpointID *TrafficEndpointID `tlv:"131"`
    LocalFTEID        *FTEID             `tlv:"21"`
}


type PFCPSessionEstablishmentResponse struct {
    NodeID                     *NodeID            `tlv:"60"`
    Cause                      *Cause             `tlv:"19"`
    OffendingIE                *OffendingIE       `tlv:"40"`
    UPFSEID                    *FSEID             `tlv:"57"`
    CreatedPDR                 *CreatedPDR                 `tlv:"8"`
    LoadControlInformation     *LoadControlInformation     `tlv:"51"`
    FailedRuleID               *FailedRuleID      `tlv:"114"`
    CreatedTrafficEndpoint     *CreatedTrafficEndpoint     `tlv:"128"`
}

type CreatedPDR struct {
    PDRID      *PacketDetectionRuleID `tlv:"56"`
    LocalFTEID *FTEID                 `tlv:"21"`
}

type RemoveTrafficEndpoint struct {
    TrafficEndpointID *TrafficEndpointID `tlv:"131"`
}

type UpdateQER struct {
	QERID              *QERID              `tlv:"109"`
	QERCorrelationID   *QERCorrelationID   `tlv:"28"`
	GateStatus         *GateStatus         `tlv:"25"`
	MaximumBitrate     *MBR                `tlv:"26"`
	GuaranteedBitrate  *GBR                `tlv:"27"`
	PacketRate         *PacketRate         `tlv:"94"`
	DLFlowLevelMarking *DLFlowLevelMarking `tlv:"97"`
	QoSFlowIdentifier  *QFI                `tlv:"124"`
	ReflectiveQoS      *RQI                `tlv:"123"`
}

type UpdateTrafficEndpoint struct {
    TrafficEndpointID *TrafficEndpointID `tlv:"131"`
    LocalFTEID        *FTEID             `tlv:"21"`
    NetworkInstance   *NetworkInstance   `tlv:"22"`
    UEIPAddress       *UEIPAddress       `tlv:"93"`
    FramedRoute       *FramedRoute       `tlv:"153"`
    FramedRouting     *FramedRouting     `tlv:"154"`
    FramedIPv6Route   *FramedIPv6Route   `tlv:"155"`
}


type PFCPSessionModificationRequest struct {
    CPFSEID                  *FSEID                          `tlv:"57"`
    RemovePDR                []*RemovePDR                             `tlv:"15"`
    RemoveFAR                []*RemoveFAR                             `tlv:"16"`
    // RemoveURR                []*RemoveURR                             `tlv:"17"`
    // RemoveQER                []*RemoveQER                             `tlv:"18"`
    // RemoveBAR                []*RemoveBAR                             `tlv:"87"`
    RemoveTrafficEndpoint    *RemoveTrafficEndpoint                   `tlv:"130"`
    CreatePDR                []*CreatePDR                             `tlv:"1"`
    CreateFAR                []*CreateFAR                             `tlv:"3"`
    // CreateURR                []*CreateURR                             `tlv:"6"`
    CreateQER                []*CreateQER                             `tlv:"7"`
    // CreateBAR                []*CreateBAR                             `tlv:"85"`
    CreateTrafficEndpoint    *CreateTrafficEndpoint                   `tlv:"127"`
    UpdatePDR                []*UpdatePDR                             `tlv:"9"`
    UpdateFAR                []*UpdateFAR                             `tlv:"10"`
    // UpdateURR                []*UpdateURR                             `tlv:"13"`
    UpdateQER                []*UpdateQER                             `tlv:"14"`
    // UpdateBAR                *UpdateBARPFCPSessionModificationRequest `tlv:"86"`
    UpdateTrafficEndpoint    *UpdateTrafficEndpoint                   `tlv:"129"`
    PFCPSMReqFlags           *PFCPSMReqFlags                 `tlv:"49"`
    // QueryURR                 []*QueryURR                              `tlv:"77"`
    UserPlaneInactivityTimer *UserPlaneInactivityTimer       `tlv:"117"`
    // QueryURRReference        *QueryURRReference              `tlv:"125"`
    TraceInformation         *TraceInformation               `tlv:"152"`
}

type RemovePDR struct {
    PDRID *PacketDetectionRuleID `tlv:"56"`
}


type RemoveFAR struct {
    FARID *FARID `tlv:"108"`
}

type PFCPSessionModificationResponse struct {
    Cause                             *Cause                               `tlv:"19"`
    OffendingIE                       *OffendingIE                         `tlv:"40"`
    CreatedPDR                        *CreatedPDR                                   `tlv:"8"`
    LoadControlInformation            *LoadControlInformation                       `tlv:"51"`
    // OverloadControlInformation        *OverloadControlInformation                   `tlv:"54"`
    // UsageReport                       []*UsageReportPFCPSessionModificationResponse `tlv:"78"`
    FailedRuleID                      *FailedRuleID                        `tlv:"114"`
    // AdditionalUsageReportsInformation *AdditionalUsageReportsInformation   `tlv:"126"`
    CreatedUpdatedTrafficEndpoint     *CreatedTrafficEndpoint                       `tlv:"128"`
}

type PFCPSessionDeletionRequest struct{}

type PFCPSessionDeletionResponse struct {
    Cause                      *Cause                           `tlv:"19"`
    OffendingIE                *OffendingIE                     `tlv:"40"`
    LoadControlInformation     *LoadControlInformation                   `tlv:"51"`
    // OverloadControlInformation *OverloadControlInformation               `tlv:"54"`
    // UsageReport                []*UsageReportPFCPSessionDeletionResponse `tlv:"79"`
}

type HeartbeatRequest struct {
    RecoveryTimeStamp *RecoveryTimeStamp `tlv:"96"`
}

type HeartbeatResponse struct {
    RecoveryTimeStamp *RecoveryTimeStamp `tlv:"96"`
}


type Header struct {
    Version         uint8
    MP              uint8
    S               uint8
    MessageType     MessageType
    MessageLength   uint16
    SEID            uint64
    SequenceNumber  uint32
    MessagePriority uint8
}

const PfcpVersion uint8 = 1

const (
    SEID_NOT_PRESENT = 0
    SEID_PRESENT     = 1
)

func (m *Message) MessageType() MessageType {
	return m.PfcpMessage.Header.MessageType
}

func (n *NodeID) ResolveNodeIdToIp() net.IP {
	switch n.NodeIdType {
	case NodeIdTypeIpv4Address, NodeIdTypeIpv6Address:
		return n.IP
	case NodeIdTypeFqdn:
		if ns, err := net.LookupHost(n.FQDN); err != nil {
			logger.Warnf("Host lookup failed: %+v", err)
			return net.IPv4zero
		} else {
			return net.ParseIP(ns[0])
		}
	default:
		return net.IPv4zero
	}
}