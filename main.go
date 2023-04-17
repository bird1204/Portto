package main

import (
	"portto/initializer"
	"portto/model"
	"portto/scheduler"
	"strconv"

	"github.com/gin-gonic/gin"
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

	server.InitializeGin()
	// Define the API endpoints.
	server.GIN.GET("/blocks", func(ctx *gin.Context) {
		// Get the limit parameter.
		limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

		// Query the database for the lastest n blocks.
		var blocks []model.Block
		server.DB.Order("block_num desc").Limit(limit).Find(&blocks)

		// Return the blocks as JSON.
		ctx.JSON(200, gin.H{
			"blocks": blocks,
		})
	})
	server.GIN.GET("/blocks/:id", func(ctx *gin.Context) {
		// Get the block id parameter.
		id := ctx.Param("id")

		// Query the database for the block with the specified id.
		var block model.Block
		server.DB.Where("id = ?", id).First(&block)

		// Return the block as JSON.
		ctx.JSON(200, gin.H{
			"block": block,
		})
	})
	server.GIN.GET("/transaction/:txHash", func(ctx *gin.Context) {
		// Get the txHash parameter.
		txHash := ctx.Param("txHash")

		// Query the database for the transaction with the specified txHash.
		var transaction model.Transaction
		server.DB.Where("hash = ?", txHash).First(&transaction)

		// Return the transaction as JSON.
		ctx.JSON(200, gin.H{
			"transaction": transaction,
		})
	})

	go scheduler.ScanBlocks(
		17067895,                  // start block
		server.DB,                 // DB
		"http://10.0.130.61:8545", // node url
	)

	// Run the server.
	server.GIN.Run(":8080")
}
