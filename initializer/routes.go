package initializer

import (
	"context"
	"portto/model"
	"portto/service"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/gin-gonic/gin"
)

func (s *Server) CreateBlocksRoute() {
	// Define the API endpoints.
	s.GIN.GET("/blocks", func(ctx *gin.Context) {
		// Get the limit parameter.
		limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

		// Query the database for the lastest n blocks.
		var blocks []model.Block
		s.DB.Order("id desc").Limit(limit).Find(&blocks)

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

	s.GIN.GET("/blocks/:id", func(ctx *gin.Context) {
		// Get the block id parameter.
		id := ctx.Param("id")

		// Query the database for the block with the specified id.
		var block model.Block
		s.DB.Where("id = ?", id).First(&block)

		// Query the database for the transactions with the specified id.
		var transactions []model.Transaction
		s.DB.Where("block_id = ?", id).Find(&transactions)

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
}

func (s *Server) CreateTransactionsRoute() {
	s.GIN.GET("/transaction/:txHash", func(ctx *gin.Context) {
		var errorMsg string

		// Get the txHash parameter.
		txHash := ctx.Param("txHash")

		// Query the database for the transaction with the specified txHash.
		var transaction model.Transaction
		s.DB.Where("hash = ?", txHash).First(&transaction)

		abiData, err := service.GetContractAbiFromFile(
			s.EthClient,
			common.HexToAddress(transaction.To), // 0xfb6916095ca1df60bb79ce92ce3ea74c37c5d359 / 0xfB666D8B64e619FAEbEcF3e1C383A0a87CB14a2e
		)
		if err != nil { // if can't get ABI from File
			errorMsg = "Error get ABI from file: " + err.Error()
		}

		// Extract event logs
		eventLogs := []types.Log{}

		receipt, err := s.EthClient.TransactionReceipt(context.Background(), common.HexToHash(transaction.Hash))
		if err != nil {
			errorMsg = errorMsg + "\n" + "Error Get Transaction Receipt: " + err.Error()
		}

		// Decode event logs
		if errorMsg == "" {
			for _, log := range receipt.Logs {
				if err := abiData.UnpackIntoInterface(&eventLogs, "Transfer", log.Data); err != nil {
					// Ignore decoding errors and move on to the next log
					continue
				}
			}
		}

		// Create responses
		response := map[string]interface{}{
			"tx_hash":       transaction.Hash,
			"from":          transaction.From,
			"to":            transaction.To,
			"nonce":         transaction.Nonce,
			"data":          transaction.Data,
			"value":         transaction.Value,
			"event_logs":    eventLogs,
			"error_message": errorMsg,
		}

		// Return the transaction as JSON.
		ctx.JSON(200, gin.H{
			"transaction": response,
		})
	})
}
