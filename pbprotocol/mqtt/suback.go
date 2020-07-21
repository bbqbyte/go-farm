package mqtt

import (
	"fmt"
)

// A SUBACK Packet is sent by the Server to the Client to confirm receipt and processing
// of a SUBSCRIBE Packet.
//
// A SUBACK Packet contains a list of return codes, that specify the maximum QoS level
// that was granted in each Subscription that was requested by the SUBSCRIBE.
type SubackPacket struct {
	header

	returnCodes []byte
}

// NewSubackPacket creates a new SUBACK packet.
func NewSubackPacket() *SubackPacket {
	msg := &SubackPacket{}
	msg.SetType(SUBACK)

	return msg
}

// String returns a string representation of the packet.
func (sp SubackPacket) String() string {
	return fmt.Sprintf("%s, Packet ID=%d, Return Codes=%v", sp.header, sp.PacketID(), sp.returnCodes)
}

// ReturnCodes returns the list of QoS returns from the subscriptions sent in the SUBSCRIBE packet.
func (sp *SubackPacket) ReturnCodes() []byte {
	return sp.returnCodes
}

// AddReturnCodes sets the list of QoS returns from the subscriptions sent in the SUBSCRIBE packet.
// An error is returned if any of the QoS values are not valid.
func (sp *SubackPacket) AddReturnCodes(ret []byte) error {
	for _, c := range ret {
		if c != QosAtMostOnce && c != QosAtLeastOnce && c != QosExactlyOnce && c != QosFailure {
			return fmt.Errorf("suback/AddReturnCode: Invalid return code %d. Must be 0, 1, 2, 0x80.", c)
		}

		sp.returnCodes = append(sp.returnCodes, c)
	}

	return nil
}

// AddReturnCode adds a single QoS return value.
func (sp *SubackPacket) AddReturnCode(ret byte) error {
	return sp.AddReturnCodes([]byte{ret})
}

func (sp *SubackPacket) Len() int {
	ml := sp.msglen()

	if err := sp.SetRemainingLength(int32(ml)); err != nil {
		return 0
	}

	return sp.header.msglen() + ml
}

func (sp *SubackPacket) Decode(src []byte) (int, error) {
	total := 0

	hn, err := sp.header.decode(src[total:])
	total += hn
	if err != nil {
		return total, err
	}

	//this.packetId = binary.BigEndian.Uint16(src[total:])
	//this.packetId = src[total : total+2]
	sp.packetId = src[total : total+2]
	total += 2

	l := int(sp.remLen) - (total - hn)
	sp.returnCodes = src[total : total+l]
	total += len(sp.returnCodes)

	for i, code := range sp.returnCodes {
		if code != 0x00 && code != 0x01 && code != 0x02 && code != 0x80 {
			return total, fmt.Errorf("suback/Decode: Invalid return code %d for topic %d", code, i)
		}
	}

	return total, nil
}

func (sp *SubackPacket) Encode(dst []byte) (int, error) {
	for i, code := range sp.returnCodes {
		if code != 0x00 && code != 0x01 && code != 0x02 && code != 0x80 {
			return 0, fmt.Errorf("suback/Encode: Invalid return code %d for topic %d", code, i)
		}
	}

	hl := sp.header.msglen()
	ml := sp.msglen()

	if len(dst) < hl+ml {
		return 0, fmt.Errorf("suback/Encode: Insufficient buffer size. Expecting %d, got %d.", hl+ml, len(dst))
	}

	if err := sp.SetRemainingLength(int32(ml)); err != nil {
		return 0, err
	}

	total := 0

	n, err := sp.header.encode(dst[total:])
	total += n
	if err != nil {
		return total, err
	}

	if copy(dst[total:total+2], sp.packetId) != 2 {
		dst[total], dst[total+1] = 0, 0
	}
	total += 2

	copy(dst[total:], sp.returnCodes)
	total += len(sp.returnCodes)

	return total, nil
}

func (sp *SubackPacket) msglen() int {
	return 2 + len(sp.returnCodes)
}
