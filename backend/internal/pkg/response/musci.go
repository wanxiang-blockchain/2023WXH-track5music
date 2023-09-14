package response

import "backend/internal/model"

type RepSingle struct {
	model.Single
	Cover model.File `json:"Cover" gorm:"foreignKey:CoverId"`
	Demo  model.File `json:"Demo" gorm:"foreignKey:DemoId"`
}

type RepFile struct {
	FileId   uint   `json:"fileId"`
	FileUrl  string `json:"fileUrl"`
	FileType string `json:"fileType"`
}

type RepNFTDetail struct {
	model.Nfts
	Single RepSingle  `json:"single" gorm:"foreignKey:RefId"`
	File   model.File `json:"file" gorm:"foreignKey:FileId"`
}
