package main

import (
	"testing"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/quantadex/distributed_quanta_bridge/registrar/Forwarder"
	"io/ioutil"
	"github.com/ethereum/go-ethereum/accounts/keystore"
)

func TestNode(t *testing.T) {
	key, _ := crypto.GenerateKey()
	auth := bind.NewKeyedTransactor(key)

	// simulate users creating the contract
	userKey, _ := crypto.GenerateKey()
	userAuth := bind.NewKeyedTransactor(userKey)

	sim := backends.NewSimulatedBackend(core.GenesisAlloc{
		userAuth.From : { Balance: big.NewInt(10000000000)} }, 500000)


	address, err := Forwarder.SubmitContract(sim, userKey, auth.From, "QDDD")
	if err != nil {
		t.Error(err)
	}

	println(address)
}

func TestKeystore(t *testing.T) {
	keyjson, err := ioutil.ReadFile("keystore/key--7cd737655dff6f95d55b711975d2a4ace32d256e")
	if err != nil {
		t.Fatal(err)
	}

	key, err := keystore.DecryptKey(keyjson, "test123")
	if err != nil {
		t.Fatalf("test: json key failed to decrypt: %v", err)
	}

	println(key.Address.Hex())
}