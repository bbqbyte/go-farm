package hashids

import (
	"fmt"
	"testing"
)

func TestHash(t *testing.T) {
	hd := NewData()
	hd.Salt = "bbqbyte@github.com"
	hd.MinLength = 30
	h, _ := NewWithData(hd)
	e, _ := h.Encode([]int{1})
	fmt.Println(e)
	d, _ := h.DecodeWithError(e)
	fmt.Println(d)
}

func TestHash1(t *testing.T) {
	hd := NewDataWithLowerCase()
	hd.Salt = "bbqbyte@github.com"
	hd.MinLength = 18
	h, _ := NewWithData(hd)
	e, _ := h.Encode([]int{1})
	fmt.Println(e)
	d, _ := h.DecodeWithError(e)
	fmt.Println(d)
}