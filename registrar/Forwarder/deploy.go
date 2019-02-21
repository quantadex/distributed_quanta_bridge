package Forwarder

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin/contracts"
	"log"
	"time"
)

// ForwarderBin is the compiled bytecode used for deploying new contracts.
const ForwarderBin = `0x608060405234801561001057600080fd5b5060405161067338038061067383398101604052805160208083015160008054600160a060020a031916600160a060020a038516179055909201805191929091610060916001919084019061011b565b507fa49a9b1337d8427ee784aeaded38ac25b248da00282d53353ef0e2dfb664504a82826040518083600160a060020a0316600160a060020a0316815260200180602001828103825283818151815260200191508051906020019080838360005b838110156100d95781810151838201526020016100c1565b50505050905090810190601f1680156101065780820380516001836020036101000a031916815260200191505b50935050505060405180910390a150506101b6565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061015c57805160ff1916838001178555610189565b82800160010185558215610189579182015b8281111561018957825182559160200191906001019061016e565b50610195929150610199565b5090565b6101b391905b80821115610195576000815560010161019f565b90565b6104ae806101c56000396000f3006080604052600436106100615763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416633c8410a281146100e15780633ef133671461016b5780636b9f96ea1461019b578063ca325469146101b0575b60408051348152905133917f5bac0d4f99f71df67fa7cebba0369126a7cb2b183bcb02b8393dbf5185ba77b6919081900360200190a26000805460405173ffffffffffffffffffffffffffffffffffffffff909116913480156108fc02929091818181858888f193505050501580156100de573d6000803e3d6000fd5b50005b3480156100ed57600080fd5b506100f66101ee565b6040805160208082528351818301528351919283929083019185019080838360005b83811015610130578181015183820152602001610118565b50505050905090810190601f16801561015d5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561017757600080fd5b5061019973ffffffffffffffffffffffffffffffffffffffff6004351661027b565b005b3480156101a757600080fd5b506101996103e4565b3480156101bc57600080fd5b506101c5610466565b6040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b60018054604080516020600284861615610100026000190190941693909304601f810184900484028201840190925281815292918301828280156102735780601f1061024857610100808354040283529160200191610273565b820191906000526020600020905b81548152906001019060200180831161025657829003601f168201915b505050505081565b604080517f70a082310000000000000000000000000000000000000000000000000000000081523060048201819052915183929160009173ffffffffffffffffffffffffffffffffffffffff8516916370a0823191602480830192602092919082900301818787803b1580156102f057600080fd5b505af1158015610304573d6000803e3d6000fd5b505050506040513d602081101561031a57600080fd5b5051905080151561032a576103de565b60008054604080517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff92831660048201526024810185905290519186169263a9059cbb926044808401936020939083900390910190829087803b1580156103a757600080fd5b505af11580156103bb573d6000803e3d6000fd5b505050506040513d60208110156103d157600080fd5b505115156103de57600080fd5b50505050565b6040805130318152905133917fa98efcd54f1f2ae5457ba3c68d7cf8974003a2bfce00f526f5624264a87bc0ea919081900360200190a26000805460405173ffffffffffffffffffffffffffffffffffffffff90911691303180156108fc02929091818181858888f19350505050158015610463573d6000803e3d6000fd5b50565b60005473ffffffffffffffffffffffffffffffffffffffff16815600a165627a7a72305820b50997e561bdb9e5bb99527e0995b254774fd7fcc650c893beeb1510ceff1cb80029`
const ForwarderBinV2 = `608060405234801561001057600080fd5b506040516104d13803806104d183398101604052805160208083015160008054600160a060020a031916600160a060020a038516179055909201805191929091610060916001919084019061011b565b507fa49a9b1337d8427ee784aeaded38ac25b248da00282d53353ef0e2dfb664504a82826040518083600160a060020a0316600160a060020a0316815260200180602001828103825283818151815260200191508051906020019080838360005b838110156100d95781810151838201526020016100c1565b50505050905090810190601f1680156101065780820380516001836020036101000a031916815260200191505b50935050505060405180910390a150506101b6565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061015c57805160ff1916838001178555610189565b82800160010185558215610189579182015b8281111561018957825182559160200191906001019061016e565b50610195929150610199565b5090565b6101b391905b80821115610195576000815560010161019f565b90565b61030c806101c56000396000f3006080604052600436106100565763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416633c8410a281146100d65780636b9f96ea14610160578063ca32546914610177575b60408051348152905133917f5bac0d4f99f71df67fa7cebba0369126a7cb2b183bcb02b8393dbf5185ba77b6919081900360200190a26000805460405173ffffffffffffffffffffffffffffffffffffffff909116913480156108fc02929091818181858888f193505050501580156100d3573d6000803e3d6000fd5b50005b3480156100e257600080fd5b506100eb6101b5565b6040805160208082528351818301528351919283929083019185019080838360005b8381101561012557818101518382015260200161010d565b50505050905090810190601f1680156101525780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561016c57600080fd5b50610175610242565b005b34801561018357600080fd5b5061018c6102c4565b6040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b60018054604080516020600284861615610100026000190190941693909304601f8101849004840282018401909252818152929183018282801561023a5780601f1061020f5761010080835404028352916020019161023a565b820191906000526020600020905b81548152906001019060200180831161021d57829003601f168201915b505050505081565b6040805130318152905133917fa98efcd54f1f2ae5457ba3c68d7cf8974003a2bfce00f526f5624264a87bc0ea919081900360200190a26000805460405173ffffffffffffffffffffffffffffffffffffffff90911691303180156108fc02929091818181858888f193505050501580156102c1573d6000803e3d6000fd5b50565b60005473ffffffffffffffffffffffffffffffffffffffff16815600a165627a7a72305820f03161a8010bca92450c47c8afca14692b54f0d917a1bcbb6d36a9147abf885c0029`

// linux
// Truffle v4.1.14 (core: 4.1.14)
// Solidity v0.4.24 (solc-js)
const ForwarderBinV3 = `608060405234801561001057600080fd5b506040516104d13803806104d183398101604052805160208083015160008054600160a060020a031916600160a060020a038516179055909201805191929091610060916001919084019061011b565b507fa49a9b1337d8427ee784aeaded38ac25b248da00282d53353ef0e2dfb664504a82826040518083600160a060020a0316600160a060020a0316815260200180602001828103825283818151815260200191508051906020019080838360005b838110156100d95781810151838201526020016100c1565b50505050905090810190601f1680156101065780820380516001836020036101000a031916815260200191505b50935050505060405180910390a150506101b6565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061015c57805160ff1916838001178555610189565b82800160010185558215610189579182015b8281111561018957825182559160200191906001019061016e565b50610195929150610199565b5090565b6101b391905b80821115610195576000815560010161019f565b90565b61030c806101c56000396000f3006080604052600436106100565763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416633c8410a281146100d65780636b9f96ea14610160578063ca32546914610177575b60408051348152905133917f5bac0d4f99f71df67fa7cebba0369126a7cb2b183bcb02b8393dbf5185ba77b6919081900360200190a26000805460405173ffffffffffffffffffffffffffffffffffffffff909116913480156108fc02929091818181858888f193505050501580156100d3573d6000803e3d6000fd5b50005b3480156100e257600080fd5b506100eb6101b5565b6040805160208082528351818301528351919283929083019185019080838360005b8381101561012557818101518382015260200161010d565b50505050905090810190601f1680156101525780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561016c57600080fd5b50610175610242565b005b34801561018357600080fd5b5061018c6102c4565b6040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b60018054604080516020600284861615610100026000190190941693909304601f8101849004840282018401909252818152929183018282801561023a5780601f1061020f5761010080835404028352916020019161023a565b820191906000526020600020905b81548152906001019060200180831161021d57829003601f168201915b505050505081565b6040805130318152905133917fa98efcd54f1f2ae5457ba3c68d7cf8974003a2bfce00f526f5624264a87bc0ea919081900360200190a26000805460405173ffffffffffffffffffffffffffffffffffffffff90911691303180156108fc02929091818181858888f193505050501580156102c1573d6000803e3d6000fd5b50565b60005473ffffffffffffffffffffffffffffffffffffffff16815600a165627a7a72305820cfe447988d6e135b2c8ed5bf64775898da12d6279b81b22eb054b57c820833690029`
const ForwarderBinV4 = `608060405234801561001057600080fd5b5060405161067338038061067383398101604052805160208083015160008054600160a060020a031916600160a060020a038516179055909201805191929091610060916001919084019061011b565b507fa49a9b1337d8427ee784aeaded38ac25b248da00282d53353ef0e2dfb664504a82826040518083600160a060020a0316600160a060020a0316815260200180602001828103825283818151815260200191508051906020019080838360005b838110156100d95781810151838201526020016100c1565b50505050905090810190601f1680156101065780820380516001836020036101000a031916815260200191505b50935050505060405180910390a150506101b6565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061015c57805160ff1916838001178555610189565b82800160010185558215610189579182015b8281111561018957825182559160200191906001019061016e565b50610195929150610199565b5090565b6101b391905b80821115610195576000815560010161019f565b90565b6104ae806101c56000396000f3006080604052600436106100615763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416633c8410a281146100e15780633ef133671461016b5780636b9f96ea1461019b578063ca325469146101b0575b60408051348152905133917f5bac0d4f99f71df67fa7cebba0369126a7cb2b183bcb02b8393dbf5185ba77b6919081900360200190a26000805460405173ffffffffffffffffffffffffffffffffffffffff909116913480156108fc02929091818181858888f193505050501580156100de573d6000803e3d6000fd5b50005b3480156100ed57600080fd5b506100f66101ee565b6040805160208082528351818301528351919283929083019185019080838360005b83811015610130578181015183820152602001610118565b50505050905090810190601f16801561015d5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561017757600080fd5b5061019973ffffffffffffffffffffffffffffffffffffffff6004351661027b565b005b3480156101a757600080fd5b506101996103e4565b3480156101bc57600080fd5b506101c5610466565b6040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b60018054604080516020600284861615610100026000190190941693909304601f810184900484028201840190925281815292918301828280156102735780601f1061024857610100808354040283529160200191610273565b820191906000526020600020905b81548152906001019060200180831161025657829003601f168201915b505050505081565b604080517f70a082310000000000000000000000000000000000000000000000000000000081523060048201819052915183929160009173ffffffffffffffffffffffffffffffffffffffff8516916370a0823191602480830192602092919082900301818787803b1580156102f057600080fd5b505af1158015610304573d6000803e3d6000fd5b505050506040513d602081101561031a57600080fd5b5051905080151561032a576103de565b60008054604080517fa9059cbb00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff92831660048201526024810185905290519186169263a9059cbb926044808401936020939083900390910190829087803b1580156103a757600080fd5b505af11580156103bb573d6000803e3d6000fd5b505050506040513d60208110156103d157600080fd5b505115156103de57600080fd5b50505050565b6040805130318152905133917fa98efcd54f1f2ae5457ba3c68d7cf8974003a2bfce00f526f5624264a87bc0ea919081900360200190a26000805460405173ffffffffffffffffffffffffffffffffffffffff90911691303180156108fc02929091818181858888f19350505050158015610463573d6000803e3d6000fd5b50565b60005473ffffffffffffffffffffffffffffffffffffffff16815600a165627a7a7230582061de18e7775520f21638941b7eec74f2ecdf7fc3d481b5a131c0ac336fecd5780029`

/**
 * SubmitContract submits the forwarding contract to the blockchain and returns
 * the smart contract address
 */
func SubmitContract(conn bind.ContractBackend, ownerKey *ecdsa.PrivateKey, trustAddress common.Address, quantaAddr string) (string, error) {
	auth := bind.NewKeyedTransactor(ownerKey)

	// Deploy a new awesome contract for the binding demo
	address, tx, forwarder, err := contracts.DeployQuantaForwarder(auth, conn, trustAddress, quantaAddr)
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
