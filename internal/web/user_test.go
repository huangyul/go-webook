package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/huangyul/go-webook/internal/service"
	svcmock "github.com/huangyul/go-webook/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserHandler_Register(t *testing.T) {
	type RegisterReq struct {
		Email           string `json:"email" binding:"required"`
		Password        string `json:"password" binding:"required"`
		ConfirmPassword string `json:"confirm_password" binding:"required"`
	}
	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) (service.UserService, service.CodeService)
		request  RegisterReq
		wantCode int
		wantBody any
	}{
		{
			name: "注册成功",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmock.NewMockUserService(ctrl)
				userSvc.EXPECT().RegisterByEmail(gomock.Any(), "test@test.com", "123abc").Return(nil)
				return userSvc, nil
			},
			request: RegisterReq{
				Email:           "test@test.com",
				Password:        "123abc",
				ConfirmPassword: "123abc",
			},
			wantCode: http.StatusOK,
			wantBody: gin.H{"msg": "user successfully registered"},
		},
		{
			name: "没有传email",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmock.NewMockUserService(ctrl)
				return userSvc, nil
			},
			request: RegisterReq{
				Email:           "",
				Password:        "123abc",
				ConfirmPassword: "123abc",
			},
			wantCode: http.StatusBadRequest,
			wantBody: map[string]string{
				"error": "Key: 'RegisterReq.Email' Error:Field validation for 'Email' failed on the 'required' tag",
			},
		},
		{
			name: "非法邮箱",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmock.NewMockUserService(ctrl)
				return userSvc, nil
			},
			request: RegisterReq{
				Email:           "test",
				Password:        "123abc",
				ConfirmPassword: "123abc",
			},
			wantCode: http.StatusBadRequest,
			wantBody: map[string]string{
				"error": "email illegal",
			},
		},
		{
			name: "密码不符合要求",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmock.NewMockUserService(ctrl)
				return userSvc, nil
			},
			request: RegisterReq{
				Email:           "test@test.com",
				Password:        "123123",
				ConfirmPassword: "123123",
			},
			wantCode: http.StatusBadRequest,
			wantBody: map[string]string{
				"error": "password illegal",
			},
		},
		{
			name: "两次密码不相等",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmock.NewMockUserService(ctrl)
				return userSvc, nil
			},
			request: RegisterReq{
				Email:           "test@test.com",
				Password:        "123abc",
				ConfirmPassword: "123abcd",
			},
			wantCode: http.StatusBadRequest,
			wantBody: map[string]string{
				"error": "password not match",
			},
		},
		{
			name: "service层报错",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmock.NewMockUserService(ctrl)
				userSvc.EXPECT().RegisterByEmail(gomock.Any(), "test@test.com", "123abc").Return(errors.New("test error"))
				return userSvc, nil
			},
			request: RegisterReq{
				Email:           "test@test.com",
				Password:        "123abc",
				ConfirmPassword: "123abc",
			},
			wantCode: http.StatusInternalServerError,
			wantBody: map[string]string{
				"error": "test error",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			userSvc, codeSvc := testCase.mock(ctrl)
			hdl := NewUserHandler(userSvc, codeSvc, nil)
			server := gin.New()
			hdl.RegisterRoutes(server)

			data, err := json.Marshal(testCase.request)
			assert.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/user/register", bytes.NewReader(data))
			req.Header.Set("Content-Type", "application/json")
			assert.NoError(t, err)

			recorder := httptest.NewRecorder()
			server.ServeHTTP(recorder, req)
			assert.Equal(t, testCase.wantCode, recorder.Code)
			wb, _ := json.Marshal(testCase.wantBody)
			assert.Equal(t, string(wb), recorder.Body.String())
		})
	}
}
