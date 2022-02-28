package conrtollers

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/pmokeev/chartographer/internal/services"
	"github.com/pmokeev/chartographer/internal/services/mocks"
	"github.com/stretchr/testify/assert"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

func createRequestString(parametersMap map[string]string) string {
	requestString := "/chartas/"
	if id, idOk := parametersMap["id"]; idOk {
		requestString += id + "/?"
	}

	if xPosition, okXPosition := parametersMap["x"]; okXPosition {
		requestString += fmt.Sprintf("x=%s&", xPosition)
	}
	if yPosition, okYPosition := parametersMap["y"]; okYPosition {
		requestString += fmt.Sprintf("y=%s&", yPosition)
	}
	if width, okWidth := parametersMap["width"]; okWidth {
		requestString += fmt.Sprintf("width=%s&", width)
	}
	if height, okHeight := parametersMap["height"]; okHeight {
		requestString += fmt.Sprintf("height=%s&", height)
	}

	return requestString
}

func TestHandler_UpdateBMP(t *testing.T) {
	type mockBehavior func(service *mock_services.MockChartographerServicer, id, xPosition, yPosition, width, height int, receivedImage []byte)

	tests := []struct {
		testName             string
		id                   int
		xPosition            int
		yPosition            int
		width                int
		height               int
		params               map[string]string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			testName:  "OK",
			id:        0,
			xPosition: 0,
			yPosition: 0,
			width:     124,
			height:    124,
			params: map[string]string{
				"id":     "0",
				"x":      "0",
				"y":      "0",
				"width":  "124",
				"height": "124",
			},
			mockBehavior: func(service *mock_services.MockChartographerServicer, id, xPosition, yPosition, width, height int, receivedImage []byte) {
				service.EXPECT().UpdateBMP(id, xPosition, yPosition, width, height, receivedImage).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: ``,
		},
		{
			testName:  "Wrong ID",
			id:        -1,
			xPosition: 0,
			yPosition: 0,
			width:     124,
			height:    124,
			params: map[string]string{
				"id":     "-1",
				"x":      "0",
				"y":      "0",
				"width":  "124",
				"height": "124",
			},
			mockBehavior: func(service *mock_services.MockChartographerServicer, id, xPosition, yPosition, width, height int, receivedImage []byte) {
			},
			expectedStatusCode:   400,
			expectedResponseBody: ``,
		},
		{
			testName:  "Negative width",
			id:        0,
			xPosition: 0,
			yPosition: 0,
			width:     -1,
			height:    124,
			params: map[string]string{
				"id":     "-1",
				"x":      "0",
				"y":      "0",
				"width":  "124",
				"height": "124",
			},
			mockBehavior: func(service *mock_services.MockChartographerServicer, id, xPosition, yPosition, width, height int, receivedImage []byte) {
			},
			expectedStatusCode:   400,
			expectedResponseBody: ``,
		},
		{
			testName:  "Negative height",
			id:        0,
			xPosition: 0,
			yPosition: 0,
			width:     0,
			height:    -1,
			params: map[string]string{
				"id":     "-1",
				"x":      "0",
				"y":      "0",
				"width":  "124",
				"height": "124",
			},
			mockBehavior: func(service *mock_services.MockChartographerServicer, id, xPosition, yPosition, width, height int, receivedImage []byte) {
			},
			expectedStatusCode:   400,
			expectedResponseBody: ``,
		},
		{
			testName:  "Negative xPosition and yPosition",
			id:        0,
			xPosition: -1,
			yPosition: -1,
			width:     124,
			height:    124,
			params: map[string]string{
				"id":     "0",
				"x":      "-1",
				"y":      "-1",
				"width":  "124",
				"height": "124",
			},
			mockBehavior: func(service *mock_services.MockChartographerServicer, id, xPosition, yPosition, width, height int, receivedImage []byte) {
				service.EXPECT().UpdateBMP(id, xPosition, yPosition, width, height, receivedImage).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: ``,
		},
		{
			testName:  "Negative xPosition and zero yPosition",
			id:        0,
			xPosition: -1,
			yPosition: 0,
			width:     124,
			height:    124,
			params: map[string]string{
				"id":     "0",
				"x":      "-1",
				"y":      "0",
				"width":  "124",
				"height": "124",
			},
			mockBehavior: func(service *mock_services.MockChartographerServicer, id, xPosition, yPosition, width, height int, receivedImage []byte) {
				service.EXPECT().UpdateBMP(id, xPosition, yPosition, width, height, receivedImage).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: ``,
		},
		{
			testName:  "Negative yPosition and zero xPosition",
			id:        0,
			xPosition: 0,
			yPosition: -1,
			width:     124,
			height:    124,
			params: map[string]string{
				"id":     "0",
				"x":      "0",
				"y":      "-1",
				"width":  "124",
				"height": "124",
			},
			mockBehavior: func(service *mock_services.MockChartographerServicer, id, xPosition, yPosition, width, height int, receivedImage []byte) {
				service.EXPECT().UpdateBMP(id, xPosition, yPosition, width, height, receivedImage).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: ``,
		},
		{
			testName:  "Positive yPosition and xPosition",
			id:        0,
			xPosition: 10,
			yPosition: 10,
			width:     124,
			height:    124,
			params: map[string]string{
				"id":     "0",
				"x":      "10",
				"y":      "10",
				"width":  "124",
				"height": "124",
			},
			mockBehavior: func(service *mock_services.MockChartographerServicer, id, xPosition, yPosition, width, height int, receivedImage []byte) {
				service.EXPECT().UpdateBMP(id, xPosition, yPosition, width, height, receivedImage).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: ``,
		},
		{
			testName:  "Negative width",
			id:        0,
			xPosition: 10,
			yPosition: 10,
			width:     -100,
			height:    0,
			params: map[string]string{
				"id":     "-1",
				"x":      "0",
				"y":      "0",
				"width":  "124",
				"height": "124",
			},
			mockBehavior: func(service *mock_services.MockChartographerServicer, id, xPosition, yPosition, width, height int, receivedImage []byte) {
			},
			expectedStatusCode:   400,
			expectedResponseBody: ``,
		},
		{
			testName:  "Negative height",
			id:        0,
			xPosition: 10,
			yPosition: 10,
			width:     0,
			height:    -100,
			params: map[string]string{
				"id":     "-1",
				"x":      "0",
				"y":      "0",
				"width":  "124",
				"height": "124",
			},
			mockBehavior: func(service *mock_services.MockChartographerServicer, id, xPosition, yPosition, width, height int, receivedImage []byte) {
			},
			expectedStatusCode:   400,
			expectedResponseBody: ``,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.testName, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			arrayToWrite := []byte{0, 1, 2, 3, 4, 5}
			buffer := bytes.NewBuffer(nil)
			writer := multipart.NewWriter(buffer)
			fw, err := writer.CreateFormFile("upload", "nullptr")
			assert.NoError(t, err)
			_, err = fw.Write(arrayToWrite)
			assert.NoError(t, err)
			err = writer.Close()
			assert.NoError(t, err)

			mockChartService := mock_services.NewMockChartographerServicer(c)
			testCase.mockBehavior(mockChartService, testCase.id, testCase.xPosition, testCase.yPosition, testCase.width, testCase.height, arrayToWrite)
			service := &services.Service{ChartographerServicer: mockChartService}
			controller := &Controller{ChartographerController: NewChartController(service)}

			targetString := createRequestString(testCase.params)

			gin.SetMode(gin.TestMode)
			router := gin.New()
			router.POST("/chartas/:id/", controller.UpdateBMP)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodPost, targetString, buffer)
			request.Header.Set("Content-Type", writer.FormDataContentType())
			router.ServeHTTP(recorder, request)

			assert.Equal(t, testCase.expectedStatusCode, recorder.Code)
			assert.Equal(t, testCase.expectedResponseBody, recorder.Body.String())
		})
	}
}
