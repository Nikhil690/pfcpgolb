package pfcpgolb

import (
	"errors"
	"net"
	"fmt"
	logger "github.com/sirupsen/logrus"

)

func NewPfcpServer(addr string) *PfcpServer {
    server := PfcpServer{Addr: addr}
    return &server
}

func (t *ConsumerTable) LoadOrStore(consumerAddr string, storeTable *TxTable) (*TxTable, bool) {
	txTable, loaded := t.m.LoadOrStore(consumerAddr, storeTable)
	return txTable.(*TxTable), loaded
}

func (pfcpServer *PfcpServer) Listen() error {
	var serverIp net.IP
	if pfcpServer.Addr == "" {
		serverIp = net.IPv4zero
	} else {
		serverIp = net.ParseIP(pfcpServer.Addr)
	}

	addr := &net.UDPAddr{
		IP:   serverIp,
		Port: PFCP_PORT,
	}

	conn, err := net.ListenUDP("udp", addr)
	pfcpServer.Conn = conn
	return err
}

func NewMessage(remoteAddr *net.UDPAddr, pfcpmessage *PFCPMessage) (msg *Message) {
	return &Message{
		RemoteAddr:  remoteAddr,
		PfcpMessage: pfcpmessage,
	}
}

func (pfcpServer *PfcpServer) ReadFrom() (*Message, error) {
	buf := make([]byte, PFCP_MAX_UDP_LEN)
	n, addr, err := pfcpServer.Conn.ReadFromUDP(buf)
	if err != nil {
		return nil, err
	}

	pfcpMsg := &PFCPMessage{}
	msg := NewMessage(addr, pfcpMsg)

	err = pfcpMsg.Unmarshal(buf[:n])
	if err != nil {
		return msg, err
	}

	if pfcpMsg.IsRequest() {
		// Todo: Implement SendingResponse type of reliable delivery
		tx, err := pfcpServer.FindTransaction(pfcpMsg, addr)
		if err != nil {
			return msg, err
		}
		if tx != nil {
			// tx != nil => Already Replied => Resend Request
			tx.EventChannel <- ReceiveEvent{
				Type:       ReceiveEventTypeResendRequest,
				RemoteAddr: addr,
				RcvMsg:     pfcpMsg,
			}
			return msg, ErrReceivedResentRequest
		} else {
			// tx == nil => New Request
			return msg, nil
		}
	} else if pfcpMsg.IsResponse() {
		tx, err := pfcpServer.FindTransaction(pfcpMsg, pfcpServer.Conn.LocalAddr().(*net.UDPAddr))
		if err != nil {
			return msg, err
		}

		tx.EventChannel <- ReceiveEvent{
			Type:       ReceiveEventTypeValidResponse,
			RemoteAddr: addr,
			RcvMsg:     pfcpMsg,
		}
	}

	return msg, nil
}

func (t *ConsumerTable) Load(consumerAddr string) (*TxTable, bool) {
	txTable, ok := t.m.Load(consumerAddr)
	if ok {
		return txTable.(*TxTable), ok
	}
	return nil, false
}

func (pfcpServer *PfcpServer) StartReqTxLifeCycle(tx *Transaction) (resMsg *Message, err error) {
	defer func() {
		// End Transaction
		rmErr := pfcpServer.RemoveTransaction(tx)
		if rmErr != nil {
			logger.Warnf("RemoveTransaction error: %+v", rmErr)
		}
	}()

	// Start Transaction
	event, err := tx.StartSendingRequest()
	if err != nil {
		return nil, err
	}
		return NewMessage(event.RemoteAddr, event.RcvMsg), nil
}

func (pfcpServer *PfcpServer) StartResTxLifeCycle(tx *Transaction) {
	// Start Transaction
	err := tx.StartSendingResponse()
	if err != nil {
		logger.Warnf("SendingResponse error: %+v", err)
		return
	}
	// End Transaction
	err = pfcpServer.RemoveTransaction(tx)
	if err != nil {
		logger.Warnf("RemoveTransaction error: %+v", err)
	}
}

func (pfcpServer *PfcpServer) RemoveTransaction(tx *Transaction) (err error) {
	logger.Traceln("In RemoveTransaction")
	consumerAddr := tx.ConsumerAddr
	txTable, _ := pfcpServer.ConsumerTable.Load(consumerAddr)

	if txTmp, exist := txTable.Load(tx.SequenceNumber); exist {
		tx = txTmp
		if tx.TxType == SendingRequest {
			logger.Debugf("Remove Request Transaction [%d]", tx.SequenceNumber)
		} else if tx.TxType == SendingResponse {
			logger.Debugf("Remove Request Transaction [%d]", tx.SequenceNumber)
		}

		txTable.Delete(tx.SequenceNumber)
	} else {
		logger.Warnln("In RemoveTransaction")
		logger.Warnln("Consumer IP: ", consumerAddr)
		logger.Warnln("Sequence number ", tx.SequenceNumber, " doesn't exist!")
		err = fmt.Errorf("remove tx error: transaction [%d] doesn't exist", tx.SequenceNumber)
	}

	logger.Traceln("End RemoveTransaction")
	return
}


func (pfcpServer *PfcpServer) WriteRequestTo(reqMsg *PFCPMessage, addr *net.UDPAddr) (resMsg *Message, err error) {
	if !reqMsg.IsRequest() {
		return nil, errors.New("not a request message")
	}

	buf, err := reqMsg.Marshal()
	if err != nil {
		return nil, err
	}

	tx := NewTransaction(reqMsg, buf, pfcpServer.Conn, addr)

	err = pfcpServer.PutTransaction(tx)
	if err != nil {
		return nil, err
	}

	return pfcpServer.StartReqTxLifeCycle(tx)
}

func (pfcpServer *PfcpServer) WriteResponseTo(resMsg *PFCPMessage, addr *net.UDPAddr) {
	if !resMsg.IsResponse() {
		logger.Warn("not a response message")
		return
	}

	buf, err := resMsg.Marshal()
	if err != nil {
		logger.Warnf("marshal error: %+v", err)
		return
	}

	tx := NewTransaction(resMsg, buf, pfcpServer.Conn, addr)

	err = pfcpServer.PutTransaction(tx)
	if err != nil {
		logger.Warnf("PutTransaction error: %+v", err)
		return
	}

	go pfcpServer.StartResTxLifeCycle(tx)
}

func (pfcpServer *PfcpServer) Close() error {
	return pfcpServer.Conn.Close()
}