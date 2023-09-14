package handler

import (
	"backend/internal/model"
	"backend/internal/pkg/request"
	"backend/internal/pkg/response"
	"backend/internal/service"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type MusicHandler interface {
	Create(ctx *gin.Context)
	List(ctx *gin.Context)
	GetDetail(ctx *gin.Context)
	Upload(ctx *gin.Context)
	Recognize(ctx *gin.Context)
	CreateNFT(ctx *gin.Context)
	GetNFTs(ctx *gin.Context)
	GetNFTDetail(ctx *gin.Context)
}

func NewMusicHandler(handler *Handler, musicService service.MusicService) MusicHandler {
	return &musicHandler{
		Handler:      handler,
		musicService: musicService,
	}
}

type musicHandler struct {
	*Handler
	musicService service.MusicService
}

func (h *musicHandler) Create(ctx *gin.Context) {
	var musicInfo request.Music

	if err := ctx.ShouldBind(&musicInfo); err != nil {
		response.HandleError(ctx, http.StatusInternalServerError, response.ErrInternalServerError, err.Error())
		return
	}

	err := h.musicService.Create(&musicInfo)

	if err != nil {
		response.HandleError(ctx, http.StatusInternalServerError, response.ErrInternalServerError, err.Error())
		return
	}

	response.HandleSuccess(ctx, nil)
}

func (h *musicHandler) List(ctx *gin.Context) {
	list, err := h.musicService.List()
	if err != nil {
		response.HandleError(ctx, http.StatusInternalServerError, response.ErrInternalServerError, err.Error())
		return
	}
	response.HandleSuccess(ctx, list)
}

func (h *musicHandler) GetDetail(ctx *gin.Context) {
	id := ctx.Param("id")
	singleId, err := strconv.Atoi(id)
	if err != nil {
		response.HandleError(ctx, http.StatusInternalServerError, response.ErrInternalServerError, err.Error())
		return
	}
	detail, err := h.musicService.GetDetail(uint(singleId))
	if err != nil {
		response.HandleError(ctx, http.StatusInternalServerError, response.ErrInternalServerError, err.Error())
		return
	}
	response.HandleSuccess(ctx, detail)
}

func (h *musicHandler) Upload(ctx *gin.Context) {

	file, _ := ctx.FormFile("file")

	src, _ := file.Open()
	defer src.Close()

	fileType := ctx.PostForm("type")
	fileName := ctx.PostForm("fileName")

	// 如果是track先临时分离
	if fileType == "track" {
		targetFilePath := "/home/ubuntu/video/" + fileName
		baseName := filepath.Base(targetFilePath)
		fileNameWithoutExt := baseName[:len(baseName)-len(filepath.Ext(baseName))]
		if err := ctx.SaveUploadedFile(file, targetFilePath); err != nil {
			response.HandleError(ctx, http.StatusInternalServerError, response.ErrInternalServerError, err.Error())
			return
		}

		// 分离MP3
		cmd := fmt.Sprintf("ffmpeg -i %v -f wav -y -vn /home/ubuntu/audio/in/%v.wav", targetFilePath, fileNameWithoutExt)
		fmt.Println(cmd)
		command := exec.Command("bash", "-c", cmd)
		_, err := command.CombinedOutput()
		if err != nil {
			response.HandleError(ctx, http.StatusInternalServerError, response.ErrInternalServerError, err.Error())
			return
		}
		filepath2 := fmt.Sprintf("/home/ubuntu/audio/in/%v.wav", fileNameWithoutExt)
		fmt.Println(filepath2)
		file2, err := os.Open(filepath2)
		defer file2.Close()
		res, err := h.musicService.Upload(file2, fileNameWithoutExt+".wav", "demo")

		if err != nil {
			response.HandleError(ctx, http.StatusInternalServerError, response.ErrInternalServerError, nil)
			return
		}
		response.HandleSuccess(ctx, res)
	} else {
		res, err := h.musicService.Upload(src, fileName, fileType)

		if err != nil {
			response.HandleError(ctx, http.StatusInternalServerError, response.ErrInternalServerError, nil)
			return
		}

		response.HandleSuccess(ctx, res)
	}
}

func (h *musicHandler) Recognize(ctx *gin.Context) {

	// 保存文件到服务器
	file, _ := ctx.FormFile("file")
	fileName := ctx.PostForm("fileName")
	targetFilePath := "/home/ubuntu/audio/re/in/" + fileName
	baseName := filepath.Base(targetFilePath)
	fileNameWithoutExt := baseName[:len(baseName)-len(filepath.Ext(baseName))]
	if err := ctx.SaveUploadedFile(file, targetFilePath); err != nil {
		response.HandleError(ctx, http.StatusInternalServerError, response.ErrInternalServerError, err.Error())
		return
	}

	// 如果不是mp3或wav需要转码
	extension := filepath.Ext(targetFilePath)
	extension = extension[1:]

	if extension != "mp3" {
		cmd := fmt.Sprintf("ffmpeg -i %v  -acodec pcm_s16le -ac 2 -ar 44100 /home/ubuntu/audio/re/out/%v.wav", targetFilePath, fileNameWithoutExt)
		fmt.Println(cmd)
		command := exec.Command("bash", "-c", cmd)
		_, err := command.CombinedOutput()
		if err != nil {
			response.HandleError(ctx, http.StatusInternalServerError, response.ErrInternalServerError, err.Error())
			return
		}
	}

	// 水印识别
	cmd := "audiowmark get " + "/home/ubuntu/audio/re/out/" + fileNameWithoutExt + ".mp3"
	fmt.Println(cmd)
	command := exec.Command("bash", "-c", cmd)
	output, err := command.CombinedOutput()
	outputStr := string(output)
	fmt.Println(outputStr)
	if err != nil || outputStr == "" {
		response.HandleError(ctx, http.StatusInternalServerError, response.ErrInternalServerError, "fail")
		return
	}

	lines := strings.Split(outputStr, "\n")

	if len(lines) > 0 {
		firstLine := lines[0]

		// 正则表达式模式，匹配第一行中的哈希值
		patternRegex := regexp.MustCompile(`pattern\s+\S+\s+([0-9a-fA-F]{32})\s+\S+`)

		// 使用正则表达式匹配
		match := patternRegex.FindStringSubmatch(firstLine)
		fmt.Print(match)

		if match != nil && len(match) == 2 {
			hash := match[1]
			id, err := h.musicService.GetNFTByHashId(hash)
			if err != nil {
				response.HandleError(ctx, http.StatusInternalServerError, response.ErrInternalServerError, "fail")
				return
			}
			response.HandleSuccess(ctx, id)
		} else {
			response.HandleError(ctx, http.StatusInternalServerError, response.ErrInternalServerError, "fail")
			return
		}
	}
}

func (h *musicHandler) CreateNFT(ctx *gin.Context) {

	id := ctx.Param("id")
	singleId, _ := strconv.Atoi(id)

	// 获取音乐的详情
	detail, _ := h.musicService.GetDetail(uint(singleId))

	// 获取文件信息
	file, err := h.musicService.GetFileDetail(detail.DemoId)

	// 根据信息生成md5
	str := fmt.Sprintf("%v%v", detail.Name, detail.UserId)
	data := []byte(str)
	hash := md5.Sum(data)
	md5Hash := hex.EncodeToString(hash[:])

	// 添加水印并上传
	cmd1 := fmt.Sprintf("audiowmark add /home/ubuntu/audio/in/%v /home/ubuntu/audio/out/%v %v", file.FileName, file.FileName, md5Hash)
	fmt.Println(cmd1)
	command1 := exec.Command("bash", "-c", cmd1)
	_, err = command1.CombinedOutput()
	if err != nil {
		response.HandleError(ctx, http.StatusInternalServerError, response.ErrInternalServerError, err.Error())
		return
	}

	fileBody, _ := os.Open("/home/ubuntu/audio/out/" + file.FileName)
	fileNameF := h.musicService.AddSuffixBeforeExtension(file.FileName, "_f")
	defer fileBody.Close()
	res, _ := h.musicService.Upload(fileBody, fileNameF, "Single")

	// 调用服务端node mint
	_ = os.Chdir("/home/ubuntu/nft")
	cmd2 := fmt.Sprintf("npx ts-node mint.ts %v", res.FileUrl)
	fmt.Println(cmd2)
	command2 := exec.Command("bash", "-c", cmd2)
	output, err := command2.CombinedOutput()
	outputStr := string(output)

	// 获取交易信息
	re := regexp.MustCompile(`blockHash:(0x[0-9a-fA-F]+) blockNumber:(\d+) transactionHash:(0x[0-9a-fA-F]+)`)
	matches := re.FindStringSubmatch(outputStr)
	if len(matches) != 4 {
		response.HandleError(ctx, http.StatusInternalServerError, response.ErrInternalServerError, "mint失败")
		return
	}
	blockHash := matches[1]
	blockNumber := matches[2]
	blockNumber2, _ := strconv.Atoi(blockNumber)
	transactionHash := matches[3]

	// 创建数据记录
	nft := model.Nfts{
		HashId:          md5Hash,
		Type:            "Single",
		FileId:          res.FileId,
		RefId:           uint(singleId),
		OwnerId:         1,
		BlockHash:       blockHash,
		TransactionHash: transactionHash,
		BlockNumber:     uint(blockNumber2),
	}

	res2, err := h.musicService.CreateNFT(&nft)

	if err != nil {
		response.HandleError(ctx, http.StatusInternalServerError, response.ErrInternalServerError, err.Error())
		return
	}
	response.HandleSuccess(ctx, res2)
}

func (h *musicHandler) GetNFTs(ctx *gin.Context) {
	id := ctx.Param("id")
	singleId, _ := strconv.Atoi(id)

	list, err := h.musicService.GetNFTs(uint(singleId))
	if err != nil {
		response.HandleError(ctx, http.StatusInternalServerError, response.ErrInternalServerError, err.Error())
		return
	}

	response.HandleSuccess(ctx, list)
}

func (h *musicHandler) GetNFTDetail(ctx *gin.Context) {
	id := ctx.Param("id")
	singleId, _ := strconv.Atoi(id)

	detail, err := h.musicService.GetNFTDetail(uint(singleId))

	if err != nil {
		response.HandleError(ctx, http.StatusInternalServerError, response.ErrInternalServerError, err.Error())
		return
	}

	response.HandleSuccess(ctx, detail)
}
