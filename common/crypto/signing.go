package crypto

import (
	"github.com/stellar/go/keypair"
	"bytes"
	"encoding/json"
	"encoding/base64"
)

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