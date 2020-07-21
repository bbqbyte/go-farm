package mqtt

// The PINGREQ Packet is sent from a Client to the Server. It can be used to:
// 1. Indicate to the Server that the Client is alive in the absence of any other
//    Control Packets being sent from the Client to the Server.
// 2. Request that the Server responds to confirm that it is alive.
// 3. Exercise the network to indicate that the Network Connection is active.

// This Packet is used in Keep Alive processing, see Section 3.1.2.10 for more details.
type PingreqPacket struct {
	DisconnectPacket
}

// NewPingreqPacket creates a new PINGREQ packet.
func NewPingreqPacket() *PingreqPacket {
	msg := &PingreqPacket{}
	msg.SetType(PINGREQ)

	return msg
}
