package cipher

import "crypto/cipher"

type CipherMode  interface {
	SetPadding(padding Padding) CipherMode

	Cipher(block cipher.Block, iv []byte) Cipher
}

type cipherMode struct {
	padding	Padding
}

type ecbCipherModel cipherMode

func NewECBMode() CipherMode {
	return &ecbCipherModel{padding:NewPKCSPadding()}
}

func (ecb *ecbCipherModel) SetPadding(padding Padding) CipherMode {
	ecb.padding = padding
	return ecb
}

func (ecb *ecbCipherModel) Cipher(block cipher.Block, iv []byte) Cipher {
	return NewBlockCipher(ecb.padding, NewECBEncrypter(block), NewECBDecrypter(block))
}

type cbcCipherModel cipherMode

func NewCBCMode() CipherMode {
	return &cbcCipherModel{padding:NewPKCSPadding()}
}

func (cbc *cbcCipherModel) SetPadding(padding Padding) CipherMode {
	cbc.padding = padding
	return cbc
}

func (cbc *cbcCipherModel) Cipher(block cipher.Block, iv []byte) Cipher {
	return NewBlockCipher(cbc.padding, cipher.NewCBCEncrypter(block, iv), cipher.NewCBCDecrypter(block, iv))
}