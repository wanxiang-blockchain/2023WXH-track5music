package model

import "gorm.io/gorm"

type Nfts struct {
	gorm.Model
	HashId          string `json:"hashId" gorm:"hash_id"`
	Type            string `json:"type" gorm:"type"`
	FileId          uint   `json:"fileId" gorm:"file_id"`
	RefId           uint   `json:"refId" gorm:"ref_id"`
	OwnerId         uint   `json:"ownerId" gorm:"owner_id"`
	BlockHash       string `json:"blockHash" gorm:"block_hash"`
	TransactionHash string `json:"transactionHash" gorm:"transaction_hash"`
	BlockNumber     uint   `json:"blockNumber" gorm:"block_number"`
}

func (s *Nfts) TableName() string {
	return "nfts"
}
