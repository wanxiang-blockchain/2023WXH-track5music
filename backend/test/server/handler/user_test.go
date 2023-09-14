package handler

import (
	"backend/internal/handler"
	"backend/internal/pkg/middleware"
	"backend/internal/pkg/request"
	jwt2 "backend/pkg/jwt"
	"backend/test/mocks/service"
	"bytes"
	"encoding/json"
	"fmt"

	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"backend/internal/model"
	"backend/pkg/config"
	"backend/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var (
	userId = "yhs6HesfgF"

	token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOiJ5aHM2SGVzZmdGIiwiZXhwIjoxNjkzOTE0ODgwLCJuYmYiOjE2ODYxMzg4ODAsImlhdCI6MTY4NjEzODg4MH0.NnFrZFgc_333a9PXqaoongmIDksNvQoHzgM_IhJM4MQ"
)
var logger *log.Logger
var hdl *handler.Handler
var jwt *jwt2.JWT
var router *gin.Engine

func TestMain(m *testing.M) {
	fmt.Println("begin")
	err := os.Setenv("APP_CONF", "../../../config/local.yml")
	if err != nil {
		fmt.Println("Setenv error", err)
	}
	conf := config.NewConfig()

	logger = log.NewLog(conf)
	hdl = handler.NewHandler(logger)

	jwt = jwt2.NewJwt(conf)
	gin.SetMode(gin.TestMode)
	router = gin.Default()
	router.Use(
		middleware.CORSMiddleware(),
		middleware.ResponseLogMiddleware(logger),
		middleware.RequestLogMiddleware(logger),
		//middleware.SignMiddleware(log),
	)

	code := m.Run()
	fmt.Println("test end")

	os.Exit(code)
}

func TestUserHandler_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	params := request.RegisterRequest{
		Username: "xxx",
		Password: "123456",
		Email:    "xxx@gmail.com",
	}

	mockUserService := mock_service.NewMockUserService(ctrl)
	mockUserService.EXPECT().Register(gomock.Any(), &params).Return(nil)

	userHandler := handler.NewUserHandler(hdl, mockUserService)
	router.POST("/register", userHandler.Register)

	paramsJson, _ := json.Marshal(params)

	resp := performRequest(router, "POST", "/register", bytes.NewBuffer(paramsJson))

	assert.Equal(t, resp.Code, http.StatusOK)
	// Add assertions for the response body if needed
}

func TestUserHandler_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	params := request.LoginRequest{
		Username: "xxx",
		Password: "123456",
	}

	mockUserService := mock_service.NewMockUserService(ctrl)
	mockUserService.EXPECT().Login(gomock.Any(), &params).Return(token, nil)

	userHandler := handler.NewUserHandler(hdl, mockUserService)
	router.POST("/login", userHandler.Login)
	paramsJson, _ := json.Marshal(params)

	resp := performRequest(router, "POST", "/login", bytes.NewBuffer(paramsJson))

	assert.Equal(t, resp.Code, http.StatusOK)
	// Add assertions for the response body if needed
}

func TestUserHandler_GetProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mock_service.NewMockUserService(ctrl)
	mockUserService.EXPECT().GetProfile(gomock.Any(), userId).Return(&model.User{
		Id:       1,
		UserId:   userId,
		Username: "xxxxx",
		Nickname: "xxxxx",
		Password: "xxxxx",
		Email:    "xxxxx@gmail.com",
	}, nil)

	userHandler := handler.NewUserHandler(hdl, mockUserService)
	router.Use(middleware.NoStrictAuth(jwt, logger))
	router.GET("/user", userHandler.GetProfile)
	req, _ := http.NewRequest("GET", "/user", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, http.StatusOK)
	// Add assertions for the response body if needed
}

func TestUserHandler_UpdateProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	params := request.UpdateProfileRequest{
		Nickname: "alan",
		Email:    "alan@gmail.com",
		Avatar:   "xxx",
	}

	mockUserService := mock_service.NewMockUserService(ctrl)
	mockUserService.EXPECT().UpdateProfile(gomock.Any(), userId, &params).Return(nil)

	userHandler := handler.NewUserHandler(hdl, mockUserService)
	router.Use(middleware.StrictAuth(jwt, logger))
	router.PUT("/user", userHandler.UpdateProfile)
	paramsJson, _ := json.Marshal(params)

	req, _ := http.NewRequest("PUT", "/user", bytes.NewBuffer(paramsJson))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, resp.Code, http.StatusOK)
	// Add assertions for the response body if needed
}

func performRequest(r http.Handler, method, path string, body *bytes.Buffer) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, body)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	return resp
}
