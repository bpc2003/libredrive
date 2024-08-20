package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"
	"log"
)

func Encrypt(key, buf []byte) (ciphertext []byte) {
	key, _ = hex.DecodeString(string(key))
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal(err)
	}

	ciphertext = make([]byte, aes.BlockSize+len(buf))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		log.Fatal(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], buf)

	return
}

func Decrypt(key, buf []byte) ([]byte, error) {
	key, _ = hex.DecodeString(string(key))
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	iv := buf[:aes.BlockSize]
	ciphertext := buf[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	stream.XORKeyStream(ciphertext, ciphertext)
	return ciphertext, nil
}
