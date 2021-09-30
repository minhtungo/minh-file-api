package cryp

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"
	"errors"
)

func Decrypt(encryptedString string, key []byte) ([]byte, error) {
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

func Encrypt(stringToEncrypt string, key []byte) ([]byte, error) {
	plaintext := []byte(stringToEncrypt)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errors.New("1")
	}

	//Create a new GCM 
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errors.New("2")
	}

	//Create a nonce. Nonce should be from GCM
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, errors.New("3")
	}

	//Encrypt the data using aesGCM.Seal
	return aesGCM.Seal(nonce, nonce, plaintext, nil), nil
}


