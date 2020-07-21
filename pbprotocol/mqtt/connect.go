package mqtt

import (
	"encoding/binary"
	"fmt"
	"regexp"
)

var clientIdRegexp *regexp.Regexp

func init() {
	// Added space for Paho compliance test
	// Added underscore (_) for MQTT C client test
	clientIdRegexp = regexp.MustCompile("^[0-9a-zA-Z _]*$")
}

// After a Network Connection is established by a Client to a Server, the first Packet
// sent from the Client to the Server MUST be a CONNECT Packet [MQTT-3.1.0-1].
//
// A Client can only send the CONNECT Packet once over a Network Connection. The Server
// MUST process a second CONNECT Packet sent from a Client as a protocol violation and
// disconnect the Client [MQTT-3.1.0-2].  See section 4.8 for information about
// handling errors.
type ConnectPacket struct {
	// fixed header
	header

	// 7: username flag
	// 6: password flag
	// 5: will retain
	// 4-3: will QoS
	// 2: will flag
	// 1: clean session
	// 0: reserved
	connectFlags byte

	protolevel byte

	keepAlive uint16 // a time interval measured in seconds

	protoName, // protocal name
	clientId, // client identifier
	willTopic,
	willMessage,
	username,
	password []byte
}

// NewConnectPacket creates a new CONNECT packet.
func NewConnectPacket() *ConnectPacket {
	cp := &ConnectPacket{}
	cp.SetType(CONNECT)

	return cp
}

// String returns a string representation of the CONNECT message
func (cp ConnectPacket) String() string {
	return fmt.Sprintf("%s, Connect Flags=%08b, Version=%d, KeepAlive=%d, Client ID=%q, Will Topic=%q, Will Message=%q, Username=%q, Password=%q",
		cp.header,
		cp.connectFlags,
		cp.ProtoLevel(),
		cp.KeepAlive(),
		cp.ClientId(),
		cp.WillTopic(),
		cp.WillMessage(),
		cp.Username(),
		cp.Password(),
	)
}

// ProtoLevel returns the the 8 bit unsigned value that represents the revision level
// of the protocol used by the Client. The value of the Protocol Level field for
// the version 3.1.1 of the protocol is 4 (0x04).
func (cp *ConnectPacket) ProtoLevel() byte {
	return cp.protolevel
}

// SetProtoLevel sets the version value of the CONNECT message
func (cp *ConnectPacket) SetProtoLevel(v byte) error {
	if _, ok := SupportedVersions[v]; !ok {
		return fmt.Errorf("connect/SetProtoLevel: Invalid protocal level number %d", v)
	}

	cp.protolevel = v

	return nil
}

// CleanSession returns the bit that specifies the handling of the Session state.
// The Client and Server can store Session state to enable reliable messaging to
// continue across a sequence of Network Connections. This bit is used to control
// the lifetime of the Session state.
func (cp *ConnectPacket) CleanSession() bool {
	return ((cp.connectFlags >> 1) & 0x1) == 1
}

// SetCleanSession sets the bit that specifies the handling of the Session state.
func (cp *ConnectPacket) SetCleanSession(v bool) {
	if v {
		cp.connectFlags |= 0x2 // 00000010
	} else {
		cp.connectFlags &= 253 // 11111101
	}
}

// WillFlag returns the bit that specifies whether a Will Message should be stored
// on the server. If the Will Flag is set to 1 this indicates that, if the Connect
// request is accepted, a Will Message MUST be stored on the Server and associated
// with the Network Connection.
func (cp *ConnectPacket) WillFlag() bool {
	return ((cp.connectFlags >> 2) & 0x1) == 1
}

// SetWillFlag sets the bit that specifies whether a Will Message should be stored
// on the server.
func (cp *ConnectPacket) SetWillFlag(v bool) {
	if v {
		cp.connectFlags |= 0x4 // 00000100
	} else {
		cp.connectFlags &= 251 // 11111011
	}
}

// WillQos returns the two bits that specify the QoS level to be used when publishing
// the Will Message.
func (cp *ConnectPacket) WillQos() byte {
	return (cp.connectFlags >> 3) & 0x3
}

// SetWillQos sets the two bits that specify the QoS level to be used when publishing
// the Will Message.
func (cp *ConnectPacket) SetWillQos(qos byte) error {
	if qos != QosAtMostOnce && qos != QosAtLeastOnce && qos != QosExactlyOnce {
		return fmt.Errorf("connect/SetWillQos: Invalid QoS level %d", qos)
	}

	cp.connectFlags = (cp.connectFlags & 231) | (qos << 3) // 231 = 11100111

	return nil
}

// WillRetain returns the bit specifies if the Will Message is to be Retained when it
// is published.
func (cp *ConnectPacket) WillRetain() bool {
	return ((cp.connectFlags >> 5) & 0x1) == 1
}

// SetWillRetain sets the bit specifies if the Will Message is to be Retained when it
// is published.
func (cp *ConnectPacket) SetWillRetain(v bool) {
	if v {
		cp.connectFlags |= 32 // 00100000
	} else {
		cp.connectFlags &= 223 // 11011111
	}
}

// UsernameFlag returns the bit that specifies whether a user name is present in the
// payload.
func (cp *ConnectPacket) UsernameFlag() bool {
	return ((cp.connectFlags >> 7) & 0x1) == 1
}

// SetUsernameFlag sets the bit that specifies whether a user name is present in the
// payload.
func (cp *ConnectPacket) SetUsernameFlag(v bool) {
	if v {
		cp.connectFlags |= 128 // 10000000
	} else {
		cp.connectFlags &= 127 // 01111111
	}
}

// PasswordFlag returns the bit that specifies whether a password is present in the
// payload.
func (cp *ConnectPacket) PasswordFlag() bool {
	return ((cp.connectFlags >> 6) & 0x1) == 1
}

// SetPasswordFlag sets the bit that specifies whether a password is present in the
// payload.
func (cp *ConnectPacket) SetPasswordFlag(v bool) {
	if v {
		cp.connectFlags |= 64 // 01000000
	} else {
		cp.connectFlags &= 191 // 10111111
	}
}

// KeepAlive returns a time interval measured in seconds. Expressed as a 16-bit word,
// it is the maximum time interval that is permitted to elapse between the point at
// which the Client finishes transmitting one Control Packet and the point it starts
// sending the next.
func (cp *ConnectPacket) KeepAlive() uint16 {
	return cp.keepAlive
}

// SetKeepAlive sets the time interval in which the server should keep the connection
// alive.
func (cp *ConnectPacket) SetKeepAlive(v uint16) {
	cp.keepAlive = v
}

// ClientId returns an ID that identifies the Client to the Server. Each Client
// connecting to the Server has a unique ClientId. The ClientId MUST be used by
// Clients and by Servers to identify state that they hold relating to this MQTT
// Session between the Client and the Server
func (cp *ConnectPacket) ClientId() []byte {
	return cp.clientId
}

// SetClientId sets an ID that identifies the Client to the Server.
func (cp *ConnectPacket) SetClientId(v []byte) error {
	if len(v) > 0 && !cp.validClientId(v) {
		return ErrIdentifierRejected
	}

	cp.clientId = v

	return nil
}

// WillTopic returns the topic in which the Will Message should be published to.
// If the Will Flag is set to 1, the Will Topic must be in the payload.
func (cp *ConnectPacket) WillTopic() []byte {
	return cp.willTopic
}

// SetWillTopic sets the topic in which the Will Message should be published to.
func (cp *ConnectPacket) SetWillTopic(v []byte) {
	cp.willTopic = v

	if len(v) > 0 {
		cp.SetWillFlag(true)
	} else if len(cp.willMessage) == 0 {
		cp.SetWillFlag(false)
	}
}

// WillMessage returns the Will Message that is to be published to the Will Topic.
func (cp *ConnectPacket) WillMessage() []byte {
	return cp.willMessage
}

// SetWillMessage sets the Will Message that is to be published to the Will Topic.
func (cp *ConnectPacket) SetWillMessage(v []byte) {
	cp.willMessage = v

	if len(v) > 0 {
		cp.SetWillFlag(true)
	} else if len(cp.willTopic) == 0 {
		cp.SetWillFlag(false)
	}
}

// Username returns the username from the payload. If the User Name Flag is set to 1,
// this must be in the payload. It can be used by the Server for authentication and
// authorization.
func (cp *ConnectPacket) Username() []byte {
	return cp.username
}

// SetUsername sets the username for authentication.
func (cp *ConnectPacket) SetUsername(v []byte) {
	cp.username = v

	if len(v) > 0 {
		cp.SetUsernameFlag(true)
	} else {
		cp.SetUsernameFlag(false)
	}
}

// Password returns the password from the payload. If the Password Flag is set to 1,
// this must be in the payload. It can be used by the Server for authentication and
// authorization.
func (cp *ConnectPacket) Password() []byte {
	return cp.password
}

// SetPassword sets the username for authentication.
func (cp *ConnectPacket) SetPassword(v []byte) {
	cp.password = v

	if len(v) > 0 {
		cp.SetPasswordFlag(true)
	} else {
		cp.SetPasswordFlag(false)
	}
}

func (cp *ConnectPacket) Len() int {
	ml := cp.msglen()

	if err := cp.SetRemainingLength(int32(ml)); err != nil {
		return 0
	}

	return cp.header.msglen() + ml
}

// For the CONNECT message, the error returned could be a ConnackReturnCode, so
// be sure to check that. Otherwise it's a generic error. If a generic error is
// returned, this Message should be considered invalid.
//
// Caller should call ValidConnackError(err) to see if the returned error is
// a Connack error. If so, caller should send the Client back the corresponding
// CONNACK message.
func (cp *ConnectPacket) Decode(src []byte) (int, error) {
	total := 0

	n, err := cp.header.decode(src[total:])
	if err != nil {
		return total + n, err
	}
	total += n

	if n, err = cp.decodeMessage(src[total:]); err != nil {
		return total + n, err
	}
	total += n

	return total, nil
}

func (cp *ConnectPacket) Encode(dst []byte) (int, error) {
	if cp.Type() != CONNECT {
		return 0, fmt.Errorf("connect/Encode: Invalid message type. Expecting %d, got %d", CONNECT, cp.Type())
	}

	_, ok := SupportedVersions[cp.protolevel]
	if !ok {
		return 0, ErrInvalidProtocolVersion
	}

	hl := cp.header.msglen()
	ml := cp.msglen()

	if len(dst) < hl+ml {
		return 0, fmt.Errorf("connect/Encode: Insufficient buffer size. Expecting %d, got %d.", hl+ml, len(dst))
	}

	if err := cp.SetRemainingLength(int32(ml)); err != nil {
		return 0, err
	}

	total := 0

	n, err := cp.header.encode(dst[total:])
	total += n
	if err != nil {
		return total, err
	}

	n, err = cp.encodeMessage(dst[total:])
	total += n
	if err != nil {
		return total, err
	}

	return total, nil
}

func (cp *ConnectPacket) encodeMessage(dst []byte) (int, error) {
	total := 0

	n, err := writeLPBytes(dst[total:], []byte(SupportedVersions[cp.protolevel]))
	total += n
	if err != nil {
		return total, err
	}

	dst[total] = cp.protolevel
	total += 1

	dst[total] = cp.connectFlags
	total += 1

	binary.BigEndian.PutUint16(dst[total:], cp.keepAlive)
	total += 2

	n, err = writeLPBytes(dst[total:], cp.clientId)
	total += n
	if err != nil {
		return total, err
	}

	if cp.WillFlag() {
		n, err = writeLPBytes(dst[total:], cp.willTopic)
		total += n
		if err != nil {
			return total, err
		}

		n, err = writeLPBytes(dst[total:], cp.willMessage)
		total += n
		if err != nil {
			return total, err
		}
	}

	// According to the 3.1 spec, it's possible that the usernameFlag is set,
	// but the username string is missing.
	if cp.UsernameFlag() && len(cp.username) > 0 {
		n, err = writeLPBytes(dst[total:], cp.username)
		total += n
		if err != nil {
			return total, err
		}
	}

	// According to the 3.1 spec, it's possible that the passwordFlag is set,
	// but the password string is missing.
	if cp.PasswordFlag() && len(cp.password) > 0 {
		n, err = writeLPBytes(dst[total:], cp.password)
		total += n
		if err != nil {
			return total, err
		}
	}

	return total, nil
}

func (cp *ConnectPacket) decodeMessage(src []byte) (int, error) {
	var err error
	n, total := 0, 0

	cp.protoName, n, err = readLPBytes(src[total:])
	total += n
	if err != nil {
		return total, err
	}

	cp.protolevel = src[total]
	total++

	if verstr, ok := SupportedVersions[cp.protolevel]; !ok {
		return total, ErrInvalidProtocolVersion
	} else if verstr != string(cp.protoName) {
		return total, ErrInvalidProtocolVersion
	}

	cp.connectFlags = src[total]
	total++

	if cp.connectFlags&0x1 != 0 {
		return total, fmt.Errorf("connect/decodeMessage: Connect Flags reserved bit 0 is not 0")
	}

	if cp.WillQos() > QosExactlyOnce {
		return total, fmt.Errorf("connect/decodeMessage: Invalid QoS level (%d) for %s message", cp.WillQos(), cp.Name())
	}

	if !cp.WillFlag() && (cp.WillRetain() || cp.WillQos() != QosAtMostOnce) {
		return total, fmt.Errorf("connect/decodeMessage: Protocol violation: If the Will Flag (%t) is set to 0 the Will QoS (%d) and Will Retain (%t) fields MUST be set to zero", cp.WillFlag(), cp.WillQos(), cp.WillRetain())
	}

	if cp.UsernameFlag() && !cp.PasswordFlag() {
		return total, fmt.Errorf("connect/decodeMessage: Username flag is set but Password flag is not set")
	}

	if len(src[total:]) < 2 {
		return 0, fmt.Errorf("connect/decodeMessage: Insufficient buffer size. Expecting %d, got %d.", 2, len(src[total:]))
	}

	cp.keepAlive = binary.BigEndian.Uint16(src[total:])
	total += 2

	cp.clientId, n, err = readLPBytes(src[total:])
	total += n
	if err != nil {
		return total, err
	}

	// If the Client supplies a zero-byte ClientId, the Client MUST also set CleanSession to 1
	if len(cp.clientId) == 0 && !cp.CleanSession() {
		return total, ErrIdentifierRejected
	}

	// The ClientId must contain only characters 0-9, a-z, and A-Z
	// We also support ClientId longer than 23 encoded bytes
	// We do not support ClientId outside of the above characters
	if len(cp.clientId) > 0 && !cp.validClientId(cp.clientId) {
		return total, ErrIdentifierRejected
	}

	if cp.WillFlag() {
		cp.willTopic, n, err = readLPBytes(src[total:])
		total += n
		if err != nil {
			return total, err
		}

		cp.willMessage, n, err = readLPBytes(src[total:])
		total += n
		if err != nil {
			return total, err
		}
	}

	// According to the 3.1 spec, it's possible that the passwordFlag is set,
	// but the password string is missing.
	if cp.UsernameFlag() && len(src[total:]) > 0 {
		cp.username, n, err = readLPBytes(src[total:])
		total += n
		if err != nil {
			return total, err
		}
	}

	// According to the 3.1 spec, it's possible that the passwordFlag is set,
	// but the password string is missing.
	if cp.PasswordFlag() && len(src[total:]) > 0 {
		cp.password, n, err = readLPBytes(src[total:])
		total += n
		if err != nil {
			return total, err
		}
	}

	return total, nil
}

func (cp *ConnectPacket) msglen() int {
	total := 0

	ver, ok := SupportedVersions[cp.protolevel]
	if !ok {
		return total
	}

	// 2 bytes protocol name length
	// n bytes protocol name
	// 1 byte protocol level
	// 1 byte connect flags
	// 2 bytes keep alive timer
	total += 2 + len(ver) + 1 + 1 + 2

	// Add the clientID length, 2 is the length prefix
	total += 2 + len(cp.clientId)

	// Add the will topic and will message length, and the length prefixes
	if cp.WillFlag() {
		total += 2 + len(cp.willTopic) + 2 + len(cp.willMessage)
	}

	// Add the username length
	// According to the 3.1 spec, it's possible that the usernameFlag is set,
	// but the user name string is missing.
	if cp.UsernameFlag() && len(cp.username) > 0 {
		total += 2 + len(cp.username)
	}

	// Add the password length
	// According to the 3.1 spec, it's possible that the passwordFlag is set,
	// but the password string is missing.
	if cp.PasswordFlag() && len(cp.password) > 0 {
		total += 2 + len(cp.password)
	}

	return total
}

// validClientId checks the client ID, which is a slice of bytes, to see if it's valid.
// Client ID is valid if it meets the requirement from the MQTT spec:
// 		The Server MUST allow ClientIds which are between 1 and 23 UTF-8 encoded bytes in length,
//		and that contain only the characters
//
//		"0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
func (cp *ConnectPacket) validClientId(cid []byte) bool {
	if cp.ProtoLevel() == 0x3 {
		return true
	}

	return clientIdRegexp.Match(cid)
}
