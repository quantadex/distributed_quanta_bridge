package key_manager

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPubKey(t *testing.T) {
	key, err := NewGrapheneKeyManager("bb2aeb9eebaaa29d79ed81699ee49a912c19c59b9350f8f8d3d81b12fa178495")
	assert.NoError(t, err, "expect no error")

	err = key.LoadNodeKeys("5Jd9vxNwWXvMnBpcVm58gwXkJ4smzWDv9ChiBXwSRkvCTtekUrx")
	assert.NoError(t, err, "expect no error")

	pub, err := key.GetPublicKey()
	assert.NoError(t, err, "expect no error")

	fmt.Println(pub)
	assert.Equal(t, pub, "QA5Bd4w7Y48aErDG442XRP5Cap3PjmJTZ2wvuof72TWKAzNznt6q")

}
