package main

import (
	"portto/initializer"
	"portto/scheduler"
)

func main() {
	server := initializer.Server{}
	// NOTE: CHANGE TO YOUR OWNED DB
	server.InitializeDB(
		"mysql",     // driver
		"username",  // Username
		"password",  // PWD
		"3306",      // Port
		"localhost", // Host
		"portto",    // DB name
	)

	// NOTE: CHANGE TO YOUR OWNED NODE
	server.InitializeEthClient("http://10.0.130.61:8545")

	server.InitializeGin()
	server.CreateBlocksRoute()
	server.CreateTransactionsRoute()

	// NOTE: update start block to a certain block
	go scheduler.ScanBlocks(
		server.EthClient, // ETH client
		17074886,         // start block
		server.DB,        // DB
	)

	// Run the server.
	server.GIN.Run(":8080")
}
