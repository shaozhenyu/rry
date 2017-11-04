package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

var (
	commonIV = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}
	key_text = "666666666666666666666666666666666"
)

func getKey(key string) ([]byte, error) {
	keyLen := len(key)
	if keyLen < 16 {
		return nil, errors.New("key can not less than 16")
	}
	arrKey := []byte(key)
	if keyLen >= 32 {
		return arrKey[:32], nil
	}
	if keyLen >= 24 {
		return arrKey[:24], nil
	}
	return arrKey[:16], nil
}

//加密模式：CFB密码反馈模式
func Encrypter(plaintext []byte) ([]byte, error) {
	key, err := getKey(key_text)
	if err != nil {
		return nil, err
	}

	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	cfb := cipher.NewCFBEncrypter(c, commonIV)
	ciphertext := make([]byte, len(plaintext))
	cfb.XORKeyStream(ciphertext, plaintext)
	//log.Printf("%s->%x\n", plaintext, ciphertext)
	return ciphertext, nil
}

//解密
func Decrypter(ciphertext []byte) ([]byte, error) {
	key, err := getKey(key_text)
	if err != nil {
		return nil, err
	}

	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	cfbdec := cipher.NewCFBDecrypter(c, commonIV)
	plaintextCopy := make([]byte, len(ciphertext))
	cfbdec.XORKeyStream(plaintextCopy, ciphertext)
	//log.Printf("%x=>%s\n", ciphertext, plaintextCopy)
	return plaintextCopy, nil
}
