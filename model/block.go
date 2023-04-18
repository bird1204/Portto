package model

import "gorm.io/gorm"

type Block struct {
	gorm.Model
	Id           uint64 `gorm:"uniqueIndex;not null"`
	Hash         string `gorm:"size:255;uniqueIndex;not null"`
	ParentHash   string `gorm:"not null"`
	Timestamp    uint64 `gorm:"not null"`
	IsStable     bool   `gorm:"not null"`
	Transactions []Transaction
}
