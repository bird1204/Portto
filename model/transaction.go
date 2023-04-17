package model

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model
	Hash     string `gorm:"size:255;uniqueIndex;not null"`
	BlockId  uint64 `gorm:"not null"`
	From     string `gorm:"not null"`
	To       string `gorm:"not null"`
	Nonce    uint64 `gorm:"not null"`
	Data     []byte `gorm:"not null"`
	Value    uint64 `gorm:"not null"`
	LogIndex uint64 `gorm:"not null"`
}
