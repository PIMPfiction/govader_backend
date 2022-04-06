package govader_backend

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jonreiter/govader"
	echo "github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type ResponseType struct {
	Negative float64 `json:"Negative"`
	Neutral  float64 `json:"Neutral"`
	Positive float64 `json:"Positive"`
	Compound float64 `json:"Compound"`
}

func TestHandler_HandleGetRequest(t *testing.T) {
	tests := []struct {
		name    string
		h       Handler
		wantErr bool
	}{
		{
			"Success",
			Handler{analyzer: govader.NewSentimentIntensityAnalyzer()},
			false,
		},
		{
			"Missing text",
			Handler{analyzer: govader.NewSentimentIntensityAnalyzer()},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Success" {
				httpRecorder := httptest.NewRecorder()
				router := echo.New()
				router.GET("/", tt.h.HandleGetRequest)
				request, err := http.NewRequest("GET", "/?text=I am happy", nil)
				assert.NoError(t, err)
				router.ServeHTTP(httpRecorder, request)
				assert.Equal(t, http.StatusOK, httpRecorder.Code)
				var response ResponseType
				body, err := ioutil.ReadAll(httpRecorder.Result().Body)
				if err != nil {
					t.Error(err)
				}
				_ = json.Unmarshal(body, &response)
				assert.Equal(t, 1, 1)
			}
			if tt.name == "Missing text" {
				httpRecorder := httptest.NewRecorder()
				router := echo.New()
				router.GET("/", tt.h.HandleGetRequest)
				request, err := http.NewRequest("GET", "/?text=", nil)
				assert.NoError(t, err)
				router.ServeHTTP(httpRecorder, request)
				assert.Equal(t, http.StatusBadRequest, httpRecorder.Code)
				var response map[string]string
				body, err := ioutil.ReadAll(httpRecorder.Result().Body)
				if err != nil {
					t.Error(err)
				}
				_ = json.Unmarshal(body, &response)
				assert.Equal(t, response["error"], "?text= parameter is required")

			}
		})
	}
}

func TestHandler_HandlePostRequest(t *testing.T) {
	tests := []struct {
		name    string
		h       Handler
		wantErr bool
	}{
		{
			"Success",
			Handler{analyzer: govader.NewSentimentIntensityAnalyzer()},
			false,
		},
		{
			"Missing text",
			Handler{analyzer: govader.NewSentimentIntensityAnalyzer()},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Success" {
				httpRecorder := httptest.NewRecorder()
				router := echo.New()
				router.POST("/", tt.h.HandlePostRequest)
				reqBody := map[string]string{"text": "I am happy"}
				jsonValue, _ := json.Marshal(reqBody)
				//requestBody := bytes.NewBuffer([]byte(`{"text": "I am happy"}`))
				request, err := http.NewRequest("POST", "/", bytes.NewBuffer(jsonValue))
				request.Header.Add("Content-Type", "application/json")
				assert.NoError(t, err)
				router.ServeHTTP(httpRecorder, request)
				assert.Equal(t, http.StatusOK, httpRecorder.Code)
				var response map[string]interface{}
				body, err := ioutil.ReadAll(httpRecorder.Result().Body)
				if err != nil {
					t.Error(err)
				}
				_ = json.Unmarshal(body, &response)
				assert.Equal(t, response["Negative"], 0.0)
			}
			if tt.name == "Missing text" {
				httpRecorder := httptest.NewRecorder()
				router := echo.New()
				router.POST("/", tt.h.HandlePostRequest)
				_ = bytes.NewBuffer([]byte(`{"text": ""}`))
				request, err := http.NewRequest("POST", "/", nil)
				assert.NoError(t, err)
				router.ServeHTTP(httpRecorder, request)
				assert.Equal(t, http.StatusBadRequest, httpRecorder.Code)
				var response map[string]interface{}
				body, err := ioutil.ReadAll(httpRecorder.Result().Body)
				if err != nil {
					t.Error(err)
				}
				_ = json.Unmarshal(body, &response)
				assert.Equal(t, response["error"], "text required")
			}
		})
	}
}
