package mqtt

// A PINGRESP Packet is sent by the Server to the Client in response to a PINGREQ
// Packet. It indicates that the Server is alive.

// This Packet is used in Keep Alive processing, see Section 3.1.2.10 for more details.
type PingrespPacket struct {
	DisconnectPacket
}

// NewPingrespPacket creates a new PINGRESP packet.
func NewPingrespPacket() *PingrespPacket {
	msg := &PingrespPacket{}
	msg.SetType(PINGRESP)

	return msg
}
