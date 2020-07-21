package cipher

import "crypto/cipher"

type ecb struct {
	block     cipher.Block
	blockSize int
}

type ecbEncrypter ecb

func (e *ecbEncrypter) BlockSize() int {
	return e.blockSize
}

func (e *ecbEncrypter) CryptBlocks(dst, src []byte) {
	if len(src) % e.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	for len(src) > 0 {
		e.block.Encrypt(dst, src[:e.blockSize])
		src = src[e.blockSize:]
		dst = dst[e.blockSize:]
	}
}

type ecbDecrypter ecb

func (e *ecbDecrypter) BlockSize() int {
	return e.blockSize
}

func (e *ecbDecrypter) CryptBlocks(dst, src []byte) {
	if len(src) % e.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	for len(src) > 0 {
		e.block.Decrypt(dst, src[:e.blockSize])
		src = src[e.blockSize:]
		dst = dst[e.blockSize:]
	}
}

func NewECBEncrypter(block cipher.Block) cipher.BlockMode {
	return &ecbEncrypter{block: block, blockSize: block.BlockSize()}
}

func NewECBDecrypter(block cipher.Block) cipher.BlockMode {
	return &ecbDecrypter{block: block, blockSize: block.BlockSize()}
}