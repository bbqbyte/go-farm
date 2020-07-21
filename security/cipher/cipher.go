package cipher

import "crypto/cipher"

type Cipher interface {
	Encrypt(src []byte) []byte
	Decrypt(src []byte) []byte
}

type blockCipher struct {
	padding Padding
	encrypt cipher.BlockMode
	decrypt cipher.BlockMode
}

func NewBlockCipher(padding Padding, encrypt, decrypt cipher.BlockMode) Cipher {
	return &blockCipher{
		encrypt:   encrypt,
		decrypt:   decrypt,
		padding:   padding,
	}
}

func (blockCipher *blockCipher) Encrypt(plaintext []byte) []byte {
	plaintext = blockCipher.padding.Padding(plaintext, blockCipher.encrypt.BlockSize())
	ciphertext := make([]byte, len(plaintext))
	blockCipher.encrypt.CryptBlocks(ciphertext, plaintext)
	return ciphertext
}

func (blockCipher *blockCipher) Decrypt(ciphertext []byte) []byte {
	plaintext := make([]byte, len(ciphertext))
	blockCipher.decrypt.CryptBlocks(plaintext, ciphertext)
	plaintext = blockCipher.padding.UnPadding(plaintext)
	return plaintext
}

// 流加密
type streamCipher struct {
	encrypt cipher.Stream
	decrypt cipher.Stream
}

func (streamCipher *streamCipher) Encrypt(plaintext []byte) []byte {
	ciphertext := make([]byte, len(plaintext))
	streamCipher.encrypt.XORKeyStream(ciphertext, plaintext)
	return ciphertext
}
func (streamCipher *streamCipher) Decrypt(ciphertext []byte) []byte {
	plaintext := make([]byte, len(ciphertext))
	streamCipher.decrypt.XORKeyStream(plaintext, ciphertext)
	return plaintext
}