package scheduler

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"portto/model"
	"sync"

	"github.com/ethereum/go-ethereum/ethclient"
	"gorm.io/gorm"
)

func ScanBlocks(client *ethclient.Client, startBlock uint64, db *gorm.DB) error {
	// Get the latest block number.
	latestBlockNumber, err := client.BlockNumber(context.Background())
	if err != nil {
		return err
	}

	// If the start block is greater than the latest block, return an error.
	if startBlock > latestBlockNumber {
		return fmt.Errorf("start block %d is greater than the latest block %d", startBlock, latestBlockNumber)
	}

	// Start a Goroutine to scan blocks in parallel.
	var wg sync.WaitGroup
	wg.Add(int(latestBlockNumber - startBlock + 1))
	for i := startBlock; i <= latestBlockNumber; i++ {
		go func(blockNum uint64) {
			defer wg.Done()

			// Get the block from the Ethereum client.
			block, err := client.BlockByNumber(context.Background(), big.NewInt(int64(blockNum)))
			if err != nil {
				log.Printf("Error getting block %d: %s", blockNum, err.Error())
				return
			}

			// Save the block information to the database.
			err = db.Create(&model.Block{
				Id:         block.Number().Uint64(),
				Hash:       block.Hash().Hex(),
				ParentHash: block.ParentHash().Hex(),
				Timestamp:  block.Time(),
				IsStable:   false, // Assume the block is unstable initially.
			}).Error
			if err != nil {
				log.Printf("Error saving block %d to DB: %s", blockNum, err.Error())
				return
			}

			// Save the transactions in the block to the database.
			for _, tx := range block.Transactions() {
				err = db.Create(&model.Transaction{
					Hash:    tx.Hash().Hex(),
					BlockId: block.Number().Uint64(),
					// From:     tx.From().Hex(),
					To:    tx.To().Hex(),
					Nonce: tx.Nonce(),
					// Value: tx.Value().String(),
					Data: tx.Data(),
				}).Error
				if err != nil {
					log.Printf("Error saving transaction %s to DB: %s", tx.Hash().Hex(), err.Error())
				}
			}

			// Mark the block as stable if it's one of the latest 20 blocks.
			if latestBlockNumber-blockNum < 20 {
				err = db.Model(&model.Block{}).Where("ID = ?", block.Number().Uint64()).Update("is_stable", true).Error
				if err != nil {
					log.Printf("Error marking block %d as stable: %s", blockNum, err.Error())
				}
			}
		}(i)
	}

	// Wait for all the Goroutines to finish scanning.
	wg.Wait()

	return nil
}
