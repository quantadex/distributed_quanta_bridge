package crypto

import (
	"github.com/tyler-smith/go-bip32"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/common"
	"github.com/haltingstate/secp256k1-go"
	"crypto/ecdsa"
)

func GenerateBip32Key(bip32Private string, index uint32) (common.Address, *ecdsa.PrivateKey, error) {
	privKey, err := bip32.B58Deserialize(bip32Private)
	key, err := privKey.NewChildKey(uint32(index))

	if err != nil {
		return common.HexToAddress("0x0"), nil, err
	}

	uncompressed := secp256k1.UncompressPubkey(key.PublicKey().Key)
	keccak := crypto.Keccak256(uncompressed[1:])
	addressPub := common.BytesToAddress(keccak[12:]) // Encode lower 160 bits/20 bytes

	pkey, err := crypto.ToECDSA(key.Key)

	return addressPub, pkey, err
	//hexPriv := hex.EncodeToString(crypto.FromECDSA(pkey))
}