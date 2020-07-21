package mqtt

// The PUBCOMP Packet is the response to a PUBREL Packet. It is the fourth and
// final packet of the QoS 2 protocol exchange.
type PubcompPacket struct {
	PubackPacket
}

// NewPubcompPacket creates a new PUBCOMP message.
func NewPubcompPacket() *PubcompPacket {
	msg := &PubcompPacket{}
	msg.SetType(PUBCOMP)

	return msg
}
