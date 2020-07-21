package mqtt

// A PUBREC Packet is the response to a PUBLISH Packet with QoS 2.
// It is the second packet of the QoS 2 protocol exchange.
type PubrecPacket struct {
	PubackPacket
}

// NewPubrecPacket creates a new PUBREC message.
func NewPubrecPacket() *PubrecPacket {
	msg := &PubrecPacket{}
	msg.SetType(PUBREC)

	return msg
}
