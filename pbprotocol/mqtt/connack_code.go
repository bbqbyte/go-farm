package mqtt

// ConnackCode is the type representing the return code in the CONNACK packet,
// returned after the initial CONNECT packet
type ConnackCode byte

const (
	// Connection accepted
	ConnectionAccepted ConnackCode = iota

	// The Server does not support the level of the MQTT protocol requested by the Client
	ErrInvalidProtocolVersion

	// The Client identifier is correct UTF-8 but not allowed by the server
	ErrIdentifierRejected

	// The Network Connection has been made but the MQTT service is unavailable
	ErrServerUnavailable

	// The data in the user name or password is malformed
	ErrBadUsernameOrPassword

	// The Client is not authorized to connect
	ErrNotAuthorized
)

// Value returns the value of the ConnackCode, which is just the byte representation
func (cc ConnackCode) Value() byte {
	return byte(cc)
}

// Desc returns the description of the ConnackCode
func (cc ConnackCode) Desc() string {
	switch cc {
	case 0:
		return "Connection accepted"
	case 1:
		return "The Server does not support the level of the MQTT protocol requested by the Client"
	case 2:
		return "The Client identifier is correct UTF-8 but not allowed by the server"
	case 3:
		return "The Network Connection has been made but the MQTT service is unavailable"
	case 4:
		return "The data in the user name or password is malformed"
	case 5:
		return "The Client is not authorized to connect"
	}

	return "Reserved for future use"
}

// Valid checks to see if the ConnackCode is valid. Currently valid codes are <= 5
func (cc ConnackCode) Valid() bool {
	return cc <= 5
}

// Error returns the corresonding error string for the ConnackCode
func (cc ConnackCode) Error() string {
	switch cc {
	case 0:
		return "Connection accepted"
	case 1:
		return "Connection Refused, unacceptable protocol version"
	case 2:
		return "Connection Refused, identifier rejected"
	case 3:
		return "Connection Refused, Server unavailable"
	case 4:
		return "Connection Refused, bad user name or password"
	case 5:
		return "Connection Refused, not authorized"
	}

	return "Unknown error"
}
