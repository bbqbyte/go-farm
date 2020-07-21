package mqtt

import (
	"encoding/binary"
	"fmt"
)

type header struct {
	// mqtt control packet type (bits 7-4).
	// flags specific to each mqtt control packet type (bits 3-0).
	mtypeflags []byte

	// remaining length = variable header len + payload len
	// 1-4 byte
	remLen int32

	// some message need packet ID, 2 byte - uint16, big-endian
	packetId []byte
}

// String returns a string representation of the message.
func (h header) String() string {
	return fmt.Sprintf("Type=%q, Flags=%08b, Remaining Length=%d", h.Type().Name(), h.Flags(), h.remLen)
}

// Name returns a string representation of the packet type. Examples include
// "PUBLISH", "SUBSCRIBE", and others. This is statically defined for each of
// the packet types and cannot be changed.
func (h *header) Name() string {
	return h.Type().Name()
}

// Desc returns a string description of the packet type. For example, a
// CONNECT message would return "Client request to connect to Server." These
// descriptions are statically defined (copied from the MQTT spec) and cannot
// be changed.
func (h *header) Desc() string {
	return h.Type().Desc()
}

// Type returns the PacketType of the packet. The retured value should be one
// of the constants defined for PacketType.
func (h *header) Type() PacketType {
	return PacketType(h.mtypeflags[0] >> 4)
}

// SetType sets the packet type of this packet. It also correctly sets the
// default flags for the packet type. It returns an error if the type is invalid.
func (h *header) SetType(ptype PacketType) error {
	if !ptype.Valid() {
		return fmt.Errorf("header/SetType: Invalid control packet type %d", ptype)
	}

	h.mtypeflags[0] = byte(ptype)<<4 | (ptype.DefaultFlags() & 0xf)

	return nil
}

// Flags returns the fixed header flags for this packet.
func (h *header) Flags() byte {
	return h.mtypeflags[0] & 0x0f
}

// RemainingLength returns the length of the non-fixed-header part of the packet.
func (h *header) RemainingLength() int32 {
	return h.remLen
}

// SetRemainingLength sets the length of the non-fixed-header part of the message.
// It returns error if the length is greater than 268435455, which is the max
// message length as defined by the MQTT spec.
func (h *header) SetRemainingLength(remlen int32) error {
	if remlen > maxRemainingLength || remlen < 0 {
		return fmt.Errorf("header/SetLength: Remaining length (%d) out of bound (max %d, min 0)", remlen, maxRemainingLength)
	}

	h.remLen = remlen

	return nil
}

func (h *header) Len() int {
	return h.msglen()
}

// PacketId returns the ID of the packet.
func (h *header) PacketID() uint16 {
	if len(h.packetId) == 2 {
		return binary.BigEndian.Uint16(h.packetId)
	}

	return 0
}

// SetPacketId sets the ID of the packet.
func (h *header) SetPacketId(v uint16) {
	// If setting to 0, nothing to do, move on
	if v == 0 {
		return
	}

	// If packetId buffer is not 2 bytes (uint16), then we allocate a new one.
	// Then we encode the packet ID into the buffer.
	if len(h.packetId) != 2 {
		h.packetId = make([]byte, 2)
	}

	// Notice we don't set the packet to be dirty when we are not allocating a new
	// buffer. In this case, it means the buffer is probably a sub-slice of another
	// slice. If that's the case, then during encoding we would have copied the whole
	// backing buffer anyway.
	binary.BigEndian.PutUint16(h.packetId, v)
}

func (h *header) encode(dst []byte) (int, error) {
	ml := h.msglen()

	if len(dst) < ml {
		return 0, fmt.Errorf("header/Encode: Insufficient buffer size. Expecting %d, got %d.", ml, len(dst))
	}

	total := 0

	if h.remLen > maxRemainingLength || h.remLen < 0 {
		return total, fmt.Errorf("header/Encode: Remaining length (%d) out of bound (max %d, min 0)", h.remLen, maxRemainingLength)
	}

	if !h.Type().Valid() {
		return total, fmt.Errorf("header/Encode: Invalid message type %d", h.Type())
	}

	// first byte. packet control type
	dst[total] = h.mtypeflags[0]
	total += 1

	// write the remaining length
	n := binary.PutUvarint(dst[total:], uint64(h.remLen))
	total += n

	return total, nil
}

// Decode reads from the io.Reader parameter until a full packet is decoded, or
// when io.Reader returns EOF or error. The first return value is the number of
// bytes read from io.Reader. The second is error if Decode encounters any problems.
func (h *header) decode(src []byte) (int, error) {
	total := 0

	mtype := h.Type()

	h.mtypeflags = src[total : total+1]

	if !h.Type().Valid() {
		return total, fmt.Errorf("header/Decode: Invalid packet type %d.", mtype)
	}

	if mtype != h.Type() {
		return total, fmt.Errorf("header/Decode: Invalid packet type %d. Expecting %d.", h.Type(), mtype)
	}

	//this.flags = src[total] & 0x0f
	if h.Type() != PUBLISH && h.Flags() != h.Type().DefaultFlags() {
		return total, fmt.Errorf("header/Decode: Invalid packet (%d) flags. Expecting %d, got %d", h.Type(), h.Type().DefaultFlags(), h.Flags())
	}

	if h.Type() == PUBLISH && !ValidQos((h.Flags()>>1)&0x3) {
		return total, fmt.Errorf("header/Decode: Invalid QoS (%d) for PUBLISH packet.", (h.Flags()>>1)&0x3)
	}

	total++

	rl, m := binary.Uvarint(src[total:])
	total += m
	h.remLen = int32(rl)

	if h.remLen > maxRemainingLength || rl < 0 {
		return total, fmt.Errorf("header/Decode: Remaining length (%d) out of bound (max %d, min 0)", h.remLen, maxRemainingLength)
	}

	if int(h.remLen) > len(src[total:]) {
		return total, fmt.Errorf("header/Decode: Remaining length (%d) is greater than remaining buffer (%d)", h.remLen, len(src[total:]))
	}

	return total, nil
}

func (h *header) msglen() int {
	total := 1

	if h.remLen <= 127 {
		total += 1
	} else if h.remLen <= 16383 {
		total += 2
	} else if h.remLen <= 2097151 {
		total += 3
	} else {
		total += 4
	}

	return total

}
