package tests

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/pmokeev/chartographer/internal/conrtollers"
	"github.com/pmokeev/chartographer/internal/services"
	mock_services "github.com/pmokeev/chartographer/tests/mocks"
	"github.com/stretchr/testify/assert"
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
			controller := &conrtollers.Controller{ChartographerController: conrtollers.NewChartController(service)}

			targetString := "/chartas"
			widthValue, widthOk := testCase.params["width"]
			if widthOk {
				targetString += fmt.Sprintf("?width=%s", widthValue)
			}
			if heightValue, heightOk := testCase.params["height"]; heightOk {
				if widthOk {
					targetString += fmt.Sprintf("&height=%s", heightValue)
				} else {
					targetString += fmt.Sprintf("?height=%s", heightValue)
				}
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
