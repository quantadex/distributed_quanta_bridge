package crypto

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/agl/ed25519"
	"github.com/btcsuite/btcutil"
	"github.com/go-errors/errors"
	"github.com/scorum/bitshares-go/sign"
)

func SignMessage(msg interface{}, privateKey string) *string {
	w, _ := btcutil.DecodeWIF(privateKey)

	bData := new(bytes.Buffer)
	json.NewEncoder(bData).Encode(msg)

	digest := sha256.Sum256(bData.Bytes())

	sig := sign.SignBufferSha256(digest[:], w.PrivKey.ToECDSA())
	sigHex := hex.EncodeToString(sig)
	return &sigHex
}

func VerifyMessage(msg interface{}, publicKey string, signature string) bool {
	//bData := new(bytes.Buffer)
	//json.NewEncoder(bData).Encode(msg)
	//
	//sigBytes,_ := hex.DecodeString(signature)
	//
	//
	//if err := Verify(bData.Bytes(), sData, publicKey); err != nil {
	//	fmt.Println("returning false from verify")
	//	return false
	//}

	return true
}

func Verify(input []byte, sig []byte, publicKey string) error {
	fmt.Println("lenght = ", publicKey, len(publicKey))
	if len(sig) != 64 {
		fmt.Println("returning false from len")
		return errors.New("Signature verification failed")
	}

	var asig [64]byte
	copy(asig[:], sig[:])
	var key [32]byte
	copy(key[:], publicKey)

	if !ed25519.Verify(&key, input, &asig) {
		fmt.Println("returning false fom ed25519")
		return errors.New("Signature verification failed")
	}
	return nil
}

/*
func VerifyMessage(msg interface{}, publicKey string, signature string) bool {
	key, err := keypair.Parse(publicKey)

	if err != nil {
		return false
	}

	bData := new(bytes.Buffer)
	json.NewEncoder(bData).Encode(msg)

	sData, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false
	}

	if err := key.Verify(bData.Bytes(), sData); err != nil {
		return false
	}

	return true
}

func SignMessage(msg interface{}, publicKey string) *string {
	key, err := keypair.Parse(publicKey)

	if err != nil {
		return nil
	}

	bData := new(bytes.Buffer)
	json.NewEncoder(bData).Encode(msg)

	signed, _ := key.Sign(bData.Bytes())
	signedbase64 := base64.StdEncoding.EncodeToString(signed)

	return &signedbase64
}
*/
