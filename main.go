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
		server.DB.Order("id desc").Limit(limit).Find(&blocks)

		// Create responses
		var response []map[string]interface{}
		for _, block := range blocks {
			item := map[string]interface{}{
				"block_num":   block.Id,
				"block_hash":  block.Hash,
				"block_time":  block.Timestamp,
				"parent_hash": block.ParentHash,
			}
			response = append(response, item)
		}

		// Return the blocks as JSON.
		ctx.JSON(200, gin.H{
			"blocks": response,
		})
	})

	server.GIN.GET("/blocks/:id", func(ctx *gin.Context) {
		// Get the block id parameter.
		id := ctx.Param("id")

		// Query the database for the block with the specified id.
		var block model.Block
		server.DB.Where("id = ?", id).First(&block)

		// Query the database for the transactions with the specified id.
		var transactions []model.Transaction
		server.DB.Where("block_id = ?", id).Find(&transactions)

		// Create responses
		var txHashes []string
		for _, tx := range transactions {
			txHashes = append(txHashes, tx.Hash)
		}

		response := map[string]interface{}{
			"block_num":    block.Id,
			"block_hash":   block.Hash,
			"block_time":   block.Timestamp,
			"parent_hash":  block.ParentHash,
			"transactions": txHashes,
		}

		// Return the block as JSON.
		ctx.JSON(200, gin.H{
			"block": response,
		})
	})
	server.GIN.GET("/transaction/:txHash", func(ctx *gin.Context) {
		// Get the txHash parameter.
		txHash := ctx.Param("txHash")

		// Query the database for the transaction with the specified txHash.
		var transaction model.Transaction
		server.DB.Where("hash = ?", txHash).First(&transaction)

		// Create responses
		response := map[string]interface{}{
			"tx_hash":    transaction.Hash,
			"from":       transaction.From,
			"to":         transaction.To,
			"nonce":      transaction.Nonce,
			"data":       transaction.Data,
			"value":      transaction.Value,
			"event_logs": "", // TO-DO: extract event log from data
		}

		// Return the transaction as JSON.
		ctx.JSON(200, gin.H{
			"transaction": response,
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
