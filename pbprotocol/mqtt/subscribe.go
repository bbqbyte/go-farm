package mqtt

import (
	"bytes"
	"fmt"
)

// The SUBSCRIBE Packet is sent from the Client to the Server to create one or more
// Subscriptions. Each Subscription registers a Client’s interest in one or more
// Topics. The Server sends PUBLISH Packets to the Client in order to forward
// Application Messages that were published to Topics that match these Subscriptions.
// The SUBSCRIBE Packet also specifies (for each Subscription) the maximum QoS with
// which the Server can send Application Messages to the Client.
type SubscribePacket struct {
	header

	topics [][]byte
	qos    []byte
}

// NewSubscribePacket creates a new SUBSCRIBE message.
func NewSubscribePacket() *SubscribePacket {
	msg := &SubscribePacket{}
	msg.SetType(SUBSCRIBE)

	return msg
}

func (sp SubscribePacket) String() string {
	msgstr := fmt.Sprintf("%s, Packet ID=%d", sp.header, sp.PacketID())

	for i, t := range sp.topics {
		msgstr = fmt.Sprintf("%s, Topic[%d]=%q/%d", msgstr, i, string(t), sp.qos[i])
	}

	return msgstr
}

// Topics returns a list of topics sent by the Client.
func (sp *SubscribePacket) Topics() [][]byte {
	return sp.topics
}

// AddTopic adds a single topic to the message, along with the corresponding QoS.
// An error is returned if QoS is invalid.
func (sp *SubscribePacket) AddTopic(topic []byte, qos byte) error {
	if !ValidQos(qos) {
		return fmt.Errorf("Invalid QoS %d", qos)
	}

	for i, t := range sp.topics {
		if bytes.Equal(t, topic) {
			sp.qos[i] = qos
			return nil
		}
	}

	sp.topics = append(sp.topics, topic)
	sp.qos = append(sp.qos, qos)

	return nil
}

// RemoveTopic removes a single topic from the list of existing ones in the message.
// If topic does not exist it just does nothing.
func (sp *SubscribePacket) RemoveTopic(topic []byte) {
	for i, t := range sp.topics {
		if bytes.Equal(t, topic) {
			sp.topics = append(sp.topics[:i], sp.topics[i+1:]...)
			sp.qos = append(sp.qos[:i], sp.qos[i+1:]...)
			break
		}
	}
}

// TopicExists checks to see if a topic exists in the list.
func (sp *SubscribePacket) TopicExists(topic []byte) bool {
	for _, t := range sp.topics {
		if bytes.Equal(t, topic) {
			return true
		}
	}

	return false
}

// TopicQos returns the QoS level of a topic. If topic does not exist, QosFailure
// is returned.
func (sp *SubscribePacket) TopicQos(topic []byte) byte {
	for i, t := range sp.topics {
		if bytes.Equal(t, topic) {
			return sp.qos[i]
		}
	}

	return QosFailure
}

// Qos returns the list of QoS current in the message.
func (sp *SubscribePacket) Qos() []byte {
	return sp.qos
}

func (sp *SubscribePacket) Len() int {
	ml := sp.msglen()

	if err := sp.SetRemainingLength(int32(ml)); err != nil {
		return 0
	}

	return sp.header.msglen() + ml
}

func (sp *SubscribePacket) Decode(src []byte) (int, error) {
	total := 0

	hn, err := sp.header.decode(src[total:])
	total += hn
	if err != nil {
		return total, err
	}

	//this.packetId = binary.BigEndian.Uint16(src[total:])
	//this.packetId = src[total : total+2]
	sp.packetId = src[total : total+2]
	total += 2

	rl := int(sp.remLen) - (total - hn)
	for rl > 0 {
		t, n, err := readLPBytes(src[total:])
		total += n
		if err != nil {
			return total, err
		}

		sp.topics = append(sp.topics, t)

		sp.qos = append(sp.qos, src[total])
		total++

		rl = rl - n - 1
	}

	if len(sp.topics) == 0 {
		return 0, fmt.Errorf("subscribe/Decode: Empty topic list")
	}

	return total, nil
}

func (sp *SubscribePacket) Encode(dst []byte) (int, error) {
	hl := sp.header.msglen()
	ml := sp.msglen()

	if len(dst) < hl+ml {
		return 0, fmt.Errorf("subscribe/Encode: Insufficient buffer size. Expecting %d, got %d.", hl+ml, len(dst))
	}

	if err := sp.SetRemainingLength(int32(ml)); err != nil {
		return 0, err
	}

	total := 0

	n, err := sp.header.encode(dst[total:])
	total += n
	if err != nil {
		return total, err
	}

	if sp.PacketID() == 0 {
		return 0, fmt.Errorf("subscribe/Encode: invalid packetid %d", sp.PacketID())
	}

	n = copy(dst[total:], sp.packetId)
	total += n

	for i, t := range sp.topics {
		n, err := writeLPBytes(dst[total:], t)
		total += n
		if err != nil {
			return total, err
		}

		dst[total] = sp.qos[i]
		total++
	}

	return total, nil
}

func (sp *SubscribePacket) msglen() int {
	// packet ID
	total := 2

	for _, t := range sp.topics {
		total += 2 + len(t) + 1
	}

	return total
}
