package conrtollers

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/pmokeev/chartographer/internal/services"
	"github.com/pmokeev/chartographer/internal/services/mocks"
	"github.com/stretchr/testify/assert"
	"image"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_CreateBMP(t *testing.T) {
	type mockBehavior func(service *mock_services.MockChartographerServicer, width, height int)

	tests := []struct {
		testName             string
		width                int
		height               int
		params               map[string]string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			testName: "OK",
			width:    800,
			height:   800,
			params:   map[string]string{"width": "800", "height": "800"},
			mockBehavior: func(service *mock_services.MockChartographerServicer, width, height int) {
				service.EXPECT().CreateBMP(width, height).Return(0, nil)
			},
			expectedStatusCode:   201,
			expectedResponseBody: `{"id":0}`,
		},
		{
			testName:             "Too big width",
			width:                20001,
			height:               800,
			params:               map[string]string{"width": "20001", "height": "800"},
			mockBehavior:         func(service *mock_services.MockChartographerServicer, width, height int) {},
			expectedStatusCode:   400,
			expectedResponseBody: ``,
		},
		{
			testName:             "Too big height",
			width:                800,
			height:               50001,
			params:               map[string]string{"width": "800", "height": "50001"},
			mockBehavior:         func(service *mock_services.MockChartographerServicer, width, height int) {},
			expectedStatusCode:   400,
			expectedResponseBody: ``,
		},
		{
			testName:             "Too big height and width",
			width:                800,
			height:               50001,
			params:               map[string]string{"width": "20001", "height": "50001"},
			mockBehavior:         func(service *mock_services.MockChartographerServicer, width, height int) {},
			expectedStatusCode:   400,
			expectedResponseBody: ``,
		},
		{
			testName:             "Negative width",
			width:                -1,
			height:               800,
			params:               map[string]string{"width": "-1", "height": "800"},
			mockBehavior:         func(service *mock_services.MockChartographerServicer, width, height int) {},
			expectedStatusCode:   400,
			expectedResponseBody: ``,
		},
		{
			testName:             "Negative height",
			width:                800,
			height:               -1,
			params:               map[string]string{"width": "800", "height": "-1"},
			mockBehavior:         func(service *mock_services.MockChartographerServicer, width, height int) {},
			expectedStatusCode:   400,
			expectedResponseBody: ``,
		},
		{
			testName:             "Negative width and height",
			width:                -1,
			height:               -1,
			params:               map[string]string{"width": "-1", "height": "-1"},
			mockBehavior:         func(service *mock_services.MockChartographerServicer, width, height int) {},
			expectedStatusCode:   400,
			expectedResponseBody: ``,
		},
		{
			testName:             "Zero width",
			width:                0,
			height:               1,
			params:               map[string]string{"width": "0", "height": "1"},
			mockBehavior:         func(service *mock_services.MockChartographerServicer, width, height int) {},
			expectedStatusCode:   400,
			expectedResponseBody: ``,
		},
		{
			testName:             "Zero height",
			width:                1,
			height:               0,
			params:               map[string]string{"width": "1", "height": "0"},
			mockBehavior:         func(service *mock_services.MockChartographerServicer, width, height int) {},
			expectedStatusCode:   400,
			expectedResponseBody: ``,
		},
		{
			testName:             "Zero width and height",
			width:                0,
			height:               0,
			params:               map[string]string{"width": "0", "height": "0"},
			mockBehavior:         func(service *mock_services.MockChartographerServicer, width, height int) {},
			expectedStatusCode:   400,
			expectedResponseBody: ``,
		},
		{
			testName:             "Width is not a integer",
			width:                1,
			height:               1,
			params:               map[string]string{"width": "helloWorld", "height": "800"},
			mockBehavior:         func(service *mock_services.MockChartographerServicer, width, height int) {},
			expectedStatusCode:   400,
			expectedResponseBody: ``,
		},
		{
			testName:             "Height is not a integer",
			width:                1,
			height:               1,
			params:               map[string]string{"width": "800", "height": "helloWorld"},
			mockBehavior:         func(service *mock_services.MockChartographerServicer, width, height int) {},
			expectedStatusCode:   400,
			expectedResponseBody: ``,
		},
		{
			testName:             "Width is negative, height is not a integer",
			width:                -1,
			height:               1,
			params:               map[string]string{"width": "-1", "height": "helloWorld"},
			mockBehavior:         func(service *mock_services.MockChartographerServicer, width, height int) {},
			expectedStatusCode:   400,
			expectedResponseBody: ``,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.testName, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockChartService := mock_services.NewMockChartographerServicer(c)
			testCase.mockBehavior(mockChartService, testCase.width, testCase.height)
			service := &services.Service{ChartographerServicer: mockChartService}
			controller := &Controller{ChartographerController: NewChartController(service)}

			targetString := "/chartas?"
			if widthValue, widthOk := testCase.params["width"]; widthOk {
				targetString += fmt.Sprintf("width=%s&", widthValue)
			}
			if heightValue, heightOk := testCase.params["height"]; heightOk {
				targetString += fmt.Sprintf("height=%s", heightValue)
			}

			gin.SetMode(gin.TestMode)
			router := gin.New()
			router.POST("/chartas", controller.CreateBMP)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodPost, targetString, nil)
			router.ServeHTTP(recorder, request)

			assert.Equal(t, testCase.expectedStatusCode, recorder.Code)
			assert.Equal(t, testCase.expectedResponseBody, recorder.Body.String())
		})
	}
}

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

type TestResponseRecorder struct {
	*httptest.ResponseRecorder
	closeChannel chan bool
}

func (r *TestResponseRecorder) CloseNotify() <-chan bool {
	return r.closeChannel
}

func (r *TestResponseRecorder) closeClient() {
	r.closeChannel <- true
}

func CreateTestResponseRecorder() *TestResponseRecorder {
	return &TestResponseRecorder{
		httptest.NewRecorder(),
		make(chan bool, 1),
	}
}

func TestHandler_GetPartBMP(t *testing.T) {
	type mockBehavior func(service *mock_services.MockChartographerServicer, id, xPosition, yPosition, width, height int)

	tests := []struct {
		testName           string
		id                 int
		xPosition          int
		yPosition          int
		width              int
		height             int
		params             map[string]string
		mockBehavior       mockBehavior
		expectedStatusCode int
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
			mockBehavior: func(service *mock_services.MockChartographerServicer, id, xPosition, yPosition, width, height int) {
				upLeft := image.Point{}
				lowRight := image.Point{X: width, Y: height}
				image := image.NewRGBA(image.Rectangle{Min: upLeft, Max: lowRight})
				service.EXPECT().GetPartBMP(id, xPosition, yPosition, width, height).Return(image, nil)
			},
			expectedStatusCode: 200,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.testName, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockChartService := mock_services.NewMockChartographerServicer(c)
			testCase.mockBehavior(mockChartService, testCase.id, testCase.xPosition, testCase.yPosition, testCase.width, testCase.height)
			service := &services.Service{ChartographerServicer: mockChartService}
			controller := &Controller{ChartographerController: NewChartController(service)}

			targetString := createRequestString(testCase.params)

			gin.SetMode(gin.TestMode)
			router := gin.New()
			router.GET("/chartas/:id/", controller.GetPartBMP)

			recorder := CreateTestResponseRecorder()
			request, _ := http.NewRequest(http.MethodGet, targetString, nil)
			router.ServeHTTP(recorder, request)

			assert.Equal(t, testCase.expectedStatusCode, recorder.Code)
		})
	}
}

func TestHandler_DeleteBMP(t *testing.T) {
	type mockBehavior func(service *mock_services.MockChartographerServicer, id int)

	tests := []struct {
		testName           string
		id                 int
		params             map[string]string
		mockBehavior       mockBehavior
		expectedStatusCode int
	}{
		{
			testName: "OK",
			id:       0,
			params: map[string]string{
				"id": "0",
			},
			mockBehavior: func(service *mock_services.MockChartographerServicer, id int) {
				service.EXPECT().DeleteBMP(id).Return(nil)
			},
			expectedStatusCode: 200,
		},
		{
			testName: "Negative id",
			id:       -1,
			params: map[string]string{
				"id": "-1",
			},
			mockBehavior:       func(service *mock_services.MockChartographerServicer, id int) {},
			expectedStatusCode: 400,
		},
		{
			testName: "ID is not a integer",
			id:       0,
			params: map[string]string{
				"id": "notInteger",
			},
			mockBehavior:       func(service *mock_services.MockChartographerServicer, id int) {},
			expectedStatusCode: 400,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.testName, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			mockChartService := mock_services.NewMockChartographerServicer(c)
			testCase.mockBehavior(mockChartService, testCase.id)
			service := &services.Service{ChartographerServicer: mockChartService}
			controller := &Controller{ChartographerController: NewChartController(service)}

			targetString := createRequestString(testCase.params)

			gin.SetMode(gin.TestMode)
			router := gin.New()
			router.DELETE("/chartas/:id/", controller.DeleteBMP)

			recorder := httptest.NewRecorder()
			request, _ := http.NewRequest(http.MethodDelete, targetString, nil)
			router.ServeHTTP(recorder, request)

			assert.Equal(t, testCase.expectedStatusCode, recorder.Code)
		})
	}
}
