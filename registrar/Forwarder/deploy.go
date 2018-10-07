package Forwarder

import (
	"log"
	"fmt"
	"time"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"crypto/ecdsa"
)

/**
 * SubmitContract submits the forwarding contract to the blockchain and returns
 * the smart contract address
 */
func SubmitContract(conn bind.ContractBackend, ownerKey *ecdsa.PrivateKey, trustAddress common.Address, quantaAddr string) (string, error) {
	auth := bind.NewKeyedTransactor(ownerKey)


	// Deploy a new awesome contract for the binding demo
	address, tx, forwarder, err := DeployForwarder(auth, conn, trustAddress, quantaAddr)
	if err != nil {
		log.Printf("Failed to deploy new token contract: %v", err)
		return "", err
	}

	fmt.Printf("Contract pending deploy: 0x%x\n", address)
	fmt.Printf("Transaction waiting to be mined: 0x%x\n\n", tx.Hash())

	time.Sleep(250 * time.Millisecond) // Allow it to be processed by the local node :P

	name, err := forwarder.QuantaAddress(&bind.CallOpts{Pending: true})
	if err != nil {
		log.Printf("Failed to retrieve pending name: %v", err)
		return "", err
	}
	fmt.Println("Pending name:", name)
	return tx.Hash().Hex(), nil
}