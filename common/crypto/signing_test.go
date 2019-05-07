package crypto

import (
	"github.com/btcsuite/btcutil"
	"github.com/stretchr/testify/assert"
	"testing"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"crypto/sha256"
	"math/big"
	"crypto/elliptic"
	"github.com/btcsuite/btcd/btcec"
	"encoding/hex"
	"github.com/btcsuite/btcutil/base58"
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

func TestGetGraphenePublicKey(t *testing.T) {
	address, err := GenerateGrapheneKeyWithSeed("pooja", "")
	assert.NoError(t, err)
	println(address)

	pubKey, err := NewGraphenePublicKeyFromString(address)
	assert.NoError(t, err)
	println(pubKey)
}

func makeKey(curve elliptic.Curve, x *big.Int, y *big.Int) *btcec.PrivateKey {
	priv,_ := btcec.PrivKeyFromBytes(curve, x.Bytes())
	return (*btcec.PrivateKey)(priv)
}

// Qx,Qy = nB * G
// Qx,Qy = nA * G
// K = nB * (Qx, Qy)
func TestEncrypt(t *testing.T) {
	digest := sha256.Sum256([]byte("BOB"))
	masterKey, err := hdkeychain.NewMaster(digest[:], &chaincfg.TestNet3Params)
	childKey, err := masterKey.Child(0)
	bobPriv, err := childKey.ECPrivKey()
	bobPub, err := childKey.ECPubKey()
	assert.NoError(t, err)

	digest = sha256.Sum256([]byte("ALICE"))
	masterKey, err = hdkeychain.NewMaster(digest[:], &chaincfg.TestNet3Params)
	childKey, err = masterKey.Child(0)
	alicePriv, err := childKey.ECPrivKey()
	alicePub, err := childKey.ECPubKey()

	s1x, s1y := bobPriv.Curve.ScalarMult(alicePub.X,alicePub.Y, bobPriv.D.Bytes())
	s2x, s2y := alicePriv.Curve.ScalarMult(bobPub.X,bobPub.Y, alicePriv.D.Bytes())

	println(s1x.String(), s1y.String())
	println(s2x.String(), s2y.String())

	k1 := makeKey(bobPriv.Curve, s1x, s1y)
	k2 := makeKey(bobPriv.Curve, s2x, s2y)

	println(len(k1.Serialize()), hex.EncodeToString(k1.Serialize()))
	println(hex.EncodeToString(k2.Serialize()))

	msg := []byte("BOB")
	cipher, _ := btcec.Encrypt(k1.PubKey(), msg)

	println(len(cipher),base58.Encode(cipher))

	plain,err := btcec.Decrypt(k1, cipher)
	println("decrypted ", string(plain),err)

	plain2,err := btcec.Decrypt(k2, cipher)
	println("decrypted ", string(plain2),err)
}