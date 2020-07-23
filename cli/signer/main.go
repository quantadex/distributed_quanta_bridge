package main

import "github.com/quantadex/distributed_quanta_bridge/cli"

func main() {
	config, _, _, _, log, secrets := cli.Setup(true)
	cli.RunSigner(config, secrets, log, config.KmPort)
	print("Started signer...")
}
