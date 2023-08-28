package app

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"user-segmentation-service/internal/models"
	"user-segmentation-service/mocks"
)

func TestSegmentHandlers(t *testing.T) {
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
			name:    "Create Segment Success",
			handler: a.createSegmentHandler,
			requestBody: models.Segment{
				Slug: "AVITO_SALE_10",
				ExpirationDate: func() time.Time {
					t, err := time.Parse(time.RFC3339, "2023-12-31T23:59:59Z")
					if err != nil {
						panic("Невозможно преобразовать строку в time.Time")
					}
					return t
				}(),
				RandomPercentage: 0.0,
			},
			mockSetup: func() {
				mockDB.EXPECT().CreateSegment("AVITO_SALE_10", gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: map[string]interface{}{
				"message": "Segment and user assignments created successfully",
			},
		},
		{
			name:    "Create Segment Error (invalid percentage)",
			handler: a.createSegmentHandler,
			requestBody: models.Segment{
				Slug: "AVITO_SALE_10",
				ExpirationDate: func() time.Time {
					t, err := time.Parse(time.RFC3339, "2023-12-31T23:59:59Z")
					if err != nil {
						panic("Невозможно преобразовать строку в time.Time")
					}
					return t
				}(),
				RandomPercentage: 110.0,
			},
			mockSetup:    func() {},
			expectedCode: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "RandomPercentage should be between 0 and 100",
			},
		},
		{
			name:    "Delete Segment Success",
			handler: a.deleteSegmentHandler,
			requestBody: models.Segment{
				Slug: "AVITO_SALE_10",
			},
			mockSetup: func() {
				mockDB.EXPECT().DeleteSegment("AVITO_SALE_10").Return(1, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: map[string]interface{}{
				"message":    "Segment deleted successfully",
				"segment_id": float64(1),
			},
		},
		{
			name:    "Delete Segment Error (segment does not exist)",
			handler: a.deleteSegmentHandler,
			requestBody: models.Segment{
				Slug: "AVITO_SALE_666",
			},
			mockSetup: func() {
				mockDB.EXPECT().DeleteSegment(
					"AVITO_SALE_666").Return(
					0, errors.New("segment with slug 'AVITO_SALE_666' does not exist"))
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "segment with slug 'AVITO_SALE_666' does not exist",
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
