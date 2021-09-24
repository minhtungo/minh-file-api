package cryp

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/md5"
	"encoding/hex"
	"io"
	"errors"
)

func Decrypt(encryptedString string, passphrase string) ([]byte, error) {
	key := []byte(CreateHash(passphrase))
	ciphertext, err := hex.DecodeString(encryptedString)
	if err != nil {
		return nil, err
	}

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	//Create a new GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	//Get the nonce size
	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
        return nil, errors.New("ciphertext too short")
    }

	//Extract the nonce from the encrypted data
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	//Decrypt the data
	return aesGCM.Open(nil, nonce, ciphertext, nil)
}

func Encrypt(stringToEncrypt string, passphrase string) ([]byte, error) {

	//Convert decode the key to bytes
	key := []byte(CreateHash(passphrase))
	plaintext := []byte(stringToEncrypt)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	//Create a new GCM - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	//https://golang.org/pkg/crypto/cipher/#NewGCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	//Create a nonce. Nonce should be from GCM
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	//Encrypt the data using aesGCM.Seal
	return aesGCM.Seal(nonce, nonce, plaintext, nil), nil
}

func CreateHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}
