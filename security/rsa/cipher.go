package rsa

import (
	"bytes"
	"crypto"
	"errors"
	"github.com/bbqbyte/go-farm/logger"
	"io/ioutil"
)

type Cipher interface {
	Encrypt(plainText []byte) ([]byte, error)
	Decrypt(cipherText []byte) ([]byte, error)
	Sign(src []byte, hash crypto.Hash) ([]byte, error)
	Verify(src []byte, sign []byte, hash crypto.Hash) error

	EncryptPri(plainText []byte) ([]byte, error)
	DecryptPub(cipherText []byte) ([]byte, error)
}

func NewCipher(key Key, padding Padding, cipherMode CipherMode, signMode SignMode) Cipher {
	return &cipher{key: key, padding: padding, cipherMode: cipherMode, sign: signMode}
}

type cipher struct {
	key        Key
	cipherMode CipherMode
	sign       SignMode
	padding    Padding
}

func (cipher *cipher) Encrypt(plainText []byte) ([]byte, error) {
	groups := cipher.padding.Padding(plainText)
	buffer := bytes.Buffer{}
	for _, plainTextBlock := range groups {
		cipherText, err := cipher.cipherMode.Encrypt(plainTextBlock, cipher.key.PublicKey())
		if err != nil {
			plog.Error("[RSA] Encrypt", log4go.Error(err))
			return nil, err
		}
		buffer.Write(cipherText)
	}
	return buffer.Bytes(), nil
}

func (cipher *cipher) Decrypt(cipherText []byte) ([]byte, error) {
	if len(cipherText) == 0 {
		return nil, errors.New("cipher can't be null")
	}
	groups := grouping(cipherText, cipher.key.Modulus())
	buffer := bytes.Buffer{}
	for _, cipherTextBlock := range groups {
		plainText, err := cipher.cipherMode.Decrypt(cipherTextBlock, cipher.key.PrivateKey())
		if err != nil {
			plog.Error("[RSA] Decrypt", log4go.Error(err))
			return nil, err
		}
		buffer.Write(plainText)
	}
	return buffer.Bytes(), nil
}

func (cipher *cipher) Sign(src []byte, hash crypto.Hash) ([]byte, error) {
	return cipher.sign.Sign(src, hash, cipher.key.PrivateKey())
}

func (cipher *cipher) Verify(src []byte, sign []byte, hash crypto.Hash) error {
	return cipher.sign.Verify(src, sign, hash, cipher.key.PublicKey())
}

// 公钥解密
func (cipher *cipher) DecryptPub(input []byte) ([]byte, error) {
	output := bytes.NewBuffer(nil)
	err := pubKeyIO(cipher.key.PublicKey(), bytes.NewReader(input), output)
	if err != nil {
		return []byte(""), err
	}
	return ioutil.ReadAll(output)
}

// 私钥加密
func (cipher *cipher) EncryptPri(input []byte) ([]byte, error) {
	output := bytes.NewBuffer(nil)
	err := priKeyIO(cipher.key.PrivateKey(), bytes.NewReader(input), output)
	if err != nil {
		return []byte(""), err
	}
	return ioutil.ReadAll(output)
}
