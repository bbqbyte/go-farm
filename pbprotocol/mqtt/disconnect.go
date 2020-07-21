package mqtt

// The DISCONNECT Packet is the final Control Packet sent from the Client to the Server.
// It indicates that the Client is disconnecting cleanly.
type DisconnectPacket struct {
	header
}

// NewDisconnectPacket creates a new DISCONNECT packet.
func NewDisconnectPacket() *DisconnectPacket {
	msg := &DisconnectPacket{}
	msg.SetType(DISCONNECT)

	return msg
}

func (dp *DisconnectPacket) Decode(src []byte) (int, error) {
	return dp.header.decode(src)
}

func (dp *DisconnectPacket) Encode(dst []byte) (int, error) {
	return dp.header.encode(dst)
}
