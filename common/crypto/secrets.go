package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"github.com/quantadex/distributed_quanta_bridge/node/common"
	"github.com/spf13/viper"
	"golang.org/x/crypto/sha3"
	"io"
	"io/ioutil"
	"strings"
)

func DecryptSecretsFile(file string, password string) (common.Secrets, error) {
	secrets := common.Secrets{}
	ciphertext, err := ioutil.ReadFile(file)
	// if our program was unable to read the file
	// print out the reason why it can't
	if err != nil {
		return secrets, err
	}

	key := []byte(strings.TrimSpace(password))
	hash := sha3.Sum256(key)
	c, err := aes.NewCipher(hash[:])

	if err != nil {
		return secrets, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return secrets, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return secrets, err
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return secrets, err
	}

	err = viper.ReadConfig(bytes.NewBuffer(plaintext))
	if err != nil {
		return secrets, err
	}

	err = viper.Unmarshal(&secrets)
	if err != nil {
		return secrets, err
	}
	return secrets, nil
}

func EncryptSecretsFile(password string, data []byte, outFile string) error {
	key := []byte(strings.TrimSpace(password))
	hash := sha3.Sum256(key)
	c, err := aes.NewCipher(hash[:])
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}

	err = ioutil.WriteFile(outFile, gcm.Seal(nonce, nonce, data, nil), 0777)
	if err != nil {
		return err
	}
	println("Successful written to " + outFile)
	return nil
}
