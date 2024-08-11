package pfcpgolb

import (
	"net"
	"time"
)

const (
	PDNTypeIpv4 uint8 = iota + 1
	PDNTypeIpv6
	PDNTypeIpv4v6
	PDNTypeNonIp
	PDNTypeEthernet
)

const (
	OuterHeaderRemovalGtpUUdpIpv4 uint8 = iota
	OuterHeaderRemovalGtpUUdpIpv6
	OuterHeaderRemovalUdpIpv4
	OuterHeaderRemovalUdpIpv6
)

const (
	SourceInterfaceAccess uint8 = iota
	SourceInterfaceCore
	SourceInterfaceSgiLanN6Lan
	SourceInterfaceCpFunction
)

const (
	DestinationInterfaceAccess uint8 = iota
	DestinationInterfaceCore
	DestinationInterfaceSgiLanN6Lan
	DestinationInterfaceCpFunction
	DestinationInterfaceLiFunction
)

const (
	OuterHeaderCreationGtpUUdpIpv4 uint16 = 1
	OuterHeaderCreationGtpUUdpIpv6 uint16 = 1 << 1
	OuterHeaderCreationUdpIpv4     uint16 = 1 << 2
	OuterHeaderCreationUdpIpv6     uint16 = 1 << 3
)

const (
	CauseRequestAccepted uint8 = 1
)

const (
	CauseRequestRejected uint8 = iota + 64
	CauseSessionContextNotFound
	CauseMandatoryIeMissing
	CauseConditionalIeMissing
	CauseInvalidLength
	CauseMandatoryIeIncorrect
	CauseInvalidForwardingPolicy
	CauseInvalidFTeidAllocationOption
	CauseNoEstablishedPfcpAssociation
	CauseRuleCreationModificationFailure
	CausePfcpEntityInCongestion
	CauseNoResourcesAvailable
	CauseServiceNotSupported
	CauseSystemFailure
)



const (
	NodeIdTypeIpv4Address uint8 = iota
	NodeIdTypeIpv6Address
	NodeIdTypeFqdn
)

const (
	GateOpen uint8 = iota
	GateClose
)

type NodeID struct {
	NodeIdType uint8 // 0x00001111
	IP         net.IP
	FQDN       string
}

type SourceInterface struct {
	InterfaceValue uint8 // 0x00001111
}

type FTEID struct {
	Chid        bool
	Ch          bool
	V6          bool
	V4          bool
	Teid        uint32
	Ipv4Address net.IP
	Ipv6Address net.IP
	ChooseId    uint8
}

type FSEID struct {
	V4          bool
	V6          bool
	Seid        uint64
	Ipv4Address net.IP
	Ipv6Address net.IP
}

type NetworkInstance struct {
	NetworkInstance string
	FQDNEncoding    bool
}

type UEIPAddress struct {
	Ipv6d                    bool
	Sd                       bool
	V4                       bool
	V6                       bool
	Ipv4Address              net.IP
	Ipv6Address              net.IP
	Ipv6PrefixDelegationBits uint8
}

type SDFFilter struct {
	Bid                     bool
	Fl                      bool
	Spi                     bool
	Ttc                     bool
	Fd                      bool
	LengthOfFlowDescription uint16
	FlowDescription         []byte
	TosTrafficClass         []byte
	SecurityParameterIndex  []byte
	FlowLabel               []byte
	SdfFilterId             uint32
}

type QFI struct {
	QFI uint8
}

type GateStatus struct {
	ULGate uint8 // 0x00001100
	DLGate uint8 // 0x00000011
}

type GBR struct {
	ULGBR uint64 // 40-bit data
	DLGBR uint64 // 40-bit data
}

type MBR struct {
	ULMBR uint64 // 40-bit data
	DLMBR uint64 // 40-bit data
}

type ApplyAction struct {
	Dupl bool
	Nocp bool
	Buff bool
	Forw bool
	Drop bool
}

type OuterHeaderRemoval struct {
	OuterHeaderRemovalDescription uint8
}

type DestinationInterface struct {
	InterfaceValue uint8 // 0x00001111
}

type OuterHeaderCreation struct {
	OuterHeaderCreationDescription uint16
	Teid                           uint32
	Ipv4Address                    net.IP
	Ipv6Address                    net.IP
	PortNumber                     uint16
}

type UserPlaneIPResourceInformation struct {
	Assosi          bool
	Assoni          bool
	Teidri          uint8 // 0x00011100
	V6              bool
	V4              bool
	TeidRange       uint8
	Ipv4Address     net.IP
	Ipv6Address     net.IP
	NetworkInstance NetworkInstance
	SourceInterface uint8 // 0x00001111
}

type Cause struct {
	CauseValue uint8
}



type RecoveryTimeStamp struct {
	RecoveryTimeStamp time.Time
}

type CPFunctionFeatures struct {
	SupportedFeatures uint8
}

type PacketDetectionRuleID struct {
	RuleId uint16
}

type Precedence struct {
	PrecedenceValue uint32
}

type ApplicationID struct {
	ApplicationIdentifier []byte
}

type URRID struct {
	UrrIdValue uint32
}

type BARID struct {
	BarIdValue uint8
}

type DownlinkDataNotificationDelay struct {
	DelayValue uint8
}

type QERID struct {
	QERID uint32
}

type FARID struct {
	FarIdValue uint32
}

type PFCPSMReqFlags struct {
	Qaurr bool
	Sndem bool
	Drobu bool
}

type ForwardingPolicy struct {
	ForwardingPolicyIdentifierLength uint8
	ForwardingPolicyIdentifier       []byte
}

type OffendingIE struct {
	TypeOfOffendingIe uint16
}

type UPFunctionFeatures struct {
	SupportedFeatures uint16
}

type ActivatePredefinedRules struct {
	PredefinedRulesName []byte
}

type TrafficEndpointID struct {
	TrafficEndpointIdValue uint8
}

type EthernetPDUSessionInformation struct {
	EthernetPDUSessionInformationdata []byte
}

type EthernetFilterID struct {
    EthernetFilterIDdata []byte
}

type EthernetFilterProperties struct {
    EthernetFilterPropertiesdata []byte
}

type MACAddress struct {
    MACAddressdata []byte
}

type Ethertype struct {
    Ethertypedata []byte
}

type CTAG struct {
    CTAGdata []byte
}

type STAG struct {
    STAGdata []byte
}

type FramedRoute struct {
    FramedRoutedata []byte
}

type FramedRouting struct {
    FramedRoutingdata []byte
}

type FramedIPv6Route struct {
    FramedIPv6Routedata []byte
}

type DuplicatingParameters struct {
    DuplicatingParametersdata []byte
}

type RedirectInformation struct {
    RedirectAddressType         uint8 // 0x00001111
    RedirectServerAddressLength uint16
    RedirectServerAddress       []byte
}

type TransportLevelMarking struct {
    TosTrafficClass []byte
}

type HeaderEnrichment struct {
    HeaderType               uint8 // 0x00011111
    LengthOfHeaderFieldName  uint8
    HeaderFieldName          []byte
    LengthOfHeaderFieldValue uint8
    HeaderFieldValue         []byte
}

type Proxying struct {
    Proxyingdata []byte
}

type QERCorrelationID struct {
    QerCorrelationIdValue uint32
}

type PacketRateTimeUnit uint8

type PacketRate struct {
    ULPR       bool
    DLPR       bool
    ULTimeUnit PacketRateTimeUnit
    MaximumUL  uint16
    DLTimeUnit PacketRateTimeUnit
    MaximumDL  uint16
}

type DLFlowLevelMarking struct {
    DLFlowLevelMarkingdata []byte
}

type RQI struct {
    RQI bool
}

type UpdateDuplicatingParameters struct {
    UpdateDuplicatingParametersdata []byte
}

type DeactivatePredefinedRules struct {
    PredefinedRulesName []byte
}

type PDNType struct {
    PdnType uint8 // 0x00000111
}

type UserPlaneInactivityTimer struct {
    UserPlaneInactivityTimerdata []byte
}

type UserID struct {
    UserIDdata []byte
}

type TraceInformation struct {
    TraceInformationdata []byte
}

type SequenceNumber struct {
    SequenceNumberdata []byte
}

type FailedRuleID struct {
    RuleIdType  uint8 // 0x00001111
    RuleIdValue []byte
}



type EthernetPacketFilter struct {
	EthernetFilterID         *EthernetFilterID         `tlv:"138"`
	EthernetFilterProperties *EthernetFilterProperties `tlv:"139"`
	MACAddress               *MACAddress               `tlv:"133"`
	Ethertype                *Ethertype                `tlv:"136"`
	CTAG                     *CTAG                     `tlv:"134"`
	STAG                     *STAG                     `tlv:"135"`
	SDFFilter                *SDFFilter                `tlv:"23"`
}

