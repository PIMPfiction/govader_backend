package govader_backend

import (
	"fmt"
	"time"

	"github.com/jonreiter/govader"
	echo "github.com/labstack/echo/v4"
)

//Request type is used to bind the request body to the handler
type RequestType struct {
	Text string `json:"text"`
}

type Handler struct {
	analyzer *govader.SentimentIntensityAnalyzer
}

// Handles get request with query parameter ?text=
// Returns a JSON object with polarity scores
func (h Handler) HandleGetRequest(c echo.Context) error {
	text := c.QueryParam("text")
	if text == "" {
		return c.JSON(400, map[string]string{"error": "?text= parameter is required"})
	}
	score := h.analyzer.PolarityScores(text)
	return c.JSON(200, score)
}

// Handles post request with body parameter text=
// Returns a JSON object with polarity scores
func (h Handler) HandlePostRequest(c echo.Context) error {
	request := new(RequestType)
	c.Bind(request)

	if request.Text == "" {
		return c.JSON(400, map[string]string{"error": "text required"})
	}
	score := h.analyzer.PolarityScores(request.Text)
	return c.JSON(200, score)
}

// Handles Health check for the server
func (h Handler) HandleHealthCheck(c echo.Context) error {
	return c.JSON(200, map[string]string{"status": "ok"})
}

// Serve function binds handler with given echo pointer
// Starts the server on the given port
// Sample usage:
// 	e := echo.New()
// 	err := Serve(e, "8080")
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Scanln()
func Serve(e *echo.Echo, portNumber string) error {
	handler := Handler{
		analyzer: govader.NewSentimentIntensityAnalyzer(),
	}
	e.GET("/", handler.HandleGetRequest)
	e.POST("/", handler.HandlePostRequest)
	e.GET("/health", handler.HandleHealthCheck)
	var errCh = make(chan bool, 1)
	go func() {
		err := e.Start(":" + portNumber)
		if err != nil {
			errCh <- false
		}
	}()
	go func() {
		time.Sleep(time.Second * 3)
		errCh <- true
	}()
	checkErr := <-errCh // blocks until either the server is started or the timeout is reached
	if !checkErr {
		return fmt.Errorf("\033[31m Port is already in use, Change Port\033[0m")
	}
	return nil

}
