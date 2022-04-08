package govader_backend

//5 minutes
import (
	"sync"
	"time"

	"github.com/jonreiter/govader"
	echo "github.com/labstack/echo/v4"
)

type RequestType struct {
	Text string `json:"text"`
}

type Handler struct {
	analyzer *govader.SentimentIntensityAnalyzer
}

func (h Handler) HandleGetRequest(c echo.Context) error {

	// var request RequestType
	// c.Bind(request)
	text := c.QueryParam("text")
	if text == "" {
		return c.JSON(400, map[string]string{"error": "?text= parameter is required"})
	}
	score := h.analyzer.PolarityScores(text)
	return c.JSON(200, score)
}

func (h Handler) HandlePostRequest(c echo.Context) error {
	request := new(RequestType)
	c.Bind(request)

	if request.Text == "" {
		return c.JSON(400, map[string]string{"error": "text required"})
	}
	score := h.analyzer.PolarityScores(request.Text)
	return c.JSON(200, score)
}

func (h Handler) HandleHealthCheck(c echo.Context) error {
	return c.JSON(200, map[string]string{"status": "ok"})
}

func Serve(e *echo.Echo, portNumber string) error {
	handler := Handler{
		analyzer: govader.NewSentimentIntensityAnalyzer(),
	}
	e.GET("/", handler.HandleGetRequest)
	e.POST("/", handler.HandlePostRequest)
	e.GET("/health", handler.HandleHealthCheck)
	wg := sync.WaitGroup{}
	wg.Add(1)
	var errCh chan error
	go func() {
		err := e.Start(":" + portNumber)
		if err != nil {
			errCh <- err
			wg.Done()
		}
	}()
	time.Sleep(3 * time.Second)
	select {
	case err := <-errCh:
		return err
	default:
		return nil
	}
}
