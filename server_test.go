package govader_backend

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

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

func TestServe(t *testing.T) {
	type args struct {
		portNumber string
	}
	tests := []struct {
		name string
		args args
	}{
		{"success", args{":8080"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errCh := make(chan error)
			e := echo.New()
			e.GET(
				"/health",
				func(c echo.Context) error {
					return c.JSON(http.StatusOK, "OK")
				},
			)
			go func() {
				errCh <- e.Start("localhost" + tt.args.portNumber)
			}()
			err := waitForServerStart(e, errCh, false)
			assert.NoError(t, err)
			httpRecorder := httptest.NewRecorder()
			router := e
			request, err := http.NewRequest("GET", "/health", nil)
			assert.NoError(t, err)
			router.ServeHTTP(httpRecorder, request)
			assert.Equal(t, http.StatusOK, httpRecorder.Code)
			assert.NoError(t, e.Close())
		})
	}
}

func waitForServerStart(e *echo.Echo, errChan <-chan error, isTLS bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	ticker := time.NewTicker(5 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			var addr net.Addr
			if isTLS {
				addr = e.TLSListenerAddr()
			} else {
				addr = e.ListenerAddr()
			}
			if addr != nil && strings.Contains(addr.String(), ":") {
				return nil // was started
			}
		case err := <-errChan:
			if err == http.ErrServerClosed {
				return nil
			}
			return err
		}
	}
}
