package service

import (
	"backend/internal/model"
	"backend/internal/pkg/request"
	"backend/internal/pkg/response"
	"backend/internal/repository"
	"io"
	"strings"
)

type MusicService interface {
	Create(info *request.Music) error
	Upload(body io.Reader, fileName, fileType string) (res *response.RepFile, err error)
	GetDetail(id uint) (detail *response.RepSingle, err error)
	List() (list []*response.RepSingle, err error)
	CreateNFT(nft *model.Nfts) (*model.Nfts, error)
	GetFileDetail(id uint) (detail *model.File, err error)
	GetNFTs(id uint) (list []*model.Nfts, err error)
	GetNFTDetail(id uint) (detail *response.RepNFTDetail, err error)
	GetNFTByHashId(hashId string) (id uint, err error)
	AddSuffixBeforeExtension(input string, suffix string) string
}

func NewMusicService(service *Service, musicRepository repository.MusicRepository) MusicService {
	return &musicService{
		Service:         service,
		musicRepository: musicRepository,
	}
}

type musicService struct {
	*Service
	musicRepository repository.MusicRepository
}

func (s *musicService) Create(info *request.Music) error {

	single := model.Single{
		Name:        info.Name,
		Intro:       info.Intro,
		Tag:         info.Tag,
		PublishDate: nil,
		UserId:      1,
		CoverId:     info.CoverId,
		DemoId:      info.DemoId,
	}

	for _, item := range info.Tracks {
		track := model.Track{
			SingleId:     0,
			Position:     item.Position,
			Label:        item.Label,
			ReleaseQty:   item.Quantity,
			ReleasePrice: item.Price,
			FileId:       item.FileId,
		}
		single.Tracks = append(single.Tracks, track)
	}

	err := s.musicRepository.CreateMusic(&single)
	return err
}

func (s *musicService) Upload(body io.Reader, fileName, fileType string) (res *response.RepFile, err error) {

	// 获得url
	url, err := s.s3.UploadFile(fileName, body)

	// 存入数据库返回文件id
	fileId, err := s.musicRepository.CreateFile(fileName, fileType, url)

	if err != nil {
		return nil, err
	}

	res = &response.RepFile{
		FileId:   fileId,
		FileUrl:  url,
		FileType: fileType,
	}

	return res, nil
}

func (s *musicService) GetFileDetail(id uint) (detail *model.File, err error) {
	detail, err = s.musicRepository.GetFile(id)
	return
}

func (s *musicService) List() (list []*response.RepSingle, err error) {
	list, err = s.musicRepository.GetSingleList()
	return
}

func (s *musicService) GetDetail(id uint) (detail *response.RepSingle, err error) {
	detail, err = s.musicRepository.GetSingleDetail(id)
	return
}

func (s *musicService) CreateNFT(nft *model.Nfts) (*model.Nfts, error) {
	detail, err := s.musicRepository.CreateNFT(nft)
	return detail, err
}

func (s *musicService) GetNFTs(id uint) (list []*model.Nfts, err error) {
	list, err = s.musicRepository.GetNFTs(id)
	return
}

func (s *musicService) GetNFTDetail(id uint) (detail *response.RepNFTDetail, err error) {
	detail, err = s.musicRepository.GetNFTDetail(id)
	return
}

func (s *musicService) GetNFTByHashId(hashId string) (id uint, err error) {
	id, err = s.musicRepository.GetNFTByHashId(hashId)
	return
}

func (s *musicService) AddSuffixBeforeExtension(input string, suffix string) string {
	// 寻找文件名中的最后一个点
	lastDotIndex := strings.LastIndex(input, ".")

	// 如果没有找到点，或者点在字符串的末尾，则直接添加后缀
	if lastDotIndex == -1 || lastDotIndex == len(input)-1 {
		return input + suffix
	}

	// 将后缀插入到点之前，然后返回新的字符串
	return input[:lastDotIndex] + suffix + input[lastDotIndex:]
}
