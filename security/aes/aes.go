package aes

import (
	"crypto/aes"
	pbciper "github.com/bbqbyte/go-farm/security/cipher"
)

func GenerateKey(key []byte, keylen int) (genKey []byte) {
	genKey = make([]byte, keylen)
	copy(genKey, key)
	for i := keylen; i < len(key); {
		for j := 0; j < keylen && i < len(key); j, i = j+1, i+1 {
			genKey[j] ^= key[i]
		}
	}
	return genKey
}

// CBC模式
// 加密
func AesCBCEncryptWithIV(src, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if iv == nil {
		iv = key[:block.BlockSize()]
	}
	pcipher := pbciper.NewCBCMode().Cipher(block, iv)
	return pcipher.Encrypt(src), nil
}

func AesCBCEncrypt(src, key []byte) ([]byte, error) {
	return AesCBCEncryptWithIV(src, key, nil)
}

// 解密
func AesCBCDecryptWithIV(encrypted []byte, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if iv == nil {
		iv = key[:block.BlockSize()]
	}
	pcipher := pbciper.NewCBCMode().Cipher(block, iv)
	return pcipher.Decrypt(encrypted), nil
}

func AesCBCDecrypt(encrypted []byte, key []byte) ([]byte, error) {
	return AesCBCDecryptWithIV(encrypted, key, nil)
}

// ECB模式
// 加密
func AesECBEncrypt(src, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	pcipher := pbciper.NewECBMode().Cipher(block, key[:block.BlockSize()])
	return pcipher.Encrypt(src), nil
}

// 解密
func AesECBDecrypt(encrypted []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	pcipher := pbciper.NewECBMode().Cipher(block, key[:block.BlockSize()])
	return pcipher.Decrypt(encrypted), nil
}
