package repository

import (
	"backend/internal/model"
	"backend/internal/pkg/response"
)

type MusicRepository interface {
	CreateMusic(single *model.Single) error
	CreateFile(fileName, fileType, uri string) (uint, error)
	GetFile(fileId uint) (detail *model.File, err error)
	GetSingleDetail(singleId uint) (detail *response.RepSingle, err error)
	GetSingleList() (list []*response.RepSingle, err error)
	CreateNFT(nft *model.Nfts) (*model.Nfts, error)
	GetNFTs(id uint) (list []*model.Nfts, err error)
	GetNFTDetail(id uint) (detail *response.RepNFTDetail, err error)
	GetNFTByHashId(hashId string) (id uint, err error)
}

func NewMusicRepository(repository *Repository) MusicRepository {
	return &musicRepository{
		Repository: repository,
	}
}

type musicRepository struct {
	*Repository
}

// CreateMusic 创建单曲信息
func (r *musicRepository) CreateMusic(info *model.Single) error {
	err := r.db.Create(info).Error
	return err
}

// CreateFile 创建文件
func (r *musicRepository) CreateFile(fileName, fileType, uri string) (uint, error) {

	file := model.File{
		Type:     fileType,
		Url:      uri,
		FileName: fileName,
	}

	err := r.db.Create(&file).Error
	if err != nil {
		return 0, err
	}

	return file.ID, nil
}

// GetFile getFile 获取文件信息
func (r *musicRepository) GetFile(fileId uint) (detail *model.File, err error) {
	err = r.db.Where("id = ?", fileId).First(&detail).Error
	return
}

// GetSingleDetail 获取单曲详情
func (r *musicRepository) GetSingleDetail(singleId uint) (detail *response.RepSingle, err error) {
	err = r.db.Where("id = ?", singleId).
		Preload("Cover").Preload("Demo").First(&detail).Error
	return
}

// GetSingleList 获取单曲列表
func (r *musicRepository) GetSingleList() (list []*response.RepSingle, err error) {
	err = r.db.Preload("Cover").Preload("Demo").Order("ID desc").Find(&list).Error
	return
}

func (r *musicRepository) CreateNFT(nft *model.Nfts) (*model.Nfts, error) {
	err := r.db.Create(nft).Error
	return nft, err
}

func (r *musicRepository) GetNFTs(id uint) (list []*model.Nfts, err error) {
	err = r.db.Limit(6).Where("ref_id = ? and type = 'Single'", id).Order("id desc").Find(&list).Error
	return
}

func (r *musicRepository) GetNFTDetail(id uint) (detail *response.RepNFTDetail, err error) {
	err = r.db.Preload("Single").Preload("File").Preload("Single.Cover").Where("id = ?", id).First(&detail).Error
	return
}

func (r *musicRepository) GetNFTByHashId(hashId string) (id uint, err error) {
	var nft model.Nfts
	err = r.db.Where("hash_id = ?", hashId).First(&nft).Error
	return nft.ID, err
}
