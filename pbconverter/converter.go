package pbconverter

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"time"
	"reflect"
)

// ToString convert the input to a string.
func ToString(obj interface{}) string {
	res := fmt.Sprintf("%v", obj)
	return res
}

// ToDeepString convert the input to a string.
func ToDeepString(x interface{}) string {
	switch y := x.(type) {

	// Handle dates with special logic
	// This needs to come above the fmt.Stringer
	// test since time.Time's have a .String()
	// method
	case time.Time:
		return y.Format("A Monday")

		// Handle type string
	case string:
		return y

		// Handle type with .String() method
	case fmt.Stringer:
		return y.String()

		// Handle type with .Error() method
	case error:
		return y.Error()

	}

	// Handle named string type
	if v := reflect.ValueOf(x); v.Kind() == reflect.String {
		return v.String()
	}

	// Fallback to fmt package for anything else like numeric types
	return fmt.Sprint(x)
}

// ToJSON convert the input to a valid JSON string
func ToJSONString(obj interface{}) ([]byte, error) {
	res, err := json.Marshal(obj)
	if err != nil {
		res = []byte("")
	}
	return res, err
}

// ToBoolean returns the boolean value represented by the string.
//
// It accepts 1, 1.0, t, T, TRUE, true, True, YES, yes, Yes,Y, y, ON, on, On,
// 0, 0.0, f, F, FALSE, false, False, NO, no, No, N,n, OFF, off, Off.
// Any other value returns an error.
func ToBoolean(val interface{}) (value bool, err error) {
	if val != nil {
		switch v := val.(type) {
		case bool:
			return v, nil
		case string:
			switch v {
			case "1", "t", "T", "true", "TRUE", "True", "YES", "yes", "Yes", "Y", "y", "ON", "on", "On":
				return true, nil
			case "0", "f", "F", "false", "FALSE", "False", "NO", "no", "No", "N", "n", "OFF", "off", "Off":
				return false, nil
			}
		case int8, int32, int64:
			strV := fmt.Sprintf("%d", v)
			if strV == "1" {
				return true, nil
			} else if strV == "0" {
				return false, nil
			}
		case float64:
			if v == 1.0 {
				return true, nil
			} else if v == 0.0 {
				return false, nil
			}
		}
		return false, fmt.Errorf("parsing %q: invalid syntax", val)
	}
	return false, fmt.Errorf("parsing <nil>: invalid syntax")
}

// ToInt8 convert the input string to an int8, or default value if the input is not an integer.
func ToInt8(str string, defaultValue int64) (int8, error) {
	res, err := toInt(str, 8, defaultValue)
	return int8(res), err
}

// ToInt16 convert the input string to an int16, or default value if the input is not an integer.
func ToInt16(str string, defaultValue int64) (int16, error) {
	res, err := toInt(str, 16, defaultValue)
	return int16(res), err
}

// ToInt32 convert the input string to an int32, or default value if the input is not an integer.
func ToInt32(str string, defaultValue int64) (int32, error) {
	res, err := toInt(str, 32, defaultValue)
	return int32(res), err
}

// ToInt64 convert the input string to an int64, or default value if the input is not an integer.
func ToInt64(str string, defaultValue int64) (int64, error) {
	return toInt(str, 64, defaultValue)
}

func toInt(str string, bitSize int, defaultValue int64) (int64, error) {
	res, err := strconv.ParseInt(str, 10, bitSize)
	if err != nil {
		res = defaultValue
	}
	return res, err
}

// ToFloat32 convert the input string to a float32, or given value if the input is not a float.
func ToFloat32(str string, defaultValue float64) (float32, error) {
	res, err := strconv.ParseFloat(str, 32)
	if err != nil {
		res = defaultValue
	}
	return float32(res), err
}

// ToFloat64 convert the input string to a float64, or given value if the input is not a float.
func ToFloat64(str string, defaultValue float64) (float64, error) {
	res, err := strconv.ParseFloat(str, 64)
	if err != nil {
		res = defaultValue
	}
	return res, err
}

// ToUint8 convert the input string to an uint8, or default value if the input is not an integer.
func ToUint8(str string, defaultValue uint64) (uint8, error) {
	res, err := toUint(str, 8, defaultValue)
	return uint8(res), err
}

// ToUint16 convert the input string to an uint16, or default value if the input is not an integer.
func ToUint16(str string, defaultValue uint64) (uint16, error) {
	res, err := toUint(str, 16, defaultValue)
	return uint16(res), err
}

// ToUint32 convert the input string to an uint32, or default value if the input is not an integer.
func ToUint32(str string, defaultValue uint64) (uint32, error) {
	res, err := toUint(str, 32, defaultValue)
	return uint32(res), err
}

// ToUint64 convert the input string to an uint64, or default value if the input is not an integer.
func ToUint64(str string, defaultValue uint64) (uint64, error) {
	return toUint(str, 64, defaultValue)
}

func toUint(str string, bitSize int, defaultValue uint64) (uint64, error) {
	res, err := strconv.ParseUint(str, 10, bitSize)
	if err != nil {
		res = defaultValue
	}
	return res, err
}

func ToBigInt(str string, base int) *big.Int {
	v := big.NewInt(0)
	v.SetString(str, base)
	return v
}

func Int64ToBytes(i int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

func BytesToInt64(buf []byte) int64 {
	return int64(binary.BigEndian.Uint64(buf))
}
