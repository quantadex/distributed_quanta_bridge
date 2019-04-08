package crypto

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestBip32(t *testing.T) {
	bipKey := "xprv9zJzAuYXn9Lrgqd6Yu1JQu3SzEFGkft41Jd9xWJmWuhwNuj9Q3qA47DG1zvca2UimhPeTAGhtMcaNFigsMnqw2F9d5NsDUy19w8Q6S2zr8N"
	pub, priv, err := GenerateBip32Key(bipKey, 100)
	assert.NoError(t,err)
	println(pub.Hex(), hex.EncodeToString(crypto.FromECDSA(priv)))
}