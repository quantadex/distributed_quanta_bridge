package crypto

import (
	"github.com/btcsuite/btcutil"
	"github.com/stretchr/testify/assert"
	"testing"
)

/*
{
 "brain_priv_key": "PALAVER CAMP TWILT BRABBLE BERIDE RIFF DAUNTON POORISH CIRCLET ENROUGH VOIDER PILOSE SHALE GOBLINE TINDER CORGE",
 "wif_priv_key": "5JyYu5DCXbUznQRSx3XT2ZkjFxQyLtMuJ3y6bGLKC3TZWPHMDxj",
 "pub_key": "QA5oEKWyjQzhvBdNCF4JufR7aVrU2bjFc9cEPFb3fthxqs1UjZtu"
}
*/

func TestSignMessage(t *testing.T) {
	msg := "some string"

	sig := SignMessage(msg, "5JyYu5DCXbUznQRSx3XT2ZkjFxQyLtMuJ3y6bGLKC3TZWPHMDxj")
	println("sig ", *sig)
	success := VerifyMessage(msg, "QA5oEKWyjQzhvBdNCF4JufR7aVrU2bjFc9cEPFb3fthxqs1UjZtu", *sig)

	if !success {
		t.Error("expect to be successful")
	}
}

func TestGetBitcoinAddressFromGraphene(t *testing.T) {
	w, _ := btcutil.DecodeWIF("5JyYu5DCXbUznQRSx3XT2ZkjFxQyLtMuJ3y6bGLKC3TZWPHMDxj")

	address, err := GetBitcoinAddressFromGraphene(w.PrivKey.PubKey())
	assert.NoError(t, err)
	println(address)
}
