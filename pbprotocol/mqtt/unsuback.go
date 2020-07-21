package mqtt

// The UNSUBACK Packet is sent by the Server to the Client to confirm receipt of an
// UNSUBSCRIBE Packet.
type UnsubackPacket struct {
	PubackPacket
}

// NewUnsubackPacket creates a new UNSUBACK packet.
func NewUnsubackPacket() *UnsubackPacket {
	msg := &UnsubackPacket{}
	msg.SetType(UNSUBACK)

	return msg
}
