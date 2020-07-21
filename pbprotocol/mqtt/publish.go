package mqtt

import (
	"fmt"
)

// A PUBLISH Control Packet is sent from a Client to a Server or from Server to a Client
// to transport an Application Message.
type PublishPacket struct {
	header

	topic   []byte
	payload []byte
}

// NewPublishMessage creates a new PUBLISH packet.
func NewPublishPacket() *PublishPacket {
	msg := &PublishPacket{}
	msg.SetType(PUBLISH)

	return msg
}

func (pp PublishPacket) String() string {
	return fmt.Sprintf("%s, Topic=%q, Packet ID=%d, QoS=%d, Retained=%t, Dup=%t, Payload=%v",
		pp.header, pp.topic, pp.packetId, pp.QoS(), pp.Retain(), pp.Dup(), pp.payload)
}

// Dup returns the value specifying the duplicate delivery of a PUBLISH Control Packet.
// If the DUP flag is set to 0, it indicates that this is the first occasion that the
// Client or Server has attempted to send this MQTT PUBLISH Packet. If the DUP flag is
// set to 1, it indicates that this might be re-delivery of an earlier attempt to send
// the Packet.
func (pp *PublishPacket) Dup() bool {
	return ((pp.Flags() >> 3) & 0x1) == 1
}

// SetDup sets the value specifying the duplicate delivery of a PUBLISH Control Packet.
func (pp *PublishPacket) SetDup(v bool) {
	if v {
		pp.mtypeflags[0] |= 0x8 // 00001000
	} else {
		pp.mtypeflags[0] &= 247 // 11110111
	}
}

// Retain returns the value of the RETAIN flag. This flag is only used on the PUBLISH
// Packet. If the RETAIN flag is set to 1, in a PUBLISH Packet sent by a Client to a
// Server, the Server MUST store the Application Message and its QoS, so that it can be
// delivered to future subscribers whose subscriptions match its topic name.
func (pp *PublishPacket) Retain() bool {
	return (pp.Flags() & 0x1) == 1
}

// SetRetain sets the value of the RETAIN flag.
func (pp *PublishPacket) SetRetain(v bool) {
	if v {
		pp.mtypeflags[0] |= 0x1 // 00000001
	} else {
		pp.mtypeflags[0] &= 254 // 11111110
	}
}

// QoS returns the field that indicates the level of assurance for delivery of an
// Application Message. The values are QosAtMostOnce, QosAtLeastOnce and QosExactlyOnce.
func (pp *PublishPacket) QoS() byte {
	return (pp.Flags() >> 1) & 0x3
}

// SetQoS sets the field that indicates the level of assurance for delivery of an
// Application Message. The values are QosAtMostOnce, QosAtLeastOnce and QosExactlyOnce.
// An error is returned if the value is not one of these.
func (pp *PublishPacket) SetQoS(v byte) error {
	if v != 0x0 && v != 0x1 && v != 0x2 {
		return fmt.Errorf("publish/SetQoS: Invalid QoS %d.", v)
	}

	pp.mtypeflags[0] = (pp.mtypeflags[0] & 249) | (v << 1) // 249 = 11111001

	return nil
}

// Topic returns the the topic name that identifies the information channel to which
// payload data is published.
func (pp *PublishPacket) Topic() []byte {
	return pp.topic
}

// SetTopic sets the the topic name that identifies the information channel to which
// payload data is published. An error is returned if ValidTopic() is falbase.
func (pp *PublishPacket) SetTopic(v []byte) error {
	if !ValidTopic(v) {
		return fmt.Errorf("publish/SetTopic: Invalid topic name (%s). Must not be empty or contain wildcard characters", string(v))
	}

	pp.topic = v

	return nil
}

// Payload returns the application message that's part of the PUBLISH message.
// Payload length = the remaining length - variable header length. maybe zero.
func (pp *PublishPacket) Payload() []byte {
	return pp.payload
}

// SetPayload sets the application message that's part of the PUBLISH message.
func (pp *PublishPacket) SetPayload(v []byte) {
	pp.payload = v
}

func (pp *PublishPacket) Len() int {
	ml := pp.msglen()

	if err := pp.SetRemainingLength(int32(ml)); err != nil {
		return 0
	}

	return pp.header.msglen() + ml
}

func (pp *PublishPacket) Decode(src []byte) (int, error) {
	total := 0

	// decode fixed header
	hn, err := pp.header.decode(src[total:])
	total += hn
	if err != nil {
		return total, err
	}

	n := 0

	pp.topic, n, err = readLPBytes(src[total:])
	total += n
	if err != nil {
		return total, err
	}

	if !ValidTopic(pp.topic) {
		return total, fmt.Errorf("publish/Decode: Invalid topic name (%s). Must not be empty or contain wildcard characters", string(pp.topic))
	}

	// The packet identifier field is only present in the PUBLISH packets where the
	// QoS level is 1 or 2
	if pp.QoS() != 0 {
		// pp.packetId = binary.BigEndian.Uint16(src[total:])
		// pp.packetId = src[total : total+2]
		pp.packetId = src[total : total+2]
		total += 2
	}

	l := int(pp.remLen) - (total - hn)
	pp.payload = src[total : total+l]
	total += len(pp.payload)

	return total, nil
}

func (pp *PublishPacket) Encode(dst []byte) (int, error) {
	if len(pp.topic) == 0 {
		return 0, fmt.Errorf("publish/Encode: Topic name is empty.")
	}

	if len(pp.payload) == 0 {
		return 0, fmt.Errorf("publish/Encode: Payload is empty.")
	}

	ml := pp.msglen()

	if err := pp.SetRemainingLength(int32(ml)); err != nil {
		return 0, err
	}

	hl := pp.header.msglen()

	if len(dst) < hl+ml {
		return 0, fmt.Errorf("publish/Encode: Insufficient buffer size. Expecting %d, got %d.", hl+ml, len(dst))
	}

	total := 0

	n, err := pp.header.encode(dst[total:])
	total += n
	if err != nil {
		return total, err
	}

	n, err = writeLPBytes(dst[total:], pp.topic)
	total += n
	if err != nil {
		return total, err
	}

	// The packet identifier field is only present in the PUBLISH packets where the QoS level is 1 or 2
	if pp.QoS() != 0 {
		if pp.PacketID() == 0 {
			//pp.SetPacketId(uint16(atomic.AddUint64(&gPacketId, 1) & 0xffff))
			return 0, fmt.Errorf("publish/Encode: invalid packetid %d when qos == 0", pp.PacketID())
		}

		n = copy(dst[total:], pp.packetId)

		total += n
	}

	copy(dst[total:], pp.payload)
	total += len(pp.payload)

	return total, nil
}

func (pp *PublishPacket) msglen() int {
	total := 2 + len(pp.topic) + len(pp.payload)
	if pp.QoS() != 0 {
		total += 2
	}

	return total
}
