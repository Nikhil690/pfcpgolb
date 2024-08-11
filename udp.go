package pfcpgolb

import (
	"errors"
	"net"
    "sync"
)

const (
    PFCP_PORT        = 8805
    PFCP_MAX_UDP_LEN = 2048
)

type Message struct {
    RemoteAddr  *net.UDPAddr
    PfcpMessage *PFCPMessage
}

type ConsumerTable struct {
	m sync.Map // map[string]pfcp.TxTable
}


type PfcpServer struct {
    Addr string
    Conn *net.UDPConn
    // Consumer Table
    // Map Consumer IP to its tx table
    ConsumerTable ConsumerTable
}

var ErrReceivedResentRequest = errors.New("received a request that is re-sent")

type ReceiveEventType uint8

type ReceiveEvent struct {
    Type       ReceiveEventType
    RemoteAddr *net.UDPAddr
    RcvMsg     *PFCPMessage
}

const (
	ReceiveEventTypeResendRequest ReceiveEventType = iota
	ReceiveEventTypeValidResponse
)


