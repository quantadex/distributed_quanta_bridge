package crypto

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/agl/ed25519"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/base58"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/gcash/bchutil"
	"github.com/juju/errors"
	chaincfg3 "github.com/ltcsuite/ltcd/chaincfg"
	"github.com/ltcsuite/ltcutil"
	"github.com/scorum/bitshares-go/sign"

	chaincfg2 "github.com/gcash/bchd/chaincfg"
)

const PREFIX = "QA"

func NewGraphenePublicKeyFromString(key string) (*btcec.PublicKey, error) {
	prefix := key[:len(PREFIX)]

	if prefix != PREFIX {
		return nil, errors.New("Wrong chain")
	}

	b58 := base58.Decode(key[len(PREFIX):])
	if len(b58) < 5 {
		return nil, errors.New(fmt.Sprintf("Invalid public key 1 %d %v", len(b58), b58))
	}

	chk1 := b58[len(b58)-4:]

	keyBytes := b58[:len(b58)-4]
	chk2, err := Ripemd160Checksum(keyBytes)
	if err != nil {
		return nil, errors.Annotate(err, "Invalid checksum")
	}

	if !bytes.Equal(chk1, chk2) {
		return nil, errors.New("Invalid public key 2")
	}

	pub, err := btcec.ParsePubKey(keyBytes, btcec.S256())
	if err != nil {
		return nil, errors.Annotate(err, "ParsePubKey??")
	}

	return pub, nil
}

func GenerateGrapheneKeyWithSeed(str string, prefix string) (string, error) {
	digest := sha256.Sum256([]byte(prefix + ":" + str))
	//digest2 := bytes.NewBuffer([]byte{0x2})
	//digest2.Write(digest[:])
	masterKey, err := hdkeychain.NewMaster(digest[:], &chaincfg.TestNet3Params)
	if err != nil {
		return "", errors.Annotate(err, "Could not get masterKey key")
	}
	childKey, err := masterKey.Child(0)
	if err != nil {
		return "", errors.Annotate(err, "Could not get child key")
	}
	pubKey, err := childKey.ECPubKey()
	if err != nil {
		return "", errors.Annotate(err, "Could not get ECPubKey")
	}

	digest2 := pubKey.SerializeCompressed()

	chk, err := Ripemd160Checksum(digest2)
	if err != nil {
		return "", err
	}
	b := append(digest2, chk...)
	pubkey := base58.Encode(b)
	return fmt.Sprintf("%s%s", PREFIX, pubkey), nil
}

func GetGraphenePublicKey(pubKey *btcec.PublicKey) (string, error) {
	buf := pubKey.SerializeCompressed()
	chk, err := Ripemd160Checksum(buf)
	if err != nil {
		return "", err
	}
	b := append(buf, chk...)
	pubkey := base58.Encode(b)
	return fmt.Sprintf("%s%s", PREFIX, pubkey), nil
}

func GetBitcoinAddressFromGraphene(pubKey *btcec.PublicKey) (*btcutil.AddressPubKey, error) {
	address, err := btcutil.NewAddressPubKey(pubKey.SerializeUncompressed(), &chaincfg.RegressionNetParams)
	if err != nil {
		return nil, err
	}
	address.SetFormat(btcutil.PKFUncompressed)
	return address, err
}

func GetLitecoinAddressFromGraphene(pubKey *btcec.PublicKey) (*ltcutil.AddressPubKey, error) {
	address, err := ltcutil.NewAddressPubKey(pubKey.SerializeUncompressed(), &chaincfg3.RegressionNetParams)
	if err != nil {
		return nil, err
	}
	//address.SetFormat(ltcutil.PKFUncompressed)
	return address, err
}

func GetBCHAddressFromGraphene(pubKey *btcec.PublicKey) (*bchutil.AddressPubKey, error) {
	address, err := bchutil.NewAddressPubKey(pubKey.SerializeUncompressed(), &chaincfg2.RegressionNetParams)
	if err != nil {
		return nil, err
	}
	//address.SetFormat(bchutil.PKFUncompressed)
	return address, err
}

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
	bData := new(bytes.Buffer)
	json.NewEncoder(bData).Encode(msg)
	digest := sha256.Sum256(bData.Bytes())

	sig, _ := hex.DecodeString(signature)
	p, _, err := btcec.RecoverCompact(btcec.S256(), sig, digest[:])

	if err != nil {
		return false
	}

	pub, err := GetGraphenePublicKey(p)
	if err != nil {
		return false
	}

	return pub == publicKey
}

func GetPublicKey(msg interface{}, signature string) (string, error) {
	bData := new(bytes.Buffer)
	json.NewEncoder(bData).Encode(msg)
	digest := sha256.Sum256(bData.Bytes())

	sig, err := hex.DecodeString(signature)
	p, _, err := btcec.RecoverCompact(btcec.S256(), sig, digest[:])

	if err != nil {
		return "", nil
	}

	pub, err := GetGraphenePublicKey(p)
	if err != nil {
		return "", nil
	}

	return pub, nil
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
