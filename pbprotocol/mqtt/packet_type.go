package mqtt

import "fmt"

// mqtt control packet type (bits 7-4).
type PacketType byte

const (
	// 0: Reserved. Forbidden.
	RESERVED PacketType = iota

	// 1: Client to Server. Client request to connect to Server.
	CONNECT

	// 2: Server to Client. Connect acknowledgment.
	CONNACK

	// 3: Client to Server or Server to Client. Publish message.
	PUBLISH

	// 4: Client to Server or Server to Client. Publish acknowledgment for.
	// Qos 1 messages
	PUBACK

	// 5: Client to Server or Server to Client. Publish received for Qos 2 messages.
	// Assured delivery part 1.
	PUBREC

	// 6: Client to Server or Server to Client. Publish release for Qos 2 messages.
	// Assured delivery part 2.
	PUBREL

	// 7: Client to Server or Server to Client. Publish complete for Qos 2 messages.
	// Assured delivery part 3.
	PUBCOMP

	// 8: Client to Server. Client subscribe request.
	SUBSCRIBE

	// 9: Server to Client. Subscribe acknowledgment.
	SUBACK

	// 10: Client to Server. Unsubscribe request.
	UNSUBSCRIBE

	// 11: Server to Client. Unsubscribe acknowledgment.
	UNSUBACK

	// 12: Client to Server. PING request.
	PINGREQ

	// 13: Server to Client. PING response.
	PINGRESP

	// 14: Client to Server. Client is disconnecting.
	DISCONNECT

	// 15: Reserved.
	RESERVED2
)

func (pkt PacketType) String() string {
	return pkt.Name()
}

// Name returns the name of the packet control type.
func (pt PacketType) Name() string {
	switch pt {
	case RESERVED:
		return "RESERVED"
	case CONNECT:
		return "CONNECT"
	case CONNACK:
		return "CONNACK"
	case PUBLISH:
		return "PUBLISH"
	case PUBACK:
		return "PUBACK"
	case PUBREC:
		return "PUBREC"
	case PUBREL:
		return "PUBREL"
	case PUBCOMP:
		return "PUBCOMP"
	case SUBSCRIBE:
		return "SUBSCRIBE"
	case SUBACK:
		return "SUBACK"
	case UNSUBSCRIBE:
		return "UNSUBSCRIBE"
	case UNSUBACK:
		return "UNSUBACK"
	case PINGREQ:
		return "PINGREQ"
	case PINGRESP:
		return "PINGRESP"
	case DISCONNECT:
		return "DISCONNECT"
	case RESERVED2:
		return "RESERVED2"
	}

	return "UNKNOWN"
}

// Desc returns the description of the packet control type.
func (pkt PacketType) Desc() string {
	switch pkt {
	case RESERVED:
		return "Reserved"
	case CONNECT:
		return "Client request to connect to Server"
	case CONNACK:
		return "Connect acknowledgement"
	case PUBLISH:
		return "Publish message"
	case PUBACK:
		return "Publish acknowledgement"
	case PUBREC:
		return "Publish received (assured delivery part 1)"
	case PUBREL:
		return "Publish release (assured delivery part 2)"
	case PUBCOMP:
		return "Publish complete (assured delivery part 3)"
	case SUBSCRIBE:
		return "Client subscribe request"
	case SUBACK:
		return "Subscribe acknowledgement"
	case UNSUBSCRIBE:
		return "Unsubscribe request"
	case UNSUBACK:
		return "Unsubscribe acknowledgement"
	case PINGREQ:
		return "PING request"
	case PINGRESP:
		return "PING response"
	case DISCONNECT:
		return "Client is disconnecting"
	case RESERVED2:
		return "Reserved2"
	}

	return "UNKNOWN"
}

// DefaultFlags returns the default flag values for the packet control type
func (pkt PacketType) DefaultFlags() byte {
	switch pkt {
	case PUBREL:
		return 2
	case SUBSCRIBE:
		return 2
	case UNSUBSCRIBE:
		return 2
	default:
		return 0
	}
}

// New creates a new packet based on the packet control type. It is a shortcut to call
// one of the New*Packet functions. If an error is returned then the packet control type
// is invalid.
func (pkt PacketType) New() (Packet, error) {
	switch pkt {
	case CONNECT:
		return NewConnectPacket(), nil
	case CONNACK:
		return NewConnackPacket(), nil
	case PUBLISH:
		return NewPublishPacket(), nil
	case PUBACK:
		return NewPubackPacket(), nil
	case PUBREC:
		return NewPubrecPacket(), nil
	case PUBREL:
		return NewPubrelPacket(), nil
	case PUBCOMP:
		return NewPubcompPacket(), nil
	case SUBSCRIBE:
		return NewSubscribePacket(), nil
	case SUBACK:
		return NewSubackPacket(), nil
	case UNSUBSCRIBE:
		return NewUnsubscribePacket(), nil
	case UNSUBACK:
		return NewUnsubackPacket(), nil
	case PINGREQ:
		return NewPingreqPacket(), nil
	case PINGRESP:
		return NewPingrespPacket(), nil
	case DISCONNECT:
		return NewDisconnectPacket(), nil
	default:
		return nil, fmt.Errorf("msgtype/NewMessage: Invalid packet type %d", pkt)
	}
}

// Valid returns a boolean indicating whether the packet control type is valid or not.
func (pkt PacketType) Valid() bool {
	return pkt > RESERVED && pkt < RESERVED2
}
