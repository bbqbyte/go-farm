package mqtt

import (
	"fmt"
)

// A PUBACK Packet is the response to a PUBLISH Packet with QoS level 1.
type PubackPacket struct {
	header
}

// NewPubackPacket creates a new PUBACK packet.
func NewPubackPacket() *PubackPacket {
	msg := &PubackPacket{}
	msg.SetType(PUBACK)

	return msg
}

func (pp PubackPacket) String() string {
	return fmt.Sprintf("%s, Packet ID=%d", pp.header, pp.packetId)
}

func (pp *PubackPacket) Len() int {
	ml := pp.msglen()

	if err := pp.SetRemainingLength(int32(ml)); err != nil {
		return 0
	}

	return pp.header.msglen() + ml
}

func (pp *PubackPacket) Decode(src []byte) (int, error) {
	total := 0

	n, err := pp.header.decode(src[total:])
	total += n
	if err != nil {
		return total, err
	}

	//this.packetId = binary.BigEndian.Uint16(src[total:])
	//pp.packetId = src[total : total+2]
	pp.packetId = src[total : total+2]
	total += 2

	return total, nil
}

func (pp *PubackPacket) Encode(dst []byte) (int, error) {
	hl := pp.header.msglen()
	ml := pp.msglen()

	if len(dst) < hl+ml {
		return 0, fmt.Errorf("puback/Encode: Insufficient buffer size. Expecting %d, got %d.", hl+ml, len(dst))
	}

	if err := pp.SetRemainingLength(int32(ml)); err != nil {
		return 0, err
	}

	total := 0

	n, err := pp.header.encode(dst[total:])
	total += n
	if err != nil {
		return total, err
	}

	if copy(dst[total:total+2], pp.packetId) != 2 {
		dst[total], dst[total+1] = 0, 0
	}
	total += 2

	return total, nil
}

func (pp *PubackPacket) msglen() int {
	// packet ID
	return 2
}
