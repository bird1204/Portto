package main

import (
	"portto/initializer"
	"portto/scheduler"
)

func main() {
	server := initializer.Server{}
	server.InitializeDB(
		"mysql",     // driver
		"root",      // Username
		"rootroot",  // PWD
		"3306",      // Port
		"localhost", // Host
		"portto",    // DB name
	)

	// node url: https://polygon-mainnet.g.alchemy.com/v2/T4msPSgqhQshJuaZ5ZxZJIU2QrksAjQt / http://10.0.130.61:8545
	server.InitializeEthClient("https://polygon-mainnet.g.alchemy.com/v2/T4msPSgqhQshJuaZ5ZxZJIU2QrksAjQt")

	server.InitializeGin()
	server.CreateBlocksRoute()
	server.CreateTransactionsRoute()

	go scheduler.ScanBlocks(
		server.EthClient, // ETH client
		17074886,         // start block
		server.DB,        // DB
	)

	// Run the server.
	server.GIN.Run(":8080")
}
