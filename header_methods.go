package pfcpgolb

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net"

	logger "github.com/sirupsen/logrus"
)

func (h *Header) Len() int {
	// Node Related Header
	if int(h.MessageType) < 50 {
		return 8
	}
	return 16
}

func (h *Header) UnmarshalBinary(data []byte) error {
	var tmpBuf uint8
	byteReader := bytes.NewReader(data)
	if err := binary.Read(byteReader, binary.BigEndian, &tmpBuf); err != nil {
		return errors.New("")
	}
	h.Version, h.MP, h.S = tmpBuf>>5, (tmpBuf&0x02)>>1, tmpBuf&0x01
	if err := binary.Read(byteReader, binary.BigEndian, &h.MessageType); err != nil {
		fmt.Printf("Binary write error: %+v", err)
	}
	if err := binary.Read(byteReader, binary.BigEndian, &h.MessageLength); err != nil {
		fmt.Printf("Binary write error: %+v", err)
	}
	if h.S&1 != 0 {
		if err := binary.Read(byteReader, binary.BigEndian, &h.SEID); err != nil {
			fmt.Printf("Binary write error: %+v", err)
		}
	}
	var snAndSpare uint32
	if err := binary.Read(byteReader, binary.BigEndian, &snAndSpare); err != nil {
		fmt.Printf("Binary write error: %+v", err)
	}

	h.SequenceNumber = snAndSpare >> 8

	if h.MP&1 != 0 {
		h.MessagePriority = uint8(snAndSpare&0x00FF) >> 4
	}
	return nil
}

func (pfcpServer *PfcpServer) FindTransaction(msg *PFCPMessage, addr *net.UDPAddr) (*Transaction, error) {
	var tx *Transaction

	fmt.Println("In FindTransaction")
	consumerAddr := addr.String()

	if msg.IsResponse() {
		txTable, exist := pfcpServer.ConsumerTable.Load(consumerAddr)
		if !exist {
			fmt.Printf("In FindTransaction")
			fmt.Printf("Can't find txTable from consumer addr: [%s]", consumerAddr)
			return nil, fmt.Errorf("FindTransaction Error: txTable not found")
		}

		seqNum := msg.Header.SequenceNumber

		tx, exist = txTable.Load(seqNum)
		if !exist {
			fmt.Printf("In FindTransaction")
			fmt.Println("Consumer Addr: ", consumerAddr)
			fmt.Printf("Can't find tx [%d] from txTable: ", seqNum)
			return nil, fmt.Errorf("FindTransaction Error: sequence number [%d] not found", seqNum)
		}
	} else if msg.IsRequest() {
		txTable, exist := pfcpServer.ConsumerTable.Load(consumerAddr)
		if !exist {
			return nil, nil
		}

		seqNum := msg.Header.SequenceNumber

		tx, exist = txTable.Load(seqNum)
		if !exist {
			return nil, nil
		}
	}
	fmt.Printf("End FindTransaction")
	return tx, nil
}

func (h *Header) MarshalBinary() (data []byte, err error) {
	var tmpbuf uint8
	buffer := new(bytes.Buffer)
	tmpbuf = h.Version<<5 | (h.MP&1)<<1 | (h.S & 1)
	if err := binary.Write(buffer, binary.BigEndian, &tmpbuf); err != nil {
		fmt.Printf("Binary write error: %+v", err)
	}
	if err := binary.Write(buffer, binary.BigEndian, &h.MessageType); err != nil {
		fmt.Printf("Binary write error: %+v", err)
	}
	if err := binary.Write(buffer, binary.BigEndian, &h.MessageLength); err != nil {
		fmt.Printf("Binary write error: %+v", err)
	}
	if h.S&1 != 0 {
		if err := binary.Write(buffer, binary.BigEndian, &h.SEID); err != nil {
			fmt.Printf("Binary write error: %+v", err)
		}
	}
	var snAndSpare uint32
	var spareAndMP uint8
	if h.MP&1 != 0 {
		spareAndMP = h.MessagePriority << 4
	} else {
		spareAndMP = 0
	}
	if h.SequenceNumber > (1<<24 - 1) {
		fmt.Printf("Sequence number must be less 24bit integer")
	}

	snAndSpare = h.SequenceNumber<<8 | uint32(spareAndMP)
	if err := binary.Write(buffer, binary.BigEndian, &snAndSpare); err != nil {
		fmt.Printf("Binary write error: %+v", err)
	}
	return buffer.Bytes(), nil
}

func (pfcpServer *PfcpServer) PutTransaction(tx *Transaction) (err error) {
	logger.Traceln("In PutTransaction")

	consumerAddr := tx.ConsumerAddr
	txTable, _ := pfcpServer.ConsumerTable.LoadOrStore(consumerAddr, &TxTable{})

	if _, exist := txTable.LoadOrStore(tx.SequenceNumber, tx); exist {
		logger.Warnln("In PutTransaction")
		logger.Warnln("Consumer Addr: ", consumerAddr)
		logger.Warnln("Sequence number ", tx.SequenceNumber, " already exist!")
		err = fmt.Errorf("Insert tx error: duplicate sequence number %d", tx.SequenceNumber)
	}

	logger.Traceln("End PutTransaction")
	return
}