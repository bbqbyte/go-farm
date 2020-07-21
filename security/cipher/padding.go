package cipher

import "bytes"

type Padding interface {
	Padding(src []byte, blockSize int) []byte

	UnPadding(src []byte) []byte
}

type padding struct{}

// PKCS5Padding & PKCS7Padding(blockSize := 16)
type pkcsPadding padding

func NewPKCSPadding() Padding {
	return &pkcsPadding{}
}

func (p *pkcsPadding) Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func (p *pkcsPadding) UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// zero padding
type zeroPadding padding

func NewZeroPadding() Padding {
	return &zeroPadding{}
}

func (p *zeroPadding) Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(0)}, padding)
	return append(ciphertext, padtext...)
}

func (p *zeroPadding) UnPadding(origData []byte) []byte {
	return bytes.TrimRightFunc(origData, func(r rune) bool {
		return r == rune(0)
	})
}
