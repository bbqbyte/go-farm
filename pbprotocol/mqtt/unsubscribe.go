package mqtt

import (
	"bytes"
	"fmt"
)

// An UNSUBSCRIBE Packet is sent by the Client to the Server, to unsubscribe from topics.
type UnsubscribePacket struct {
	header

	topics [][]byte
}

// NewUnsubscribePacket creates a new UNSUBSCRIBE packet.
func NewUnsubscribePacket() *UnsubscribePacket {
	msg := &UnsubscribePacket{}
	msg.SetType(UNSUBSCRIBE)

	return msg
}

func (up UnsubscribePacket) String() string {
	msgstr := fmt.Sprintf("%s", up.header)

	for i, t := range up.topics {
		msgstr = fmt.Sprintf("%s, Topic%d=%s", msgstr, i, string(t))
	}

	return msgstr
}

// Topics returns a list of topics sent by the Client.
func (up *UnsubscribePacket) Topics() [][]byte {
	return up.topics
}

// AddTopic adds a single topic to the message.
func (up *UnsubscribePacket) AddTopic(topic []byte) {
	if up.TopicExists(topic) {
		return
	}

	up.topics = append(up.topics, topic)
}

// RemoveTopic removes a single topic from the list of existing ones in the message.
// If topic does not exist it just does nothing.
func (up *UnsubscribePacket) RemoveTopic(topic []byte) {
	for i, t := range up.topics {
		if bytes.Equal(t, topic) {
			up.topics = append(up.topics[:i], up.topics[i+1:]...)
			break
		}
	}
}

// TopicExists checks to see if a topic exists in the list.
func (up *UnsubscribePacket) TopicExists(topic []byte) bool {
	for _, t := range up.topics {
		if bytes.Equal(t, topic) {
			return true
		}
	}

	return false
}

func (up *UnsubscribePacket) Len() int {
	ml := up.msglen()

	if err := up.SetRemainingLength(int32(ml)); err != nil {
		return 0
	}

	return up.header.msglen() + ml
}

// Decode reads from the io.Reader parameter until a full message is decoded, or
// when io.Reader returns EOF or error. The first return value is the number of
// bytes read from io.Reader. The second is error if Decode encounters any problems.
func (up *UnsubscribePacket) Decode(src []byte) (int, error) {
	total := 0

	hn, err := up.header.decode(src[total:])
	total += hn
	if err != nil {
		return total, err
	}

	//this.packetId = binary.BigEndian.Uint16(src[total:])
	//this.packetId = src[total : total+2]
	up.packetId = src[total : total+2]
	total += 2

	rl := int(up.remLen) - (total - hn)
	for rl > 0 {
		t, n, err := readLPBytes(src[total:])
		total += n
		if err != nil {
			return total, err
		}

		up.topics = append(up.topics, t)
		rl = rl - n - 1
	}

	if len(up.topics) == 0 {
		return 0, fmt.Errorf("unsubscribe/Decode: Empty topic list")
	}

	return total, nil
}

// Encode returns an io.Reader in which the encoded bytes can be read. The second
// return value is the number of bytes encoded, so the caller knows how many bytes
// there will be. If Encode returns an error, then the first two return values
// should be considered invalid.
// Any changes to the message after Encode() is called will invalidate the io.Reader.
func (up *UnsubscribePacket) Encode(dst []byte) (int, error) {
	hl := up.header.msglen()
	ml := up.msglen()

	if len(dst) < hl+ml {
		return 0, fmt.Errorf("unsubscribe/Encode: Insufficient buffer size. Expecting %d, got %d.", hl+ml, len(dst))
	}

	if err := up.SetRemainingLength(int32(ml)); err != nil {
		return 0, err
	}

	total := 0

	n, err := up.header.encode(dst[total:])
	total += n
	if err != nil {
		return total, err
	}

	if up.PacketID() == 0 {
		return 0, fmt.Errorf("subscribe/Encode: invalid packetid %d", up.PacketID())
	}

	n = copy(dst[total:], up.packetId)
	total += n

	for _, t := range up.topics {
		n, err := writeLPBytes(dst[total:], t)
		total += n
		if err != nil {
			return total, err
		}
	}

	return total, nil
}

func (up *UnsubscribePacket) msglen() int {
	// packet ID
	total := 2

	for _, t := range up.topics {
		total += 2 + len(t)
	}

	return total
}
