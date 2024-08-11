package pfcpgolb

import (
	"fmt"
	"github.com/Nikhil690/pfcpgolb/tlv"
)

func (m *PFCPMessage) Unmarshal(data []byte) error {
	if err := m.Header.UnmarshalBinary(data); err != nil {
		return fmt.Errorf("pfcp: unmarshal msg failed: %s", err)
	}

	// Check Message Length field in header
	if int(m.Header.MessageLength) != len(data)-4 {
		return fmt.Errorf("Message Length Incorrect: Expected %d, got %d", m.Header.MessageLength, len(data)-4)
	}
	switch m.Header.MessageType {
	case PFCP_HEARTBEAT_REQUEST:
		Body := HeartbeatRequest{}
		if err := tlv.Unmarshal(data[m.Header.Len():], &Body); err != nil {
			return err
		}
		m.Body = Body
	case PFCP_HEARTBEAT_RESPONSE:
		Body := HeartbeatResponse{}
		if err := tlv.Unmarshal(data[m.Header.Len():], &Body); err != nil {
			return err
		}
		m.Body = Body
	case PFCP_ASSOCIATION_SETUP_REQUEST:
		Body := PFCPAssociationSetupRequest{}
		if err := tlv.Unmarshal(data[m.Header.Len():], &Body); err != nil {
			return err
		}
		m.Body = Body
	case PFCP_ASSOCIATION_SETUP_RESPONSE:
		Body := PFCPAssociationSetupResponse{}
		if err := tlv.Unmarshal(data[m.Header.Len():], &Body); err != nil {
			return err
		}
		m.Body = Body
	case PFCP_ASSOCIATION_RELEASE_REQUEST:
		Body := PFCPAssociationReleaseRequest{}
		if err := tlv.Unmarshal(data[m.Header.Len():], &Body); err != nil {
			return err
		}
		m.Body = Body
	case PFCP_ASSOCIATION_RELEASE_RESPONSE:
		Body := PFCPAssociationReleaseResponse{}
		if err := tlv.Unmarshal(data[m.Header.Len():], &Body); err != nil {
			return err
		}
		m.Body = Body
	case PFCP_SESSION_ESTABLISHMENT_REQUEST:
		Body := PFCPSessionEstablishmentRequest{}
		if err := tlv.Unmarshal(data[m.Header.Len():], &Body); err != nil {
			return err
		}
		m.Body = Body
	case PFCP_SESSION_ESTABLISHMENT_RESPONSE:
		Body := PFCPSessionEstablishmentResponse{}
		if err := tlv.Unmarshal(data[m.Header.Len():], &Body); err != nil {
			return err
		}
		m.Body = Body
	case PFCP_SESSION_MODIFICATION_REQUEST:
		Body := PFCPSessionModificationRequest{}
		if err := tlv.Unmarshal(data[m.Header.Len():], &Body); err != nil {
			return err
		}
		m.Body = Body
	case PFCP_SESSION_MODIFICATION_RESPONSE:
		Body := PFCPSessionModificationResponse{}
		if err := tlv.Unmarshal(data[m.Header.Len():], &Body); err != nil {
			return err
		}
		m.Body = Body
	case PFCP_SESSION_DELETION_REQUEST:
		Body := PFCPSessionDeletionRequest{}
		if err := tlv.Unmarshal(data[m.Header.Len():], &Body); err != nil {
			return err
		}
		m.Body = Body
	case PFCP_SESSION_DELETION_RESPONSE:
		Body := PFCPSessionDeletionResponse{}
		if err := tlv.Unmarshal(data[m.Header.Len():], &Body); err != nil {
			return err
		}
		m.Body = Body
	default:
		return fmt.Errorf("pfcp: unmarshal msg type %d not supported", m.Header.MessageType)
	}
	return nil
}

func (message *PFCPMessage) IsRequest() (IsRequest bool) {
	switch message.Header.MessageType {
	case PFCP_HEARTBEAT_REQUEST:
		IsRequest = true
	case PFCP_PFD_MANAGEMENT_REQUEST:
		IsRequest = true
	case PFCP_ASSOCIATION_SETUP_REQUEST:
		IsRequest = true
	case PFCP_ASSOCIATION_UPDATE_REQUEST:
		IsRequest = true
	case PFCP_ASSOCIATION_RELEASE_REQUEST:
		IsRequest = true
	case PFCP_NODE_REPORT_REQUEST:
		IsRequest = true
	case PFCP_SESSION_SET_DELETION_REQUEST:
		IsRequest = true
	case PFCP_SESSION_ESTABLISHMENT_REQUEST:
		IsRequest = true
	case PFCP_SESSION_MODIFICATION_REQUEST:
		IsRequest = true
	case PFCP_SESSION_DELETION_REQUEST:
		IsRequest = true
	case PFCP_SESSION_REPORT_REQUEST:
		IsRequest = true
	default:
		IsRequest = false
	}

	return
}

func (message *PFCPMessage) IsResponse() (IsResponse bool) {
	IsResponse = false
	switch message.Header.MessageType {
	case PFCP_HEARTBEAT_RESPONSE:
		IsResponse = true
	case PFCP_PFD_MANAGEMENT_RESPONSE:
		IsResponse = true
	case PFCP_ASSOCIATION_SETUP_RESPONSE:
		IsResponse = true
	case PFCP_ASSOCIATION_UPDATE_RESPONSE:
		IsResponse = true
	case PFCP_ASSOCIATION_RELEASE_RESPONSE:
		IsResponse = true
	case PFCP_NODE_REPORT_RESPONSE:
		IsResponse = true
	case PFCP_SESSION_SET_DELETION_RESPONSE:
		IsResponse = true
	case PFCP_SESSION_ESTABLISHMENT_RESPONSE:
		IsResponse = true
	case PFCP_SESSION_MODIFICATION_RESPONSE:
		IsResponse = true
	case PFCP_SESSION_DELETION_RESPONSE:
		IsResponse = true
	case PFCP_SESSION_REPORT_RESPONSE:
		IsResponse = true
	default:
		IsResponse = false
	}

	return
}
