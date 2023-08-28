package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"user-segmentation-service/internal/models"
	"user-segmentation-service/mocks"
)

func TestUserHandlers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockInterface(ctrl)
	a := &App{db: mockDB}

	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		handler      gin.HandlerFunc
		requestBody  interface{}
		mockSetup    func()
		expectedCode int
		expectedBody map[string]interface{}
	}{
		{
			name:    "Create User Success",
			handler: a.createUserHandler,
			requestBody: models.User{
				Name: "John",
			},
			mockSetup: func() {
				mockDB.EXPECT().CreateUser("John").Return(int64(1), nil)
			},
			expectedCode: http.StatusCreated,
			expectedBody: map[string]interface{}{
				"message": "User created successfully",
				"user_id": float64(1),
			},
		},
		{
			name:    "Delete User Success",
			handler: a.deleteUserHandler,
			requestBody: models.DeleteUserRequest{
				UserId: 1,
			},
			mockSetup: func() {
				mockDB.EXPECT().DeleteUser(1).Return(1, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: map[string]interface{}{
				"message": "User deleted successfully",
				"user_id": float64(1),
			},
		},
		{
			name:    "Delete User Error",
			handler: a.deleteUserHandler,
			requestBody: models.DeleteUserRequest{
				UserId: 12,
			},
			mockSetup: func() {
				mockDB.EXPECT().DeleteUser(12).Return(0, errors.New("user with ID 12 does not exist"))
			},
			expectedCode: http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"error": "user with ID 12 does not exist",
			},
		},
		{
			name:    "Update User Segments Success",
			handler: a.updateUserSegmentsHandler,
			requestBody: models.UpdateSegmentsRequest{
				UserId: 1,
				Add: []models.Segment{
					{
						Slug: "AVITO_SALE_10",
						ExpirationDate: func() time.Time {
							t, err := time.Parse(time.RFC3339, "2023-12-31T23:59:59Z")
							if err != nil {
								panic("Невозможно преобразовать строку в time.Time")
							}
							return t
						}(),
					},
					{
						Slug: "AVITO_SALE_20",
						ExpirationDate: func() time.Time {
							t, err := time.Parse(time.RFC3339, "2023-12-31T23:59:59Z")
							if err != nil {
								panic("Невозможно преобразовать строку в time.Time")
							}
							return t
						}(),
					},
				},
				Remove: []string{},
			},
			mockSetup: func() {
				mockDB.EXPECT().UpdateUserSegments(1, gomock.Any(), gomock.Any()).Return(1, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: map[string]interface{}{
				"message": "User segments updated successfully",
				"user_id": float64(1),
			},
		},
		{
			name:    "Update User Segments Error (user does not exist)",
			handler: a.updateUserSegmentsHandler,
			requestBody: models.UpdateSegmentsRequest{
				UserId: 13,
				Add: []models.Segment{
					{
						Slug: "AVITO_SALE_10",
						ExpirationDate: func() time.Time {
							t, err := time.Parse(time.RFC3339, "2023-12-31T23:59:59Z")
							if err != nil {
								panic("Невозможно преобразовать строку в time.Time")
							}
							return t
						}(),
					},
					{
						Slug: "AVITO_SALE_20",
						ExpirationDate: func() time.Time {
							t, err := time.Parse(time.RFC3339, "2023-12-31T23:59:59Z")
							if err != nil {
								panic("Невозможно преобразовать строку в time.Time")
							}
							return t
						}(),
					},
				},
				Remove: []string{},
			},
			mockSetup: func() {
				mockDB.EXPECT().UpdateUserSegments(13, gomock.Any(), gomock.Any()).Return(0, errors.New("user with ID '13' does not exist"))
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "user with ID '13' does not exist",
			},
		},
		{
			name:    "Update User Segments Error (segment does not exist)",
			handler: a.updateUserSegmentsHandler,
			requestBody: models.UpdateSegmentsRequest{
				UserId: 1,
				Add: []models.Segment{
					{
						Slug: "AVITO_SALE_666",
						ExpirationDate: func() time.Time {
							t, err := time.Parse(time.RFC3339, "2023-12-31T23:59:59Z")
							if err != nil {
								panic("Невозможно преобразовать строку в time.Time")
							}
							return t
						}(),
					},
					{
						Slug: "AVITO_SALE_20",
						ExpirationDate: func() time.Time {
							t, err := time.Parse(time.RFC3339, "2023-12-31T23:59:59Z")
							if err != nil {
								panic("Невозможно преобразовать строку в time.Time")
							}
							return t
						}(),
					},
				},
				Remove: []string{},
			},
			mockSetup: func() {
				mockDB.EXPECT().UpdateUserSegments(1, gomock.Any(), gomock.Any()).Return(0, errors.New("segment with slug 'AVITO_SALE_120' does not exist"))
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "segment with slug 'AVITO_SALE_120' does not exist",
			},
		},
		{
			name:    "Get User Segment Success",
			handler: a.getUserSegmentsHandler,
			requestBody: models.UserSegmentsRequest{
				UserId: 1,
			},
			mockSetup: func() {
				mockDB.EXPECT().GetUserSegments(1).Return(1, []string{"AVITO_SALE_10", "AVITO_SALE_20"}, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: map[string]interface{}{
				"segments": []interface{}{
					"AVITO_SALE_10",
					"AVITO_SALE_20",
				},
				"user_id": float64(1),
			},
		},
		{
			name:    "Get User Segment Error (user does not exist)",
			handler: a.getUserSegmentsHandler,
			requestBody: models.UserSegmentsRequest{
				UserId: 13,
			},
			mockSetup: func() {
				mockDB.EXPECT().GetUserSegments(13).Return(0, nil, errors.New("user with ID '13' does not exist"))
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "user with ID '13' does not exist",
			},
		},
		{
			name:    "Get User Report Success",
			handler: a.getUserReportHandler,
			requestBody: models.ReportRequest{
				UserId:    1,
				YearMonth: "2023-08",
			},
			mockSetup: func() {
				mockDB.EXPECT().GetUserReport(1, "2023-08").Return("http://localhost:8080/user/report/user_1_report_2023-08.csv", nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: map[string]interface{}{
				"download_link": "http://localhost:8080/user/report/user_1_report_2023-08.csv",
				"message":       "Report generated successfully",
			},
		},
		{
			name:    "Get User Report Error (user does not exist)",
			handler: a.getUserReportHandler,
			requestBody: models.ReportRequest{
				UserId:    13,
				YearMonth: "2023-08",
			},
			mockSetup: func() {
				mockDB.EXPECT().GetUserReport(13, "2023-08").Return("", errors.New("user with ID '13' does not exist"))
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "user with ID '13' does not exist",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assertion := assert.New(t)
			if tc.mockSetup != nil {
				tc.mockSetup()
			}

			requestData, _ := json.Marshal(tc.requestBody)
			r := httptest.NewRequest("POST", "/", bytes.NewBuffer(requestData))
			w := httptest.NewRecorder()

			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = r

			tc.handler(ctx)

			assertion.Equal(tc.expectedCode, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assertion.NoError(err)
			assertion.Equal(tc.expectedBody, response)
		})
	}
}
