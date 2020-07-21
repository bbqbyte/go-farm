package mqtt

import "fmt"

// The CONNACK Packet is the packet sent by the Server in response to a CONNECT Packet
// received from a Client. The first packet sent from the Server to the Client MUST
// be a CONNACK Packet [MQTT-3.2.0-1].
//
// If the Client does not receive a CONNACK Packet from the Server within a reasonable
// amount of time, the Client SHOULD close the Network Connection. A "reasonable" amount
// of time depends on the type of application and the communications infrastructure.
type ConnackPacket struct {
	header

	sessionPresent bool
	returnCode     ConnackCode
}

// NewConnackPacket creates a new CONNACK message
func NewConnackPacket() *ConnackPacket {
	msg := &ConnackPacket{}
	msg.SetType(CONNACK)

	return msg
}

// String returns a string representation of the CONNACK packet
func (cp ConnackPacket) String() string {
	return fmt.Sprintf("%s, Session Present=%t, Return code=%q\n", cp.header, cp.sessionPresent, cp.returnCode)
}

// SessionPresent returns the session present flag value
func (cp *ConnackPacket) SessionPresent() bool {
	return cp.sessionPresent
}

// SetSessionPresent sets the value of the session present flag
func (cp *ConnackPacket) SetSessionPresent(v bool) {
	if v {
		cp.sessionPresent = true
	} else {
		cp.sessionPresent = false
	}
}

// ReturnCode returns the return code received for the CONNECT packet. The return
// type is an error
func (cp *ConnackPacket) ReturnCode() ConnackCode {
	return cp.returnCode
}

func (cp *ConnackPacket) SetReturnCode(ret ConnackCode) {
	cp.returnCode = ret
}

func (cp *ConnackPacket) Len() int {
	ml := cp.msglen()

	if err := cp.SetRemainingLength(int32(ml)); err != nil {
		return 0
	}

	return cp.header.msglen() + ml
}

func (cp *ConnackPacket) Decode(src []byte) (int, error) {
	total := 0

	n, err := cp.header.decode(src)
	total += n
	if err != nil {
		return total, err
	}

	b := src[total]

	if b&254 != 0 {
		return 0, fmt.Errorf("connack/Decode: Bits 7-1 in Connack Acknowledge Flags byte (1) are not 0")
	}

	cp.sessionPresent = b&0x1 == 1
	total++

	b = src[total]

	// Read return code
	if b > 5 {
		return 0, fmt.Errorf("connack/Decode: Invalid CONNACK return code (%d)", b)
	}

	cp.returnCode = ConnackCode(b)
	total++

	return total, nil
}

func (cp *ConnackPacket) Encode(dst []byte) (int, error) {
	// CONNACK remaining length fixed at 2 bytes
	hl := cp.header.msglen()
	// variable header 2 bytes
	ml := cp.msglen()

	if len(dst) < hl+ml {
		return 0, fmt.Errorf("connack/Encode: Insufficient buffer size. Expecting %d, got %d.", hl+ml, len(dst))
	}

	if err := cp.SetRemainingLength(int32(ml)); err != nil {
		return 0, err
	}

	total := 0

	n, err := cp.header.encode(dst[total:])
	total += n
	if err != nil {
		return 0, err
	}

	if cp.sessionPresent {
		dst[total] = 1
	}
	total++

	if cp.returnCode > 5 {
		return total, fmt.Errorf("connack/Encode: Invalid CONNACK return code (%d)", cp.returnCode)
	}

	dst[total] = cp.returnCode.Value()
	total++

	return total, nil
}

func (cp *ConnackPacket) msglen() int {
	return 2
}
