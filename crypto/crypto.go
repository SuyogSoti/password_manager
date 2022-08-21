package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/suyogsoti/password_manager/auth"
)

func getKeyFromPassword(password string) []byte {
	diff := 32 - len(password)
	if diff > 0 {
		password = password + strings.Repeat("a", diff)
	}
	return []byte(password)
}

func Encrypt(c *gin.Context, sitePswd string) (string, error) {
	block, err := aes.NewCipher(getKeyFromPassword(auth.GetCredentials(c).Password))
	if err != nil {
		return "", fmt.Errorf("could not create cipher: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("could not generate nonce: %w", err)
	}

	ciphertext := aesGCM.Seal(nonce, nonce, []byte(sitePswd), nil)
	return fmt.Sprintf("%x", ciphertext), nil
}


func Decrypt(c *gin.Context, encryptedString string) (string, error) {

	key := getKeyFromPassword(auth.GetCredentials(c).Password)
	enc, _ := hex.DecodeString(encryptedString)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	//Create a new GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	//Get the nonce size
	nonceSize := aesGCM.NonceSize()

	//Extract the nonce from the encrypted data
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]

	//Decrypt the data
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s", plaintext), nil
}
