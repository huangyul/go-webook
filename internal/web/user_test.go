package web

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/huangyul/go-blog/internal/service"
	mocksvc "github.com/huangyul/go-blog/internal/service/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserHandler_Signup(t *testing.T) {
	type Req struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirm_password"`
	}
	tests := []struct {
		name     string
		before   func(ctrl *gomock.Controller) (service.UserService, service.CodeService)
		req      Req
		wantCode int
		wantBody any
	}{
		{
			name: "success",
			before: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := mocksvc.NewMockUserService(ctrl)
				userSvc.EXPECT().Signup(gomock.Any(), "123123@qq.com", "123123").Return(nil)
				return userSvc, nil
			},
			req: Req{
				Email:           "123123@qq.com",
				Password:        "123123",
				ConfirmPassword: "123123",
			},
			wantCode: http.StatusOK,
			wantBody: struct{}{},
		},
		{
			name: "success",
			before: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				return nil, nil
			},
			req: Req{
				Password:        "123123",
				ConfirmPassword: "123123",
			},
			wantCode: http.StatusOK,
			wantBody: struct{}{},
		},
		{
			name: "success",
			before: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				return nil, nil
			},
			req: Req{
				Email:           "123123@qq.com",
				Password:        "123123",
				ConfirmPassword: "1231233",
			},
			wantCode: http.StatusOK,
			wantBody: struct{}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userSvc, codeSvc := tt.before(ctrl)

			server := gin.Default()
			uHdl := NewUserHandler(userSvc, codeSvc)
			uHdl.RegisterRoutes(server)

			data, err := json.Marshal(tt.req)
			assert.NoError(t, err)

			request := httptest.NewRequest(http.MethodPost, "/user/signup", bytes.NewReader(data))
			request.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()

			server.ServeHTTP(recorder, request)
			assert.Equal(t, http.StatusOK, recorder.Code)

		})
	}
}
