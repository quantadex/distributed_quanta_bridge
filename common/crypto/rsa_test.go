package crypto

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"encoding/hex"
)

func TestRsa(t *testing.T) {
	priv, err := LoadSSHKey("/Users/quoc/.ssh/id_rsa")
	assert.NoError(t, err)

	res := EncryptWithPublicKey([]byte("blah"), &priv.PublicKey)
	decrypted := DecryptWithPrivateKey(res, priv)
	println(len(res),hex.EncodeToString(res),string(decrypted))
}
