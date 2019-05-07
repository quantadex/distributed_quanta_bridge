package crypto

import (
	"io/ioutil"
	"log"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/rand"
	"encoding/pem"
	"crypto/x509"
)

func LoadSSHKey(privateKey string) (*rsa.PrivateKey, error){
	key, err := ioutil.ReadFile(privateKey)
	if err != nil {
		log.Fatalf("Unable to read private key: %v", err)
		return nil, err
	}
	return BytesToPrivateKey(key)
}

// EncryptWithPublicKey encrypts data with public key
func EncryptWithPublicKey(msg []byte, pub *rsa.PublicKey) []byte {
	hash := sha512.New()
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, pub, msg, nil)
	if err != nil {
		log.Fatalf("encrypt error:",err.Error())
	}
	return ciphertext
}

// BytesToPrivateKey bytes to private key
func BytesToPrivateKey(priv []byte) (*rsa.PrivateKey,error) {
	block, _ := pem.Decode(priv)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		log.Println("is encrypted pem block")
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			log.Fatalf("decrypt pem error:",err.Error())
			return nil, err
		}
	}
	key, err := x509.ParsePKCS1PrivateKey(b)
	if err != nil {
		log.Fatalf("parse private key error:", err.Error())
		return nil, err
	}
	return key,nil
}


// DecryptWithPrivateKey decrypts data with private key
func DecryptWithPrivateKey(ciphertext []byte, priv *rsa.PrivateKey) []byte {
	hash := sha512.New()
	plaintext, err := rsa.DecryptOAEP(hash, rand.Reader, priv, ciphertext, nil)
	if err != nil {
		log.Fatalf("decrypt error:", err.Error())
	}
	return plaintext
}