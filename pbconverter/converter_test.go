package pbconverter

import (
	"testing"
	"log"
)

func TestToInt8(t *testing.T) {
	tests := []string {"125", "-123", "abcdef", "200"}
	expected := []int8{125, -123, 0, 0}

	for i := 0; i < len(tests); i++ {
		result, _ := ToInt8(tests[i], 0)
		if result != expected[i] {
			t.Log("Case ", i, ": expected ", expected[i], " when result is ", result)
			t.FailNow()
		}
	}
}

func TestToInt64(t *testing.T) {
	tests := []string {"", "1000", "-123", "abcdef", "100000000000000000000000000000000000000000000"}
	expected := []int64{0, 1000, -123, 0, 0}

	for i := 0; i < len(tests); i++ {
		result, _ := ToInt64(tests[i], 0)
		if result != expected[i] {
			t.Log("Case ", i, ": expected ", expected[i], " when result is ", result)
			t.FailNow()
		}
	}
}

func TestToBigInt(t *testing.T) {
	res := ToBigInt("123", 10)
	log.Printf("%v", res)
}