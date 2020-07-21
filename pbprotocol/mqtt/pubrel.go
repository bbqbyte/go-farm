package mqtt

// A PUBREL Packet is the response to a PUBREC Packet. It is the third packet of the
// QoS 2 protocol exchange.
type PubrelPacket struct {
	PubackPacket
}

// NewPubrelPacket creates a new PUBREL packet.
func NewPubrelPacket() *PubrelPacket {
	msg := &PubrelPacket{}
	msg.SetType(PUBREL)

	return msg
}
